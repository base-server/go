package config

import "time"

type LongPollingServer struct {
	Address         string        `json:"address"`
	Timeout         int           `json:"timeout"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`

	SubscriptionURI string `json:"subscription_uri"`
	PublishURI      string `json:"publish_uri"`

	FilePersistorInfo struct {
		Use                     bool   `json:"use"`
		FileName                string `json:"file_name"`
		WriteBufferSize         int    `json:"write_buffer_size"`
		WriteFlushPeriodSeconds int    `json:"write_flush_period_seconds"`
	} `json:"file_persistor_info"`

	Log struct {
		Level  string `json:"level"`
		Output string `json:"output"`
		File   struct {
			Name          string `json:"name"`
			ExtensionName string `json:"extensionName"`
			AddDate       bool   `json:"addDate"`
		} `json:"file"`
		WithCallerInfo bool `json:"withCallerInfo"`
	} `json:"log"`
}
