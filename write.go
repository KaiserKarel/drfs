package drfs

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Write buffer p to the file. Write will return len(p), nil or an the written bytes, error. Writes are performed
// asynchronously for each bucket within the file. On error, each write after the error is rolled back. If rollbacks
// error, the file must first be recovered by verifying the index integrity and index of each bucket.
func (f *File) Write(p []byte) (int, error) {
	return f.WriteCtx(context.Background(), p)
}

// Write using the provided context for API calls.
func (f *File) WriteCtx(ctx context.Context, p []byte) (int, error) {
	var n int
	for n <= len(p) {
		a, err := f.WriteBatch(context.TODO(), p[n:])
		n += a
		if err != nil || a == 0 {
			return n, err
		}
	}
	return n, nil
}

// WriteBatch writes up to FileOptions.NumThreads * EffectiveReplySize bytes to the drfs file.
func (f *File) WriteBatch(ctx context.Context, p []byte) (int, error) {
	var numbuckets = len(f.index.Buckets)
	var errs = make([]error, numbuckets)
	var segments = slice(p, EffectiveReplySize)
	var written = make([]int, numbuckets)
	var incompleteWrite = false

	grp := sync.WaitGroup{}
	for i := 0; i < numbuckets && i < len(segments); i++ {
		segment := segments[i]
		i := i
		bucket := f.writers.Get()
		grp.Add(1)
		go func() {
			n, err := bucket.WriteCtx(context.TODO(), f.service, f.file.Id, p[segment.lower:segment.upper])
			errs[i] = err
			written[i] = n
			if bucket.Header.Capacity > 0 {
				incompleteWrite = true
			}
			grp.Done()
		}()
	}
	grp.Wait()

	if incompleteWrite {
		f.writers.Ring = f.writers.Prev()
	}

	for k, err := range errs {
		err := err
		if err == nil {
			continue
		}
		// rollback all writes from this error
		// rollback ring to write that first errored
		f.writers.Ring = f.writers.Move(-k)
		bucket := f.writers.Get()

		grp, _ := errgroup.WithContext(context.TODO())
		for i := k; i < len(errs); i++ {
			written[i] = 0
			grp.Go(func() error {
				errRB := bucket.Rollback(context.TODO(), f.service, f.file.Id)
				if err != nil {
					return fmt.Errorf("unable to write: %w [rollback status: %s]", err, errRB)
				}
				return nil
			})
		}
		err = grp.Wait()
		return sum(written), err // an error here is a catastrophic failure.
	}

	return sum(written), nil
}
