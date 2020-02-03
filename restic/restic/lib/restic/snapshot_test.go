package restic_test

import (
	"testing"
	"time"

	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	rtest "github.com/kaiserkarel/drfs/restic/restic/lib/test"
)

func TestNewSnapshot(t *testing.T) {
	paths := []string{"/home/foobar"}

	_, err := restic.NewSnapshot(paths, nil, "foo", time.Now())
	rtest.OK(t, err)
}
