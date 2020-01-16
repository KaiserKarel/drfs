package drfs

import (
	"context"
	"fmt"
	"io"
	"os"
)

// Copy a local file to newly created drfs.File.
func Copy(src, dst string, service Service) error {
	file, err := CreateFileCtx(context.TODO(), service, dst, FileOptions{NumThreads: 20})
	if err != nil {
		return fmt.Errorf("unable to create drive file: %w", err)
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("unable to open src: %v", err)
	}

	_, err = io.Copy(file, f)
	if err != nil {
		return fmt.Errorf("unable to copy file to drive: %v", err)
	}
	return nil
}

// // Download a file from drfs to a local file.
// func Download(src, dst string, service Service) error {
// 	file, err := OpenCtx(context.TODO(), service, src, ThreadOption{PageSize: 100})
// 	if err != nil {
// 		return fmt.Errorf("unable to open drive file: %w", err)
// 	}
//
// 	f, err := os.Open(dst)
// 	if err != nil {
// 		return fmt.Errorf("unable to open src: %v", err)
// 	}
//
// 	_, err = io.Copy(f, file)
// 	if err != nil {
// 		return fmt.Errorf("unable to copy file to drive: %v", err)
// 	}
// 	return nil
// }
