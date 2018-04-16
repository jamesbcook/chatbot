
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
	mkdir bin/active-plugins
	mkdir bin/background-plugins
	cp chatbot-plugins/cryptocurrency/main.go bin/active-plugins/crypto-api.go
	cp chatbot-plugins/help/main.go bin/active-plugins/help.go
	cp chatbot-plugins/reddit/main.go bin/active-plugins/reddit.go
	cp chatbot-plugins/weather/main.go bin/active-plugins/weather.go
	cp chatbot-plugins/hibp/account/main.go bin/active-plugins/hibp-account.go
	cp chatbot-plugins/hibp/password/main.go bin/active-plugins/hibp-password.go
	cp chatbot-plugins/media/giphy/main.go bin/active-plugins/giphy.go
	cp chatbot-plugins/media/direct/main.go bin/active-plugins/media.go
	cp chatbot-plugins/media/direct/main.go bin/active-plugins/media.go
	cp chatbot-plugins/screenshot/main.go bin/active-plugins/screenshot.go
	cp chatbot-plugins/virustotal/hash/main.go bin/active-plugins/virustotal.go

	cp chatbot-plugins/chatlog/plain/main.go bin/background-plugins/log.go
	cp chatbot-plugins/auth/team/main.go bin/background-plugins/auth.go

plugin-build:
	@echo "+ $@"
	go build --buildmode=plugin -o bin/active-plugins/crypto-api.so bin/active-plugins/crypto-api.go
	rm bin/active-plugins/crypto-api.go
	go build --buildmode=plugin -o bin/active-plugins/help.so bin/active-plugins/help.go
	rm bin/active-plugins/help.go
	go build --buildmode=plugin -o bin/active-plugins/reddit.so bin/active-plugins/reddit.go
	rm bin/active-plugins/reddit.go
	go build --buildmode=plugin -o bin/active-plugins/weather.so bin/active-plugins/weather.go
	rm bin/active-plugins/weather.go
	go build --buildmode=plugin -o bin/active-plugins/hibp-account.so bin/active-plugins/hibp-account.go
	rm bin/active-plugins/hibp-account.go
	go build --buildmode=plugin -o bin/active-plugins/hibp-password.so bin/active-plugins/hibp-password.go
	rm bin/active-plugins/hibp-password.go
	go build --buildmode=plugin -o bin/active-plugins/giphy.so bin/active-plugins/giphy.go
	rm bin/active-plugins/giphy.go
	go build --buildmode=plugin -o bin/active-plugins/media.so bin/active-plugins/media.go
	rm bin/active-plugins/media.go
	go build --buildmode=plugin -o bin/active-plugins/screenshot.so bin/active-plugins/screenshot.go
	rm bin/active-plugins/screenshot.go
	go build --buildmode=plugin -o bin/active-plugins/virustotal.so bin/active-plugins/virustotal.go
	rm bin/active-plugins/virustotal.go

	go build --buildmode=plugin -o bin/background-plugins/auth.so bin/background-plugins/auth.go
	rm bin/background-plugins/auth.go
	go build --buildmode=plugin -o bin/background-plugins/log.so bin/background-plugins/log.go
	rm bin/background-plugins/log.go

clean:
	rm -r bin/

all: lint vet test build plugin-setup plugin-build