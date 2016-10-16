package command

import (
	"flag"
	"fmt"
	"io"
	"log"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	input "github.com/tcnksm/go-input"
)

const (
	DefaultConfigName = ".dutyme.json"
)

var Debug bool

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	OutStream io.Writer
	ErrStream io.Writer

	UI *input.UI
}

func (m *Meta) NewFlagSet(name, usage string) *flag.FlagSet {
	flags := flag.NewFlagSet(name, flag.ContinueOnError)
	flags.SetOutput(m.ErrStream)
	flags.Usage = func() {
		fmt.Fprintf(m.OutStream, usage+"\n")
	}
	return flags
}

func (m *Meta) ConfigPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, DefaultConfigName), nil
}

func Debugf(format string, args ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}
