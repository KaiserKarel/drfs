package drive

import (
	"container/ring"
	"context"
	"sync"

	"github.com/kaiserkarel/drfs"
	"golang.org/x/oauth2/google"

	"golang.org/x/sync/errgroup"
	"google.golang.org/api/option"

	"golang.org/x/time/rate"
	"google.golang.org/api/drive/v3"
)

//nolint
var DefaultScopes = []string{
	"https://www.googleapis.com/auth/drive",
}

// Service implements drfs.Service using a ring of clients to alternate the source of API calls,
// allowing for a larger effective rate limit.
type Service struct {
	Limit   *rate.Limiter
	mu      *sync.Mutex
	ring    *clientRing
	clients []*Client
}

// NewService constructs a a Service consisting of len(credentials) clients || 1 client. If
// no clients are provided; google.FindDefaultCredentials is used to obtain the credentials.
// No secret will be obtained, and thus the call to Emails() will return nil, although the
// service contains a single valid client.
func NewService(ctx context.Context, credentials ...Credential) (*Service, error) {
	if len(credentials) == 0 {
		cred, err := google.FindDefaultCredentials(ctx, DefaultScopes...)
		if err != nil {
			return nil, err
		}
		credentials = []Credential{{
			Cred:   cred,
			Secret: Secret{},
		}}
	}

	var clients = make([]*Client, len(credentials))

	grp, ctx := errgroup.WithContext(ctx)

	for i := 0; i < len(credentials); i++ {
		credential := credentials[i]
		i := i
		grp.Go(func() error {
			service, err := drive.NewService(ctx, option.WithCredentials(credential.Cred))
			if err != nil {
				return err
			}
			clients[i] = &Client{
				Limiter: rate.NewLimiter(MaxUserLimit, DefaultUserBurst),
				service: service,
				Secret:  credential.Secret,
				i:       i,
			}
			return nil
		})
	}

	err := grp.Wait()
	if err != nil {
		return nil, err
	}

	r := ring.New(len(clients))
	for _, client := range clients {
		r.Value = client
		r = r.Next()
	}
	return &Service{
		Limit:   rate.NewLimiter(TotalLimit, DefaultTotalBurst),
		mu:      &sync.Mutex{},
		ring:    &clientRing{r},
		clients: clients,
	}, nil
}

type clientRing struct {
	ring *ring.Ring
}

// Returns the emails used by the service accounts. If no credentials were provided, No emails are
// returned.
func (s *Service) Emails() []string {
	var emails []string
	for _, client := range s.clients {
		email := client.Secret.ClientEmail
		if email == "" {
			continue
		}
		emails = append(emails, email)
	}
	return emails
}

// Requests N tokens from the global rate limiter, then obtains the next client and requests N tokens from
// that client too. If the context is cancelled an error is returned.
func (s *Service) Take(ctx context.Context, n int) (drfs.Client, error) {
	err := s.Limit.WaitN(ctx, n)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	client := s.ring.next()
	s.mu.Unlock()
	err = client.Limiter.WaitN(ctx, n)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *clientRing) next() *Client {
	cl := c.ring.Value.(*Client)
	c.ring = c.ring.Next()
	return cl
}
