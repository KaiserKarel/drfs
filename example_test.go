package drfs_test

import (
	"context"
	"fmt"
	"github.com/kaiserkarel/gdfs"
	"google.golang.org/api/drive/v3"
	"io"
	"os"
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

func ExampleFile_Copy_To_Gdfs() {
	var fileName = "examples/lorem.txt"

	// This test takes about 150s, thus we don't run automatically run it.
	if !(os.Getenv("LONG_TESTS") == "TRUE") {
		fmt.Println("Skipping ExampleFile_Copy_To_Gdfs")
		return
	}

	ctx := context.Background()
	driveService, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	lorem, err := os.Open("lorem.txt")
	if err != nil {
		panic(fmt.Sprintf("unable to open lorem.txt: %v", err))
	}

	gdfsFile, err := drfs.Open(ctx, driveService, fileName)
	if err != nil {
		panic(fmt.Sprintf("unable to open file: %v", err))
	}

	_, err = io.Copy(gdfsFile, lorem)
	if err != nil {
		panic(fmt.Sprintf("unable to copy lorem to drfs: %v", err))
	}
	// Output:
}
