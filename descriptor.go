package drfs

import (
	"path"
)

type descriptor struct {
	fileName    string
	commentName string
}

func parse(address string) (descriptor, error) {
	filename, commentname := path.Split(address)
	return descriptor{
		fileName:    filename,
		commentName: commentname,
	}, nil
}
