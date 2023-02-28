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
	mu         sync.RWMutex
	cachePool  CachePool
	lifetime   time.Duration
	cleanUpJob *time.Ticker
}

// New creates a new cache with the specified lifetime (in seconds).
func New(lifetime time.Duration) *InMemoryCacheProvider {
	return &InMemoryCacheProvider{
		sync.RWMutex{},
		CachePool{},
		lifetime,
		nil,
	}
}

// cleanUpCache runs a single cleanup process
// over the cache.
func cleanUpCache(provider *InMemoryCacheProvider) {
	provider.mu.RLock()
	for _, cachePool := range provider.cachePool {
		// Asynchronously process pool clean up to prevent locking
		//the entire cache longer than necessary.
		cachePool := cachePool
		go func() {
			var invalidItemKeys []string

			// step 1: collect
			cachePool.mu.RLock()
			for key, value := range cachePool.entries {
				if value.invalid {
					invalidItemKeys = append(invalidItemKeys, key)

					continue
				}
			}
			cachePool.mu.RUnlock()

			// step 2: drop
			cachePool.mu.Lock()
			for _, invalidItem := range invalidItemKeys {
				delete(cachePool.entries, invalidItem)
			}
			cachePool.mu.Unlock()
		}()
	}
	provider.mu.RUnlock()
}

// UseGarbageCollector configures a periodic job
// that throws out items from the cache that are invalid,
// to prevent unnecessary memory usage.
//
// Note that currently there is no method to stop the garbage collector
// once started. However, this is fine for the applications needs.
func (provider *InMemoryCacheProvider) UseGarbageCollector() bool {
	if nil != provider.cleanUpJob {
		return false
	}

	ticker := time.NewTicker(provider.lifetime)
	provider.cleanUpJob = ticker

	go func() {
		for range ticker.C {
			go cleanUpCache(provider)
		}
	}()

	return true
}

// Get an Item from the cache, if there is a valid one.
// The function will return nil if there is no valid cache entry.
// A valid cache entry is present when:
//  1. for the given type and key an item can be found.
//  2. the found items lifetime is not exceeded
func (provider *InMemoryCacheProvider) Get(key string, t reflect.Type) (interface{}, bool) {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[t]
	lifetime := provider.lifetime
	provider.mu.RUnlock()

	if !ok || nil == typedCache {
		return nil, false
	}

	typedCache.mu.RLock()
	item, ok := typedCache.entries[key]
	typedCache.mu.RUnlock()

	if !ok || nil == item {
		return nil, false
	}

	item.mu.RLock()
	defer item.mu.RUnlock()

	if item.invalid {
		return nil, false
	}

	timeDiff := time.Since(item.since)
	if timeDiff >= lifetime {
		return nil, false
	}

	return item.value, true
}

// Update adds and item to the cache or updates it.
func (provider *InMemoryCacheProvider) Update(key string, t reflect.Type, value interface{}) error {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[t]
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

	return nil
}

// Invalidate manually invalidates the cache item behind
// the supplied key, if there is a cache item.
func (provider *InMemoryCacheProvider) Invalidate(key string, t reflect.Type) bool {
	provider.mu.RLock()
	typedCache, ok := provider.cachePool[t]
	provider.mu.RUnlock()

	if !ok || nil == typedCache {
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
