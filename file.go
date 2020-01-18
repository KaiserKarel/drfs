package drfs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"github.com/google/uuid"

	"golang.org/x/sync/errgroup"
	"google.golang.org/api/drive/v3"
)

// maxReplySize defines the maximum number of bytes Drive allows in the content of
// single reply, encoded in UTF-8.
const MaxReplySize = 4096

// EffectiveReplySize is the size used per reply, as leading spaces are removed by Drive.
const EffectiveReplySize = MaxReplySize - 2

type FileHeader struct {
	FileOptions `json:"o"`
}

func (f FileHeader) MustMarshall() []byte {
	p, err := json.Marshal(f)
	if err != nil {
		panic(fmt.Errorf("marshaling fileheader failed: %w", err))
	}
	return p
}

func FileHeaderFromJSON(p io.Reader) (*FileHeader, error) {
	dec := json.NewDecoder(p)
	dec.DisallowUnknownFields()

	header := &FileHeader{}
	err := dec.Decode(header)
	return header, err
}

type FileOptions struct {
	NumThreads int
}

func (f *FileOptions) setDefaults() {
	if f.NumThreads == 0 {
		f.NumThreads = 512
	}
}

func OpenCtx(ctx context.Context, file *drive.File, service Service) (*File, error) {
	index, err := IndexFromFile(context.TODO(), service, file)
	if err != nil {
		return nil, fmt.Errorf("unable to index file: %w", err)
	}

	var writerlist = make([]*Thread, len(index.Buckets))
	copy(writerlist, index.Buckets)
	sort.Sort(byLengthThenNumber(writerlist))

	return &File{
		file:    file,
		index:   *index,
		writers: newThreadRing(writerlist),
		readers: newThreadRing(index.Buckets),
		service: service,
	}, nil
}

func CreateFileCtx(ctx context.Context, service Service, fileName string, options FileOptions) (*File, error) {
	options.setDefaults()

	var fileheader = FileHeader{FileOptions: options}
	var buckets = make([]*Thread, options.NumThreads)

	client, err := service.Take(context.TODO(), 2)
	if err != nil {
		return nil, err
	}

	file, err := client.FilesService().
		Create(&drive.File{Name: fileName}).
		Fields("id").
		Context(context.TODO()).
		Do()
	if err != nil {
		return nil, err
	}

	emails := service.Emails()
	if len(emails) > 1 {
		err = ensurePermissionsSet(ctx, client, emails, file.Id)
		if err != nil {
			return nil, fmt.Errorf("unable to set permissions: %w", err)
		}
	}

	grp, ctx := errgroup.WithContext(context.TODO())
	comments := make([]*drive.Comment, options.NumThreads)

	// create the file header itself.
	grp.Go(func() error {
		return retry(ctx, func() error {
			client, err := service.Take(context.TODO(), 1)
			if err != nil {
				return err
			}

			header := FileHeader{options}
			_, err = client.CommentsService().
				Create(file.Id, &drive.Comment{Content: string(header.MustMarshall())}).
				Context(context.TODO()).Fields("id").
				Do()
			return err
		})
	})

	// create individual threads.
	for i := 0; i < options.NumThreads; i++ {

		i := i
		grp.Go(func() error {
			return retry(ctx, func() error {
				client, err := service.Take(context.TODO(), 1)
				if err != nil {
					return err
				}

				header := &ThreadHeader{
					Number: i,
					UUID:   uuid.New(),
				}
				comment, err := client.CommentsService().
					Create(file.Id, &drive.Comment{Content: string(header.MustMarshall())}).
					Context(context.TODO()).Fields("id").
					Do()
				if err != nil {
					return err
				}

				comments[i] = comment
				buckets[i] = &Thread{
					FileID:    file.Id,
					CommentID: comment.Id,
					Header:    *header,
					service:   service,
					cursor:    0,
					ri:        0,
					replies:   nil,
					oldState:  nil,
				}
				return nil
			})
		})
	}

	err = grp.Wait()
	if err != nil {
		client, limitErr := service.Take(context.TODO(), 1)
		if limitErr != nil {
			return nil, fmt.Errorf("%w (%s)", err, limitErr)
		}

		deleteErr := client.FilesService().
			Delete(file.Id).
			Fields("id").
			Context(context.TODO()).
			Do()
		if deleteErr != nil {
			return nil, fmt.Errorf("%w (%s)", err, deleteErr)
		}
		return nil, err
	}

	return &File{
		file: file,
		index: Index{
			Header:  fileheader,
			Buckets: buckets,
		},
		writers: newThreadRing(buckets),
		readers: newThreadRing(buckets),
		service: service,
	}, nil
}

type File struct {
	file    *drive.File
	index   Index
	writers *threadRing
	readers *threadRing
	service Service
}

func (f *File) Service() Service {
	return f.service
}

func (f *File) Index() Index {
	return f.index
}

type bounds struct {
	lower int
	upper int
}
