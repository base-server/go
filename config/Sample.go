package config

import "github.com/heaven-chp/common-library-go/json"

type Sample struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`
}

func (this *Sample) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
