package cmd

import "database/sql"

type Tag string
type MigrationCommand string
type MigrationOption string
type CustomMode string

type Handler func(option MigrationOption, tx *sql.Tx, data string) error
