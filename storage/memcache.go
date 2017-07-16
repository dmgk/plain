// Memcache backing store.
package storage

import (
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/appengine/memcache"
)

type memcacheStore struct{}

func NewMemcache() Store {
	return &memcacheStore{}
}

func (s *memcacheStore) Add(ctx context.Context, content string) (string, error) {
	key, err := s.nextKey(ctx)
	if err != nil {
		return "", err
	}
	item := &memcache.Item{
		Key:        key,
		Value:      []byte(content),
		Expiration: DefaultTTL,
	}
	return key, s.translateError(memcache.Add(ctx, item))
}

func (s *memcacheStore) Get(ctx context.Context, key string) (string, error) {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return "", s.translateError(err)
	}
	return string(item.Value), nil
}

func (s *memcacheStore) Delete(ctx context.Context, key string) error {
	return s.translateError(memcache.Delete(ctx, key))
}

func (s *memcacheStore) Expire(ctx context.Context) error {
	return nil // nop
}

const (
	keyKey = "_paste_id"
)

func (s *memcacheStore) nextKey(ctx context.Context) (string, error) {
	nextVal, err := memcache.IncrementExisting(ctx, keyKey, 1)

	if err != nil {
		if err != memcache.ErrCacheMiss {
			return "", s.translateError(err)
		}

		// no key was stored yet, add a new key with initial value
		nextVal = 0
		item := &memcache.Item{
			Key:        keyKey,
			Value:      []byte(strconv.FormatUint(nextVal, 10)),
			Expiration: 0, // no expiration
		}

		if err := memcache.Add(ctx, item); err != nil {
			if err != memcache.ErrNotStored {
				return "", s.translateError(err)
			}

			// other goroutine beat us to it, retry getting the next key value
			return s.nextKey(ctx)
		}
	}

	return strconv.FormatUint(nextVal, 36), nil
}

func (s *memcacheStore) translateError(err error) error {
	if err == memcache.ErrCacheMiss {
		return ErrNotFound
	}
	return err
}
