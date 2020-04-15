package main

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"log"
	"os"
	"strings"
	"time"
)

func check(err error)  {
	if err != nil {
		log.Fatalf("%s", err)
	}
}

func is500(err error) bool  {
	if err == nil {
		return false
	}

	var apiError googleapi.Error
	if errors.Is(err, &apiError); apiError.Code == 500 {
		return true
	}

	if strings.Contains(err.Error(), "500") {
		return true
	}
	return false
}

func main() {
	service, err := drive.NewService(context.Background())
	if err != nil {
		log.Fatalf("unable to init drive service: %s", err)
	}

	fileService := drive.NewFilesService(service)
	commentService := drive.NewCommentsService(service)
	repliesService := drive.NewRepliesService(service)


	f, err := fileService.Create(&drive.File{Name: time.Now().String()}).Fields("*").Do()
	check(err)
	fmt.Printf("fileID: %s\n", f.Id)

	c, err := commentService.Create(f.Id, &drive.Comment{Content: "."}).Fields("*").Do()
	check(err)
	fmt.Printf("commentID: %s\n", c.Id)

	sleep := 0

	fname := fmt.Sprintf("../data/replies_%d_%s.csv", sleep, time.Now().String())
	file, err := os.OpenFile(fname, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err)
	defer file.Close()

	_, err = file.WriteString(f.Id+"\n")
	check(err)

	_, err = file.WriteString(c.Id+"\n")
	check(err)
	for i := 0; i < 100000; i++ {
		r, err := repliesService.Create(f.Id, c.Id, &drive.Reply{Content: "."}).Fields("*").Do()
		if is500(err) {
			time.Sleep(time.Duration(sleep)*time.Second)
			continue
		}
		check(err)
		_, err =
			file.WriteString(fmt.Sprintf("%s, %s, %s\n", r.Id, r.CreatedTime, r.ModifiedTime))
		check(err)
		check(file.Sync())
		time.Sleep(time.Duration(sleep)*time.Second)
	}
}
