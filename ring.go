package drfs

import (
	"container/ring"
)

func newThreadRing(bs []*Thread) *threadRing {
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
	return &threadRing{r}
}

type threadRing struct {
	*ring.Ring
}

func (b *threadRing) Get() *Thread {
	val := b.Peek()
	b.Next()
	return val
}

func (b *threadRing) Peek() *Thread {
	return b.Value.(*Thread)
}

func (b *threadRing) Next() {
	b.Ring = b.Ring.Next()
}
