package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// formulaTmpl is temaplate for homebrew formula
var formulaTmpl = `require "formula"

class {{ .Name | Title }} < Formula
  homepage "https://github.com/tcnksm/{{ .Name }}"
  version 'v{{ .Version }}'

  url "https://github.com/tcnksm/{{ .Name }}/releases/download/v{{ .Version }}/darwin_amd64.zip"
  sha256 "{{ .Sha256 }}"

  def install
    bin.install '{{ .Name }}'
  end
end
`

func main() {
	os.Exit(_main())
}

func _main() int {
	if len(os.Args) != 3 {
		log.Println("Usage: go run main.go VERSION FILE")
		return 0
	}

	name := "dutyme"
	version := os.Args[1]
	file := os.Args[2]

	file, err := filepath.Abs(file)
	if err != nil {
		log.Println(err)
		return 1
	}

	f, err := os.Open(file)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println(err)
		return 1
	}
	checkSum := sha256.Sum256(buf)

	tmpl, err := template.New("formula").Funcs(template.FuncMap{
		"Title": strings.Title,
	}).Parse(formulaTmpl)
	if err != nil {
		log.Fatal(err)
	}

	if err := tmpl.Execute(os.Stdout, struct {
		Name, Version, Sha256 string
	}{
		Name:    name,
		Version: version,
		Sha256:  fmt.Sprintf("%x", checkSum),
	}); err != nil {
		log.Println(err)
		return 1
	}

	return 0
}
