package drfs

import (
	"context"
	"os"
	"time"

	"google.golang.org/api/drive/v3"
)

// FileInfo is a complete description of a File
type FileInfo interface {
	ID() string
	QuotaBytesUsed() int64

	os.FileInfo
}

// Fstat returns the full file stats.
func (f *File) Fstat() (FileInfo, error) {
	client, err := f.service.Take(context.Background(), 1)
	if err != nil {
		return nil, err
	}

	refresh, err := client.FilesService().Get(f.file.Id).Fields("*").Do()
	if err != nil {
		return nil, err
	}

	f.file = refresh

	return &stat{
		fileID:   f.file.Id,
		fileName: f.file.Name,
		size:     f.size(),
		modtime:  f.modTime(),
		sys:      f.file,
	}, nil
}

// Stat returns file stats, mimicking the os API
func (f *File) Stat() (os.FileInfo, error) {
	return &stat{
		fileID:  f.file.Name,
		size:    f.size(),
		modtime: f.modTime(),
		sys:     f.file,
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
