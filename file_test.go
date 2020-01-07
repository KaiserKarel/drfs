package drfs

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImplementsReader(t *testing.T) {
	assert.Implements(t, (*io.Reader)(nil), &File{}, "File should implement io.Reader")
}

func TestImplementsWriter(t *testing.T) {
	assert.Implements(t, (*io.Writer)(nil), &File{}, "File should implement io.Writer")
}
