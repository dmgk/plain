.PHONY: all
all: plain

SOURCES := $(shell find . -type f -name \*.go)

plain: $(SOURCES)
	go build ./...

server: plain
	goapp serve appengine

deploy:
	goapp deploy -application plain-im -version 1 appengine
