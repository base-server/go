package config

type SocketServer struct {
	Address        string `json:"address"`
	ClientPoolSize int    `json:"client_pool_size"`

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
