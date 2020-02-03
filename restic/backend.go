package restic

import (
	"context"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	"io"
)

type Backend struct {

}

func (b *Backend) Location() string  {
	return "drfs"
}

func (b *Backend) Test(ctx context.Context, h restic.Handle) (bool, error) {
	return false, nil
}

func (b *Backend) Remove(ctx context.Context, h restic.Handle) error  {
	return nil
}

func (b *Backend) Close() error  {
	return nil
}

func (b *Backend) Save(ctx context.Context, h restic.Handle, rd restic.RewindReader) error  {
	return nil
}

func (b *Backend) Load(ctx context.Context, h restic.Handle, length int, offset int64, fn func(rd io.Reader) error) error  {
	return nil
}

func (b *Backend) Stat(ctx context.Context, h restic.Handle) (restic.FileInfo, error)  {
	return restic.FileInfo{}, nil
}

func (b *Backend) List(ctx context.Context, t restic.FileType, fn func(restic.FileInfo) error) error  {
	return nil
}

func (b *Backend) IsNotExist(err error) bool  {
	return true
}

func (b *Backend) Delete(ctx context.Context) error  {
	return nil
}