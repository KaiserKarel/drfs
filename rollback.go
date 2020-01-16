package drfs

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
)

// ErrNoRollback is returned if rollback is not possible or if it fails. Failing rollbacks are a catastrophic failure.
// The best bet is to save the current and previous bucket states; wait a day for rate limits to regenerate and attempt
// to restore the index and trim buckets.
var ErrNoRollback = fmt.Errorf("rollback not possible")

// RollbackCtx returns a bucket from state new to old by deleting and updating replies. At most there should be a length difference
// of 1 between the buckets. (There is no API for searching reply by number, thus deleting between two arbitrary replies
// is expensive).
func RollbackCtx(ctx context.Context, s Service, fileID string, commentID string, old ThreadHeader, new ThreadHeader) error {
	if new.Length-old.Length > 1 {
		panic("headers should differ by max length 1")
	}

	// remove a created reply
	if old.Tail != new.Tail {
		service, err := s.Take(ctx, 2)
		if err != nil {
			return err
		}

		err = service.RepliesService().
			Delete(fileID, commentID, new.Tail).
			Fields("*").
			Context(ctx).
			Do()
		if err != nil {
			return err
		}

		_, err = service.CommentsService().
			Update(fileID, commentID, &drive.Comment{Content: string(old.MustMarshall())}).
			Fields("*").
			Context(ctx).
			Do()
		return err
	}

	if old.Capacity-new.Capacity < 0 {
		panic("capacity difference < 0 for append rollback")
	}

	service, err := s.Take(ctx, 3)
	if err != nil {
		return err
	}

	// update an appended piece of data
	end := old.Capacity - new.Capacity
	reply, err := service.RepliesService().
		Get(fileID, commentID, old.Tail).
		Context(ctx).
		Fields("*").
		Do()
	if err != nil {
		return err
	}

	reply.Content = reply.Content[:end]
	_, err = service.RepliesService().
		Update(fileID, commentID, old.Tail, reply).
		Fields("*").
		Context(ctx).
		Do()
	if err != nil {
		return err
	}
	_, err = service.CommentsService().
		Update(fileID, commentID, &drive.Comment{Content: string(old.MustMarshall())}).
		Fields("*").
		Context(ctx).
		Do()
	return err
}
