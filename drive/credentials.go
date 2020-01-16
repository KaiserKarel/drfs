package drive

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/oauth2/google"
)

// Credential combines the content of a JSON secret file and obtained credentials from the google API.
type Credential struct {
	Cred   *google.Credentials
	Secret Secret
}

// CredentialsFromDirectory walks the given directory, searching for files which have the suffix .json
// and unmarshal these into Credentials and Secrets.
func CredentialsFromDirectory(ctx context.Context, dir string) ([]Credential, error) {
	var creds []Credential
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(info.Name()) == ".json" {
			b, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			cred, err := google.CredentialsFromJSON(ctx, b, DefaultScopes...)
			if err != nil {
				return err
			}

			var secret Secret
			err = json.Unmarshal(b, &secret)
			if err != nil {
				return err
			}

			creds = append(creds, Credential{
				Cred:   cred,
				Secret: secret,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return creds, nil
}
