BIN := goma
VERSION := $$(make -s show-version)
CURRENT_REVISION := $(shell git rev-parse --short HEAD)
BUILD_LDFLAGS := "-s -w -X main.revision=$(CURRENT_REVISION)"
GOBIN ?= $(shell go env GOPATH)/bin

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags=$(BUILD_LDFLAGS) -o $(BIN) cmd/goma/main.go

.PHONY: install
install:
	go install -ldflags=$(BUILD_LDFLAGS) cmd/goma/main.go

.PHONY: show-version
show-version: $(GOBIN)/gobump
	@gobump show -r cmd/goma

$(GOBIN)/gobump:
	@go install github.com/x-motemen/gobump/cmd/gobump@latest

.PHONY: cross
cross: $(GOBIN)/goxz
	CGO_ENABLED=0 goxz -n $(BIN) -pv=v$(VERSION) -build-ldflags=$(BUILD_LDFLAGS) ./cmd/goma

$(GOBIN)/goxz:
	go install github.com/Songmu/goxz/cmd/goxz@latest

.PHONY: clean
clean:
	rm -rf $(BIN) goxz
	go clean

.PHONY: upload
upload: $(GOBIN)/ghr
	ghr "v$(VERSION)" goxz

$(GOBIN)/ghr:
	go install github.com/tcnksm/ghr@latest
