# QuickSearch

## Contents

+ [Introduction](#introduction)
+ [Getting Started](#getting-started)
+ [API Reference](#api-reference)
+ [Run or build from source](#run-or-build-from-source)
+ [Tests](#tests)

### Introduction

Quicksearch is a lightweight search engine which deploys and runs as a single binary. It's inspired
by [Zinc](https://github.com/zinclabs/zinc) but it uses [Bleve](https://github.com/blevesearch/bleve) as underlying
indexing library and supports chinese by default. It just supports local storage now.

### Getting Started

If you have installed [golang](https://golang.google.cn/learn/), just run following:

 ```bash
// if go version >= 1.16
go install github.com/feimingxliu/quicksearch/cmd/quicksearch@latest
// or `go get -u github.com/feimingxliu/quicksearch/cmd/quicksearch` for go < 1.16
 ```

This will generate binary in your `$GOPATH/bin`, note this **does not install UI**. Or you can get the prebuilt binary
from [Releases](https://github.com/feimingxliu/quicksearch/releases) *which includes UI*.

To run the quicksearch. Copy the example config.

```bash
wget https://raw.githubusercontent.com/feimingxliu/quicksearch/master/configs/config.yaml
// run quicksearch
quicksearch -c config.yaml
```

Or use [Docker](https://docs.docker.com/get-docker/) **(includes UI)**

```bash
docker run -d -p 9200:9200 feimingxliu/quicksearch
```

Quicksearch server will listen on [:9200](http://localhost:9200) by default, open [:9200](http://localhost:9200)
to view the UI.
If you enable the `http.auth` in [config.yaml](configs/config.yaml), login with `admin:admin` by default.
You can change the username `http.auth.username` directly in [config.yaml](configs/config.yaml), but change password follows below
```sh
go run github.com/feimingxliu/quicksearch/cmd/bcrypt -p $YourPassword
```
the above command will generate the bcrypt password hash, copy it to the `http.auth.password` in [config.yaml](configs/config.yaml).
In this case, **note** the `http.auth.enabled` must be **true**.

### API Reference

#### Index API

+ *Create Index*

```
POST /<index>
{
    "settings": <Index Settings>,
    "mappings": <Index Mappings>
}
```

The request body can be ignored which use default.
`<Index Settings>` is an object which contains index's setting

```
{
    "number_of_shards": int 
}
```

`<Index Mappings>` is an object which defines index's mapping

```

{
	"types": {
		"Document Type": <Documnet Mapping>, 
		....
		"Document Type": <Documnet Mapping>
	}
	"default_mapping": <Documnet Mapping>,
	"type_field": string,
	"default_type": string,
	"default_analyzer": string 
}

```

`<Documnet Mapping>`  is an object which defines document's mapping

```

{
	"disabled": bool,	# disable this documnet mapping
	"properties": {
			"name": <Documnet Mapping>,	# this enables nested json
			...
			"name": <Documnet Mapping>
		},
	"fields": <Field Mapping>,
	"default_analyzer": string 
}
```

`<Field Mapping>` is an array which defines field level mapping

```

[
	# you can define one more field mapping for one field, that's why its array {
	"type": string,	# support "keyword", "text", "datetime", "number", "boolean", "geopoint", "IP"
	"analyzer": string,	# specifies the name of the analyzer to use for this field
	"store": bool,	# indicates whether to store field values in the index
	"index": bool	# indicates whether to analyze the field }
]

```

+ *Update Index Mapping*

```
PUT /<index>/_mapping
<Index Mapping>
```

+ *Get Index Detail*

```
GET /<index>
```

+ *Open Index*

```
POST /<index>/_open
```

+ *Close Index*

```
POST /<index>/_close
```

+ *Clone Index*

```
POST /<index>/_clone/<cloned index>
```

+ *List Indices*

```
GET /_all
```

+ *Delete Index*

```
DELETE /<index>
```

#### Document API

+ *Index Document*

```
POST /<index>/_doc
<document json object>
# or with custom documnet id
POST /<index>/_doc/<docID>
<document json object>
```

  If index a document with same docID, the newer one will cover old fully.

+ *Bulk*

```
POST /_bulk 
POST /<index>/_bulk
<Action Line>
<optional document json object>
......
<Action Line>
<optional document json object>
```

`<Action Line>` is an object defines which operation to execute.

```
{ 
	# <Action> can be `create`, `delete`, `index`, `update`
	<Action>: {
		"_index": string,
		"_id": string 
	} 
}
```

+ *Update Document*

```
PUT /<index>/_doc/<docID>
{
	"fieldName": any 
}
```

  This can update part fields of document.

+ *Get Document*

```
GET /<index>/_doc/<docID>
```

+ *Delete Document*

```
DELETE /<index>/_doc/<docID>
```

#### Search API

```
POST /<index>/_search
GET /<index>/_search
POST /_search 
GET /_search 
{ 
	"query": <Query>,
	"size": int,
	"from": int,
	"highlight": []string, # fields to highlight
	"fields": []string,
	"facets": {
		<facet name>: {
            "size": int,
            "field": string,
            "numeric_ranges": [
                {
                    "name": string,
                    "min": float64,
                    "max": float64
                }
            ],
            "date_ranges": [
                {
                    "name": string,
                    "start": datetime, # RFC3339
                    "end": datetime # RFC3339 
                }
            ]
         } 
     },
     "explain": bool,
     "sort": []sting,
     "includeLocations": bool,
     "search_after": []sting,
     "search_before": []string 
}
```

`<Query>` indicates different query, see [Queries](http://blevesearch.com/docs/Query/).

+ *QueryStringQuery* is the simplest query for search,
  see  [full query language specification](http://blevesearch.com/docs/Query-String-Query/).

```
{
  "query": string,
  "boost": float64 
}
```

+ *TermQuery* 

  A term query is the simplest possible query. It performs an exact match in the index for the provided term.

  Most of the time users should use a Match Query instead.

```
{
  "term": string,
  "field": string
}
```

+ *MatchQuery*

  A match query is like a term query, but the input text is analyzed first. An attempt is made to use the same analyzer that was used when the field was indexed.

  The match query can optionally perform fuzzy matching. If the fuzziness parameter is set to a non-zero integer the analyzed text will be matched with the specified level of fuzziness. Also, the prefix_length parameter can be used to require that the term also have the same prefix of the specified length.

```
{
  "match": string,
  "field": string,
  "analyzer": string,
  "boost": float64,
  "prefix_length": int,
  "fuzziness": int,
  "operator": string # "and" or "or"
}
```

+ *PhraseQuery*

  A phrase query searches for terms occurring in the specified position and offsets.

  The phrase query is performing an exact term match for all the phrase constituents. If you want the phrase to be analyzed, consider using the Match Phrase Query instead.

```
{
  "terms": []string,
  "field": string,
  "boost": float64
}
```

+ *MatchPhraseQuery*

  The match phrase query is like the phrase query, but the input text is analyzed and a phrase query is built with the terms resulting from the analysis. 

```
{
  "match_phrase": string,
  "analyzer": string,
  "field": string,
  "boost": float64
}
```

+ *PrefixQuery*

  The prefix query finds documents containing terms that start with the provided prefix. 

```
{
  "prefix": string,
  "field": string,
  "boost": float64
}
```

+ *FuzzyQuery*

  A fuzzy query is a term query that matches terms within a specified edit distance (Levenshtein distance). Also, you can optionally specify that the term must have a matching prefix of the specified length. 

```
{
  "term": string,
  "field": string,
  "boost": float64,
  "prefix_length": int,
  "fuzziness": int
}
```

+ *ConjunctionQuery*

  The conjunction query is a compound query. Result documents must satisfy all of the child queries. 

```
{
  "conjuncts": []<Query>,
  "boost": 1
}
```

+ *DisjunctionQuery*

  The disjunction query is a compound query. Result documents must satisfy a configurable `min` number of child queries. By default this `min` is set to 1. 

```
{
  "disjuncts": []<Query>,
  "boost": 1,
  "min": 1
}
```

+ *BooleanQuery*

  The boolean query is useful combination of conjunction and disjunction queries. The query takes three lists of queries:

  - must - result documents must satisfy all of these queries
  - should - result documents should satisfy at least `minShould` of these queries
  - must not - result documents must not satisfy any of these queries

  The `minShould` value is configurable, defaults to 0.

```
{
  "must": <Query>, # must be *ConjunctionQuery* or *DisjunctionQuery*
  "should": <Query>, # must be *ConjunctionQuery* or *DisjunctionQuery*
  "must_not": <Query>, # must be *ConjunctionQuery* or *DisjunctionQuery*
  "boost": 1
}
```

+ *NumericRangeQuery*

  The numeric range query finds documents containing a numeric value in the specified field within the specified range. You can omit one endpoint, but not both. The `inclusiveMin` and `inclusiveMax` properties control whether or not the end points are included or excluded. 

```
{
  "min": float64,
  "max": float64,
  "inclusive_min": bool,
  "inclusive_max": bool,
  "field": string,
  "boost": float64
}
```

+ *DateRangeQuery*

  The date range query finds documents containing a date value in the specified field within the specified range. You can omit one endpoint, but not both. The inclusiveStart and inclusiveEnd properties control whether or not the end points are included or excluded. 

```
{
  "start": datetime,
  "end": datetime,
  "inclusive_start": bool,
  "inclusive_end": bool,
  "field": string,
  "boost": float64
}
```

+ *MatchAllQuery*

  The match all query will match all documents in the index.

```
{
  "match_all": {},
  "boost": float64
}
```

+ *MatchNoneQuery*

  The match none query will not match any documents in the index. 

```
{
  "match_none": {},
  "boost": float64
}
```

+ *DocIDQuery*

  The doc ID query will match only documents that contain one of the supplied document identifiers. 

```
{
  "ids": []string,
  "boost": 1
}
```

### Run or build from source

To run the `quicksearch` from source, clone the repo firstly.

```sh
git clone --recurse-submodules git@github.com:feimingxliu/quicksearch.git
# or use 'git clone --recurse-submodules https://github.com/feimingxliu/quicksearch.git' if you don't set SSH key.
```

Then download the dependencies.

```sh
cd quicksearch && go mod tidy -compat=1.17 # go version >= 1.17
```

Build the frontend

```sh
yarn --cwd web/web && yarn --cwd web/web build
```

Run the following command to start the `quicksearch`.

```sh
go run github.com/feimingxliu/quicksearch/cmd/quicksearch -c configs/config.yaml
```

Or build the binary like this:

```sh
go build -o bin/quicksearch github.com/feimingxliu/quicksearch/cmd/quicksearch
```

Run binary:

```sh
bin/quicksearch -c configs/config.yaml
```

### Tests

The tests use some testdata which stores with [git-lfs](https://git-lfs.github.com/). After you have installed
the [git-lfs](https://git-lfs.github.com/), you can run

```sh
git lfs pull
```

in the project root to fetch the large test file.

Then run

```sh
go run github.com/feimingxliu/quicksearch/cmd/xmltojson
```

The above command will generate `test/testdata/zhwiki-20220601-abstract.json`, you can open it to see the content.

In the end, just run all the tests by

```sh
go test -timeout 0 ./...
```

If everything works well, an `ok` will appear at the end of output.