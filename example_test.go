package drfs_test

import (
	"context"
	"log"

	"github.com/kaiserkarel/drfs"
	"github.com/kaiserkarel/drfs/drive"
)

func ExampleFileV2_Write_Read() {
	ctx := context.Background()
	service, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	file, err := drfs.CreateFileCtx(ctx, service, "test", drfs.FileOptions{NumThreads: 5})
	if err != nil {
		log.Fatalf("error while creating file: %v", err)
	}

	_, err = file.Write([]byte("hello world"))
	if err != nil {
		log.Fatalf("error while writing to file: %v", err)
	}

	buf := make([]byte, 100)
	n, err := file.Read(buf)
	if err != nil {
		log.Fatalf("error while reading from file: %v read: %d", err, n)
	}

	log.Println(string(buf))

	err = drfs.Remove(file)
	if err != nil {
		log.Fatalf("error while removing file: %v", err)
	}
	// Output:
}

// func ExampleFile_Write() {
// 	ctx := context.Background()
// 	driveService, err := drive.NewService(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	var data = []byte("Hello there.")
// 	var fileName = "examples/fileWrite"
//
// 	file, err := drfs.Open(ctx, driveService, fileName)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to open file: %v", err))
// 	}
//
// 	_, err = file.Write(data)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to write to file: %v", err))
// 	}
//
// 	var buf = make([]byte, len(data))
// 	_, err = file.Read(buf)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to read from file: %v", err))
// 	}
//
// 	err = drfs.Remove(file)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to remove file: %v", err))
// 	}
// 	fmt.Println(string(buf))
// 	// Output:
// 	// Hello there.
// }
//
// func ExampleFile_Reopen() {
// 	ctx := context.Background()
// 	driveService, err := drive.NewService(ctx)
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	var data = []byte("Written earlier")
// 	var data2 = []byte("Written now")
//
// 	var fileName = "examples/fileReopen"
//
// 	oldFile, err := drfs.Open(ctx, driveService, fileName)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to open file: %v", err))
// 	}
//
// 	_, err = oldFile.Write(data)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to write to file: %v", err))
// 	}
//
// 	newFile, err := drfs.Open(ctx, driveService, fileName)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to open file: %v", err))
// 	}
//
// 	_, err = newFile.Write(data2)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to write to file: %v", err))
// 	}
//
// 	var buf1 = make([]byte, len(data))
// 	_, err = newFile.Read(buf1)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to read from file: %v", err))
// 	}
//
// 	var buf2 = make([]byte, len(data2))
// 	_, err = newFile.Read(buf2)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to read from file: %v", err))
// 	}
//
// 	err = drfs.Remove(newFile)
// 	if err != nil {
// 		panic(fmt.Sprintf("unable to remove file: %v", err))
// 	}
//
// 	fmt.Println(string(buf1))
// 	fmt.Println(string(buf2))
// 	// Output:
// 	// Written earlier
// 	// Written now
// }
