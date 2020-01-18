package drfs

import (
	"context"
	"io"
	"sync"
)

func (f *File) Read(p []byte) (int, error) {
	return f.ReadCtx(context.Background(), p)
}

func (f *File) ReadCtx(ctx context.Context, p []byte) (int, error) {
	var n int
	for n <= len(p) {
		a, err := f.ReadBatch(context.TODO(), p[n:])
		n += a
		if err != nil || a == 0 {
			return n, err
		}
	}
	return n, io.EOF
}

func (f *File) ReadBatch(ctx context.Context, p []byte) (int, error) {
	var numbuckets = len(f.index.Buckets)
	var segments = slice(p, EffectiveReplySize)
	var read = make([]int, numbuckets)
	var errs = make([]error, numbuckets)
	var grp sync.WaitGroup

	for i := 0; i < numbuckets && i < len(segments); i++ {
		i := i
		segment := segments[i]
		bucket := f.readers.Get()

		grp.Add(1)
		go func() {
			defer grp.Done()
			n, err := bucket.ReadCtx(context.TODO(), p[segment.lower:segment.upper])
			read[i] = n
			errs[i] = err
		}()
	}
	grp.Wait()

	var total int
	for i, n := range read {
		total += n

		if err := errs[i]; err != nil {
			return total, err
		}
	}
	return total, nil
}
