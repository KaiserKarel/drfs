package drfs

import (
	"net"
)

// IPGenerator defines the interface for generating IP addresses, which are used to increase the rate limit in
// read operations by using the mock IP address. IPGenerator should not just randomly generate addresses, but
// use consistent zones and preferably rotate a set of ~1000 IP addresses to avoid detection.
type IPGenerator interface {
	Generate() *net.IPAddr
}
