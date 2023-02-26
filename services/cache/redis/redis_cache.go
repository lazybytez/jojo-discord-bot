/*
 * JOJO Discord Bot - An advanced multi-purpose discord bot
 * Copyright (C) 2022 Lazy Bytez (Elias Knodel, Pascal Zarrad)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/url"
	"reflect"
	"sync"
	"time"
)

// GoRedisCacheProvider is a cache provider that allows to store
// cache entries in Redis. It uses a basic go-redis redis.Client to
// communicate with Redis.
type GoRedisCacheProvider struct {
	mu       sync.RWMutex
	cacheDsn string
	lifetime time.Duration
	context  context.Context
	client   *redis.Client
}

// New creates a new cache with the specified lifetime (in seconds) and given redis DSN.
func New(cacheDsn string, lifetime time.Duration) (*GoRedisCacheProvider, error) {
	cacheUrl, err := url.Parse(cacheDsn)

	if nil != err {
		return nil, err
	}

	user := cacheUrl.User.Username()
	password, hasPassword := cacheUrl.User.Password()
	if !hasPassword {
		user = ""
		password = cacheUrl.User.Username()
	}

	cacheProvider := &GoRedisCacheProvider{
		mu:       sync.RWMutex{},
		cacheDsn: cacheDsn,
		lifetime: lifetime,
		client: redis.NewClient(&redis.Options{
			Username: user,
			Password: password,
			Addr:     cacheUrl.Host,
		}),
		context: context.Background(),
	}

	err = cacheProvider.getClient().Ping(cacheProvider.context).Err()
	if nil != err {
		return nil, err
	}

	return cacheProvider, nil
}

// Get an Item from the cache, if there is a valid one.
// The function will return nil if there is no valid cache entry.
// A valid cache entry is present when:
//  1. for the given type and key an item can be found.
//  2. the found items lifetime is not exceeded
func (cache *GoRedisCacheProvider) Get(key string, t interface{}) (interface{}, bool) {
	result, err := cache.getClient().Get(cache.getContext(), computeCacheKey(key, t)).Result()

	if err != nil {
		return nil, false
	}

	var prototype interface{}
	switch typed := t.(type) {
	case reflect.Type:
		prototype = reflect.New(typed).Interface()
	default:
		prototype = reflect.New(reflect.TypeOf(t)).Interface()
	}

	err = json.Unmarshal([]byte(result), &prototype)

	if err != nil {
		return nil, false
	}

	return prototype, true
}

// Update adds and item to the cache or updates it.
func (cache *GoRedisCacheProvider) Update(key string, t reflect.Type, value interface{}) error {
	object, err := json.Marshal(value)
	if nil != err {
		return err
	}

	result := cache.getClient().Set(cache.getContext(), computeCacheKey(key, t), string(object), cache.lifetime)

	return result.Err()
}

// Invalidate manually invalidates the cache item behind
// the supplied key, if there is a cache item.
func (cache *GoRedisCacheProvider) Invalidate(key string, t reflect.Type) bool {
	result := cache.getClient().Del(cache.getContext(), computeCacheKey(key, t))

	return result.Val() >= 1
}

// computeCacheKey creates a cache key from the passed key and value (t) that has been passed.
// If the passed value is not of type reflect.Type, it will be converted to it.
func computeCacheKey(key string, t interface{}) string {
	switch tType := t.(type) {
	case reflect.Type:
		return computeCacheKeyFromKeyAndType(key, tType)
	default:
		return computeCacheKeyFromKeyAndType(key, reflect.TypeOf(t))
	}
}

// computeCacheKeyFromKeyAndType creates a new cache key from a key and a type.
// The format is "PackagePath_TypeName_Key".
func computeCacheKeyFromKeyAndType(key string, t reflect.Type) string {
	return fmt.Sprintf("%s_%s_%s", t.PkgPath(), t.Name(), key)
}

// getClient provides synchronized access to the client of a GoRedisCacheProvider.
func (cache *GoRedisCacheProvider) getClient() *redis.Client {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.client
}

// getContext provides synchronized access to the context of a GoRedisCacheProvider.
func (cache *GoRedisCacheProvider) getContext() context.Context {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	return cache.context
}
