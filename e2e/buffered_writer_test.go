package tests

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/kaiserkarel/drfs"
	dros "github.com/kaiserkarel/drfs/os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBufferedWriter(t *testing.T) {
	dest, err := dros.Open(fmt.Sprintf("TestBufferedWriter_%s", time.Now().String()))
	require.NoError(t, err)

	src, err := os.Open("../testdata/lorem.txt")
	require.NoError(t, err)

	buf := drfs.NewBufferedWriter(dest)
	_, err = io.Copy(buf, src)
	require.NoError(t, err)

	assert.NoError(t, buf.Flush())
}
