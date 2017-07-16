// Google Datastore backing store.
package storage

import (
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/memcache"
)

type gdsStore struct{}

func NewGDS() Store {
	return &gdsStore{}
}

const (
	pasteKind       = "Paste"
	counterKind     = "Counter"
	counterShards   = 10
	counterCacheKey = "Counter"
)

type paste struct {
	Content   string `datastore:",noindex"`
	CreatedAt time.Time
}

type counter struct {
	Count int `datastore:",noindex"`
}

func (s *gdsStore) Add(ctx context.Context, content string) (string, error) {
	key, err := s.nextKey(ctx)
	if err != nil {
		return "", s.translateError(err)
	}
	p := paste{
		Content:   content,
		CreatedAt: time.Now(),
	}
	if _, err := datastore.Put(ctx, key, &p); err != nil {
		return "", s.translateError(err)
	}
	return key.StringID(), nil
}

func (s *gdsStore) Get(ctx context.Context, key string) (string, error) {
	var p paste
	err := datastore.Get(ctx, datastore.NewKey(ctx, pasteKind, key, 0, nil), &p)
	if err != nil {
		return "", s.translateError(err)
	}
	return p.Content, nil
}

func (s *gdsStore) Delete(ctx context.Context, key string) error {
	err := datastore.Delete(ctx, datastore.NewKey(ctx, pasteKind, key, 0, nil))
	return s.translateError(err)
}

func (s *gdsStore) Expire(ctx context.Context) error {
	q := datastore.NewQuery(pasteKind).
		Filter("CreatedAt <", time.Now().Add(-DefaultTTL)).
		KeysOnly()
	keys, err := q.GetAll(ctx, nil)
	if err != nil {
		return err
	}
	return datastore.DeleteMulti(ctx, keys)
}

// Strongly consistent sharded counter implementation.
// Based on https://cloud.google.com/appengine/articles/sharding_counters
// but with lookup by keys for strong consistency.

func (s *gdsStore) nextKey(ctx context.Context) (*datastore.Key, error) {
	// increment counter in a randomly selected shard
	err := datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		var c counter
		shard := rand.Intn(counterShards)
		key := datastore.NewKey(ctx, counterKind, strconv.FormatInt(int64(shard), 10), 0, nil)
		err := datastore.Get(ctx, key, &c)
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		c.Count++
		_, err = datastore.Put(ctx, key, &c)
		return err
	}, nil)

	if err != nil {
		return nil, err
	}

	// bump cached total value
	total, err := memcache.IncrementExisting(ctx, counterCacheKey, 1)

	if err != nil {
		if err != memcache.ErrCacheMiss {
			return nil, err
		}

		// cache miss, get and cache total count
		keys := make([]*datastore.Key, counterShards)
		for i := 0; i < counterShards; i++ {
			keys[i] = datastore.NewKey(ctx, counterKind, strconv.FormatInt(int64(i), 10), 0, nil)
		}
		cc := make([]*counter, counterShards)
		err := datastore.GetMulti(ctx, keys, cc)
		if err != nil {
			if err, ok := err.(appengine.MultiError); ok {
				for _, e := range err {
					if e != nil && e != datastore.ErrNoSuchEntity {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}
		for _, c := range cc {
			if c != nil {
				total += uint64(c.Count)
			}
		}

		item := memcache.Item{
			Key:        counterCacheKey,
			Value:      []byte(strconv.FormatUint(total, 10)),
			Expiration: 5 * time.Minute,
		}
		if err := memcache.Add(ctx, &item); err != nil {
			if err != memcache.ErrNotStored {
				return nil, err
			}

			// other goroutine beat us to it, retry getting the next key value
			return s.nextKey(ctx)
		}
	}

	return datastore.NewKey(ctx, pasteKind, strconv.FormatUint(total, 36), 0, nil), nil
}

func (s *gdsStore) translateError(err error) error {
	if err == datastore.ErrNoSuchEntity {
		return ErrNotFound
	}
	return err
}
