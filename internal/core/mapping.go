package core

import (
	"fmt"
	"github.com/blevesearch/bleve/v2"
	imapping "github.com/blevesearch/bleve/v2/mapping"
	"github.com/feimingxliu/quicksearch/pkg/util/maps"
	"strings"
)

type IndexMapping struct {
	TypeMapping     map[string]*DocumentMapping `json:"types,omitempty" mapstructure:"types"`
	DefaultMapping  *DocumentMapping            `json:"default_mapping" mapstructure:"default_mapping"`
	TypeField       *string                     `json:"type_field" mapstructure:"type_field"`
	DefaultType     *string                     `json:"default_type" mapstructure:"default_type"`
	DefaultAnalyzer *string                     `json:"default_analyzer" mapstructure:"default_analyzer"` // standard
}

type DocumentMapping struct {
	Disabled        bool                        `json:"disabled" mapstructure:"disabled"`
	Properties      map[string]*DocumentMapping `json:"properties,omitempty" mapstructure:"properties"`
	Fields          []*FieldMapping             `json:"fields,omitempty" mapstructure:"fields"`
	DefaultAnalyzer string                      `json:"default_analyzer,omitempty" mapstructure:"default_analyzer"`
}

type FieldMapping struct {
	Type     string  `json:"type,omitempty" mapstructure:"type"`         // support "keyword", "text", "datetime", "number", "boolean", "geopoint", "IP"
	Analyzer *string `json:"analyzer,omitempty" mapstructure:"analyzer"` // Analyzer specifies the name of the analyzer to use for this field.
	Store    *bool   `json:"store,omitempty" mapstructure:"store"`       // Store indicates whether to store field values in the index.
	Index    *bool   `json:"index,omitempty" mapstructure:"index"`       // Store indicates whether to analyze the field.
}

func BuildIndexMappingFromMap(m map[string]interface{}) (*IndexMapping, error) {
	if m == nil {
		return nil, nil
	}
	im := new(IndexMapping)
	if err := maps.MapToStruct(m, im); err != nil {
		return nil, err
	} else {
		return im, nil
	}
}

func buildIndexMapping(im *IndexMapping) (imapping.IndexMapping, error) {
	indexMapping := bleve.NewIndexMapping()
	if im == nil {
		return indexMapping, nil
	}
	typeMapping := make(map[string]*imapping.DocumentMapping, len(im.TypeMapping))
	for t, dm := range im.TypeMapping {
		if idm, err := buildDocumentMapping(dm); err != nil {
			return nil, err
		} else {
			typeMapping[t] = idm
		}
	}
	indexMapping.TypeMapping = typeMapping
	if ddm, err := buildDocumentMapping(im.DefaultMapping); err != nil {
		return nil, err
	} else {
		indexMapping.DefaultMapping = ddm
	}
	if im.TypeField != nil {
		indexMapping.TypeField = *im.TypeField
	}
	if im.DefaultType != nil {
		indexMapping.DefaultType = *im.DefaultType
	}
	if im.DefaultAnalyzer != nil {
		if err := setDefaultAnalyzerForMapping(indexMapping, *im.DefaultAnalyzer, nil); err != nil {
			return nil, err
		}
	}
	return indexMapping, nil
}

func buildDocumentMapping(dm *DocumentMapping) (*imapping.DocumentMapping, error) {
	documentMapping := bleve.NewDocumentMapping()
	if dm == nil {
		return documentMapping, nil
	}
	properties := make(map[string]*imapping.DocumentMapping, len(dm.Properties))
	for name, vdm := range dm.Properties {
		if rdm, err := buildDocumentMapping(vdm); err != nil {
			return nil, err
		} else {
			properties[name] = rdm
		}
	}
	fields := make([]*imapping.FieldMapping, 0, len(dm.Fields))
	for _, fm := range dm.Fields {
		if ifm, err := buildFieldMapping(fm); err != nil {
			return nil, err
		} else {
			fields = append(fields, ifm)
		}
	}
	documentMapping.Enabled = !dm.Disabled
	documentMapping.Properties = properties
	documentMapping.Fields = fields
	documentMapping.DefaultAnalyzer = dm.DefaultAnalyzer
	return documentMapping, nil
}

func buildFieldMapping(fm *FieldMapping) (*imapping.FieldMapping, error) {
	fieldMapping := new(imapping.FieldMapping)
	fm.Type = strings.ToLower(fm.Type)
	switch fm.Type {
	case "keyword":
		fieldMapping = bleve.NewKeywordFieldMapping()
	case "text":
		fieldMapping = bleve.NewTextFieldMapping()
	case "datetime":
		fieldMapping = bleve.NewDateTimeFieldMapping()
	case "number":
		fieldMapping = bleve.NewNumericFieldMapping()
	case "boolean":
		fieldMapping = bleve.NewBooleanFieldMapping()
	case "geopoint":
		fieldMapping = bleve.NewGeoPointFieldMapping()
	case "ip":
		fieldMapping = bleve.NewIPFieldMapping()
	default:
		return nil, fmt.Errorf("unknown field type [%s]", fm.Type)
	}
	if fm.Analyzer != nil {
		fieldMapping.Analyzer = *fm.Analyzer
	}
	if fm.Store != nil {
		fieldMapping.Store = *fm.Store
	}
	if fm.Index != nil {
		fieldMapping.Index = *fm.Index
	}
	return fieldMapping, nil
}

func setDefaultAnalyzerForMapping(mapping *imapping.IndexMappingImpl, analyzer string, config map[string]interface{}) error {
	if mapping == nil {
		mapping = bleve.NewIndexMapping()
	}
	switch strings.ToLower(analyzer) {
	case "default", "standard":
		mapping.DefaultAnalyzer = "standard"
		return nil
	//case "jieba":
	//	err = mapping.AddCustomTokenizer("gojieba",
	//		map[string]interface{}{
	//			"dict_path":      "",
	//			"hmm_path":       "",
	//			"user_dict_path": "",
	//			"idf":            "",
	//			"stop_words":     "",
	//			"type":           "gojieba",
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	err = mapping.AddCustomAnalyzer("gojieba",
	//		map[string]interface{}{
	//			"type":      "gojieba",
	//			"tokenizer": "gojieba",
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	mapping.DefaultAnalyzer = "gojieba"
	//	return nil
	//case "gse":
	//	err = mapping.AddCustomTokenizer("gse",
	//		map[string]interface{}{
	//			"type":       "gse",
	//			"dict_path":  "",
	//			"stop_words": "",
	//			"alpha":      false,
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	err = mapping.AddCustomAnalyzer("gse",
	//		map[string]interface{}{
	//			"type":      "gse",
	//			"tokenizer": "gse",
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	mapping.DefaultAnalyzer = "gse"
	//	return nil
	//case "sego":
	//	err = mapping.AddCustomTokenizer("sego",
	//		map[string]interface{}{
	//			"type":      "sego",
	//			"dict_path": "",
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	err = mapping.AddCustomAnalyzer("sego",
	//		map[string]interface{}{
	//			"type":      "sego",
	//			"tokenizer": "sego",
	//		},
	//	)
	//	if err != nil {
	//		return err
	//	}
	//	mapping.DefaultAnalyzer = "sego"
	//	return nil
	default:
		mapping.DefaultAnalyzer = "standard"
		return nil
	}
}
