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

package entity_manager_mock

import (
	"github.com/lazybytez/jojo-discord-bot/services"
	"github.com/stretchr/testify/mock"
)

type EntityManagerMock struct {
	mock.Mock
}

func (em *EntityManagerMock) RegisterEntity(entityType interface{}) error {
	call := em.Called(entityType)

	return call.Error(0)
}

func (em *EntityManagerMock) DB() services.DatabaseAccess {
	call := em.Called()

	switch v := call.Get(0).(type) {
	case services.DatabaseAccess:
		return v
	default:
		return nil
	}
}

func (em *EntityManagerMock) Logger() services.Logger {
	call := em.Called()

	switch v := call.Get(0).(type) {
	case services.Logger:
		return v
	default:
		return nil
	}
}
