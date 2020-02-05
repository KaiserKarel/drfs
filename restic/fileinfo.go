package restic

import (
	"errors"
	"fmt"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	"strings"
)

const Prefix = "restic::"

var ErrNoParse = errors.New("incorrect format")

// Convert a restic.Handle into a Google Drive file name, preserving the
// file type.
func abs(h restic.Handle) string  {
	return fmt.Sprintf("%s_%s", h.Type, h.Name)
}

// Canonicalize a Google Drive file name into a restic Handle.
// Also validates the handle.
func canon(abs string) (*restic.Handle, error)  {
	splitted := strings.SplitN(abs, "::", 3)
	if len(splitted) != 3 {
		return nil, ErrNoParse
	}

	if splitted[0] != "" {

	}

	handle :=  &restic.Handle{
		Type: restic.FileType(splitted[0]),
		Name: splitted[1],
	}

	if err := handle.Valid(); err != nil {
		return nil, err
	}
	return handle, nil
}