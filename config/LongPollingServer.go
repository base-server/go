package config

import "github.com/heaven-chp/common-library-go/json"

type LongPollingServer struct {
	Address         string `json:"address"`
	Timeout         int    `json:"timeout"`
	ShutdownTimeout uint64 `json:"shutdownTimeout"`

	SubscriptionURI string `json:"subscription_uri"`
	PublishURI      string `json:"publish_uri"`

	FilePersistorInfo struct {
		Use                     bool   `json:"use"`
		FileName                string `json:"file_name"`
		WriteBufferSize         int    `json:"write_buffer_size"`
		WriteFlushPeriodSeconds int    `json:"write_flush_period_seconds"`
	} `json:"file_persistor_info"`

	Log struct {
		Level           string `json:"level"`
		OutputPath      string `json:"output_path"`
		FileNamePrefix  string `json:"file_name_prefix"`
		PrintCallerInfo bool   `json:"print_caller_info"`
		ChannelSize     int    `json:"channel_size"`
	} `json:"log"`
}

func (this *LongPollingServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
