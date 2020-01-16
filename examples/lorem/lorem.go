package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/kaiserkarel/drfs/drive"

	"github.com/kaiserkarel/drfs"
	"github.com/udhos/equalfile"
)

func main() {
	var localSrc = "lorem.txt"
	var localDest = "lorem.download.txt"
	var fileID = "1Ux4VckG-RUvle-wLQB_v65Q7bhShNh65"

	// log.Println("Searching for credentials.")
	// creds, err := drive.CredentialsFromDirectory(context.TODO(), "secrets")
	// if err != nil {
	// 	log.Fatalf("unable to read credentials: %w", err)
	// }

	// log.Printf("Using %d different credentials.", len(creds))
	service, err := drive.NewService(context.Background())
	if err != nil {
		log.Fatalf("unable to init drive.Service: %w", err)
	}

	file, err := createFile(service)
	if err != nil {
		log.Fatalf("unable to create file:  %s: %s. Perhaps create it first and update this example.", fileID, err)
	}

	// s, err := service.Take(context.TODO(), 1)
	// if err != nil {
	// 	log.Fatalf("unable to take service: %s", err)
	// }

	// driveFile, err := s.FilesService().Get(fileID).Fields("*").Do()
	// if err != nil {
	// 	log.Fatalf("unable to create file:  %s: %s. Perhaps create it first and update this example.", fileID, err)
	// }

	// file, err := drfs.OpenCtx(context.TODO(), driveFile, service)
	// if err != nil {
	// 	log.Fatalf("unable to index file: %s", err)
	// }

	stat, _ := file.Stat()
	log.Printf("initial local stats: %+v", stat)

	log.Println("starting upload")
	err = upload(localSrc, file)
	if err != nil {
		log.Fatalf("unable to upload: %s", err)
	}

	stat, _ = file.Stat()
	log.Printf("final local stats: %+v", stat)

	log.Println("opening dest")
	f, err := os.OpenFile(localDest, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Fatalf("unable to open src: %v", err)
	}

	log.Println("downloading")
	_, err = io.Copy(f, file)
	if err != nil {
		log.Fatalf("unable to copy file to drive: %v", err)
	}

	cmp := equalfile.New(make([]byte, 10000), equalfile.Options{
		Debug:         false,
		ForceFileRead: true,
		MaxSize:       10000000,
	})

	log.Println("closing local dest")
	err = f.Close()
	if err != nil {
		log.Fatalf("unable to close localDest: %v", err)
	}

	ok, err := cmp.CompareFile(localSrc, localDest)
	if err != nil {
		log.Fatalf("unable to compare files: %w", err)
	}

	if !ok {
		log.Fatalf("files do not match!")
	}

	log.Println("files match!")
}

func createFile(service drfs.Service) (*drfs.File, error) {
	var name = "lorem_" + time.Now().String()

	file, err := drfs.CreateFileCtx(context.TODO(), service, name, drfs.FileOptions{NumThreads: 4})
	if err != nil {
		return nil, err
	}

	stat, _ := file.Stat()

	log.Printf("created file! id: %s\n", stat.Name())
	return file, nil
}

func upload(local string, file *drfs.File) error {
	l, err := os.Open(local)
	if err != nil {
		return err
	}
	n, err := io.Copy(file, l)
	log.Printf("copied %d bytes", n)
	return err
}
