package storage

import (
	"errors"
	"time"

	"golang.org/x/net/context"
)

const DefaultTTL = 24 * time.Hour

var ErrNotFound = errors.New("store: item not found")

type Store interface {
	Add(ctx context.Context, content string) (string, error)
	Get(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	Expire(ctx context.Context) error
}
