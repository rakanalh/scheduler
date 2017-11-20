package task

import "github.com/stretchr/testify/mock"

// CallbackMock is used for testing Task
type CallbackMock struct {
	mock.Mock
}

// CallNoArgs is a dummy function which accepts no arguments
func (m *CallbackMock) CallNoArgs() {
	m.Called()
}

// CallWithArgs is a dummy function which accepts two arguments
func (m *CallbackMock) CallWithArgs(arg1 string, arg2 bool) {
	m.Called(arg1, arg2)
}

// CallWithChan is a dummy function to test failure with non-hashable parameters
func (m *CallbackMock) CallWithChan(channel chan bool) {
	m.Called(channel)
}
