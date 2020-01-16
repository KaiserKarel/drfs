package drive_test

import (
	"testing"

	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/drive"
	"github.com/stretchr/testify/assert"
)

func TestImplementsClient(t *testing.T) {
	assert.Implements(t, (*drfs.Client)(nil), &drive.Client{}, "Client should implement drfs.Client")
}
