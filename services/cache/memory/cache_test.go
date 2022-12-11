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
	cache *InMemoryCacheProvider
}

func (suite *CacheTestSuite) SetupTest() {
	suite.cache = &InMemoryCacheProvider{
		sync.RWMutex{},
		CachePool{},
		10 * time.Second,
		nil,
	}
}

func (suite *CacheTestSuite) TestNew() {
	tables := []struct {
		lifetime time.Duration
		expected time.Duration
	}{
		{5 * time.Second, 5 * time.Second},
		{5 * time.Minute, 5 * time.Minute},
	}

	for _, table := range tables {
		result := New(table.lifetime).lifetime

		suite.Equal(
			table.expected,
			result,
			"Arguments: %v",
			table.lifetime)
	}
}

func TestCache(t *testing.T) {
	suite.Run(t, new(CacheTestSuite))
}
