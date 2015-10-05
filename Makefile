# Makefile for gotgbot
EXEC := tgbot
FILE := $(GOPATH)/bin/$(EXEC)

SOURCES := \
	tgbot/*.go \
	support/types/*.go \
	support/loader/*.go \
	support/help/*.go \
	support/utils/*.go \
	misc/*.go \
	scholar/*.go \
	chinese/*.go \
	script/*.go \
	pictures/*.go \

.PHONY: clean all test

clean:
	rm -r $(FILE)

$(FILE): $(SOURCES)
	cd $(EXEC) && go install

all: $(FILE)

test: $(FILE)
	$(FILE) test/test.json
