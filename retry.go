package drfs

import (
	"context"
	"errors"
	"time"

	"google.golang.org/api/googleapi"

	"github.com/cenkalti/backoff/v4"
)

// Retry an operation using exponential backoff. If the operation returns a googleapi.Error with code 500 or 404, the
// operation is retried as well.
func retry(ctx context.Context, operation backoff.Operation) error {
	return tryUntil(operation, backoff.WithContext(backoff.NewExponentialBackOff(), ctx), checkErr)
}

// Check if an error is googleapi.Error.Code 500 or 404
func checkErr(err error) bool {
	var apiError googleapi.Error
	if errors.Is(err, &apiError) {
		return apiError.Code != 500 && apiError.Code != 404 // 404 might occur because of eventual consistency
	}
	return false
}

// Adapted from github.com/cenkalti/backoff/v4 to allow control over error checking.
func tryUntil(operation backoff.Operation, b backoff.BackOffContext, f func(error) bool) error {
	var err error
	var next time.Duration
	t := &defaultTimer{}
	ctx := b.Context()
	b.Reset()

	for {
		if err = operation(); err == nil {
			return nil
		}

		if f(err) {
			return err
		}

		if next = b.NextBackOff(); next == backoff.Stop {
			return err
		}

		t.Start(next)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C():
		}
	}
}

// defaultTimer implements backoff.Timer interface using time.Timer
type defaultTimer struct {
	timer *time.Timer
}

// C returns the timers channel which receives the current time when the timer fires.
func (t *defaultTimer) C() <-chan time.Time {
	return t.timer.C
}

// Start starts the timer to fire after the given duration
func (t *defaultTimer) Start(duration time.Duration) {
	if t.timer == nil {
		t.timer = time.NewTimer(duration)
	} else {
		t.timer.Reset(duration)
	}
}

// Stop is called when the timer is not used anymore and resources may be freed.
func (t *defaultTimer) Stop() {
	if t.timer != nil {
		t.timer.Stop()
	}
}
