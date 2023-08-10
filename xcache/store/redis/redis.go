// Package redis
/*
 * @Date: 2023-07-20 09:48:51
 * @LastEditTime: 2023-07-20 09:48:52
 * @Description:
 */
package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	store *redis.Client
}

// New
// @param conf
// @date 2022-07-02 08:12:14
func New(conf *redis.Options) *Store {
	cli := redis.NewClient(conf)

	return &Store{store: cli}
}

// NewWithDb
// @param tx
// @date 2022-07-02 08:12:12
func NewWithDb(tx *redis.Client) *Store {
	return &Store{store: tx}
}

// Set
// @param ctx
// @param key
// @param value
// @param ttl
// @date 2022-07-02 08:12:11
func (r *Store) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return r.store.Set(ctx, key, value, ttl).Err()
}

// Get
// @param ctx
// @param key
// @date 2022-07-02 08:12:09
func (r *Store) Get(ctx context.Context, key string) ([]byte, error) {
	return r.store.Get(ctx, key).Bytes()
}

// RemoveFromTag
// @param ctx
// @param tag
// @date 2022-07-02 08:12:08
func (r *Store) RemoveFromTag(ctx context.Context, tag string) error {
	keys, err := r.store.SMembers(ctx, tag).Result()
	if err != nil {
		return err
	}

	return r.store.Del(ctx, keys...).Err()
}

func (r *Store) RemoveFromKey(ctx context.Context, key string) error {
	return r.store.Del(ctx, key).Err()
}

// SaveTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) SaveTagKey(ctx context.Context, tag, key string) error {
	return r.store.SAdd(ctx, tag, key).Err()
}

// RemoveTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) RemoveTagKey(ctx context.Context, tag, key string) error {
	return r.store.SRem(ctx, tag, key).Err()
}

// MemberTagKey
// @param ctx
// @param tag
// @param key
// @date 2022-07-02 08:12:05
func (r *Store) MemberTagKey(ctx context.Context, tag, key string) (bool, error) {
	return r.store.SIsMember(ctx, tag, key).Result()
}
