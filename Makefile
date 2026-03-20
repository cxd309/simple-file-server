BINARY  := simple-file-server
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"
OUTDIR  := dist

TARGETS := \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: build
build:
	@mkdir -p $(OUTDIR)
	@$(foreach target,$(TARGETS), \
		$(eval OS   := $(word 1,$(subst /, ,$(target)))) \
		$(eval ARCH := $(word 2,$(subst /, ,$(target)))) \
		$(eval EXT  := $(if $(filter windows,$(OS)),.exe,)) \
		echo "Building $(OS)/$(ARCH)..."; \
		GOOS=$(OS) GOARCH=$(ARCH) go build $(LDFLAGS) -o $(OUTDIR)/$(BINARY)-$(OS)-$(ARCH)$(EXT) . ; \
	)

.PHONY: clean
clean:
	rm -rf $(OUTDIR)

.PHONY: fmt
fmt:
	go mod tidy
	gofmt -w .
	golangci-lint run
	dprint fmt