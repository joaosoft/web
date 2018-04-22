package gomapper

// appConfig ...
type appConfig struct {
	GoMapper goMapperConfig `json:"gomapper"`
}

// goMapperConfig ...
type goMapperConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
