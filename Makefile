all: test build

GIT_DESCRIBE=$(shell git describe)
GIT_DESCRIBE_LONG=$(shell git describe --long)
GIT_DESCRIBE_NO_V=$(shell git describe | sed 's/^v//g')

clean:
	rm -rf ./bin

test:
	go test ./...

_manpage:
	asciidoctor -b manpage -a version="$(GIT_DESCRIBE)" -D $(OUTPUT_MAN) dist/uformat.adoc

manpage:
	OUTPUT_MAN=bin/ $(MAKE) _manpage

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
		cp bin/uformat.1 bin/release/$(OS)-$(ARCH)/uformat.1;\
	fi

	tar -cvzf bin/release/bin/$(OS)-$(ARCH).tar.gz -C bin/release $(OS)-$(ARCH)

_macrelease:
	mkdir -p ./bin/release/darwin-universal

	lipo -create -output bin/release/darwin-universal/uformat bin/release/darwin-amd64/uformat bin/release/darwin-arm64/uformat

	rm bin/release/bin/darwin-*
	cp bin/uformat.1 bin/release/darwin-universal/uformat.1
	tar -cvzf bin/release/bin/darwin-universal.tar.gz -C bin/release darwin-universal

_linuxrelease:
	ARCH=$(ARCH) VERSION=$(GIT_DESCRIBE_NO_V) UFORMAT_BIN=bin/release/linux-$(ARCH)/uformat UFORMAT_MAN=bin/release/linux-$(ARCH)/uformat.1 nfpm pkg --config dist/nfpm.yaml --packager deb --target bin/release/bin
	ARCH=$(ARCH) VERSION=$(GIT_DESCRIBE_NO_V) UFORMAT_BIN=bin/release/linux-$(ARCH)/uformat UFORMAT_MAN=bin/release/linux-$(ARCH)/uformat.1 nfpm pkg --config dist/nfpm.yaml --packager rpm --target bin/release/bin

release: clean manpage
	mkdir -p ./bin/release/bin

	OS=darwin ARCH=amd64 $(MAKE) _release
	OS=darwin ARCH=arm64 $(MAKE) _release
	if [[ "$(shell uname -s)" == "Darwin" ]]; then\
		$(MAKE) _macrelease;\
	fi

	OS=linux ARCH=amd64 $(MAKE) _release
	OS=linux ARCH=arm64 $(MAKE) _release
	if $(shell which nfpm); then\
		ARCH=amd64 $(MAKE) _linuxrelease;\
		ARCH=arm64 $(MAKE) _linuxrelease;\
	fi

	OS=windows ARCH=amd64 $(MAKE) _release
	OS=windows ARCH=arm64 $(MAKE) _release
