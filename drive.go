package drfs

import (
	"context"
	"errors"
	"os"

	"google.golang.org/api/googleapi"

	"golang.org/x/sync/errgroup"

	"google.golang.org/api/drive/v3"
)

var (
	ErrNotExist = os.ErrNotExist
	errFound    = errors.New("") // sentinel error to exist pagination.
)

const (
	// Maximum number of pages returned in a pagination request.
	MaxPages = 100
)

func findFile(ctx context.Context, s Service, fileName string, fields ...googleapi.Field) (*drive.File, error) {
	service, err := s.Take(ctx, 1)
	if err != nil {
		return nil, err
	}

	var result *drive.File
	err = service.FilesService().List().PageSize(MaxPages).Context(ctx).Fields(fields...).Pages(ctx, func(list *drive.FileList) error {
		for _, file := range list.Files {
			if file.Name == fileName {
				result = file
				return errFound
			}
		}
		_, err := s.Take(ctx, 1)
		if err != nil {
			return err
		}
		return nil
	})

	if errors.Is(err, errFound) {
		return result, nil
	}
	return nil, ErrNotExist
}

func ensurePermissionsSet(ctx context.Context, client Client, emails []string, fileID string) error {

	permissions := make(map[string]struct{})
	err := client.PermissionsService().List(fileID).Pages(ctx, func(list *drive.PermissionList) error {
		for _, perm := range list.Permissions {
			permissions[perm.EmailAddress] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return err
	}

	var permissionsToCreate []string

	for _, email := range emails {
		if _, ok := permissions[email]; !ok {
			permissionsToCreate = append(permissionsToCreate, email)
		}
	}

	grp, ctx := errgroup.WithContext(ctx)
	for _, email := range permissionsToCreate {
		email := email
		grp.Go(func() error {
			_, err := client.PermissionsService().
				Create(fileID, &drive.Permission{
					EmailAddress: email,
					Role:         "commenter",
					Type:         "user"}).
				SendNotificationEmail(false).
				Fields("").Context(ctx).Do()
			if err != nil {
				return err
			}
			return nil
		})
	}
	return grp.Wait()
}
