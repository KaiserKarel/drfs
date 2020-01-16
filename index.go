package drfs

import (
	"context"
	"errors"
	"sort"
	"strings"

	"google.golang.org/api/drive/v3"
)

var (
	ErrMissingFileHeader = errors.New("fileheader missing")
)

// Index describes the structure of a file.
type Index struct {
	Header  FileHeader
	Buckets []*Bucket
}

// IndexFromFile queries the buckets from a file to generate an Index.
func IndexFromFile(ctx context.Context, s Service, file *drive.File) (*Index, error) {
	var fileheader *FileHeader
	var buckets []*Bucket

	service, err := s.Take(ctx, 6) // 512 comments is the default per drfsFile. 100 pages per pagination means at most
	// it will take 6 calls in the paginator
	if err != nil {
		return nil, err
	}

	err = service.CommentsService().List(file.Id).Fields("*").PageSize(MaxPages).Pages(ctx, func(list *drive.CommentList) error {
		for _, comment := range list.Comments {
			payload := strings.NewReader(comment.Content)
			threadheader, err := ThreadHeaderFromJSON(payload)
			if err != nil {
				// possibly the file header. Check if we already encountered it. If so error anyway, else try to decode.
				if fileheader != nil {
					return err
				}

				payload := strings.NewReader(comment.Content)
				fileheader, err = FileHeaderFromJSON(payload)
				if err != nil {
					return err
				}
				continue
			}

			buckets = append(buckets, &Bucket{
				CommentID: comment.Id,
				Header:    *threadheader,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if fileheader == nil {
		return nil, ErrMissingFileHeader
	}

	sort.Sort(byHeaderNumber(buckets))

	return &Index{
		Header:  *fileheader,
		Buckets: buckets,
	}, nil
}

type byHeaderNumber []*Bucket

func (s byHeaderNumber) Len() int {
	return len(s)
}
func (s byHeaderNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byHeaderNumber) Less(i, j int) bool {
	return s[i].Header.Number < s[j].Header.Number
}
