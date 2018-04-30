package gomapper

// appConfig ...
type appConfig struct {
	GoMapper MapperConfig `json:"gomapper"`
}

// MapperConfig ...
type MapperConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
