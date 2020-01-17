package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	drfs "github.com/kaiserkarel/drfs/os"
	"github.com/kaiserkarel/drfs/recovery"
	"github.com/stretchr/testify/require"
)

func TestAppendWrites(t *testing.T) {
	var fileName = fmt.Sprintf("TestAppendWrites%s", time.Now().String())

	file, err := drfs.Open(fileName)
	require.NoError(t, err, "cannot open drfs file")

	var payload = "hello payload!"
	for i := 0; i < 10; i++ {
		_, err := file.Write([]byte(payload))
		require.NoError(t, err, "write should work")
	}

	stat, err := file.Stat()
	require.NoError(t, err, "should stat")
	rStat, err := recovery.Stats(file)
	require.NoError(t, err, "should recovery.Stat")

	if !assert.Equal(t, int(rStat.Size()), 10*len(payload)) ||
		!assert.Equal(t, rStat.Size(), stat.Size(), "sizes should match") {
		buf := make([]byte, len(payload))
		var err error
		var r int
		var n int
		for err == nil {
			r, err = file.Read(buf)
			n += r
			log.Print("checking")
			assert.Equal(t, payload, string(buf), "fail at %d", n)
		}
	}
}
