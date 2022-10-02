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

package database

import (
	"github.com/lazybytez/jojo-discord-bot/api/entities"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type EntityManagerMock struct {
	mock.Mock
}

func (e *EntityManagerMock) RegisterEntity(entityType interface{}) error {
	call := e.Called(entityType)

	return call.Error(0)
}

func (e *EntityManagerMock) Create(entity interface{}) error {
	call := e.Called(entity)

	return call.Error(0)
}

func (e *EntityManagerMock) Save(entity interface{}) error {
	call := e.Called(entity)

	return call.Error(0)
}

func (e *EntityManagerMock) UpdateEntity(entityContainer interface{}, column string, value interface{}) error {
	call := e.Called(entityContainer, column, value)

	return call.Error(0)
}

func (e *EntityManagerMock) DeleteEntity(entityContainer interface{}) error {
	call := e.Called(entityContainer)

	return call.Error(0)
}

func (e *EntityManagerMock) GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error {
	call := e.Called(entityContainer, conditions)

	return call.Error(0)
}

func (e *EntityManagerMock) GetLastEntity(entityContainer interface{}, conditions ...interface{}) error {
	call := e.Called(entityContainer, conditions)

	return call.Error(0)
}

func (e *EntityManagerMock) GetEntities(entities interface{}, conditions ...interface{}) error {
	call := e.Called(entities, conditions)

	return call.Error(0)
}

func (e *EntityManagerMock) WorkOn(entityContainer interface{}) *gorm.DB {
	call := e.Called(entityContainer)

	gormDB := call.Get(0)

	switch typed := gormDB.(type) {
	case *gorm.DB:
		return typed
	default:
		return nil
	}
}

func (e *EntityManagerMock) Guilds() *entities.GuildEntityManager {
	call := e.Called()

	guildEntityManager := call.Get(0)

	switch typed := guildEntityManager.(type) {
	case *entities.GuildEntityManager:
		return typed
	default:
		return nil
	}
}

func (e *EntityManagerMock) GlobalComponentStatus() *entities.GlobalComponentStatusEntityManager {
	call := e.Called()

	gcsemEntityManager := call.Get(0)

	switch typed := gcsemEntityManager.(type) {
	case *entities.GlobalComponentStatusEntityManager:
		return typed
	default:
		return nil
	}
}

func (e *EntityManagerMock) RegisteredComponent() *entities.RegisteredComponentEntityManager {
	call := e.Called()

	rcemEntityManager := call.Get(0)

	switch typed := rcemEntityManager.(type) {
	case *entities.RegisteredComponentEntityManager:
		return typed
	default:
		return nil
	}
}

func (e *EntityManagerMock) GuildComponentStatus() *entities.GuildComponentStatusEntityManager {
	call := e.Called()

	gcemEntityManager := call.Get(0)

	switch typed := gcemEntityManager.(type) {
	case *entities.GuildComponentStatusEntityManager:
		return typed
	default:
		return nil
	}
}
