package drive_test

import (
	"testing"

	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/drive"
	"github.com/stretchr/testify/assert"
)

func TestImplementsService(t *testing.T) {
	assert.Implements(t, (*drfs.Service)(nil), &drive.Service{}, "Service should implement drfs.Service")
}
