package opener

import (
	"context"
)

type Driver interface {
	Open(ctx context.Context, target string) error
}
