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
	"github.com/lazybytez/jojo-discord-bot/api/log"
	"sync"
	"time"
)

// LogPrefix is the prefix for log messages by cache functions and routines
const LogPrefix = "Cache"

// Cache is a simple struct that
// is used to create cache instances that allow to cache data.
//
// The lifetime specifies how often items are refreshed.
// Note that the cache considers the time an item has been pulled
// and not the last access!
//
// Note that automated cleanup of no longer used values can be
// enabled using EnableAutoCleanup. It is important that
// DisableAutoCleanup is called when the cleanup should stop,
// to enable the cache to be collected by the GC.
// Without disabling the automated cleanup, a reference to the cache
// will be hold infinite.
type Cache[I comparable, C any] struct {
	// registeredComponentCache is a map that holds component
	// codes and their reference.
	//
	// It acts as a cache to prevent exhaustive database calls
	// for something that is required during the entire application
	// lifetime and does not change.
	cache       map[I]*Item[*C]
	lock        sync.RWMutex
	lifetime    int64
	autoCleanup bool
}

// Item is a single cache item that holds
// the value of the item and the time the value has been pulled.
// Additionally, there is a lastAccessed value that stores the last access to
// the cache item. It is used to determine whether the
// item can be cleaned up.
type Item[C any] struct {
	value        C
	since        time.Time
	lastAccessed time.Time
}

// AutoCleanup is the interface that defines the public methods
// that allow to enable and disable automated cleanup for caches.
type AutoCleanup interface {
	EnableAutoCleanup(interval time.Duration)
	DisableAutoCleanup()
}

// New creates a new cache with the specified lifetime (in milliseconds).
func New[I comparable, C any](lifetime time.Duration) *Cache[I, C] {
	return &Cache[I, C]{
		map[I]*Item[*C]{},
		sync.RWMutex{},
		lifetime.Milliseconds(),
		false,
	}
}

// Get an Item from the cache, if there is one with a valid lifetime
func Get[I comparable, C any](cache *Cache[I, C], key I) (*C, bool) {
	cache.lock.RLock()
	item, ok := cache.cache[key]
	cache.lock.RUnlock()

	if !ok {
		return nil, false
	}

	timeDiff := time.Since(item.since).Milliseconds()
	if cache.lifetime != 0 && timeDiff > cache.lifetime {
		return nil, false
	}

	cache.lock.Lock()
	item.lastAccessed = time.Now()
	cache.lock.Unlock()

	return item.value, true
}

// Update adds or updates an item in the cache.
func Update[I comparable, C any](cache *Cache[I, C], key I, value *C) {
	cache.lock.Lock()
	defer cache.lock.Unlock()

	item, ok := cache.cache[key]

	if !ok {
		item = &Item[*C]{}

		cache.cache[key] = item
	}

	item.value = value
	item.since = time.Now()
}

// EnableAutoCleanup enables the automated cleanup of
// cache items that exceeded their lifetime since their last access.
//
// Be sure to always call DisableAutoCleanup when a cache is no longer
// needed. Else the occupied resources of the cache might be allocated forever!
//
// Note that this does not block a graceful exit.
// Enabled auto-cleanups can be ignored for application shutdown
// and interrupting a cleanup should never lead to data loss as data is only
// cached for read.
func (c *Cache[I, C]) EnableAutoCleanup(interval time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.lifetime < 1 {
		log.Warn(LogPrefix, "It makes no sense to start a cleanup routine "+
			"for infinite lifetime caches!")

		return
	}

	if c.autoCleanup {
		log.Warn(LogPrefix, "Something tried to register two cleanup tasks for the same cache!")

		return
	}

	c.autoCleanup = true
	go autoCleanupRoutine(c, interval)
}

// autoCleanupRoutine is the function used to start go routines
// that enable automated cleanup of a cache.
func autoCleanupRoutine[I comparable, C any](c *Cache[I, C], interval time.Duration) {
	for {
		c.lock.RLock()
		if !c.autoCleanup {
			break
		}

		var markedForDelete []I

		// Stage 1: Collect old items
		for key, entry := range c.cache {
			timeDiff := time.Since(entry.since).Milliseconds()

			if timeDiff > c.lifetime {
				markedForDelete = append(markedForDelete, key)
			}
		}
		c.lock.RUnlock()

		// We don't want to read lock if we can avoid
		if len(markedForDelete) == 0 {
			time.Sleep(interval)
			continue
		}

		// Stage 2: Delete old items
		c.lock.Lock()
		for _, marked := range markedForDelete {
			delete(c.cache, marked)
		}
		c.lock.Unlock()

		if len(markedForDelete) > 0 {
			log.Info(LogPrefix, "Cleaned up %v items from cache", len(markedForDelete))
		}

		time.Sleep(interval)
	}
}

// DisableAutoCleanup disables the automated cleanup by notifying the
// cleanup routine that it should stop cleaning up things.
func (c *Cache[I, C]) DisableAutoCleanup() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.autoCleanup = false
}
