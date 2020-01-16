package drive

import (
	"golang.org/x/time/rate"
	"google.golang.org/api/drive/v3"
)

type Client struct {
	Secret  Secret
	Limiter *rate.Limiter
	service *drive.Service
	i       int
}

func (s *Client) FilesService() *drive.FilesService {
	return drive.NewFilesService(s.service)
}

func (s *Client) RepliesService() *drive.RepliesService {
	return drive.NewRepliesService(s.service)
}

func (s *Client) CommentsService() *drive.CommentsService {
	return drive.NewCommentsService(s.service)
}

func (s *Client) DrivesService() *drive.DrivesService {
	return drive.NewDrivesService(s.service)
}

func (s *Client) PermissionsService() *drive.PermissionsService {
	return drive.NewPermissionsService(s.service)
}
