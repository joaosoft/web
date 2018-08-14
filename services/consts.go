package services

const (
	DefaultURL = "http://localhost:8001"

	TagMigrationUp   MigrationTag = "-- migrate up"
	TagMigrationDown MigrationTag = "-- migrate down"

	MigrationCmd        MigrationCommand = "migrate"
	MigrationNumberCmd  MigrationCommand = "number"
	MigrationOptionUp   MigrationOption  = "up"
	MigrationOptionDown MigrationOption  = "down"
)
