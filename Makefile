BOT_BINARY=bot
WEB_BINARY=web

JS_FILES = $(shell find static/src/ -type f -name '*.js')

BOT_SOURCEDIR=./cmd/bot/
BOT_SOURCES := $(shell find $(BOT_SOURCEDIR) -name '*.go')

.PHONY: all
all: bot web

bot: $(BOT_SOURCES)
	go build -o ${BOT_BINARY} $(BOT_SOURCES)

web: cmd/webserver/web.go static
	go build -o ${WEB_BINARY} cmd/webserver/web.go

npm: static/package.json
	cd static && npm install .

gulp: $(JS_FILES)
	cd static && gulp dist

.PHONY: static
static: npm gulp

.PHONY: clean
clean:
	rm -r ${BOT_BINARY} ${WEB_BINARY} static/dist/
