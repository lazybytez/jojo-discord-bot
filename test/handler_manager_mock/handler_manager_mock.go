package handler_manager_mock

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

// HandlerManagerMock is a custom handler manager embedding
// mock.Mock and allows to do expectations on handler management methods.
type HandlerManagerMock struct {
	mock.Mock
}

func (h *HandlerManagerMock) RegisterSimpleMessageHandler(name string, handler func(session *discordgo.Session, create *discordgo.MessageCreate), messages ...string) (string, bool) {
	result := h.Called(name, handler, messages)

	return result.String(0), result.Bool(1)
}

func (h *HandlerManagerMock) RegisterComplexMessageHandler(name string, handler func(session *discordgo.Session, create *discordgo.MessageCreate), messages ...string) (string, bool) {
	result := h.Called(name, handler, messages)

	return result.String(0), result.Bool(1)
}

func (h *HandlerManagerMock) Register(name string, handler interface{}) (string, bool) {
	result := h.Called(name, handler)

	return result.String(0), result.Bool(1)
}

func (h *HandlerManagerMock) RegisterOnce(name string, handler interface{}) (string, bool) {
	result := h.Called(name, handler)

	return result.String(0), result.Bool(1)
}

func (h *HandlerManagerMock) Unregister(name string) error {
	result := h.Called(name)

	return result.Error(0)
}

func (h *HandlerManagerMock) AddDecorator(name string, decorator interface{}) bool {
	result := h.Called(name, decorator)

	return result.Bool(1)
}

func (h *HandlerManagerMock) UnregisterAll() {
	h.Called()
}
