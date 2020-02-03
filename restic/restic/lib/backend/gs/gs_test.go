package gs_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/kaiserkarel/drfs/restic/restic/lib/backend"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/gs"
	"github.com/kaiserkarel/drfs/restic/restic/lib/backend/test"
	"github.com/kaiserkarel/drfs/restic/restic/lib/errors"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	rtest "github.com/kaiserkarel/drfs/restic/restic/lib/test"
)

func newGSTestSuite(t testing.TB) *test.Suite {
	tr, err := backend.Transport(backend.TransportOptions{})
	if err != nil {
		t.Fatalf("cannot create transport for tests: %v", err)
	}

	return &test.Suite{
		// do not use excessive data
		MinimalData: true,

		// NewConfig returns a config for a new temporary backend that will be used in tests.
		NewConfig: func() (interface{}, error) {
			gscfg, err := gs.ParseConfig(os.Getenv("RESTIC_TEST_GS_REPOSITORY"))
			if err != nil {
				return nil, err
			}

			cfg := gscfg.(gs.Config)
			cfg.ProjectID = os.Getenv("RESTIC_TEST_GS_PROJECT_ID")
			cfg.Prefix = fmt.Sprintf("test-%d", time.Now().UnixNano())
			return cfg, nil
		},

		// CreateFn is a function that creates a temporary repository for the tests.
		Create: func(config interface{}) (restic.Backend, error) {
			cfg := config.(gs.Config)

			be, err := gs.Create(cfg, tr)
			if err != nil {
				return nil, err
			}

			exists, err := be.Test(context.TODO(), restic.Handle{Type: restic.ConfigFile})
			if err != nil {
				return nil, err
			}

			if exists {
				return nil, errors.New("config already exists")
			}

			return be, nil
		},

		// OpenFn is a function that opens a previously created temporary repository.
		Open: func(config interface{}) (restic.Backend, error) {
			cfg := config.(gs.Config)
			return gs.Open(cfg, tr)
		},

		// CleanupFn removes data created during the tests.
		Cleanup: func(config interface{}) error {
			cfg := config.(gs.Config)

			be, err := gs.Open(cfg, tr)
			if err != nil {
				return err
			}

			return be.Delete(context.TODO())
		},
	}
}

func TestBackendGS(t *testing.T) {
	defer func() {
		if t.Skipped() {
			rtest.SkipDisallowed(t, "restic/backend/gs.TestBackendGS")
		}
	}()

	vars := []string{
		"GOOGLE_APPLICATION_CREDENTIALS",
		"RESTIC_TEST_GS_PROJECT_ID",
		"RESTIC_TEST_GS_REPOSITORY",
	}

	for _, v := range vars {
		if os.Getenv(v) == "" {
			t.Skipf("environment variable %v not set", v)
			return
		}
	}

	t.Logf("run tests")
	newGSTestSuite(t).RunTests(t)
}

func BenchmarkBackendGS(t *testing.B) {
	vars := []string{
		"GOOGLE_APPLICATION_CREDENTIALS",
		"RESTIC_TEST_GS_PROJECT_ID",
		"RESTIC_TEST_GS_REPOSITORY",
	}

	for _, v := range vars {
		if os.Getenv(v) == "" {
			t.Skipf("environment variable %v not set", v)
			return
		}
	}

	t.Logf("run tests")
	newGSTestSuite(t).RunBenchmarks(t)
}