package drfs

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"sync/atomic"

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
	log.Printf("requested: %d", len(p))
	var n int
	for n <= len(p) {
		if len(p[n:]) == 0 {
			return n, nil
		}
		a, err := f.WriteBatch(context.TODO(), p[n:])
		log.Printf("WriteCtx: a %d n %d err%s", a, n, err)
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
	var written = make([]int, numbuckets)
	var incompleteWrite = newAtomicCheck()
	grp := &sync.WaitGroup{}

	var put = func(thread *Thread, buf []byte, i int) {
		grp.Add(1)
		log.Printf("b %d payload: %v", thread.Header.Number, len(buf))

		go func() {
			var err error

			if thread.Capacity() == 0 {
				err = thread.Put(context.TODO(), buf)
			} else {
				err = thread.Update(context.TODO(), buf)
			}

			errs[i] = err
			written[i] = len(buf)
			if thread.Header.Capacity > 0 {
				incompleteWrite.set()
			}
			grp.Done()
		}()
	}

	last := f.writers.Peek()
	var offset int
	var skip int
	if last.Capacity() > 0 {
		skip = 1
		offset = min(len(p), last.Capacity())
		f.writers.Next()
		put(last, p[:offset], 0)
	}

	var remaining = p[offset:]
	var segments = slice(remaining, EffectiveReplySize)

	for i := skip; i < numbuckets && i < len(segments); i++ {
		payload := remaining[segments[i-skip].lower:segments[i-skip].upper]
		put(f.writers.Get(), payload, i)
	}
	grp.Wait()

	if incompleteWrite.isSet() {
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

func slice(p []byte, size int) []bounds {
	var b []bounds
	length := len(p)
	segmentCount := int(math.Ceil(float64(length) / float64(size)))
	var start, stop int
	for i := 0; i < segmentCount; i += 1 {
		start = i * size
		stop = start + size
		if stop > length {
			stop = length
		}
		b = append(b, bounds{lower: start, upper: stop})
	}
	return b
}

func sum(p []int) int {
	var res int
	for _, i := range p {
		res += i
	}
	return res
}

type atomicCheck struct {
	flag int32
}

func newAtomicCheck() *atomicCheck {
	return &atomicCheck{0}
}

func (c *atomicCheck) isSet() bool {
	return c.flag > 0
}

func (c *atomicCheck) set() {
	if !atomic.CompareAndSwapInt32(&c.flag, 0, 1) {
		panic("setting an already set check!") // indicates programming error within this library
	}
}
