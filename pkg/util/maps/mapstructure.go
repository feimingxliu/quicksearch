package maps

import "github.com/mitchellh/mapstructure"

// MapToStruct converts a map to a struct, 'toStruct' must be a pointer.
func MapToStruct(fromMap map[string]interface{}, toStruct interface{}) error {
	return mapstructure.Decode(fromMap, toStruct)
}
