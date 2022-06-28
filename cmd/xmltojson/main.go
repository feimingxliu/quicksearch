package main

import (
	"compress/gzip"
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"github.com/go-creed/sat"
	"log"
	"net/url"
	"os"
)

var (
	inputXmlGz = "test/testdata/zhwiki-20220601-abstract.xml.gz"
	outputJson = "test/testdata/zhwiki-20220601-abstract.json"
)

func init() {
	inputXmlGz = *flag.String("input", inputXmlGz, "the input xml file")
	outputJson = *flag.String("output", outputJson, "the output json file")
	flag.Parse()
}

func main() {
	log.Println("loading xml.gz source file ......")
	inputFile, err := os.OpenFile(inputXmlGz, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()
	log.Println("Creating bulk json output file ......")
	outputFile, err := os.OpenFile(outputJson, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()
	log.Println("Processing ......")
	r, err := gzip.NewReader(inputFile)
	if err != nil {
		log.Fatalln(err)
	}
	doc := new(document)
	dicter := sat.DefaultDict()
	var buf []byte
	var count uint
	xmlDecoder := xml.NewDecoder(r)
	// use stream to process.
	for {
		token, _ := xmlDecoder.Token()
		if token == nil {
			break
		}
		switch token := token.(type) {
		case xml.StartElement:
			if token.Name.Local == "doc" {
				count++
				err = xmlDecoder.DecodeElement(doc, &token)
				if err != nil {
					log.Fatalln(err)
				}
				doc.ID = uuid.GetUUID()
				doc.URL, _ = url.QueryUnescape(doc.URL)
				doc.URL = dicter.Read(doc.URL)
				doc.Title = dicter.Read(doc.Title)
				doc.Text = dicter.Read(doc.Text)
				buf, _ = json.Marshal(doc)
				_, _ = outputFile.Write(buf)
				_, _ = outputFile.WriteString("\n")
				fmt.Printf("Total proccessed item: %d\r", count)
			}
		}
	}
	err = outputFile.Sync()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Complete successfully!")
}

type document struct {
	ID    string `json:"id"`
	Title string `xml:"title" json:"title"`
	URL   string `xml:"url" json:"url"`
	Text  string `xml:"abstract" json:"text"`
}
