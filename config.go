package web

type AppConfig struct {
	Server ServerConfig `json:"Server"`
	Client ClientConfig `json:"Client"`
}
