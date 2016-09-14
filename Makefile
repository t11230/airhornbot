BOT_BINARY=bot
WEB_BINARY=auth

JS_FILES = $(shell find static/src/ -type f -name '*.js')

BOT_SOURCEDIR=./cmd/bot/
BOT_SOURCES := $(shell find $(BOT_SOURCEDIR) -name '*.go')

.PHONY: all
all: bot auth

bot: $(BOT_SOURCES)
	go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`" -o ${BOT_BINARY} $(BOT_SOURCES)

auth: cmd/authserver/auth.go
	go build -ldflags "-X main.buildstamp `date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.githash `git rev-parse HEAD`" -o ${WEB_BINARY} cmd/authserver/auth.go

npm: static/package.json
	cd static && npm install .

gulp: $(JS_FILES)
	cd static && gulp dist

.PHONY: static
static: npm gulp

.PHONY: clean
clean:
	rm -r ${BOT_BINARY} ${WEB_BINARY} static/dist/
