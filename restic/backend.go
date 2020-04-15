package restic

import (
	"context"
	"errors"
	"fmt"
	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	"github.com/kaiserkarel/qstring"
	"google.golang.org/api/drive/v3"
	"io"
	"log"
	"os"
)

const (
	ResticFileType = "resticFileType"
)

type Backend struct {
	backend.Layout
	Config

	service drfs.Service
}

// Location describes the type of backend, including a description of the
// Drive account used.
func (b *Backend) Location() string  {
	return "drfs"
}

// Test checks if the given file exists
func (b *Backend) Test(ctx context.Context, h restic.Handle) (bool, error) {
	f, err := b.getFile(context.Background(), h)
	if b.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return f != nil, nil
}

// Removes a file described by the absolute name of h.
func (b *Backend) Remove(ctx context.Context, h restic.Handle) error  {
	client, err := b.service.Take(context.Background(), 1)
	if err != nil {
		return err
	}

	var query = qstring.Name().Equals(h.Name).And().Properties().Has(ResticFileType, string(h.Type))
	err = client.FilesService().List().Q(query.String()).Pages(context.Background(), func(list *drive.FileList) error {
		for _, f := range list.Files {
			client, err := b.service.Take(context.Background(), 1)
			if err != nil {
				return err
			}

			err = client.FilesService().Delete(f.Id).Context(context.Background()).Do()
			if err != nil {
				return err
			}

			if list.IncompleteSearch {
				_, err := b.service.Take(context.Background(), 1)
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
	panic("Close")
	return nil
}

// Save appends the date of a given rd to the file described by the handle.
func (b *Backend) Save(ctx context.Context, h restic.Handle, rd restic.RewindReader) error  {
	if err := h.Valid(); err != nil {
		return err
	}

	file, err := b.getOrCreateFile(context.Background(), h)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, rd)
	return err
}

// Load runs fn with a reader that yields the contents of the file at h at the
// given offset. If length is larger than zero, only a portion of the file
// is read.
//
// The function fn may be called multiple times during the same Load invocation
// and therefore must be idempotent.
func (b *Backend) Load(ctx context.Context, h restic.Handle, length int, offset int64, fn func(rd io.Reader) error) error  {
	return backend.DefaultLoad(context.Background(), h, length, offset, b.openReader, fn)
}

func (b *Backend) openReader(ctx context.Context, h restic.Handle, length int, offset int64) (io.ReadCloser, error)  {

	file, err := b.getOrCreateFile(context.Background(), h)
	if err != nil {
		return nil, err
	}

	if offset > 0 {
		_, err = file.Seek(offset, 0)
		if err != nil {
			return nil, err
		}
	}

	if length > 0 {
		return backend.LimitReadCloser(file, int64(length)), nil
	}
	return file, nil
}



// Stat returns information about the File identified by h.
func (b *Backend) Stat(ctx context.Context, h restic.Handle) (restic.FileInfo, error)  {
	file, err := b.getFile(context.Background(), h)
	if err != nil {
		return restic.FileInfo{}, err
	}

	stat, err := file.Fstat()
	if err != nil {
		return restic.FileInfo{}, err
	}

	return restic.FileInfo{
		Size: stat.Size(),
		Name: stat.Name(),
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

	client, err := b.service.Take(context.Background(), 1)
	if err != nil {
		return err
	}

	err = client.FilesService().List().Q(query).Pages(context.Background(), func(list *drive.FileList) error {
		for _, file := range list.Files {
			f, err := drfs.OpenCtx(context.Background(), file, b.service)
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
			_, err := b.service.Take(context.Background(), 1)
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

func (b *Backend) getFile(ctx context.Context, h restic.Handle) (*drfs.File, error)  {
	filename := b.Layout.Filename(h)

	client, err := b.service.Take(ctx, 1)
	if err != nil {
		return nil, err
	}

	query := qstring.Name().Equals(filename)

	resp, err := client.FilesService().List().Q(query.String()).Fields("*").Context(ctx).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Files) > 1 {
		return nil, errors.New("multiple files found for given filename")
	}

	if len(resp.Files) == 0 {
		return nil, drfs.ErrNotFound
	}

	return drfs.OpenCtx(ctx, resp.Files[0], b.service)
}

func (b *Backend) createFile(ctx context.Context, h restic.Handle) (*drfs.File, error)  {
	filename := b.Layout.Filename(h)
	log.Printf("creating file: %s", filename)

	f, err := drfs.CreateFileCtx(
		context.Background(), b.service, filename, drfs.FileOptions{
			NumThreads:1,
			Properties: map[string]string{
				ResticFileType: string(h.Type),
			},
		},)

	stat, _ := f.Fstat()
	os.Stdout.WriteString(	fmt.Sprintf("created file: %s", stat.ID()))

	if err != nil {
		return nil, err
	}
	return f, nil
}

func (b *Backend) getOrCreateFile(ctx context.Context, h restic.Handle) (*drfs.File, error)  {
	f, err := b.getFile(context.Background(), h)
	if errors.Is(err, drfs.ErrNotFound) {
		f, err = b.createFile(context.Background(), h)
	}

	if err != nil {
		return nil, err
	}
	return f, nil
}