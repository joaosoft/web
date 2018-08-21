package service

import "time"

type FileInfo struct {
	FullName string
	Name     string
	Size     int64
	ModTime  time.Time
}

type Event struct {
	File      string
	Operation Operation
}

type Operation string
