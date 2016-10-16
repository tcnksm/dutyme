package command

import "io"

// Meta contain the meta-option that nearly all subcommand inherited.
type Meta struct {
	OutStream io.Writer
	ErrStream io.Writer
	InStream  io.Reader
}
