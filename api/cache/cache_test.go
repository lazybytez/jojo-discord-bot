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
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tables := []struct {
		lifetime time.Duration
		expected int64
	}{
		{0, 0},
		{5 * time.Second, 5000},
		{5 * time.Minute, 300 * 1000},
	}

	for _, table := range tables {
		if cache := New[string, string](table.lifetime); cache.lifetime != table.expected {
			t.Errorf("output of cache new was incorrect with %v, got: %v, want: %v.",
				table.lifetime,
				cache.lifetime,
				table.expected)
		}
	}
}

func TestGet(t *testing.T) {
	str := "some value"

	tables := []struct {
		key          string
		value        *string
		lifetime     time.Duration
		delaySeconds time.Duration
		expected     *string
		expectedOK   bool
	}{
		{"test1", &str, 0, 0, &str, true},
		{"test2", &str, 0, 2, &str, true},
		{"test3", &str, 3, 1, &str, true},
		{"test4", &str, 1, 2, nil, false},
	}

	for _, table := range tables {
		duration := table.lifetime * time.Second

		testCache := &Cache[string, string]{
			map[string]*Item[*string]{},
			sync.RWMutex{},
			duration.Milliseconds(),
			nil,
		}

		testCache.cache[table.key] = &Item[*string]{
			table.value,
			time.Now(),
			time.Now(),
		}

		time.Sleep(table.delaySeconds * time.Second)

		if r, ok := Get(testCache, table.key); r != table.expected || ok != table.expectedOK {
			t.Errorf("output of cache get with key \"%v\" and "+
				"value \"%v\" was incorrect, got: %v & %v, want: %v & %v.",
				table.key,
				table.value,
				r,
				ok,
				table.expected,
				table.expectedOK)
		}
	}
}

func TestCache_EnableAutoCleanup(t *testing.T) {
	str := "some value"

	tables := []struct {
		key          string
		value        *string
		lifetime     time.Duration
		delaySeconds time.Duration
		expectedOK   bool
	}{
		{"test1", &str, 0, 0, true},
		{"test2", &str, 0, 2, true},
		{"test3", &str, 3, 1, true},
		{"test4", &str, 1, 2, false},
	}

	for _, table := range tables {
		duration := table.lifetime * time.Second

		testCache := &Cache[string, string]{
			map[string]*Item[*string]{},
			sync.RWMutex{},
			duration.Milliseconds(),
			nil,
		}
		testCache.EnableAutoCleanup(500 * time.Millisecond)

		testCache.cache[table.key] = &Item[*string]{
			table.value,
			time.Now(),
			time.Now(),
		}

		time.Sleep(table.delaySeconds * time.Second)

		if _, ok := testCache.cache[table.key]; ok != table.expectedOK {
			t.Errorf("output of cache get with key \"%v\" and "+
				"value \"%v\" was incorrect, got: %v, want: %v.",
				table.key,
				table.value,
				ok,
				table.expectedOK)
		}
	}
}

func TestUpdate(t *testing.T) {
	const testKey = "test"
	testDataA := "A"
	testDataB := "B"

	cache := &Cache[string, string]{
		map[string]*Item[*string]{},
		sync.RWMutex{},
		0,
		nil,
	}

	if r, ok := cache.cache[testKey]; ok {
		t.Errorf("got a result when no result was expected, got: %v & %v, want: %v & %v.",
			r,
			ok,
			nil,
			false)
	}

	Update(cache, testKey, &testDataA)

	if r, ok := cache.cache[testKey]; ok && r.value != &testDataA {
		t.Errorf("got wrong result, got: %v & %v, want: %v & %v.",
			r,
			ok,
			testDataA,
			true)
	}

	Update(cache, testKey, &testDataB)

	if r, ok := cache.cache[testKey]; ok && r.value != &testDataB {
		t.Errorf("got wrong result, got: %v & %v, want: %v & %v.",
			r,
			ok,
			testDataB,
			true)
	}
}
