package config

type CloudEventsServer struct {
	Address         string `json:"address" `
	ShutdownTimeout string `json:"shutdownTimeout"`
}
