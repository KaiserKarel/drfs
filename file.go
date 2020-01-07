package drfs

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
)

type File struct {
	docID     string
	commentID string

	service *drive.Service
	file    *drive.File
	comment *drive.Comment
	reply   *drive.Reply
	replIds []string
	cursor  int
}

// maxReplySize defines the maximum number of bytes Drive allows in the content of
// single reply, encoded in UTF-8.
const maxReplySize = 4096

// Open a file by creating or getting an existing file from a comment thread.
func Open(ctx context.Context, service *drive.Service, docID string) (*File, error) {
	desc, err := parse(docID)
	if err != nil {
		return nil, err
	}

	file, err := ensureFileExists(ctx, service, desc.fileName, "*")
	if err != nil {
		return nil, err
	}

	comment, err := ensureCommentExists(ctx, service, file.Id, desc.commentName, "*")
	if err != nil {
		return nil, err
	}

	var replyIDs = make([]string, len(comment.Replies))
	for i, r := range comment.Replies {
		replyIDs[i] = r.Id
	}

	var reply *drive.Reply
	if len(comment.Replies) > 0 {
		reply = comment.Replies[0]
	}

	return &File{docID: file.Id, commentID: comment.Id, service: service, file: file, comment: comment, replIds: replyIDs, reply: reply}, nil
}

func Remove(f *File) error {
	commentService := drive.NewCommentsService(f.service)
	return commentService.Delete(f.file.Id, f.commentID).Fields("").Do()
}

func (f *File) Read(p []byte) (int, error) {
	replyService := drive.NewRepliesService(f.service)
	replyID := f.replIds[f.cursor]

	reply, err := replyService.Get(f.docID, f.commentID, replyID).Fields("*").Do()
	if err != nil {
		return 0, err
	}

	var wi int
	for i := 0; i < len(reply.Content) && i < len(reply.Content); i++ {
		p[wi] = []byte(reply.Content)[i]
		wi++
	}
	f.cursor += 1
	return wi, nil
}

func (f *File) Write(data []byte) (int, error) {
	replyService := drive.NewRepliesService(f.service)

	for _, chunk := range split(data, maxReplySize) {
		reply, err := replyService.Create(f.docID, f.commentID, &drive.Reply{Content: string(chunk)}).Fields("id").Do()
		if err != nil {
			return 0, fmt.Errorf("unable to create new reply: %w", err)
		}
		f.replIds = append(f.replIds, reply.Id)
		f.reply = reply
	}
	return len(data), nil
}

func split(buf []byte, lim int) [][]byte {
	var chunk []byte
	chunks := make([][]byte, 0, len(buf)/lim+1)
	for len(buf) >= lim {
		chunk, buf = buf[:lim], buf[lim:]
		chunks = append(chunks, chunk)
	}
	if len(buf) > 0 {
		chunks = append(chunks, buf[:len(buf)])
	}
	return chunks
}
