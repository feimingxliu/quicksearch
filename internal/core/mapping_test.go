package core

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	"github.com/feimingxliu/quicksearch/pkg/util/json"
	"testing"
)

const rawJson = `
{
    "input": {
      "type": "container"
    },
	"log": {
      "offset": 36554067,
      "file": {
        "path": "/var/lib/docker/containers/a41811b2f50326ba5b8c5fb14e79e0d89b714a9970df5fffebdcede02aa7141e/a41811b2f50326ba5b8c5fb14e79e0d89b714a9970df5fffebdcede02aa7141e-json.log"
      }
    },
    "message": "I0707 09:57:35.410132       1 log.go:184] http: TLS handshake error from 100.127.21.195:50315: read tcp 192.168.33.193:10004->100.127.21.195:50315: read: connection timed out",
    "timestamp": "2022-07-07T01:57:35.410Z"
}`

func TestBuildIndexMappingFromMap(t *testing.T) {
	rawMapping := `{
    "default_mapping": {
     "properties": {
      "input": {
    "properties": {
     "type": {
    "fields": [
     {
      "type": "keyword"
     }
    ]
   }
    }
   },
      "log": {
    "properties": {
     "file": {
    "properties": {
     "path": {
    "fields": [
     {
      "type": "text"
     }
    ]
   }
    }
   },
     "offset": {
    "fields": [
     {
      "type": "number"
     }
    ]
   }
    }
   },
      "message": {
    "fields": [
     {
      "type": "text"
     }
    ]
   },
      "timestamp": {
    "fields": [
     {
      "type": "datetime"
     }
    ]
   }
     }
    },
    "type_field": "_type",
    "default_type": "_default",
    "default_analyzer": ""
   }`
	m := make(map[string]interface{})
	if err := json.Unmarshal([]byte(rawMapping), &m); err != nil {
		t.Fatal(err)
	}
	im, err := BuildIndexMappingFromMap(m)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("internal IndexMapping", im)
	indexMapping, err := buildIndexMapping(im)
	if err != nil {
		t.Fatal(err)
	}
	json.Print("built bleve index mapping", indexMapping)

	mapping := bleve.NewIndexMapping()
	docMapping := bleve.NewDocumentMapping()
	inputFieldMapping := bleve.NewDocumentMapping()
	inputFieldMapping.AddFieldMappingsAt("type", bleve.NewKeywordFieldMapping())
	docMapping.AddSubDocumentMapping("input", inputFieldMapping)
	logFieldMapping := bleve.NewDocumentMapping()
	logFieldMapping.AddFieldMappingsAt("offset", bleve.NewNumericFieldMapping())
	fileFieldMapping := bleve.NewDocumentMapping()
	fileFieldMapping.AddFieldMappingsAt("path", bleve.NewTextFieldMapping())
	logFieldMapping.AddSubDocumentMapping("file", fileFieldMapping)
	docMapping.AddSubDocumentMapping("log", logFieldMapping)
	docMapping.AddFieldMappingsAt("message", bleve.NewTextFieldMapping())
	docMapping.AddFieldMappingsAt("timestamp", bleve.NewDateTimeFieldMapping())
	mapping.DefaultMapping = docMapping
	json.Print("api bleve mapping", mapping)

	m1, _ := json.Marshal(indexMapping)
	m2, _ := json.Marshal(mapping)
	fmt.Println("Is built mapping equals to api generated ?", string(m1) == string(m2))
}
