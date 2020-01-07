package drfs

import (
	"context"
	"errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"os"
)

var (
	ErrNotExist = os.ErrNotExist
	errFound    = errors.New("")
)

func ensureFileExists(ctx context.Context, service *drive.Service, fileName string, fields ...googleapi.Field) (*drive.File, error) {
	fileService := drive.NewFilesService(service)

	file, err := findFile(ctx, service, fileName, fields...)
	if err != nil {
		file, err = fileService.Create(&drive.File{Name: fileName}).Context(ctx).Fields(fields...).Do()
		if err != nil {
			return nil, err
		}
	}
	return file, nil
}

func ensureCommentExists(ctx context.Context, service *drive.Service, fileID, commentName string, fields ...googleapi.Field) (*drive.Comment, error) {
	commentsService := drive.NewCommentsService(service)

	comment, err := findComment(ctx, service, fileID, commentName, fields...)
	if err != nil {
		comment, err = commentsService.Create(fileID, &drive.Comment{Content: commentName, ForceSendFields: []string{"Content"}}).Context(ctx).Fields(fields...).Do()
		if err != nil {
			return nil, err
		}
	}
	return comment, nil
}

func findFile(ctx context.Context, service *drive.Service, fileName string, fields ...googleapi.Field) (*drive.File, error) {
	if len(fields) == 0 {
		fields = []googleapi.Field{"*"}
	}

	fileService := drive.NewFilesService(service)

	var result *drive.File
	err := fileService.List().Context(ctx).Fields(fields...).Pages(ctx, func(list *drive.FileList) error {
		for _, file := range list.Files {
			if file.Name == fileName {
				result = file
				return errFound
			}
		}
		return ErrNotExist
	})

	if errors.Is(err, errFound) {
		return result, nil
	}
	return nil, err
}

func findComment(ctx context.Context, service *drive.Service, fileID, commentName string, fields ...googleapi.Field) (*drive.Comment, error) {
	if len(fields) == 0 {
		fields = []googleapi.Field{"*"}
	}

	commentService := drive.NewCommentsService(service)
	var result *drive.Comment
	err := commentService.List(fileID).Context(ctx).Fields(fields...).Pages(ctx, func(list *drive.CommentList) error {
		for _, comment := range list.Comments {
			if comment.Content == commentName {
				result = comment
				return errFound
			}
		}
		return ErrNotExist
	})

	if errors.Is(err, errFound) {
		return result, nil
	}
	return nil, err
}
