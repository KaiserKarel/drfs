package tests

import (
	"context"
	"testing"

	"google.golang.org/api/drive/v3"
)

// TestLeadingSpaceDissapears was created since sometimes spaces would disappear when comparing two files. This was due
// to Drive removing leading spaces in comments. As a result the total size of a reply available to us is 4095 bytes;
// or the drfs.EffectiveReplySize.
func TestLeadingEndingSpaceDissapears(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test")
	}
	service, err := drive.NewService(context.Background())
	if err != nil {
		t.Fatalf("unable to init drive service: %s", err)
	}

	fileService := drive.NewFilesService(service)
	commentService := drive.NewCommentsService(service)
	repliesService := drive.NewRepliesService(service)

	f, err := fileService.Create(&drive.File{Name: "TestLeadingSpaceDissapears"}).Fields("*").Do()
	if err != nil {
		t.Fatalf("cannot create file: %s", err)
	}
	t.Log("fileID: %", f.Id)

	comment, err := commentService.Create(f.Id, &drive.Comment{Content: "Test"}).Fields("*").Do()
	if err != nil {
		t.Fatalf("cannot create comment: %s", err)
	}

	var content = " hello "

	reply, err := repliesService.Create(f.Id, comment.Id, &drive.Reply{Content: content}).Fields("*").Do()
	if err != nil {
		t.Fatalf("cannot create reply: %s", err)
	}

	if reply.Content != "hello" {
		t.Errorf("whitespace was not stripped! expected: [%s], but got: [%s]", content, reply.Content)
	}
}
