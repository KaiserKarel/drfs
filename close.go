package drfs

// Close is currently a noop, but might later be used to flush buffers or
// commit the index.
func (f *File) Close() error {
	return nil
}
