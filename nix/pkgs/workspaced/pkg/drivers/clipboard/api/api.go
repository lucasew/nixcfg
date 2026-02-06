package api

import (
	"context"
	"io"
)

type Driver interface {
	WriteImageReader(ctx context.Context, r io.Reader) error
	WriteText(ctx context.Context, text string) error
}
