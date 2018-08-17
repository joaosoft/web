package cmd

const (
	CmdMigrate MigrationCommand = "migrate"
	CmdNumber  MigrationCommand = "number"

	OptionUp   MigrationOption = "up"
	OptionDown MigrationOption = "down"

	FileTagMigrate     Tag = "-- migrate %s"
	FileTagMigrateUp   Tag = "-- migrate up"
	FileTagMigrateDown Tag = "-- migrate down"

	FileTagCustom     Tag = "-- %s %s"
	FileTagCustomUp   Tag = "-- %s up"
	FileTagCustomDown Tag = "-- %s down"
)
