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
      -output "out/${VERSION}/{{.Dir}}_{{.OS}}_{{.Arch}}"
	cd out/${VERSION} && shasum -a 256 * > SHASUMS && cat SHASUMS

ghr: 
	ghr v$(VERSION) out/$(VERSION)/

# TODO(tcnksm): Include tcnksm/homebrew-dutyme into this repository.
brew:
	go run release/main.go $(VERSION) out/$(VERSION)/dutyme_darwin_amd64 

test: vet
	go test -v $$(go list ./... | grep -v 'vendor/' )

vet:
	@go tool vet -all $$(ls -d */ | grep -v vendor)

lint:
	@go get github.com/golang/lint/golint
	golint ./...
