package manidator

import (
	"bytes"
	"io"
	"strings"
)

// Dator provides an interface for an object that can provide the last line of output from, say, a command or event
// stream.
type Dator interface {
	io.WriteCloser

	// GetName returns the name of the Dator.
	GetName() string

	// GetLastLine returns the last line of output from the Dator, with leading and trailing spaces trimmed off.
	GetLastLine() string

	// IsClosed indicates that the Dator is closed and will provide no additional output.
	IsClosed() bool
}

// DefaultDator is the default implementation of the Dator interface.
type DefaultDator struct {
	*bytes.Buffer

	// name is the Dator's name.
	name string

	// closed indicates when the Dator has no more output to provide.
	closed bool
}

// NewDator creates a new DefaultDator with the specified name.
func NewDator(name string) Dator {
	return &DefaultDator{
		Buffer: bytes.NewBuffer(nil),
		name:   name,
	}
}

// GetName returns the name of the Dator.
func (this *DefaultDator) GetName() string {
	return this.name
}

// GetLastLine returns the last line of output from the Dator, with leading and trailing spaces trimmed off.
func (this *DefaultDator) GetLastLine() string {
	line := strings.TrimSpace(this.String())
	if pos := strings.LastIndex(line, "\n"); pos != -1 {
		line = line[pos:]
	}

	return strings.TrimSpace(line)
}

// IsClosed indicates that the Dator is closed and will provide no additional output.
func (this *DefaultDator) IsClosed() bool {
	return this.closed
}

// Close marks the Dator as closed.
func (this *DefaultDator) Close() error {
	this.closed = true
	return nil
}
