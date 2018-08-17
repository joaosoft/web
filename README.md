# dbmigration
[![Build Status](https://travis-ci.org/joaosoft/dbmigration.svg?branch=master)](https://travis-ci.org/joaosoft/dbmigration) | [![codecov](https://codecov.io/gh/joaosoft/dbmigration/branch/master/graph/badge.svg)](https://codecov.io/gh/joaosoft/dbmigration) | [![Go Report Card](https://goreportcard.com/badge/github.com/joaosoft/dbmigration)](https://goreportcard.com/report/github.com/joaosoft/dbmigration) | [![GoDoc](https://godoc.org/github.com/joaosoft/dbmigration?status.svg)](https://godoc.org/github.com/joaosoft/dbmigration)

A simple database migration tool to integrate in your projects

###### If i miss something or you have something interesting, please be part of this project. Let me know! My contact is at the end.

## With support for
* Migration Up and Down options
* Custom Tags with handler to Up and Down options
* Postgres

## Dependecy Management 
>### Dep

Project dependencies are managed using Dep. Read more about [Dep](https://github.com/golang/dep).
* Install dependencies: `dep ensure`
* Update dependencies: `dep ensure -update`


>### Go
```
go get github.com/joaosoft/dbmigration
```

## Usage 
This examples are available in the project at [dbmigration/Makefile](https://github.com/joaosoft/dbmigration/tree/master/Makefile)
```
import (
	github.com/joaosoft/dbmigration
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
```


Configuration file
```
{
  "dbmigration": {
    "host": "localhost:8001",
    "path": "schema/db/postgres/example",
    "db": {
      "driver": "postgres",
      "datasource": "postgres://postgres:postgres@localhost:5432?sslmode=disable"
    },
    "log": {
      "level": "info"
    }
  },
  "manager": {
    "log": {
      "level": "error"
    }
  }
}
```

Migration file example
```
-- migrate up
CREATE TABLE dbmigration.test1();

-- custom up
teste do joao A
teste do joao B



-- migrate down
DROP TABLE dbmigration.test1;

-- custom down
teste do joao 1
teste do joao 2
```

Migration commands
```
// migrate up all migrations
dbmigration -migrate up

// migrate up 2 migrations
dbmigration -migrate up -number 2

// migrate down one migration
dbmigration -migrate down

// migrate down 2 migration
dbmigration -migrate down -number 2

// migrate down all migration
dbmigration -migrate down -number -1
```

> Administration
>> Get a migration (GET)
```
http://localhost:8001/api/v1/migrations/<migration_name>
```
>> Get migrations (GET)
```
http://localhost:8001/api/v1/migrations
```
>> Create a migration (POST)
```
http://localhost:8001/api/v1/migrations
```
with body:
```
{
	"id_migration": "<migration_name",
	"executed_at": "2018-06-02 13:11:37"
}
```
>> Delete a migration (DELETE)
```
http://localhost:8001/api/v1/migrations/<migration_name>
```
>> Delete migrations (DELETE)
```
http://localhost:8001/api/v1/migrations
```


## Known issues

## Follow me at
Facebook: https://www.facebook.com/joaosoft

LinkedIn: https://www.linkedin.com/in/jo%C3%A3o-ribeiro-b2775438/

##### If you have something to add, please let me know joaosoft@gmail.com
