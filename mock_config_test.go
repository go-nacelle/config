// DO NOT EDIT
// Code generated automatically by github.com/efritz/go-mockgen
// $ go-mockgen github.com/efritz/zubrin -i Config -o mock_config_test.go -f

package zubrin

import (
	"sync"
)

type MockConfig struct {
	LoadFunc     func(interface{}, ...TagModifier) error
	histLoad     []ConfigLoadParamSet
	MustLoadFunc func(interface{}, ...TagModifier)
	histMustLoad []ConfigMustLoadParamSet
	mutex        sync.RWMutex
}
type ConfigLoadParamSet struct {
	Arg0 interface{}
	Arg1 []TagModifier
}
type ConfigMustLoadParamSet struct {
	Arg0 interface{}
	Arg1 []TagModifier
}

func NewMockConfig() *MockConfig {
	m := &MockConfig{}
	m.LoadFunc = m.defaultLoadFunc
	m.MustLoadFunc = m.defaultMustLoadFunc
	return m
}
func (m *MockConfig) Load(v0 interface{}, v1 ...TagModifier) error {
	m.mutex.Lock()
	m.histLoad = append(m.histLoad, ConfigLoadParamSet{v0, v1})
	m.mutex.Unlock()
	return m.LoadFunc(v0, v1...)
}
func (m *MockConfig) LoadFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histLoad)
}
func (m *MockConfig) LoadFuncCallParams() []ConfigLoadParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histLoad
}

func (m *MockConfig) MustLoad(v0 interface{}, v1 ...TagModifier) {
	m.mutex.Lock()
	m.histMustLoad = append(m.histMustLoad, ConfigMustLoadParamSet{v0, v1})
	m.mutex.Unlock()
	m.MustLoadFunc(v0, v1...)
}
func (m *MockConfig) MustLoadFuncCallCount() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return len(m.histMustLoad)
}
func (m *MockConfig) MustLoadFuncCallParams() []ConfigMustLoadParamSet {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.histMustLoad
}

func (m *MockConfig) defaultLoadFunc(v0 interface{}, v1 ...TagModifier) error {
	return nil
}
func (m *MockConfig) defaultMustLoadFunc(v0 interface{}, v1 ...TagModifier) {
	return
}
