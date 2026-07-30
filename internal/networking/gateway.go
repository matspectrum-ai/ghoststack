package networking

import "context"

type Gateway interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() string
}
