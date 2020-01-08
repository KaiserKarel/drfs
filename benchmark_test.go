package drfs_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/kaiserkarel/drfs"
	"google.golang.org/api/drive/v3"
)

func init() {
	var err error
	service, err = drive.NewService(context.Background())
	if err != nil {
		log.Fatalf("unable to init drive.Service: %w", err)
	}
}

var service *drive.Service

func BenchmarkFile_Read_Small(b *testing.B) {
	b.StopTimer()
	file, err := drfs.Open(context.Background(), service, "benchmarks/BenchmarkFile_Read_Small")
	if err != nil {
		b.Fatalf("unable to open file: %w", err)
	}

	var data = []byte("Small")

	b.StartTimer()
	for n := 0; n < b.N; n++ {
		_, err = file.Write(data)
		if err != nil {
			b.Fatalf("encountered error while writing: %w", err)
		}
	}
	b.StopTimer()
	drfs.Remove(file)
}

func BenchmarkFile_Lorem_Copy(b *testing.B) {
	b.StopTimer()
	var fileName = "examples/lorem.txt"

	ctx := context.Background()
	driveService, err := drive.NewService(ctx)
	if err != nil {
		panic(err)
	}

	lorem, err := os.Open("lorem.txt")
	if err != nil {
		panic(fmt.Sprintf("unable to open lorem.txt: %v", err))
	}

	drfsFile, err := drfs.Open(ctx, driveService, fileName)
	if err != nil {
		panic(fmt.Sprintf("unable to open file: %v", err))
	}

	b.StartTimer()
	_, err = io.Copy(drfsFile, lorem)
	if err != nil {
		panic(fmt.Sprintf("unable to copy lorem to drfs: %v", err))
	}
	b.StopTimer()
	drfs.Remove(drfsFile)
}
