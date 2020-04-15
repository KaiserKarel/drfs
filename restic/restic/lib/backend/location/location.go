// Package location implements parsing the restic repository location from a string.
package location

import (
	drfs "github.com/kaiserkarel/drfs/restic"
	"strings"

	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/azure"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/b2"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/gs"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/local"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/rclone"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/rest"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/s3"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/sftp"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/swift"
	"github.com/kaiserkarel/drfs/restic/restic/lib/errors"
)

// Location specifies the location of a repository, including the method of
// access and (possibly) credentials needed for access.
type Location struct {
	Scheme string
	Config interface{}
}

type parser struct {
	scheme string
	parse  func(string) (interface{}, error)
}

// parsers is a list of valid config parsers for the backends. The first parser
// is the fallback and should always be set to the local backend.
var parsers = []parser{
	{"b2", b2.ParseConfig},
	{"local", local.ParseConfig},
	{"sftp", sftp.ParseConfig},
	{"s3", s3.ParseConfig},
	{"gs", gs.ParseConfig},
	{"azure", azure.ParseConfig},
	{"swift", swift.ParseConfig},
	{"rest", rest.ParseConfig},
	{"rclone", rclone.ParseConfig},
	{"drfs", drfs.ParseConfig},
}

func isPath(s string) bool {
	if strings.HasPrefix(s, "../") || strings.HasPrefix(s, `..\`) {
		return true
	}

	if strings.HasPrefix(s, "/") || strings.HasPrefix(s, `\`) {
		return true
	}

	if len(s) < 3 {
		return false
	}

	// check for drive paths
	drive := s[0]
	if !(drive >= 'a' && drive <= 'z') && !(drive >= 'A' && drive <= 'Z') {
		return false
	}

	if s[1] != ':' {
		return false
	}

	if s[2] != '\\' && s[2] != '/' {
		return false
	}

	return true
}

// Parse extracts repository location information from the string s. If s
// starts with a backend name followed by a colon, that backend's Parse()
// function is called. Otherwise, the local backend is used which interprets s
// as the name of a directory.
func Parse(s string) (u Location, err error) {
	scheme := extractScheme(s)
	u.Scheme = scheme

	for _, parser := range parsers {
		if parser.scheme != scheme {
			continue
		}

		u.Config, err = parser.parse(s)
		if err != nil {
			return Location{}, err
		}

		return u, nil
	}

	// if s is not a path or contains ":", it's ambiguous
	if !isPath(s) && strings.ContainsRune(s, ':') {
		return Location{}, errors.New("invalid backend\nIf the repo is in a local directory, you need to add a `local:` prefix")
	}

	u.Scheme = "local"
	u.Config, err = local.ParseConfig("local:" + s)
	if err != nil {
		return Location{}, err
	}

	return u, nil
}

func extractScheme(s string) string {
	data := strings.SplitN(s, ":", 2)
	return data[0]
}
