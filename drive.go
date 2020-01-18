package drfs

import (
	"context"

	"golang.org/x/sync/errgroup"

	"google.golang.org/api/drive/v3"
)

const (
	// Maximum number of pages returned in a pagination request.
	MaxPages = 100
)

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
