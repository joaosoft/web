package services

import (
	"time"
)

type Migration struct {
	IdMigration string    `json:"id_migration"`
	User        string    `json:"user"`
	ExecutedAt  time.Time `json:"executed_at"`
}

type ListMigration []*Migration
