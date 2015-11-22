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
	barcode/*.go \
	channels/gank/*.go \

.PHONY: clean all test fmt

clean:
	rm -r $(FILE)

fmt:
	gofmt -w -s ./

$(FILE): $(SOURCES)
	cd $(EXEC) && go install

all: $(FILE)

test: $(FILE)
	$(FILE) test/test.json
