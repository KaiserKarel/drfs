package drfs

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/api/drive/v3"
)

const padding = "1"

type Bucket struct {
	CommentID string       `json:"c"`
	Header    ThreadHeader `json:"h"`

	cursor   int // current read position within reply
	ri       int // index of current reply being read.
	replies  *drive.ReplyList
	oldState *ThreadHeader
	modTime  time.Time
}

// Write buffer p to bucket, writing at most 4095 bytes. WriteCtx decides between appending to an existing reply or
// creating a new one. The ThreadHeader is updated as well; returning the old ThreadHeader.
func (b *Bucket) WriteCtx(ctx context.Context, s Service, fileID string, p []byte) (int, error) {
	defer func() {
		b.modTime = time.Now()
	}()
	old := b.Header
	var payload string
	var newHeader *ThreadHeader
	var err error
	var written int

	log.Printf("bucket %d has cap %d", b.Header.Number, b.Header.Capacity)
	if b.Header.Capacity == 0 {
		data := string(p[:min(EffectiveReplySize, len(p))])
		payload = padding + data
		if len(payload) == MaxReplySize-1 {
			payload += padding
		}
		err = retry(ctx, func() error {
			newHeader, err = CreateReply(ctx, s, fileID, *b, &drive.Reply{Content: payload})
			return err
		})
		written = len(data)
	} else {
		payload = string(p[:min(b.Header.Capacity, len(p))])
		err = retry(ctx, func() error {
			newHeader, err = AppendToReply(ctx, s, fileID, *b, payload)
			return err
		})
		written = len(payload)
	}

	if newHeader != nil {
		b.Header = *newHeader
		b.oldState = &old
	}

	log.Printf("bucket %d wrote %d cap %d", b.Header.Number, written, b.Header.Capacity)
	if err != nil {
		return 0, err
	}
	return written, nil
}

// Rollback to the previous state. This is quite a desperate operation which may leave the file
// in an inconsistent state if API calls fail.
func (b *Bucket) Rollback(ctx context.Context, service Service, fileID string) error {
	if b.oldState == nil {
		return ErrNoRollback
	}
	return RollbackCtx(ctx, service, fileID, b.CommentID, *b.oldState, b.Header)
}

// Read at most maxReplySize bytes into buffer p. Starts reading at the first reply of a thread, incrementing the reply every
// maxReplySize bytes. ReadCtx expects replies to be always filled. Replies are fetched in bulk and cached locally.
func (b *Bucket) ReadCtx(ctx context.Context, s Service, fileID string, p []byte) (int, error) {
	// Initial fetch
	if b.replies == nil {
		err := retry(ctx, func() error {
			service, err := s.Take(ctx, 1)
			if err != nil {
				return err
			}
			b.replies, err = service.RepliesService().
				List(fileID, b.CommentID).
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
	if b.ri > len(b.replies.Replies) {
		err := retry(ctx, func() error {
			service, err := s.Take(ctx, 1)
			if err != nil {
				return err
			}
			replies, err := service.RepliesService().
				List(fileID, b.CommentID).
				PageToken(b.replies.NextPageToken).
				Fields("id", "content").
				Context(ctx).
				Do()

			if err != nil {
				return err
			}
			b.replies = replies
			b.ri = 0
			return nil
		})
		if err != nil {
			return 0, err
		}
	}

	if len(b.replies.Replies) == 0 || b.ri == len(b.replies.Replies) {
		return 0, io.EOF
	}

	reply := b.replies.Replies[b.ri]

	var content []byte
	if len(reply.Content) == MaxReplySize {
		content = []byte(reply.Content)[1 : EffectiveReplySize+1]
	} else {
		content = []byte(reply.Content)[1:]
	}

	copy(p, content[b.cursor:])
	written := min(len(content[b.cursor:]), len(p))

	b.cursor += written
	if b.cursor == EffectiveReplySize {
		b.cursor = 0
		b.ri++
	}
	return written, nil
}

func (b *Bucket) Reply() drive.Reply {
	var i = b.ri
	if b.ri >= len(b.replies.Replies) {
		i -= 1
	}
	r := b.replies.Replies[i]
	return *r
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// byLengthThenNumber defines a sorting interface for first sorting on shortest length; then by ascending bucket number.
type byLengthThenNumber []*Bucket

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
