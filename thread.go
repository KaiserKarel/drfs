package drfs

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"google.golang.org/api/drive/v3"

	"github.com/google/uuid"
)

type ThreadOption struct {
	// PageSize used in read requests. Defaults to 20. 0 is illegal.
	PageSize int64
}

type ThreadHeader struct {
	Number   int       `json:"n"`
	Length   int64     `json:"l"`
	Tail     string    `json:"t"`
	Capacity int       `json:"c"`
	UUID     uuid.UUID `json:"u"`
}

func (t ThreadHeader) MustMarshall() []byte {
	p, err := json.Marshal(t)
	if err != nil {
		panic(fmt.Errorf("marshaling theadheader failed: %w", err))
	}
	return p
}

func ThreadHeaderFromJSON(p io.Reader) (*ThreadHeader, error) {
	dec := json.NewDecoder(p)
	dec.DisallowUnknownFields()

	header := &ThreadHeader{}
	err := dec.Decode(header)
	return header, err
}

const padding = "1"

type Thread struct {
	FileID    string
	CommentID string       `json:"c"`
	Header    ThreadHeader `json:"h"`

	service  Service
	cursor   int // current read position within reply
	ri       int // index of current reply being read.
	replies  *drive.ReplyList
	oldState *ThreadHeader
	modTime  time.Time
}

func (t *Thread) Capacity() int {
	return t.Header.Capacity
}

func (t *Thread) Update(ctx context.Context, p []byte) error {
	if t.Header.Capacity == 0 {
		return errors.New("no capacity")
	}

	if t.Header.Capacity < len(p) {
		return errors.New("insufficient capacity")
	}

	old := t.Header

	payload := string(p[:min(t.Header.Capacity, len(p))])
	var header *ThreadHeader
	err := retry(ctx, func() error {
		newHeader, err := AppendToReply(ctx, t.service, t.FileID, *t, payload)
		header = newHeader
		if err != nil {
			return err
		}
		return nil
	})
	if header != nil {
		t.Header = *header
		t.oldState = &old
	}
	return err
}

func (t *Thread) Put(ctx context.Context, p []byte) error {
	old := t.Header
	data := string(p[:min(EffectiveReplySize, len(p))])
	payload := padding + data + padding

	var header *ThreadHeader
	err := retry(ctx, func() error {
		newHeader, err := CreateReply(ctx, t.service, t.FileID, *t, &drive.Reply{Content: payload})
		header = newHeader
		return err
	})
	if header != nil {
		t.Header = *header
		t.oldState = &old
	}
	return err
}

// Rollback to the previous state. This is quite a desperate operation which may leave the file
// in an inconsistent state if API calls fail.
func (t *Thread) Rollback(ctx context.Context, service Service, fileID string) error {
	if t.oldState == nil {
		return ErrNoRollback
	}
	return RollbackCtx(ctx, service, fileID, t.CommentID, *t.oldState, t.Header)
}

// Read at most maxReplySize bytes into buffer p. Starts reading at the first reply of a thread, incrementing the reply every
// maxReplySize bytes. ReadCtx expects replies to be always filled. Replies are fetched in bulk and cached locally.
func (t *Thread) ReadCtx(ctx context.Context, p []byte) (int, error) {
	// Initial fetch
	if t.replies == nil {
		err := retry(ctx, func() error {
			client, err := t.service.Take(ctx, 1)
			if err != nil {
				return err
			}
			t.replies, err = client.RepliesService().
				List(t.FileID, t.CommentID).
				Context(ctx).
				Fields("*").
				Do()
			return err
		})
		if err != nil {
			return 0, err
		}
	}

	// fetch new bulk of replies
	if t.ri > len(t.replies.Replies) {
		err := retry(ctx, func() error {
			client, err := t.service.Take(ctx, 1)
			if err != nil {
				return err
			}
			replies, err := client.RepliesService().
				List(t.FileID, t.CommentID).
				PageToken(t.replies.NextPageToken).
				Fields("id", "content").
				Context(ctx).
				Do()

			if err != nil {
				return err
			}
			t.replies = replies
			t.ri = 0
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	if len(t.replies.Replies) == 0 || t.ri == len(t.replies.Replies) {
		return 0, io.EOF
	}

	reply := t.replies.Replies[t.ri]

	content := []byte(reply.Content)[1 : len(reply.Content)-1]
	copy(p, content[t.cursor:])
	read := min(len(content[t.cursor:]), len(p))

	t.cursor += read
	if t.cursor == EffectiveReplySize {
		t.cursor = 0
		t.ri++
	}
	return read, nil
}

func (t *Thread) Reply() drive.Reply {
	var i = t.ri
	if t.ri >= len(t.replies.Replies) {
		i -= 1
	}
	r := t.replies.Replies[i]
	return *r
}

// Create a new reply and update the ThreadHeader. The new ThreadHeader is returned.
//
// This function does not actually alter the reply or bucket, making it possible to retry this with exponential backoff.
func CreateReply(ctx context.Context, s Service, fileID string, bucket Thread, reply *drive.Reply) (*ThreadHeader, error) {
	service, err := s.Take(ctx, 2)
	if err != nil {
		return nil, err
	}

	r, err := service.RepliesService().
		Create(fileID, bucket.CommentID, reply).
		Context(ctx).
		Fields("*").
		Do()
	if err != nil {
		return nil, fmt.Errorf("create reply: %w", err)
	}

	bucket.Header.Capacity = MaxReplySize - len(reply.Content)
	bucket.Header.Tail = r.Id
	bucket.Header.Length++

	_, err = service.CommentsService().
		Update(fileID, bucket.CommentID, &drive.Comment{Content: string(bucket.Header.MustMarshall())}).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("update comment: %w", err)
	}
	return &bucket.Header, nil
}

// AppendToReply adds the content to the buckets tail. The caller should ensure that the content of the new
// reply does not exceed EffectiveReplySize.
func AppendToReply(ctx context.Context, s Service, fileID string, bucket Thread, content string) (*ThreadHeader, error) {
	service, err := s.Take(ctx, 3)
	if err != nil {
		return nil, err
	}

	reply, err := service.RepliesService().
		Get(fileID, bucket.CommentID, bucket.Header.Tail).
		Fields("*").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("get reply: %w", err)
	}

	reply.Content = reply.Content[:len(reply.Content)-1] + content + padding

	// TODO remove this once certain no 1 of errors are present
	if len(reply.Content) > MaxReplySize {
		return nil, fmt.Errorf("reply exceeded max size: %d", len(reply.Content))
	}

	bucket.Header.Capacity = MaxReplySize - len(reply.Content)

	_, err = service.RepliesService().
		Update(fileID, bucket.CommentID, bucket.Header.Tail, reply).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return nil, fmt.Errorf("update reply: %w", err)
	}

	_, err = service.CommentsService().
		Update(fileID, bucket.CommentID, &drive.Comment{Content: string(bucket.Header.MustMarshall())}).
		Fields("*").
		Context(ctx).
		Do()

	if err != nil {
		return &bucket.Header, fmt.Errorf("update comment: %w", err)
	}
	return &bucket.Header, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// byLengthThenNumber defines a sorting interface for first sorting on shortest length; then by ascending bucket number.
type byLengthThenNumber []*Thread

func (s byLengthThenNumber) Len() int {
	return len(s)
}
func (s byLengthThenNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLengthThenNumber) Less(i, j int) bool {
	if s[i].Header.Length == s[j].Header.Length {
		return s[i].Header.Number < s[j].Header.Number
	}
	return s[i].Header.Length < s[j].Header.Length
}
