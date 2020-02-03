// +build windows

package fs

import (
	"os"

	"github.com/kaiserkarel/drfs/restic/restic/lib/errors"
)

// DeviceID extracts the device ID from an os.FileInfo object by casting it
// to syscall.Stat_t
func DeviceID(fi os.FileInfo) (deviceID uint64, err error) {
	return 0, errors.New("Device IDs are not supported on Windows")
}
