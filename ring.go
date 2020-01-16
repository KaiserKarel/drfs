package drfs

import (
	"container/ring"
)

func newBucketRing(bs []*Bucket) *bucketRing {
	if len(bs) < 1 {
		panic("0 length rings are invalid")
	}

	r := ring.New(len(bs))
	for _, b := range bs {
		if b == nil {
			panic("bucket cannot be nil")
		}
		r.Value = b
		r = r.Next()
	}
	return &bucketRing{r}
}

type bucketRing struct {
	*ring.Ring
}

func (b *bucketRing) Get() *Bucket {
	val := b.Value.(*Bucket)
	b.Ring = b.Next()
	return val
}
