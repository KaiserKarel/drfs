package drfs

import (
	"bufio"
)

// NewBufferedWriter returns a buffered writer capable of at least filling every single thread each buffer; meaning
// that the maximum upload limit becomes dependent on the rate limiter.
func NewBufferedWriter(f *File) *bufio.Writer {
	return bufio.NewWriterSize(f, EffectiveReplySize*len(f.Index().Buckets))
}
