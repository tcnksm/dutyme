package command

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	input "github.com/tcnksm/go-input"
)

const (
	EnvToken = "PD_SERVICE_KEY"

	EnvDebug = "DUTYME_DEBUG"
	EnvTrace = "DUTYME_TRACE"
)

const (
	ExitCodeOK = iota
	ExitCodeError
)

const (
	DefaultConfigName = ".dutyme.json"
)

var (
	Debug bool
	Trace bool
)

func init() {
	if v := os.Getenv(EnvDebug); len(v) != 0 {
		Debug = true
	}

	if v := os.Getenv(EnvTrace); len(v) != 0 {
		Trace = true
	}
}

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

func (m *Meta) AskToken() (string, error) {
	fmt.Fprintf(m.OutStream, `To use dutyme command, you need a PagerDuty API v2 token.
The token must have full access to read, write, update, and delete.

Only account administrators have the ability to generate token.
(See more about token on official doc https://goo.gl/VPvlwB)

You can set the token via %q env var (then you can
skip the following input). Or after running first overriding, you can
save login info on config file.

`, EnvToken)

	query := "Input PagerDuty API token"
	return m.UI.Ask(query, &input.Options{
		Required:  true,
		Loop:      true,
		HideOrder: true,
		Mask:      true,
	})
}

// Debugf prints debug information when debug env var
// has non-empty value. If not, it does nothing.
func Debugf(format string, args ...interface{}) {
	if Debug {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// Trace prints pkg/errors stack trace information when trace env
// var has non-empty value. If not, it does nothing.
func TracePrint(w io.Writer, err error) {
	if Trace {
		fmt.Fprintf(w, "\n[TRACE] %+v\n", err)
	}
}
