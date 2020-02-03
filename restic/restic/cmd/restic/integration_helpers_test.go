package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/kaiserkarel/drfs/restic/restic/lib/options"
	"github.com/kaiserkarel/drfs/restic/restic/lib/repository"
	"github.com/kaiserkarel/drfs/restic/restic/lib/restic"
	rtest "github.com/kaiserkarel/drfs/restic/restic/lib/test"
)

type dirEntry struct {
	path string
	fi   os.FileInfo
	link uint64
}

func walkDir(dir string) <-chan *dirEntry {
	ch := make(chan *dirEntry, 100)

	go func() {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return nil
			}

			name, err := filepath.Rel(dir, path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
				return nil
			}

			ch <- &dirEntry{
				path: name,
				fi:   info,
				link: nlink(info),
			}

			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Walk() error: %v\n", err)
		}

		close(ch)
	}()

	// first element is root
	_ = <-ch

	return ch
}

func isSymlink(fi os.FileInfo) bool {
	mode := fi.Mode() & (os.ModeType | os.ModeCharDevice)
	return mode == os.ModeSymlink
}

func sameModTime(fi1, fi2 os.FileInfo) bool {
	switch runtime.GOOS {
	case "darwin", "freebsd", "openbsd", "netbsd":
		if isSymlink(fi1) && isSymlink(fi2) {
			return true
		}
	}

	same := fi1.ModTime().Equal(fi2.ModTime())
	if !same && (runtime.GOOS == "darwin" || runtime.GOOS == "openbsd") {
		// Allow up to 1μs difference, because macOS <10.13 cannot restore
		// with nanosecond precision and the current version of Go (1.9.2)
		// does not yet support the new syscall. (#1087)
		mt1 := fi1.ModTime()
		mt2 := fi2.ModTime()
		usecDiff := (mt1.Nanosecond()-mt2.Nanosecond())/1000 + (mt1.Second()-mt2.Second())*1000000
		same = usecDiff <= 1 && usecDiff >= -1
	}
	return same
}

// directoriesEqualContents checks if both directories contain exactly the same
// contents.
func directoriesEqualContents(dir1, dir2 string) bool {
	ch1 := walkDir(dir1)
	ch2 := walkDir(dir2)

	changes := false

	var a, b *dirEntry
	for {
		var ok bool

		if ch1 != nil && a == nil {
			a, ok = <-ch1
			if !ok {
				ch1 = nil
			}
		}

		if ch2 != nil && b == nil {
			b, ok = <-ch2
			if !ok {
				ch2 = nil
			}
		}

		if ch1 == nil && ch2 == nil {
			break
		}

		if ch1 == nil {
			fmt.Printf("+%v\n", b.path)
			changes = true
		} else if ch2 == nil {
			fmt.Printf("-%v\n", a.path)
			changes = true
		} else if !a.equals(b) {
			if a.path < b.path {
				fmt.Printf("-%v\n", a.path)
				changes = true
				a = nil
				continue
			} else if a.path > b.path {
				fmt.Printf("+%v\n", b.path)
				changes = true
				b = nil
				continue
			} else {
				fmt.Printf("%%%v\n", a.path)
				changes = true
			}
		}

		a, b = nil, nil
	}

	if changes {
		return false
	}

	return true
}

type dirStat struct {
	files, dirs, other uint
	size               uint64
}

func isFile(fi os.FileInfo) bool {
	return fi.Mode()&(os.ModeType|os.ModeCharDevice) == 0
}

// dirStats walks dir and collects stats.
func dirStats(dir string) (stat dirStat) {
	for entry := range walkDir(dir) {
		if isFile(entry.fi) {
			stat.files++
			stat.size += uint64(entry.fi.Size())
			continue
		}

		if entry.fi.IsDir() {
			stat.dirs++
			continue
		}

		stat.other++
	}

	return stat
}

type testEnvironment struct {
	base, cache, repo, mountpoint, testdata string
	gopts                                   GlobalOptions
}

// withTestEnvironment creates a test environment and returns a cleanup
// function which removes it.
func withTestEnvironment(t testing.TB) (env *testEnvironment, cleanup func()) {
	if !rtest.RunIntegrationTest {
		t.Skip("integration tests disabled")
	}

	repository.TestUseLowSecurityKDFParameters(t)
	restic.TestDisableCheckPolynomial(t)

	tempdir, err := ioutil.TempDir(rtest.TestTempDir, "restic-test-")
	rtest.OK(t, err)

	env = &testEnvironment{
		base:       tempdir,
		cache:      filepath.Join(tempdir, "cache"),
		repo:       filepath.Join(tempdir, "repo"),
		testdata:   filepath.Join(tempdir, "testdata"),
		mountpoint: filepath.Join(tempdir, "mount"),
	}

	rtest.OK(t, os.MkdirAll(env.mountpoint, 0700))
	rtest.OK(t, os.MkdirAll(env.testdata, 0700))
	rtest.OK(t, os.MkdirAll(env.cache, 0700))
	rtest.OK(t, os.MkdirAll(env.repo, 0700))

	env.gopts = GlobalOptions{
		Repo:     env.repo,
		Quiet:    true,
		CacheDir: env.cache,
		ctx:      context.Background(),
		password: rtest.TestPassword,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
		extended: make(options.Options),
	}

	// always overwrite global options
	globalOptions = env.gopts

	cleanup = func() {
		if !rtest.TestCleanupTempDirs {
			t.Logf("leaving temporary directory %v used for test", tempdir)
			return
		}
		rtest.RemoveAll(t, tempdir)
	}

	return env, cleanup
}