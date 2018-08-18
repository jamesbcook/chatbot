
lint:
	@echo "+ $@"
	@golint ./... | tee /dev/stderr

vet:
	@echo "+ $@"
	@go vet $(shell go list ./...) | tee /dev/stderr

test:
	@echo "+ $@"
	@go test -cover $(shell go list ./... | grep -vE '(cmd)')

all: lint vet test