package main

import (
	"encoding/xml"
	"flag"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"github.com/feimingxliu/quicksearch/pkg/util/uuid"
	"log"
	"os"
)

var (
	inputXml   = "test/testdata/zhwiki-20220601-abstract.xml"
	outputJson = "test/testdata/zhwiki-20220601-abstract.json"
)

func init() {
	inputXml = *flag.String("input", inputXml, "the input xml file")
	outputJson = *flag.String("output", outputJson, "the output json file")
	flag.Parse()
}

func main() {
	log.Println("loading xml source file ......")
	inputFile, err := os.OpenFile(inputXml, os.O_RDONLY, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer inputFile.Close()
	xmlDecoder := xml.NewDecoder(inputFile)
	dump := struct {
		Documents []document `xml:"doc"`
	}{}
	log.Println("Decoding xml into struct ......")
	err = xmlDecoder.Decode(&dump)
	docs := dump.Documents
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Adding additional ID field ......")
	for i := range docs {
		docs[i].ID = uuid.GetUUID()
	}
	log.Println("Creating json output file ......")
	outputFile, err := os.OpenFile(outputJson, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer outputFile.Close()
	jsonEncoder := json.NewEncoder(outputFile)
	jsonEncoder.SetIndent("", " ")
	log.Println("Writing data into json file ......")
	err = jsonEncoder.Encode(docs)
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
