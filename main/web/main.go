package main

import (
	"db-migration/services/web"
)

func main() {
	m, err := web.NewService()
	if err != nil {
		panic(err)
	}

	if err := m.Start(); err != nil {
		panic(err)
	}
}
