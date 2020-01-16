package drfs

import (
	"time"
)

func (f *File) modTime() time.Time {
	var mod time.Time
	for _, b := range f.index.Buckets {
		if b.modTime.After(mod) {
			mod = b.modTime
		}
	}
	return mod
}
