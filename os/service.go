package os

import (
	"context"
	"os"
	"sync"

	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/drive"
)

var service drfs.Service
var serviceErr error

const (
	DRFS_CREDS = "DRFS_APPLICATION_CREDENTIALS"
)

var once = sync.Once{}

func ensure() error {
	once.Do(func() {
		dir, ok := os.LookupEnv(DRFS_CREDS)
		if ok {
			serviceFromDir(dir)
		} else {
			defaultService()
		}
	})
	return serviceErr
}

func serviceFromDir(dir string) {
	creds, err := drive.CredentialsFromDirectory(context.Background(), dir)
	if err != nil {
		serviceErr = err
	}

	service, serviceErr = drive.NewService(context.Background(), creds...)
}

func defaultService() {
	service, serviceErr = drive.NewService(context.Background())
}
