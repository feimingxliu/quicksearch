package core

import (
	"context"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/errors"
	"github.com/feimingxliu/quicksearch/pkg/util/slices"
	"golang.org/x/sync/errgroup"
	"sort"
	"time"
)

//NewSearchOption returns a new empty SearchOption.
func NewSearchOption() *SearchOption {
	return new(SearchOption)
}

type SearchOption struct {
	query   string
	timeout time.Duration
	topN    int
}

func (o *SearchOption) SetQuery(query string) *SearchOption {
	o.query = query
	return o
}

func (o *SearchOption) SetTimeout(timeout time.Duration) *SearchOption {
	o.timeout = timeout
	return o
}

func (o *SearchOption) SetTopN(topN int) *SearchOption {
	o.topN = topN
	return o
}

//Search search for docs with options. TODO: optimize the hit score calculation.
func (index *Index) Search(option *SearchOption) *SearchResult {
	if option == nil {
		return nil
	}
	if err := index.Open(); err != nil {
		return &SearchResult{Error: fmt.Sprintf("%+v", errors.WithStack(err))}
	}
	var (
		startTime = time.Now()                                                      //record search starts.
		isTimeout = false                                                           //if search timeout.
		keywords  = index.tokenizer.KeywordsWeight(option.query, len(option.query)) //keyword with weight in query.
		addition  = slices.DifferenceStr(index.tokenizer.Tokenize(option.query),
			index.tokenizer.Keywords(option.query, len(option.query))) //secondary tokens in query.
		factor     = float64(len(addition)) / float64(len(keywords)+len(addition)) //used to compute score.
		docIDScore = make(map[string]float64)                                      //store unique doc IDs along with their scores.
		pctx       context.Context                                                 //parent context.
		cancel     context.CancelFunc                                              //cancel func.
		c          = make(chan *idScore, len(keywords)*10)                         //chan between docID producer and consumer.
		done       = make(chan struct{}, len(keywords))                            //to inform consumer that a producer has done.
	)
	defer func() {
		close(c)
		close(done)
	}()
	if option.timeout > 0 {
		pctx, cancel = context.WithTimeout(context.Background(), option.timeout)
	} else {
		pctx, cancel = context.Background(), func() {}
	}
	defer cancel()
	g, ctx := errgroup.WithContext(pctx)
	//fetch all related docs' ids.
	for i := range keywords {
		//because i is a temp var.
		j := i
		g.Go(func() error {
			defer func() {
				done <- struct{}{}
			}()
			ids, err := index.GetIDsByKeyword(keywords[j].Word)
			if err != nil {
				return err
			}
			var idx int
			for {
				select {
				case <-ctx.Done():
					isTimeout = true
					return ctx.Err()
				default:
					if idx < len(ids) {
						is := &idScore{
							id:    ids[idx],
							score: keywords[j].Weight,
						}
						c <- is
						idx++
					} else {
						return nil
					}
				}
			}
		})
	}
	//receive docs from chan.
	if len(keywords) > 0 {
		g.Go(func() error {
			doneCount := 0
		Receive:
			for {
				select {
				case is := <-c:
					docIDScore[is.id] += is.score
				case <-done:
					doneCount++
				case <-ctx.Done():
					isTimeout = true
					return ctx.Err()
				default:
					if doneCount == len(keywords) {
						break Receive
					}
				}
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return &SearchResult{Error: fmt.Sprintf("%+v", errors.WithStack(err))}
	}
	//fetch all docs.
	hits := make([]*Hit, 0, len(docIDScore))
	var rerr error
	for docID, score := range docIDScore {
		doc, err := index.RetrieveDocument(docID)
		if err == nil && doc != nil {
			for _, ad := range addition {
				if slices.ContainsStr(doc.KeyWords, ad) {
					score += score / float64(len(keywords)) * factor
				}
			}
			hits = append(hits, &Hit{
				Index:     doc.IndexName,
				Type:      "_doc",
				ID:        doc.ID,
				Score:     score,
				Timestamp: doc.Timestamp,
				Source:    doc.Source,
			})
		} else {
			if err != nil {
				rerr = fmt.Errorf("RetrieveDocument(%s): err: %s\n", docID, err.Error())
			}
		}
	}
	//sort hits in descending order by score.
	sort.Sort(HitSlice(hits))
	searchResult := &SearchResult{
		Took: time.Now().Sub(startTime).String(),
		TimedOut: func() bool {
			if isTimeout {
				return true
			} else {
				return false
			}
		}(),
		MaxScore: func() float64 {
			if len(hits) > 0 {
				return hits[0].Score
			}
			return 0
		}(),
		Hits: Hits{
			Total: Total{Value: len(hits)},
		},
		Error: func() string {
			if rerr != nil {
				return rerr.Error()
			}
			return ""
		}(),
	}
	if option.topN > 0 && option.topN < len(hits) {
		searchResult.Hits.Hits = hits[:option.topN]
	} else {
		searchResult.Hits.Hits = hits
	}
	return searchResult
}

type idScore struct {
	id    string
	score float64
}
