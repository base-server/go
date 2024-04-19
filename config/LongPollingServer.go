package config

type LongPollingServer struct {
	Address         string `json:"address"`
	TimeoutSeconds  int    `json:"timeoutSeconds"`
	ShutdownTimeout string `json:"shutdownTimeout"`

	SubscriptionURI string `json:"subscriptionUri"`
	PublishURI      string `json:"publishUri"`

	FilePersistorInfo struct {
		Use                     bool   `json:"use"`
		FileName                string `json:"fileName"`
		WriteBufferSize         int    `json:"writeBufferSize"`
		WriteFlushPeriodSeconds int    `json:"writeFlushPeriodSeconds"`
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
