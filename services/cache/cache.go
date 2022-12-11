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

package cache

import (
	"github.com/lazybytez/jojo-discord-bot/services/cache/memory"
	"reflect"
	"time"
)

// Mode specifies what cache implementation should be used
type Mode string

// Available cache modes
const (
	ModeMemory Mode = "memory"
	ModeRedis  Mode = "redis"
)

// Dsn is the information necessary to connect
// to a remote server like redis.
// It is necessary for some cache implementations.
type Dsn string

// Provider specifies the interface that different cache implementations must provide.
type Provider interface {
	Get(key string, t interface{}) (interface{}, bool)
	Update(key string, t reflect.Type, value interface{})
	Invalidate(key string, t reflect.Type) bool
}

// cache is the currently active Provider that manages the underlying
// cache implementation.
var cache Provider

// Init initializes the caching system.
// The cache implementation used is chosen by the supplied Mode.
func Init(mode Mode, lifetime time.Duration, _ Dsn) {
	switch mode {
	case ModeMemory:
		inMemoryCache := memory.New(lifetime)
		inMemoryCache.UseGarbageCollector()
		cache = inMemoryCache
	case ModeRedis:
		panic("Cache mode redis is not implemented yet!")
	default:
		inMemoryCache := memory.New(lifetime)
		inMemoryCache.UseGarbageCollector()
		cache = inMemoryCache
	}
}

// Get returns a value from the cache.
// When there is no valid cache item available,
// the function will return the passed parameter t.
func Get[T any](key string, t T) (T, bool) {
	result, ok := cache.Get(key, t)

	switch v := result.(type) {
	case T:
		return v, ok
	default:
		return t, ok
	}
}

// Update adds or updates the given value to the cache with the given key.
func Update[T any](key string, value T) {
	cache.Update(key, reflect.TypeOf(value), value)
}

// Invalidate the cache item with the given key and type.
// The function returns whether an item has been invalidated or not.
func Invalidate[T any](key string, t T) bool {
	return cache.Invalidate(key, reflect.TypeOf(t))
}
