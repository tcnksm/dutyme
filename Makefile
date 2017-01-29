VERSION = $(shell grep 'Version string' version.go | sed -E 's/.*"(.+)"$$/\1/')
COMMIT = $(shell git describe --always)
PACKAGES = $(shell go list ./... | grep -v '/vendor/')

EXTERNAL_TOOLS = \
	github.com/mitchellh/gox \
	github.com/tcnksm/ghr

default: test

# install external tools for this project
bootstrap:
	@for tool in $(EXTERNAL_TOOLS) ; do \
		echo "Installing $$tool" ; \
        go get -v $$tool; \
    done

install:
	go install -ldflags "-X main.GitCommit=$(COMMIT)"	

build: 
	go build -ldflags "-X main.GitCommit=$(COMMIT)" -o bin/dutyme

xbuild: 
	@if [ -d "out/$(VERSION)" ]; then rm -fr out; fi
	gox \
      -ldflags "-X main.GitCommit=$(COMMIT)" \
      -parallel=4 \
      -os="darwin linux windows" \
      -arch="amd64" \
      -output "out/${VERSION}/{{.OS}}_{{.Arch}}/{{.Dir}}"

package: xbuild
	@if [ -d "out/$(VERSION)/dist" ]; then rm -fr "out/$(VERSION)/dist"; fi

	@mkdir -p "out/$(VERSION)/dist"
	@for P in `find ./out/$(VERSION) -mindepth 1 -maxdepth 1 -type d`; do \
		PLATFORM_NAME=$$(basename $$P); \
		if [ $$PLATFORM_NAME = "dist" ]; then continue; fi; \
		pushd $$P && zip $$PLATFORM_NAME.zip ./* && mv $$PLATFORM_NAME.zip ../dist/. && popd; \
	done

	@pushd out/$(VERSION)/dist && shasum -a 256 * > SHASUMS && popd

upload: 
	ghr v$(VERSION) out/$(VERSION)/dist

# TODO(tcnksm): Include tcnksm/homebrew-dutyme into this repository.
brew:
	go run release/main.go $(VERSION) out/$(VERSION)/dist/darwin_amd64.zip > ../homebrew-dutyme/dutyme.rb

PHONY: build xbuild install package brew upload

test-all: vet lint test test-race

test:
	go test -v -parallel=4 ${PACKAGES}

test-race:
	go test -v -race ${PACKAGES}

vet:
	go vet ${PACKAGES}

lint:
	@go get github.com/golang/lint/golint
	go list ./... | grep -v vendor | xargs -n1 golint

cover:
	@go get golang.org/x/tools/cmd/cover
	go test -coverprofile=cover.out
	go tool cover -html cover.out
	rm cover.out

.PHONY: test test-race test-all vet lint cover