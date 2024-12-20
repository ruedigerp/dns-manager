BINARY_NAME=dns-manager

GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get

# Ensure linker embeds versioning information
VERSION=${shell git describe --tags $(git rev-list --tags --max-count=1)}
COMMIT_HASH=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
BUILD_TIMESTAMP=$(shell date)
LDFLAGS=-s -w -extldflags= \
  -X 'dns-manager/core.GitCommitHash=$(COMMIT_HASH)' \
  -X 'dns-manager/core.Branch=$(GIT_BRANCH)' \
  -X 'dns-manager/core.BuildTimestamp=$(BUILD_TIMESTAMP)' \
  -X 'dns-manager/core.Ver=$(VERSION)' \
  -X 'dns-manager/core.Agent=dns-manager' \
  -X 'dns-manager/core.Stage=prod'

all: darwin_arm64 darwin_amd64 linux_amd64 linux_386 linux_arm64 linux_arm windows_amd64 windows_386

tarballs: all
	for file in builds/*; do \
		tar -czvf builds/`basename "$$file"`.tar.gz -C builds `basename "$$file"`; \
	done

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)-$(VERSION)-darwin-amd64
	rm -f $(BINARY_NAME)-$(VERSION)-darwin-386
	rm -f $(BINARY_NAME)-$(VERSION)-linux-amd64
	rm -f $(BINARY_NAME)-$(VERSION)-linux-386
	rm -f $(BINARY_NAME)-$(VERSION)-linux-arm64
	rm -f $(BINARY_NAME)-$(VERSION)-linux-arm
	rm -f $(BINARY_NAME)-$(VERSION)-windows-amd64
	rm -f $(BINARY_NAME)-$(VERSION)-windows-386

test:
	$(GOTEST) -v ./...

darwin_arm64:
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-darwin-arm64 -v 

darwin_amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-darwin-amd64 -v 

linux_amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-linux-amd64 -v

linux_386:
	GOOS=linux GOARCH=386 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-linux-386 -v

linux_arm64:
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-linux-arm64 -v

linux_arm:
	GOOS=linux GOARCH=arm GOARM=7 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-linux-arm -v

windows_amd64:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-windows-amd64 -v

windows_386:
	GOOS=windows GOARCH=386 $(GOBUILD) -ldflags="$(LDFLAGS)" -o builds/$(BINARY_NAME)-$(VERSION)-windows-386 -v

deps:
	$(GOGET) ./...

.PHONY: