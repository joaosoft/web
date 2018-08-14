run:
	go run ./bin/launcher/main.go

build:
	mkdir -p bin/cmd/
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o bin/cmd/dbmigration main/cmd/main.go

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*

migrate-up:
	./bin/cmd/dbmigration -migrate up

migrate-down:
	./bin/cmd/dbmigration -migrate down