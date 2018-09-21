package web

type AppConfig struct {
	Server ServerConfig `json:"server"`
	Client ClientConfig `json:"client"`
}
