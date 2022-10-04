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

package db

import (
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type DatabaseAccessMock struct {
	mock.Mock
}

func (dba *DatabaseAccessMock) RegisterEntity(entityType interface{}) error {
	call := dba.Called(entityType)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) Create(entity interface{}) error {
	call := dba.Called(entity)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) Save(entity interface{}) error {
	call := dba.Called(entity)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) UpdateEntity(entityContainer interface{}, column string, value interface{}) error {
	call := dba.Called(entityContainer, column, value)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) DeleteEntity(entityContainer interface{}) error {
	call := dba.Called(entityContainer)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) GetFirstEntity(entityContainer interface{}, conditions ...interface{}) error {
	call := dba.Called(entityContainer, conditions)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) GetLastEntity(entityContainer interface{}, conditions ...interface{}) error {
	call := dba.Called(entityContainer, conditions)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) GetEntities(entities interface{}, conditions ...interface{}) error {
	call := dba.Called(entities, conditions)

	return call.Error(0)
}

func (dba *DatabaseAccessMock) WorkOn(entityContainer interface{}) *gorm.DB {
	call := dba.Called(entityContainer)

	gormDB := call.Get(0)

	switch typed := gormDB.(type) {
	case *gorm.DB:
		return typed
	default:
		return nil
	}
}

func (dba *DatabaseAccessMock) DB() *gorm.DB {
	call := dba.Called()

	gormDB := call.Get(0)

	switch typed := gormDB.(type) {
	case *gorm.DB:
		return typed
	default:
		return nil
	}
}
