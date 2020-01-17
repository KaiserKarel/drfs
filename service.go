package drfs

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
)

// Service is the interface for obtaining and configuring clients drive clients.
type Service interface {
	// Take provides the api context and approximate number of calls that will be made using
	// the returned client. The service may use a rate limiter to block on Take.
	Take(ctx context.Context, n int) (Client, error)
	Emails() []string
}

// A client is a single API client used to access Drive. Once a client is obtained through the service, no rate limiting
// should be enforced.
type Client interface {
	FilesService() *drive.FilesService
	RepliesService() *drive.RepliesService
	CommentsService() *drive.CommentsService
	DrivesService() *drive.DrivesService
	PermissionsService() *drive.PermissionsService
}

// Create a new reply and update the ThreadHeader. The new ThreadHeader is returned.
//
// This function does not actually alter the reply or bucket, making it possible to retry this with exponential backoff.
func CreateReply(ctx context.Context, s Service, fileID string, bucket Bucket, reply *drive.Reply) (*ThreadHeader, error) {
	service, err := s.Take(ctx, 2)
	if err != nil {
		return nil, err
	}

	r, err := service.RepliesService().
		Create(fileID, bucket.CommentID, reply).
		Context(ctx).
		Fields("*").
		Do()
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}

	bucket.Header.Capacity = MaxReplySize - len(reply.Content)
	bucket.Header.Tail = r.Id
	bucket.Header.Length++

	_, err = service.CommentsService().
		Update(fileID, bucket.CommentID, &drive.Comment{Content: string(bucket.Header.MustMarshall())}).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}
	return &bucket.Header, nil
}

// AppendToReply adds the content to the buckets tail. The caller should ensure that the content of the new
// reply does not exceed EffectiveReplySize.
func AppendToReply(ctx context.Context, s Service, fileID string, bucket Bucket, content string) (*ThreadHeader, error) {
	service, err := s.Take(ctx, 3)
	if err != nil {
		return nil, err
	}

	reply, err := service.RepliesService().
		Get(fileID, bucket.CommentID, bucket.Header.Tail).
		Fields("*").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("get reply: %w", err)
	}

	reply.Content = reply.Content[:len(reply.Content)-1] + content + padding

	// TODO remove this once certain no 1 of errors are present
	if len(reply.Content) > MaxReplySize {
		return nil, fmt.Errorf("reply exceeded max size: %d", len(reply.Content))
	}

	bucket.Header.Capacity = MaxReplySize - len(reply.Content)

	_, err = service.RepliesService().
		Update(fileID, bucket.CommentID, bucket.Header.Tail, reply).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("update reply: %w", err)
	}

	_, err = service.CommentsService().
		Update(fileID, bucket.CommentID, &drive.Comment{Content: string(bucket.Header.MustMarshall())}).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return &bucket.Header, fmt.Errorf("update comment: %w", err)
	}
	return &bucket.Header, nil
}
