package drfs

import (
	"context"

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
