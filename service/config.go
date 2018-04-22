package gomapper

// appConfig ...
type appConfig struct {
	GoMapper GoMapperConfig `json:"gomapper"`
}

// GoMapperConfig ...
type GoMapperConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
