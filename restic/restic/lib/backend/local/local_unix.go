// +build !windows

package local

import (
	"os"

	"github.com/kaiserkarel/drfs/restic/restic/lib/fs"
)

// set file to readonly
func setNewFileMode(f string, mode os.FileMode) error {
	return fs.Chmod(f, mode)
}
