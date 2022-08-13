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
	"sync"
	"time"
)

// Cache is a simple struct that
// is used to create cache instances that allow to cache data.
//
// The lifetime specifies how often items are refreshed.
// Note that the cache considers the time an item has been pulled
// and not the last access!
type Cache[I comparable, C any] struct {
	// registeredComponentCache is a map that holds component
	// codes and their reference.
	//
	// It acts as a cache to prevent exhaustive database calls
	// for something that is required during the entire application
	// lifetime and does not change.
	cache    map[I]*Item[*C]
	lock     sync.RWMutex
	lifetime int64
}

// Item is a single cache item that holds
// the value of the item and the time the value has been pulled.
type Item[C any] struct {
	value C
	since time.Time
}

// New creates a new cache with the specified lifetime (in milliseconds).
func New[I comparable, C any](lifetime time.Duration) *Cache[I, C] {
	return &Cache[I, C]{
		map[I]*Item[*C]{},
		sync.RWMutex{},
		lifetime.Milliseconds(),
	}
}

// Get an Item from the cache, if there is one with a valid lifetime
func Get[I comparable, C any](cache *Cache[I, C], key I) (*C, bool) {
	cache.lock.RLock()
	defer cache.lock.RUnlock()
	item, ok := cache.cache[key]

	if !ok {
		return nil, false
	}

	timeDiff := time.Since(item.since).Milliseconds()
	if cache.lifetime != 0 && timeDiff > cache.lifetime {
		return nil, false
	}

	return item.value, true
}

// Update adds or updates an item in the cache.
func Update[I comparable, C any](cache *Cache[I, C], key I, value *C) {
	item, ok := cache.cache[key]

	if !ok {
		item = &Item[*C]{}

		cache.lock.Lock()
		cache.cache[key] = item
		cache.lock.Unlock()
	}

	item.value = value
	item.since = time.Now()
}
