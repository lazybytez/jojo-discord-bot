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

package memory

import (
	"reflect"
	"sync"
	"time"
)

// Item is a single cache item that holds
// the value of the item and the time the value has been pulled.
type Item struct {
	mu      sync.RWMutex
	value   interface{}
	since   time.Time
	invalid bool
}

// CacheEntries is a map that holds cache items.
type CacheEntries map[string]*Item

// Cache holds cache entries and a lock to allow
// manage concurrent access.
type Cache struct {
	mu      sync.RWMutex
	entries CacheEntries
}

// CachePool holds type specific Cache instances.
// This structure allows to optimize locking of the cache.
// As every type has its own Cache instance, the entire
// CachePool does not need to be locked longer than necessary,
// as large part of data is in the dedicated Cache instances.
type CachePool map[reflect.Type]*Cache

// InMemoryCacheProvider provides the ability to cache data in the RAM.
type InMemoryCacheProvider struct {
	mu        sync.RWMutex
	cachePool CachePool
	lifetime  time.Duration
}

// New creates a new cache with the specified lifetime (in seconds).
func New(lifetime time.Duration) *InMemoryCacheProvider {
	return &InMemoryCacheProvider{
		sync.RWMutex{},
		CachePool{},
		lifetime,
	}
}

// Get an Item from the cache, if there is a valid one.
// The function will return nil if there is no valid cache entry.
// A valid cache entry is present when:
//  1. for the given type and key an item can be found.
//  2. the found items lifetime is not exceeded
func (provider *InMemoryCacheProvider) Get(key string, t interface{}) interface{} {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[reflect.TypeOf(t)]
	lifetime := provider.lifetime
	provider.mu.RUnlock()

	if !ok || nil == typedCache {
		return nil
	}

	typedCache.mu.RLock()
	item, ok := typedCache.entries[key]
	typedCache.mu.RUnlock()

	if !ok || nil == item {
		return nil
	}

	item.mu.RLock()
	defer item.mu.RUnlock()

	if item.invalid {
		return nil
	}

	timeDiff := time.Since(item.since)
	if timeDiff >= lifetime {
		return nil
	}

	return item.value
}

// Update adds and item to the cache or updates it.
func (provider *InMemoryCacheProvider) Update(key string, t reflect.Type, value interface{}) {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[reflect.TypeOf(t)]
	provider.mu.RUnlock()

	if !ok || nil == typedCache {
		typedCache = &Cache{
			sync.RWMutex{},
			CacheEntries{},
		}

		provider.mu.Lock()
		provider.cachePool[t] = typedCache
		provider.mu.Unlock()
	}

	typedCache.mu.RLock()
	item, ok := typedCache.entries[key]
	typedCache.mu.RUnlock()
	if !ok || nil == item {
		item = &Item{}

		typedCache.mu.Lock()
		typedCache.entries[key] = item
		typedCache.mu.Unlock()
	}

	item.mu.Lock()
	defer item.mu.Unlock()

	item.value = value
	item.since = time.Now()
}

// Invalidate manually invalidates the cache item behind
// the supplied key, if there is a cache item.
func (provider *InMemoryCacheProvider) Invalidate(key string, t reflect.Type) bool {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[t]
	provider.mu.RUnlock()

	if nil == typedCache {
		return false
	}

	typedCache.mu.RLock()
	item, ok := typedCache.entries[key]
	typedCache.mu.RUnlock()
	if !ok || nil == item {
		return false
	}

	item.mu.Lock()
	defer item.mu.Unlock()
	item.invalid = true

	return true
}
