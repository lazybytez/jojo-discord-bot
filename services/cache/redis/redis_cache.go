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
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"net/url"
	"reflect"
	"time"
)

const (
	// localCacheCount is the max number of keys cached in the local short-lived cache
	localCacheCount = 10000
	// localCacheTtl is the defaultLifetime of a single cache key in the local cache
	localCacheTtl = 30 * time.Second
)

// GoRedisCacheProvider is a cache provider that allows to store
// cache entries in Redis. It uses a basic go-redis redis.Client to
// communicate with Redis.
type GoRedisCacheProvider struct {
	// defaultLifetime is the lifetime set by default for
	// cache entries stored in Redis.
	defaultLifetime time.Duration
	// cache is the cache.Cache instance used to store and retrieve values.
	cache *cache.Cache
	// client is the Redis client that is used.
	// The client used by cache and client should be always the same.
	client redis.UniversalClient
}

// New creates a new cache with the specified defaultLifetime (in seconds) and given redis DSN.
func New(cacheDsn string, lifetime time.Duration) (*GoRedisCacheProvider, error) {
	client, err := createClientFromDsn(cacheDsn)
	if err != nil {
		return nil, err
	}

	cacheProvider := &GoRedisCacheProvider{
		defaultLifetime: lifetime,
		cache: cache.New(&cache.Options{
			Redis:      client,
			LocalCache: cache.NewTinyLFU(localCacheCount, localCacheTtl),
			// Unfortunately MessagePack has a bug when decoding zero time values
			// nested in an interface.
			// Issue: https://github.com/vmihailenco/msgpack/issues/332
			// For now, we use JSON as a workaround.
			Marshal:   json.Marshal,
			Unmarshal: json.Unmarshal,
		}),
		client: client,
	}

	return cacheProvider, nil
}

// createClientFromDsn creates a new Redis client instance using the passed
// DSN. The function supports a DSN made out of a URL with a hostname or ip, a port
// and a username and password. If only one of username/password is provided, the value is used as
// the password.
func createClientFromDsn(cacheDsn string) (*redis.Client, error) {
	cacheUrl, err := url.Parse(cacheDsn)
	if nil != err {
		return nil, err
	}

	user, password := extractCredentialsFromDsnUrl(cacheUrl)

	client := redis.NewClient(&redis.Options{
		Username: user,
		Password: password,
		Addr:     cacheUrl.Host,
	})

	return client, nil
}

// extractCredentialsFromDsnUrl extracts the username (if set) and password from the passed
// URL instance.
// The function will use the username as password, when only one of username and password has
// been specified.
func extractCredentialsFromDsnUrl(cacheUrl *url.URL) (string, string) {
	user := cacheUrl.User.Username()
	password, hasPassword := cacheUrl.User.Password()
	if !hasPassword {
		user = ""
		password = cacheUrl.User.Username()
	}

	return user, password
}

// CheckRedisReachable checks if the configured Redis instance is reachable
// by sending a ping command.
func (grc *GoRedisCacheProvider) CheckRedisReachable() error {
	return grc.client.Ping(context.TODO()).Err()
}

// Get an Item from the cache, if there is a valid one.
// The function will return nil if there is no valid cache entry.
// A valid cache entry is present when:
//  1. for the given type and key an item can be found.
//  2. the found items defaultLifetime is not exceeded
func (grc *GoRedisCacheProvider) Get(key string, t reflect.Type) (interface{}, bool) {
	prototype := reflect.New(t).Interface()

	err := grc.cache.Get(context.TODO(), computeCacheKeyFromKeyAndType(key, t), prototype)

	if err != nil {
		return nil, false
	}

	prototype = reflect.Indirect(reflect.ValueOf(prototype).Elem()).Interface()

	return prototype, true
}

// Update adds and item to the cache or updates it.
func (grc *GoRedisCacheProvider) Update(key string, t reflect.Type, value interface{}) error {
	return grc.cache.Set(&cache.Item{
		Ctx:   context.TODO(),
		Key:   computeCacheKeyFromKeyAndType(key, t),
		Value: value,
		TTL:   grc.defaultLifetime,
	})
}

// Invalidate manually invalidates the cache item behind
// the supplied key, if there is a cache item.
func (grc *GoRedisCacheProvider) Invalidate(key string, t reflect.Type) bool {
	return grc.cache.Delete(context.TODO(), computeCacheKeyFromKeyAndType(key, t)) == nil
}

// Shutdown closes the Redis client attached to the cache instance.
func (grc *GoRedisCacheProvider) Shutdown() {
	_ = grc.client.Close()
}

// computeCacheKeyFromKeyAndType creates a new cache key from a key and a type.
// The format is "PackagePath_TypeName_Key".
func computeCacheKeyFromKeyAndType(key string, t reflect.Type) string {
	return fmt.Sprintf("%s_%s_%s", t.PkgPath(), t.Name(), key)
}
