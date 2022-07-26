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

package handler_manager_mock

import (
	"github.com/lazybytez/jojo-discord-bot/api"
	"github.com/stretchr/testify/mock"
)

// HandlerManagerMock is a custom handler manager embedding
// mock.Mock and allows to do expectations on handler management methods.
type HandlerManagerMock struct {
	mock.Mock
}

func (h *HandlerManagerMock) Register(name string, handler interface{}) (api.HandlerName, bool) {
	result := h.Called(name, handler)

	return api.HandlerName(result.String(0)), result.Bool(1)
}

func (h *HandlerManagerMock) RegisterOnce(name string, handler interface{}) (api.HandlerName, bool) {
	result := h.Called(name, handler)

	return api.HandlerName(result.String(0)), result.Bool(1)
}

func (h *HandlerManagerMock) Unregister(name api.HandlerName) error {
	result := h.Called(name)

	return result.Error(0)
}

func (h *HandlerManagerMock) UnregisterAll() {
	h.Called()
}
