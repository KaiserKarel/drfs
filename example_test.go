package drfs_test

import (
	"context"
	"fmt"
	// "io"
	// "os"

	"google.golang.org/api/drive/v3"

	"github.com/kaiserkarel/drfs"
)

func ExampleFile_Write() {
	ctx := context.Background()
	driveService, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	var data = []byte("Hello there.")
	var fileName = "examples/fileWrite"

	file, err := drfs.Open(ctx, driveService, fileName)
	if err != nil {
		panic(fmt.Sprintf("unable to open file: %v", err))
	}

	_, err = file.Write(data)
	if err != nil {
		panic(fmt.Sprintf("unable to write to file: %v", err))
	}

	var buf = make([]byte, len(data))
	_, err = file.Read(buf)
	if err != nil {
		panic(fmt.Sprintf("unable to read from file: %v", err))
	}

	err = drfs.Remove(file)
	if err != nil {
		panic(fmt.Sprintf("unable to remove file: %v", err))
	}
	fmt.Println(string(buf))
	// Output:
	// Hello there.
}

func ExampleFile_Reopen() {
	ctx := context.Background()
	driveService, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	var data = []byte("Written earlier")
	var data2 = []byte("Written now")

	var fileName = "examples/fileReopen"

	oldFile, err := drfs.Open(ctx, driveService, fileName)
	if err != nil {
		panic(fmt.Sprintf("unable to open file: %v", err))
	}

	_, err = oldFile.Write(data)
	if err != nil {
		panic(fmt.Sprintf("unable to write to file: %v", err))
	}

	newFile, err := drfs.Open(ctx, driveService, fileName)
	if err != nil {
		panic(fmt.Sprintf("unable to open file: %v", err))
	}

	_, err = newFile.Write(data2)
	if err != nil {
		panic(fmt.Sprintf("unable to write to file: %v", err))
	}

	var buf1 = make([]byte, len(data))
	_, err = newFile.Read(buf1)
	if err != nil {
		panic(fmt.Sprintf("unable to read from file: %v", err))
	}

	var buf2 = make([]byte, len(data2))
	_, err = newFile.Read(buf2)
	if err != nil {
		panic(fmt.Sprintf("unable to read from file: %v", err))
	}

	err = drfs.Remove(newFile)
	if err != nil {
		panic(fmt.Sprintf("unable to remove file: %v", err))
	}

	fmt.Println(string(buf1))
	fmt.Println(string(buf2))
	// Output:
	// Written earlier
	// Written now
}
