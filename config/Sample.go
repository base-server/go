package config

import "github.com/heaven-chp/common-library-go/json"

type Sample struct {
	Field1 int    `json:"field1"`
	Field2 string `json:"field2"`
}

func (this *Sample) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
