GIT_VERSION=`git log --pretty=format:"%h" -1`
BIN_VERSION=`cat version.txt`

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
	env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.gitCommit=${GIT_VERSION} -X main.binVersion=${BIN_VERSION}" -o bin/bot cmd/main.go

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
	cp chatbot-plugins/shodan/main.go bin/active-plugins/shodan.go
	cp chatbot-plugins/urlshorten/google/main.go bin/active-plugins/urlshorten.go
	cp chatbot-plugins/nmap/main.go bin/active-plugins/nmap.go
	cp chatbot-plugins/strawpoll/main.go bin/active-plugins/strawpoll.go
	cp chatbot-plugins/remindme/main.go bin/active-plugins/remindme.go
	cp chatbot-plugins/api/main.go bin/active-plugins/api.go

	cp chatbot-plugins/chatlog/plain/main.go bin/background-plugins/log.go
	cp chatbot-plugins/auth/team/main.go bin/background-plugins/auth.go
	cp chatbot-plugins/ratelimit/main.go bin/background-plugins/ratelimit.go

plugin-build:
	@echo "+ $@"
	go build --buildmode=plugin -o bin/active-plugins/crypto-api.so bin/active-plugins/crypto-api.go
	go build --buildmode=plugin -ldflags "-X main.version=${BIN_VERSION}" -o bin/active-plugins/help.so bin/active-plugins/help.go
	go build --buildmode=plugin -o bin/active-plugins/reddit.so bin/active-plugins/reddit.go
	go build --buildmode=plugin -o bin/active-plugins/weather.so bin/active-plugins/weather.go
	go build --buildmode=plugin -o bin/active-plugins/hibp-account.so bin/active-plugins/hibp-account.go
	go build --buildmode=plugin -o bin/active-plugins/hibp-password.so bin/active-plugins/hibp-password.go
	go build --buildmode=plugin -o bin/active-plugins/giphy.so bin/active-plugins/giphy.go
	go build --buildmode=plugin -o bin/active-plugins/media.so bin/active-plugins/media.go
	go build --buildmode=plugin -o bin/active-plugins/screenshot.so bin/active-plugins/screenshot.go
	go build --buildmode=plugin -o bin/active-plugins/virustotal.so bin/active-plugins/virustotal.go
	go build --buildmode=plugin -o bin/active-plugins/shodan.so bin/active-plugins/shodan.go
	go build --buildmode=plugin -o bin/active-plugins/urlshorten.so bin/active-plugins/urlshorten.go
	go build --buildmode=plugin -o bin/active-plugins/nmap.so bin/active-plugins/nmap.go
	go build --buildmode=plugin -o bin/active-plugins/strawpoll.so bin/active-plugins/strawpoll.go
	go build --buildmode=plugin -o bin/active-plugins/remindme.so bin/active-plugins/remindme.go
	go build --buildmode=plugin -o bin/active-plugins/api.so bin/active-plugins/api.go

	go build --buildmode=plugin -o bin/background-plugins/auth.so bin/background-plugins/auth.go
	go build --buildmode=plugin -o bin/background-plugins/log.so bin/background-plugins/log.go
	go build --buildmode=plugin -o bin/background-plugins/ratelimit.so bin/background-plugins/ratelimit.go

	rm bin/active-plugins/*.go
	rm bin/background-plugins/*.go

clean:
	rm -r bin/

all: lint vet test build plugin-setup plugin-build