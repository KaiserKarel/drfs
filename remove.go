package drfs

import (
	"context"
)

// Remove deletes a file from Drive.
func Remove(file *File) error {
	return RemoveCtx(context.Background(), file)
}

// RemoveCtx deletes a file from Drive using the provided context.
func RemoveCtx(ctx context.Context, file *File) error {
	service, err := file.service.Take(ctx, 1)
	if err != nil {
		return err
	}

	return service.FilesService().
		Delete(file.file.Id).
		Fields("*").
		Context(ctx).
		Do()
}
