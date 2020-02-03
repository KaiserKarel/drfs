package restic_test

import (
	drestic "github.com/kaiserkarel/drfs/restic"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackendImplementsResticBackend(t *testing.T) {
	assert.Implements(t, (*restic.Backend)(nil), &drestic.Backend{}, "File should implement io.Reader")
}
