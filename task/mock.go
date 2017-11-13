package task

import "github.com/stretchr/testify/mock"

type CallbackMock struct {
	mock.Mock
}

func (m *CallbackMock) CallNoArgs() {
	m.Called()
}

func (m *CallbackMock) CallWithArgs(arg1 string, arg2 bool) {
	m.Called(arg1, arg2)
}

func (m *CallbackMock) CallWithChan(channel chan bool) {
	m.Called(channel)
}
