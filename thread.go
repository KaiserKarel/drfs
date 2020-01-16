package drfs

import (
	"encoding/json"
	"fmt"
	"io"

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
		panic(fmt.Errorf("marshalling theadheader failed: %w", err))
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
