package task

import "github.com/stretchr/testify/mock"

type CallbackMock struct {
	mock.Mock
}

func (m *CallbackMock) CallNoArgs() {
	m.Called()
}

func (m *CallbackMock) CallWithArgs(arg1 string, arg2 bool) (string, bool) {
	args := m.Called(arg1, arg2)
	return args.String(0), args.Bool(1)
}
