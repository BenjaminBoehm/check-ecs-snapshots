.PHONY: lint vet

build:
	go build -o check_ecs_snapshots main.go
lint:
	go fmt $(go list ./... | grep -v /vendor/)
vet:
	go vet $(go list ./... | grep -v /vendor/)
test:
	go test -v -cover ./...
