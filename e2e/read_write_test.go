package tests

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	"unsafe"

	"github.com/kaiserkarel/drfs/os"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestReadWrite is used to identify the cause of a bug where it
// seemed file length was correct, but read text was forward shifted
// by ~ 16 characters
func TestReadWrite(t *testing.T) {
	var fileName = fmt.Sprintf("TestReadWrite_%s", time.Now().String())

	f, err := os.Open(fileName)
	require.NoError(t, err)

	for i := 0; i < 40940/30; i++ {
		var payload = []byte(randStr(30))
		var buf = make([]byte, len(payload))
		n, err := f.Write(payload)
		assert.NoError(t, err, "i: %d written: %d", i, n)

		r, err := f.Read(buf)
		assert.Equal(t, r, n, "i: %d written: %d read:", i, n, r)

		assert.Equal(t, payload, buf)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
var src = rand.NewSource(time.Now().UnixNano())

func randStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
