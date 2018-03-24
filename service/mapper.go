package gomapper

// GoMapper ...
type GoMapper struct{}

// NewMapper ...
func NewMapper(options ...GoMapperOption) *GoMapper {
	gomapper := &GoMapper{}

	gomapper.Reconfigure(options...)

	return gomapper
}