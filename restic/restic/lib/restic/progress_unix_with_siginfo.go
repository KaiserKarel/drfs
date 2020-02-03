// +build darwin freebsd netbsd openbsd dragonfly

package restic

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiserkarel/drfs/restic/restic/lib/debug"
)

func init() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1)
	signal.Notify(c, syscall.SIGINFO)
	go func() {
		for s := range c {
			debug.Log("Signal received: %v\n", s)
			forceUpdateProgress <- true
		}
	}()
}
