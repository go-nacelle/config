// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/zubrin -i Logger -o mock_logger_test.go -f

package zubrin

import "sync"

type MockLogger struct {
	PrintfFunc func(string, ...interface{})
	histPrintf []LoggerPrintfParamSet
	mutex      sync.RWMutex
}
type LoggerPrintfParamSet struct {
	Arg0 string
	Arg1 []interface{}
}

func NewMockLogger() *MockLogger {
	m := &MockLogger{}
	m.PrintfFunc = m.defaultPrintfFunc
	return m
}
func (m *MockLogger) Printf(v0 string, v1 ...interface{}) {
	m.mutex.Lock()
	m.histPrintf = append(m.histPrintf, LoggerPrintfParamSet{v0, v1})
	m.mutex.Unlock()
	m.PrintfFunc(v0, v1...)
}
func (m *MockLogger) PrintfFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histPrintf)
}
func (m *MockLogger) PrintfFuncCallParams() []LoggerPrintfParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histPrintf
}

func (m *MockLogger) defaultPrintfFunc(v0 string, v1 ...interface{}) {
	return
}
