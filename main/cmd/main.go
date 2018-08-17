package main

import (
	"db-migration/services/cmd"
	"flag"

	"fmt"

	"database/sql"
)

func main() {
	var cmdMigrate string
	var cmdNumber int

	flag.StringVar(&cmdMigrate, string(cmd.CmdMigrate), string(cmd.OptionUp), "Runs the specified command. Valid options are: `"+string(cmd.OptionUp)+"`, `"+string(cmd.OptionDown)+"`.")
	flag.IntVar(&cmdNumber, string(cmd.CmdNumber), 0, "Runs the specified command.")
	flag.Parse()

	m, err := cmd.NewService()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}

	m.AddTag("custom", CustomHandler)
	if executed, err := m.Execute(cmd.MigrationOption(cmdMigrate), cmdNumber); err != nil {
		panic(err)
	} else {
		fmt.Printf("EXECUTED: %d", executed)
	}
}

func CustomHandler(option cmd.MigrationOption, tx *sql.Tx, data string) error {
	fmt.Printf("\nexecuting with option '%s' and data '%s", option, data)
	return nil
}
