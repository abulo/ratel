package basep2c

import "google.golang.org/grpc/balancer"

// P2c ...
type P2c interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a item.
	Add(interface{})
}
