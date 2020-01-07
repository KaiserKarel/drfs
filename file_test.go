package drfs

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestImplementsReader(t *testing.T) {
	assert.Implements(t, (*io.Reader)(nil), &File{}, "File should implement io.Reader")
}

func TestImplementsWriter(t *testing.T) {
	assert.Implements(t, (*io.Writer)(nil), &File{}, "File should implement io.Writer")
}
