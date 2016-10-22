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

test: vet
	go test -v $$(go list ./... | grep -v 'vendor/' )

vet:
	@go tool vet -all $$(ls -d */ | grep -v vendor)

lint:
	@go get github.com/golang/lint/golint
	golint ./...
