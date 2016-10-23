VERSION ?= $$(grep 'Version' version.go | sed -E 's/.*"(.+)"$$/\1/')
COMMIT ?= $$(git describe --always)

default: test

deps:
	go get github.com/mitchellh/gox
	go get -v .

build: deps
	go build -ldflags "-X main.GitCommit=$(COMMIT)" -o bin/dutyme

xbuild: deps
	@if [ -d "out/$(VERSION)" ]; then rm -fr out; fi
	gox \
      -ldflags "-X main.GitCommit=$(COMMIT)" \
      -parallel=3 \
      -os="darwin linux windows" \
      -arch="amd64" \
      -output "out/${VERSION}/{{.OS}}_{{.Arch}}/{{.Dir}}"

package: #xbuild
	@if [ -d "out/$(VERSION)/dist" ]; then rm -fr "out/$(VERSION)/dist"; fi

	mkdir -p "out/$(VERSION)/dist"
	@for P in `find ./out/$(VERSION) -mindepth 1 -maxdepth 1 -type d`; do \
		PLATFORM_NAME=$$(basename $$P); \
		if [ $$PLATFORM_NAME = "dist" ]; then continue; fi; \
		pushd $$P && zip $$PLATFORM_NAME.zip ./* && mv $$PLATFORM_NAME.zip ../dist/. && popd; \
	done

	pushd out/$(VERSION)/dist && shasum -a 256 * > SHASUMS && popd

ghr: 
	ghr v$(VERSION) out/$(VERSION)/dist

# TODO(tcnksm): Include tcnksm/homebrew-dutyme into this repository.
brew:
	go run release/main.go $(VERSION) out/$(VERSION)/dist/darwin_amd64.zip > ../homebrew-dutyme/dutyme.rb

test: vet
	go test -v $$(go list ./... | grep -v 'vendor/' )

vet:
	@go tool vet -all $$(ls -d */ | grep -v vendor)

lint:
	@go get github.com/golang/lint/golint
	golint ./...
