package mocks

import (
	"github.com/stretchr/testify/mock"
)

type Logger struct {
	mock.Mock
}

func (m *Logger) Debug(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *Logger) Info(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *Logger) Warn(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *Logger) Error(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *Logger) Fatal(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}
