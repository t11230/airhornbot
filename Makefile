BOT_BINARY=bot
WEB_BINARY=auth

JS_FILES = $(shell find static/src/ -type f -name '*.js')

BOT_SOURCEDIR=./cmd/bot/
BOT_SOURCES := $(shell find $(BOT_SOURCEDIR) -name '*.go')

.PHONY: all
all: bot auth

bot: $(BOT_SOURCES)
	go build -o ${BOT_BINARY} $(BOT_SOURCES)

auth: cmd/authserver/auth.go static
	go build -o ${WEB_BINARY} cmd/authserver/auth.go

npm: static/package.json
	cd static && npm install .

gulp: $(JS_FILES)
	cd static && gulp dist

.PHONY: static
static: npm gulp

.PHONY: clean
clean:
	rm -r ${BOT_BINARY} ${WEB_BINARY} static/dist/
