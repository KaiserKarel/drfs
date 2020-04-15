package drfs

import (
	"context"
	"golang.org/x/sync/errgroup"
)

func (f *File) jumpToStart(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	for _, thread := range f.index.Buckets {
		grp.Go(func() error {
			return thread.jumpToStart(ctx)
		})
	}

	err := grp.Wait()
	if err != nil {
		return err
	}

	f.readers = newThreadRing(f.index.Buckets)
	return nil
}
