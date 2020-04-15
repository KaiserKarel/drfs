package drfs

import (
	"context"
	"errors"
	"io"
)

// Seek sets the read offset for subsequent reads. Seeking outside of the file length will set the
// offset to the EOF. Seek does not yet allow to increment the write cursor.
func (f *File) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		return f.seekStart(offset)
	case io.SeekCurrent:
		return f.seekCurrent(offset)
	case io.SeekEnd:
		return 0, errors.New("SeekEnd not supported")
	default:
		panic("invalid whence")
	}
}

func (f *File) seekCurrent(offset int64) (int64, error) {
	if offset < 0 {
		return 0, errors.New("seeking to offset before start of file")
	}

	// compute offset for each thread
	var offsets = make([]int64, len(f.index.Buckets))
	for {
		for i, _ := range offsets {
			if offset - EffectiveReplySize < 0 {
				offsets[i] = offsets[i] + offset
				offset = 0
				break
			}

			offset -= EffectiveReplySize
			offsets[i] = offsets[i] + EffectiveReplySize
		}
		if offset == 0 {
			break
		}
	}

	var pos int64
	for i := 0; i < len(f.index.Buckets); i++  {
		r := f.readers.Get()
		n, err := r.seekCurrent(offsets[i])
		pos += n
		if err != nil {
			return pos, err
		}
	}
	return pos, nil
}

func (f *File) seekStart(offset int64) (int64, error) {
	ctx := context.Background()
	err := f.jumpToStart(ctx)
	if err != nil {
		return 0, err
	}
	return f.seekCurrent(offset)
}


