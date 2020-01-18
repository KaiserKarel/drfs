package drfs

// size computes the amount of bytes using the index for computation.
func (f *File) size() int64 {
	var size int64
	for _, b := range f.index.Buckets {
		size += b.size()
	}
	return size
}

// size computes the amount of bytes in the bucket using the header information.
func (t *Thread) size() int64 {
	return t.Header.Length*EffectiveReplySize - int64(t.Header.Capacity)
}
