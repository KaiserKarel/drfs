package os

import (
	"context"
	"fmt"
	"os"

	"github.com/kaiserkarel/drfs"
)

const DefaultNumThreads = 48

// Open either creates or opens the file by filename; first searching for that specific file through the list api.
// It errors if more or less than 1 file(s) are found.
func Open(fileName string) (*drfs.File, error) {
	err := ensure()
	if err != nil {
		return nil, err
	}

	client, err := service.Take(context.Background(), 1)
	if err != nil {
		return nil, err
	}

	resp, err := client.FilesService().List().Q(fmt.Sprintf("name = '%s'", fileName)).Fields("*").Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Files) == 0 {
		file, err := drfs.CreateFileCtx(context.Background(), service, fileName, drfs.FileOptions{NumThreads: DefaultNumThreads})
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	if len(resp.Files) > 1 {
		return nil, fmt.Errorf("multiple files match name: %s", fileName)
	}

	return drfs.OpenCtx(context.Background(), resp.Files[0], service)
}

func OpenFile(fileName string, _ int, _ os.FileMode) (*drfs.File, error) {
	return Open(fileName)
}
