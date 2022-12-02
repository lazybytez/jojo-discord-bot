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
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CacheTestSuite struct {
	suite.Suite
}

func (suite *CacheTestSuite) TestNew() {
	tables := []struct {
		lifetime time.Duration
		expected int64
	}{
		{0, 0},
		{5 * time.Second, (5 * time.Second).Milliseconds()},
		{5 * time.Minute, (5 * time.Minute).Milliseconds()},
	}

	for _, table := range tables {
		result := New[string, string](table.lifetime).lifetime

		suite.Equal(
			table.expected,
			result,
			"Arguments: %v",
			table.lifetime)
	}
}

func (suite *CacheTestSuite) TestGeWithHit() {
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
			false,
		}

		currentTime := time.Now().Add(-table.delaySeconds * time.Second)
		testCache.cache[table.key] = &Item[*string]{
			table.value,
			currentTime,
			currentTime,
		}

		resultItem, resultOk := Get(testCache, table.key)

		suite.Equalf(table.expectedOK, resultOk, "(ok check) Arguments: %v, %v", table.key, table.value)
		suite.Equalf(table.expected, resultItem, "(item check) Arguments: %v, %v", table.key, table.value)
	}
}

func (suite *CacheTestSuite) TestGeWithMiss() {
	firstKey := "my_item"

	testCache := &Cache[string, string]{
		map[string]*Item[*string]{},
		sync.RWMutex{},
		0,
		false,
	}

	// 1. Check with empty cache
	firstResultItem, firstResultOk := testCache.cache[firstKey]
	suite.False(firstResultOk, "(ok check) Arguments: %v", firstKey)
	suite.Nil(firstResultItem, "(item check) Arguments: %v", firstKey)

	// 1. Check with some value in cache
	testValue := "you expected a value, but it was me, DIO!"
	testCache.cache["some_key"] = &Item[*string]{
		&testValue,
		time.Now(),
		time.Now(),
	}

	secondResultItem, secondResultOk := Get(testCache, firstKey)
	suite.False(secondResultOk, "(ok check) Arguments: %v", firstKey)
	suite.Nil(secondResultItem, "(item check) Arguments: %v", firstKey)
}

func (suite *CacheTestSuite) TestCache_EnableAutoCleanupWithSuccess() {
	str := "some value"

	tables := []struct {
		key          string
		value        *string
		lifetime     time.Duration
		delaySeconds time.Duration
		expectedOK   bool
	}{
		{"test3", &str, 3, 1, true},
		{"test4", &str, 1, 2, false},
	}

	for _, table := range tables {
		duration := table.lifetime * time.Second

		testCache := &Cache[string, string]{
			map[string]*Item[*string]{},
			sync.RWMutex{},
			duration.Milliseconds(),
			false,
		}

		// Run every ms. With 10ms delay, cleanup should never fail during test
		err := testCache.EnableAutoCleanup(1 * time.Millisecond)
		suite.NoError(
			err,
			"Got an error when enabling auto cleanup, when no error was expected! Arguments:",
			table.key, table.value)

		testCache.lock.Lock()
		currentTime := time.Now().Add(-(table.delaySeconds * time.Second))
		testCache.cache[table.key] = &Item[*string]{
			table.value,
			currentTime,
			currentTime,
		}
		testCache.lock.Unlock()

		time.Sleep(50 * time.Millisecond) // Cleanup is async task, ensure its done

		testCache.lock.RLock()
		_, resultOk := testCache.cache[table.key]
		testCache.lock.RUnlock()

		suite.Equalf(table.expectedOK, resultOk, "(Arguments: %v, %v", table.key, table.value)
	}
}

func (suite *CacheTestSuite) TestCache_EnableAutoCleanupWithZeroOrLowerLifetime() {
	tables := []struct {
		lifetime time.Duration
	}{
		{0},
		{-1},
		{-10},
	}

	for _, table := range tables {
		duration := table.lifetime * time.Second

		testCache := &Cache[string, string]{
			map[string]*Item[*string]{},
			sync.RWMutex{},
			duration.Milliseconds(),
			false,
		}

		// Run every ms. With 10ms delay, cleanup should never fail during test
		err := testCache.EnableAutoCleanup(1 * time.Millisecond)

		testCache.lock.Lock()
		testCache.autoCleanup = false
		testCache.lock.Unlock()

		suite.Error(
			err,
			"Got an no error when enabling auto cleanup with zero or lower lifetime! Arguments:",
			table.lifetime)
	}
}

func (suite *CacheTestSuite) TestCache_EnableAutoCleanupWithMultipleCalls() {
	const lifetime = 10 * time.Millisecond

	testCache := &Cache[string, string]{
		map[string]*Item[*string]{},
		sync.RWMutex{},
		lifetime.Milliseconds(),
		false,
	}

	// Run every ms. With 10ms delay, cleanup should never fail during test
	result1 := testCache.EnableAutoCleanup(1 * time.Millisecond)
	suite.Nil(result1, "Did not expected an error on first auto cleanup enable!")

	result2 := testCache.EnableAutoCleanup(1 * time.Millisecond)
	suite.Error(
		result2,
		"Got an no error when enabling auto cleanup multiple times!")

	testCache.lock.Lock()
	testCache.autoCleanup = false
	testCache.lock.Unlock()
}

func (suite *CacheTestSuite) TestCache_DisableAutoCleanup() {
	str := "some value"

	tables := []struct {
		key          string
		value        *string
		lifetime     time.Duration
		delaySeconds time.Duration
		expectedOK   bool
	}{
		{"test3", &str, 3, 1, true},
		{"test4", &str, 1, 2, true},
	}

	for _, table := range tables {
		duration := table.lifetime * time.Second

		testCache := &Cache[string, string]{
			map[string]*Item[*string]{},
			sync.RWMutex{},
			duration.Milliseconds(),
			false,
		}

		err := testCache.EnableAutoCleanup(1 * time.Millisecond)
		suite.NoError(
			err,
			"Got an error when enabling auto cleanup, when no error was expected! Arguments:",
			table.key, table.value)

		testCache.lock.Lock()
		currentTime := time.Now().Add(-(table.delaySeconds * time.Second))
		testCache.cache[table.key] = &Item[*string]{
			table.value,
			currentTime,
			currentTime,
		}
		testCache.lock.Unlock()

		testCache.DisableAutoCleanup()

		time.Sleep(200 * time.Millisecond) // Cleanup is async task, ensure its done

		testCache.lock.RLock()
		_, resultOk := testCache.cache[table.key]
		testCache.lock.RUnlock()

		suite.Equalf(table.expectedOK, resultOk, "(Arguments: %v, %v", table.key, table.value)
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
		false,
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

func TestCache(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
