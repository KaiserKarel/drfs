package restic

import (
	"context"
	"errors"
	"github.com/kaiserkarel/drfs/drive"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend"
	"github.com/kaiserkarel/drfs/restic/restic/lib/options"
	"strings"
)

const defaultLayout = "default"

func Open(cfg Config) (*Backend, error)  {
	service, err := drive.NewService(context.Background())
	if err != nil {
		return nil, err
	}

	l, err := backend.ParseLayout(
		&backend.LocalFilesystem{},
		cfg.Layout,
		defaultLayout,
		cfg.Path,
		)

	if err != nil {
		return nil, err
	}


	return &Backend{
		Config: cfg,
		Layout: l,
		service:service}, nil
}

func init() {
	options.Register("drfs", Config{})
}

type Config struct {
	Path string `option:"path" help:"Path to the repository."`
	Layout string `option:"layout" help:"use this backend directory layout (default: auto-detect)"`
}

// ParseConfig parses a local backend config.
func ParseConfig(cfg string) (interface{}, error) {
	if !strings.HasPrefix(cfg, "drfs:") {
		return nil, errors.New(`invalid format, prefix "local" not found`)
	}

	return Config{Path: cfg[5:]}, nil
}