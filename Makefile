# Convenience tasks
.PHONY: test lint run

GO ?= go

test:
	$(GO) test -race ./...

lint:
	golangci-lint run

run:
	@if [ -z "$(DAY)" ]; then echo "Usage: make run DAY=01"; exit 1; fi; \
	cd day$(DAY) && $(GO) run ./...
