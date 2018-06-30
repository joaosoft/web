package mapper

// appConfig ...
type appConfig struct {
	GoMapper MapperConfig `json:"mapper"`
}

// MapperConfig ...
type MapperConfig struct {
	Log struct {
		Level string `json:"level"`
	} `json:"log"`
}
