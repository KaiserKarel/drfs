package recovery

import (
	"context"
	"os"
	"time"

	"google.golang.org/api/drive/v3"

	"github.com/kaiserkarel/drfs"
)

type FileInfo interface {
	drfs.FileInfo
}

// Stats queries the drive APIs to correctly obtain file info, instead of relying on the indexes
func Stats(file *drfs.File) (os.FileInfo, error) {
	s, err := file.Fstat()
	if err != nil {
		return nil, err
	}
	client, err := file.Service().Take(context.Background(), 1)
	if err != nil {
		return nil, err
	}
	index, err := drfs.IndexFromFile(context.Background(), file.Service(), s.Sys().(*drive.File))
	if err != nil {
		return nil, err
	}

	var length int64

	for _, b := range index.Buckets {
		err = client.RepliesService().
			List(s.ID(), b.CommentID).
			Fields("*").
			PageSize(100).
			Pages(context.Background(), func(list *drive.ReplyList) error {
				for _, reply := range list.Replies {
					if reply.Deleted {
						panic("a deleted reply!")
					}
					length += int64(len([]byte(reply.Content)) - 2)
				}
				return nil
			})
		if err != nil {
			return nil, err
		}
	}

	return &stat{
		fileID:         s.ID(),
		fileName:       s.Name(),
		size:           length,
		quotaBytesUsed: s.QuotaBytesUsed(),
		modtime:        s.ModTime(),
		sys:            s.Sys().(*drive.File),
	}, nil
}

type stat struct {
	fileID         string
	fileName       string
	size           int64
	quotaBytesUsed int64
	modtime        time.Time
	sys            *drive.File
}

func (s *stat) ID() string {
	return s.fileID
}

func (s *stat) Name() string {
	return s.fileName
}

func (s *stat) Size() int64 {
	return s.size
}

func (s *stat) Mode() os.FileMode {
	return os.ModeIrregular
}

func (s *stat) ModTime() time.Time {
	return s.modtime
}

func (s *stat) IsDir() bool {
	return false
}

func (s *stat) Sys() interface{} {
	return s.sys
}

func (s *stat) QuotaBytesUsed() int64 {
	return s.quotaBytesUsed
}
