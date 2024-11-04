all: test build

GIT_DESCRIBE=$(shell git describe)
GIT_DESCRIBE_LONG=$(shell git describe --long)

clean:
	rm -r ./bin

test:
	go test ./...

_build:
	go build -ldflags "-X 'main.Version=$(GIT_DESCRIBE)' -X 'main.VersionLong=$(GIT_DESCRIBE_LONG)'" -o $(OUTPUT_FILE) ./cmd/uformat

build:
	OUTPUT_FILE=bin/uformat $(MAKE) _build

_release:
	mkdir -p ./bin/release/$(OS)-$(ARCH)

	if [[ "$(OS)" == "windows" ]]; then\
		GOOS=$(OS) GOARCH=$(ARCH) OUTPUT_FILE=bin/release/$(OS)-$(ARCH)/uformat.exe $(MAKE) _build;\
	else\
		GOOS=$(OS) GOARCH=$(ARCH) OUTPUT_FILE=bin/release/$(OS)-$(ARCH)/uformat $(MAKE) _build;\
	fi

	tar -cvzf bin/release/bin/$(OS)-$(ARCH).tar.gz -C bin/release $(OS)-$(ARCH)

_macrelease:
	mkdir -p ./bin/release/darwin-universal

	lipo -create -output bin/release/darwin-universal/uformat bin/release/darwin-amd64/uformat bin/release/darwin-arm64/uformat

	rm bin/release/bin/darwin-*
	tar -cvzf bin/release/bin/darwin-universal.tar.gz -C bin/release darwin-universal

release: clean
	mkdir -p ./bin/release/bin

	OS=darwin ARCH=amd64 $(MAKE) _release
	OS=darwin ARCH=arm64 $(MAKE) _release
	if [[ "$(shell uname -s)" == "Darwin" ]]; then\
		$(MAKE) _macrelease;\
	fi

	OS=linux ARCH=amd64 $(MAKE) _release
	OS=linux ARCH=arm64 $(MAKE) _release
	OS=windows ARCH=amd64 $(MAKE) _release
	OS=windows ARCH=arm64 $(MAKE) _release
