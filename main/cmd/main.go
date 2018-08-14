package main

import (
	"db-migration/services"
	"db-migration/services/cmd"
	"flag"

	_ "github.com/lib/pq"
)

func main() {
	var migrateOption string
	var migrateNumber int

	flag.StringVar(&migrateOption, string(services.MigrationCmd), string(services.MigrationOptionUp), "Runs the specified command. Valid options are: `"+string(services.MigrationOptionUp)+"`, `"+string(services.MigrationOptionDown)+"`.")
	flag.IntVar(&migrateNumber, string(services.MigrationNumberCmd), 0, "Runs the specified command.")
	flag.Parse()

	m, err := cmd.NewService()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}

	if err := m.Execute(services.MigrationOption(migrateOption), migrateNumber); err != nil {
		panic(err)
	}
}
