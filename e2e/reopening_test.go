package tests

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	drfs "github.com/kaiserkarel/drfs/os"
	"github.com/kaiserkarel/drfs/recovery"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"testing"
)

func TestReopeningResultsInSameFileALR(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	reopenE2E(t, "../testdata/alr.gen.txt")
}

func TestReopeningResultsInSameFileShort(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	reopenE2E(t, "../testdata/lorem_short.txt")
}

func TestReopeningResultsInSameFileMedium(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	reopenE2E(t, "../testdata/lorem_medium.txt")
}

func TestReopeningResultsInSameFileLorem(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}

	reopenE2E(t, "../testdata/lorem.txt")
}

func reopenE2E(t *testing.T, src string) {
	var fileName = fmt.Sprintf("TestReopeningResultsInSameIndex_%s", time.Now().String())

	file, err := drfs.Open(fileName)
	require.NoError(t, err, "cannot open drfs file")

	fstat, err := file.Fstat()
	require.NoError(t, err, "fstat should work")

	assert.LessOrEqual(t, fstat.QuotaBytesUsed(), int64(0), "new file should consume no quota")

	lorem, err := os.Open(src)
	require.NoError(t, err, "unable to open lorem")

	_, err = io.Copy(file, lorem)
	require.NoError(t, err)

	lorem.Close()

	fstat, err = file.Fstat()
	require.NoError(t, err, "second fstat should work")
	assert.LessOrEqual(t, fstat.QuotaBytesUsed(), int64(0), "quota should remain 0 after writes")

	file2, err := drfs.Open(fileName)
	require.NoError(t, err, "cannot reopen drfs file")

	assert.Equal(t, file2.Index().Header, file.Index().Header, "index headers should be equal")

	buckets1 := file.Index().Buckets
	buckets2 := file2.Index().Buckets

	assert.Len(t, buckets2, len(buckets1), "bucket lengths should be equal")

	for i := 0; i < len(buckets1) && i < len(buckets2); i++ {
		assert.Equal(t, buckets2[i].Header, buckets1[i].Header, "headers should ne equal")
		assert.Equal(t, buckets2[i].CommentID, buckets1[i].CommentID, "commentIDs should ne equal")
	}

	stat1, _ := file.Stat()
	stat2, _ := file2.Stat()
	assert.Equal(t, stat2.Name(), stat1.Name(), "file.Stats.Name should be equal")
	assert.Equal(t, stat2.Size(), stat1.Size(), "file.Size should be equal")

	// check if local size bookkeeping matches recovery
	rStats, err := recovery.Stats(file)
	require.NoError(t, err, "recovery stats should work")

	assert.Equal(t, rStats.Size(), stat2.Size(), "index size should match recovery size")

	lorem, err = os.Open(src)
	require.NoError(t, err, "should open lorem")

	log.Printf("diffing")
	diff(t, lorem, file)
	lorem.Close()
}

const size = 4094

// https://groups.google.com/forum/#!topic/golang-nuts/keG78hYt1I0
func readComplete(r io.Reader, b []byte) (int, error) {
	var (
		n, _n int
		err   error
	)
	for _n = 0; err == nil && n < size; n += _n {
		_n, err = r.Read(b[n:])
		log.Printf("%d, %d, %s", _n, n, err)
	}
	return n, err
}

// https://groups.google.com/forum/#!topic/golang-nuts/keG78hYt1I0
// Diff compares the contents of two io.Readers.
// The return value of identical is true if and only if there are no errors
// in reading r1 and r2 (io.EOF excluded) and r1 and r2 are
// byte-for-byte identical.
func diff(t *testing.T, exp, val io.Reader) {
	var (
		input = [...]io.Reader{exp, val}

		n    [2]int
		errs [2]error
		buf  [2][size]byte
	)
	for {
		for i, r := range input {
			n[i], errs[i] = readComplete(r, buf[i][:])
			if errs[i] != nil && errs[i] != io.EOF {
				require.NoError(t, errs[i], "expected no error during read")
			}
		}

		log.Println("comparing")
		require.Equal(t, string(buf[0][:n[0]]), string(buf[1][:n[1]]))

		if errs[0] == io.EOF {
			assert.Error(t, errs[1], "second error should EOF too")
			return
		}

		if errs[1] == io.EOF {
			assert.Error(t, errs[0], "first error should EOF too")
			return
		}
	}
}
