package power

import (
	"context"
)

type Driver interface {
	Lock(ctx context.Context) error
	Logout(ctx context.Context) error
	Suspend(ctx context.Context) error
	Hibernate(ctx context.Context) error
	Reboot(ctx context.Context) error
	Shutdown(ctx context.Context) error
}
