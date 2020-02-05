package restic

import (
	"context"
	"errors"
	"fmt"
	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	"github.com/kaiserkarel/qstring"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
)

const (
	ResticFileType = "resticFileType"
)

type Backend struct {
	service drfs.Service
}

// Location describes the type of backend, including a description of the
// Drive account used.
func (b *Backend) Location() string  {
	return "drfs"
}

// Test checks if the given file exists
func (b *Backend) Test(ctx context.Context, h restic.Handle) (bool, error) {
	client, err := b.service.Take(ctx, 1)
	if err != nil {
		return false, err
	}
	f, err := client.FilesService().Get(abs(h)).Fields("id").Context(ctx).Do()
	if err != nil {
		return false, err
	}
	return f != nil, nil
}

// Removes a file described by the absolute name of h.
func (b *Backend) Remove(ctx context.Context, h restic.Handle) error  {
	client, err := b.service.Take(ctx, 1)
	if err != nil {
		return err
	}

	var query = qstring.Name().Equals(h.Name).And().Properties().Has(ResticFileType, string(h.Type))
	err = client.FilesService().List().Q(query.String()).Pages(ctx, func(list *drive.FileList) error {
		for _, f := range list.Files {
			client, err := b.service.Take(ctx, 1)
			if err != nil {
				return err
			}

			err = client.FilesService().Delete(f.Id).Context(ctx).Do()
			if err != nil {
				return err
			}

			if list.IncompleteSearch {
				_, err := b.service.Take(ctx, 1)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	return err
}

// Close is a noop.
func (b *Backend) Close() error  {
	return nil
}

// Save appends the date of a given rd to the file described by the handle.
func (b *Backend) Save(ctx context.Context, h restic.Handle, rd restic.RewindReader) error  {
	return nil
}

// Load runs fn with a reader that yields the contents of the file at h at the
// given offset. If length is larger than zero, only a portion of the file
// is read.
//
// The function fn may be called multiple times during the same Load invocation
// and therefore must be idempotent.
//
// Implementations are encouraged to use backend.DefaultLoad
func (b *Backend) Load(ctx context.Context, h restic.Handle, length int, offset int64, fn func(rd io.Reader) error) error  {
	return nil
}

// Stat returns information about the File identified by h.
func (b *Backend) Stat(ctx context.Context, h restic.Handle) (restic.FileInfo, error)  {
	var query = qstring.Name().Equals(h.Name).And().Properties().Has(ResticFileType, string(h.Type)).String()

	client, err := b.service.Take(ctx, 1)
	if err != nil {
		return restic.FileInfo{}, err
	}

	f, err := client.FilesService().List().Q(query).Context(ctx).Do()
	if err != nil {
		return restic.FileInfo{}, err
	}

	if len(f.Files) == 0 {
		return restic.FileInfo{}, os.ErrNotExist
	}

	if len(f.Files) > 1 {
		return restic.FileInfo{}, fmt.Errorf("encountered multiple files matching the query: %d", len(f.Files))
	}

	file, err := drfs.OpenCtx(ctx, f.Files[0], b.service)
	if err != nil {
		return restic.FileInfo{}, err
	}

	stat, err := file.Fstat()
	if err != nil {
		return restic.FileInfo{}, err
	}

	return restic.FileInfo{
		Size: stat.Size(),
		Name: h.Name,
	}, nil
}

// List runs fn for each file in the backend which has the type t. When an
// error occurs (or fn returns an error), List stops and returns it.
//
// The function fn is called exactly once for each file during successful
// execution and at most once in case of an error.
//
// The function fn is called in the same Goroutine that List() is called
// from.
func (b *Backend) List(ctx context.Context, t restic.FileType, fn func(restic.FileInfo) error) error  {
	var query = qstring.Properties().Has(ResticFileType, string(t)).String()

	client, err := b.service.Take(ctx, 1)
	if err != nil {
		return err
	}

	err = client.FilesService().List().Q(query).Pages(ctx, func(list *drive.FileList) error {
		for _, file := range list.Files {
			f, err := drfs.OpenCtx(ctx, file, b.service)
			if err != nil {
				return err
			}

			stat, err := f.Stat()
			if err != nil {
				return err
			}

			err = fn(restic.FileInfo{
				Size: stat.Size(),
				Name: file.Name,
			})
			if err != nil {
				return err
			}
		}
		if list.IncompleteSearch {
			_, err := b.service.Take(ctx, 1)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// IsNotExist returns true if the error was caused by a non-existing file
// in the backend.
func (b *Backend) IsNotExist(err error) bool  {
	return errors.Is(err, drfs.ErrNotFound)
}

// Delete removes all data in the backend. This deletes every instance of a
// Google Drive file which has a properly canonicalized filename.
func (b *Backend) Delete(ctx context.Context) error  {
	return errors.New("backend delete unimplemented")
}