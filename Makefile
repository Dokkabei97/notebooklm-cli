BINARY := nlm
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X github.com/jmk/notebooklm-cli/cmd.Version=$(VERSION) -X github.com/jmk/notebooklm-cli/cmd.CommitSHA=$(COMMIT)"

.PHONY: build install clean test e2e e2e-basic release-dry

build:
	go build $(LDFLAGS) -o $(BINARY) .

install:
	go install $(LDFLAGS) .

clean:
	rm -f $(BINARY)
	rm -rf dist/

release-dry:
	goreleaser release --snapshot --clean

test:
	go test ./internal/...

e2e-basic: build
	NLM_BINARY=$(PWD)/$(BINARY) go test ./e2e/ -v -run "TestVersion|TestHelp|TestSubcommand|TestCompletion|TestUnknown|TestUseRequires|Test.*Aliases|TestGlobalFlags" -count=1

e2e: build
	NLM_BINARY=$(PWD)/$(BINARY) go test ./e2e/ -v -count=1 -timeout 120s
