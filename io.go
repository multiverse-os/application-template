package application

import (
	"io"
)

type IO struct {
	Output *io.Writer
	Error  *io.Writer
	Input  *io.Reader
}
