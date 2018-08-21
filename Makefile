run:
	go run ./main/main.go

fmt:
	go fmt ./...

vet:
	go vet ./*

gometalinter:
	gometalinter ./*