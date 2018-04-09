
lint:
	@echo "+ $@"
	@golint ./... | tee /dev/stderr

vet:
	@echo "+ $@"
	@go vet $(shell go list ./...) | tee /dev/stderr

test:
	@echo "+ $@"
	@go test -cover $(shell go list ./... | grep -vE '(cmd)')

build:
	@echo "+ $@"
	env GOOS=linux GOARCH=amd64 go build -o bin/bot cmd/main.go

plugin-setup: 
	@echo "+ $@"
	mkdir active-plugins
	mkdir background-plugins
	cp chat-bot-plugins/cryptocurrency/main.go active-plugins/crypto-api.go
	cp chat-bot-plugins/help/main.go active-plugins/help.go
	cp chat-bot-plugins/reddit/main.go active-plugins/reddit.go
	cp chat-bot-plugins/rta/main.go active-plugins/rta.go
	cp chat-bot-plugins/weather/main.go active-plugins/weather.go
	cp chat-bot-plugins/leaksrus/main.go active-plugins/leaksrus.go
	cp chat-bot-plugins/hibp/account/main.go active-plugins/hibp-account.go
	cp chat-bot-plugins/hibp/password/main.go active-plugins/hibp-password.go
	cp chat-bot-plugins/media/giphy/main.go active-plugins/giphy.go
	cp chat-bot-plugins/media/direct/main.go active-plugins/media.go

	cp chat-bot-plugins/chatlog/plain/main.go background-plugins/log.go
	cp chat-bot-plugins/auth/team/main.go background-plugins/auth.go

plugin-build:
	@echo "+ $@"
	go build --buildmode=plugin -o active-plugins/crypto-api.so active-plugins/crypto-api.go
	rm active-plugins/crypto-api.go
	go build --buildmode=plugin -o active-plugins/help.so active-plugins/help.go
	rm active-plugins/help.go
	go build --buildmode=plugin -o active-plugins/reddit.so active-plugins/reddit.go
	rm active-plugins/reddit.go
	go build --buildmode=plugin -o active-plugins/rta.so active-plugins/rta.go
	rm active-plugins/rta.go
	go build --buildmode=plugin -o active-plugins/weather.so active-plugins/weather.go
	rm active-plugins/weather.go
	go build --buildmode=plugin -o active-plugins/leaksrus.so active-plugins/leaksrus.go
	rm active-plugins/leaksrus.go
	go build --buildmode=plugin -o active-plugins/hibp-account.so active-plugins/hibp-account.go
	rm active-plugins/hibp-account.go
	go build --buildmode=plugin -o active-plugins/hibp-password.so active-plugins/hibp-password.go
	rm active-plugins/hibp-password.go
	go build --buildmode=plugin -o active-plugins/giphy.so active-plugins/giphy.go
	rm active-plugins/giphy.go
	go build --buildmode=plugin -o active-plugins/media.so active-plugins/media.go
	rm active-plugins/media.go
	go build --buildmode=plugin -o background-plugins/auth.so background-plugins/auth.go
	rm background-plugins/auth.go
	go build --buildmode=plugin -o background-plugins/log.so background-plugins/log.go
	rm background-plugins/log.go

clean:
	rm -r bin/
	rm -r active-plugins/
	rm -r background-plugins/

all: lint vet test build plugin-setup plugin-build