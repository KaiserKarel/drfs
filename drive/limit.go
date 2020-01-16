package drive

import (
	"golang.org/x/time/rate"
)

// MaxUserLimit is the quota of queries per second (1000 calls per 100 seconds)
const MaxUserLimit rate.Limit = 1000 / 100

// TotalLimit is the quota of queries for the entire project. (no matter the number of users)
const TotalLimit rate.Limit = 10000 / 100

// DefaultUserBurst is the token burst allowed per client.
const DefaultUserBurst = 10

// DefaultTotalBurst is the maximum requested burst allowed.
const DefaultTotalBurst = DefaultUserBurst
