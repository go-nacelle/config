// Code generated by go-mockgen 0.1.0; DO NOT EDIT.

package mocks

import (
	structtag "github.com/fatih/structtag"
	config "github.com/go-nacelle/config"
	"reflect"
	"sync"
)

// MockTagModifier is a mock implementation of the TagModifier interface
// (from the package github.com/go-nacelle/config) used for unit testing.
type MockTagModifier struct {
	// AlterFieldTagFunc is an instance of a mock function object
	// controlling the behavior of the method AlterFieldTag.
	AlterFieldTagFunc *TagModifierAlterFieldTagFunc
}

// NewMockTagModifier creates a new mock of the TagModifier interface. All
// methods return zero values for all results, unless overwritten.
func NewMockTagModifier() *MockTagModifier {
	return &MockTagModifier{
		AlterFieldTagFunc: &TagModifierAlterFieldTagFunc{
			defaultHook: func(reflect.StructField, *structtag.Tags) error {
				return nil
			},
		},
	}
}

// NewMockTagModifierFrom creates a new mock of the MockTagModifier
// interface. All methods delegate to the given implementation, unless
// overwritten.
func NewMockTagModifierFrom(i config.TagModifier) *MockTagModifier {
	return &MockTagModifier{
		AlterFieldTagFunc: &TagModifierAlterFieldTagFunc{
			defaultHook: i.AlterFieldTag,
		},
	}
}

// TagModifierAlterFieldTagFunc describes the behavior when the
// AlterFieldTag method of the parent MockTagModifier instance is invoked.
type TagModifierAlterFieldTagFunc struct {
	defaultHook func(reflect.StructField, *structtag.Tags) error
	hooks       []func(reflect.StructField, *structtag.Tags) error
	history     []TagModifierAlterFieldTagFuncCall
	mutex       sync.Mutex
}

// AlterFieldTag delegates to the next hook function in the queue and stores
// the parameter and result values of this invocation.
func (m *MockTagModifier) AlterFieldTag(v0 reflect.StructField, v1 *structtag.Tags) error {
	r0 := m.AlterFieldTagFunc.nextHook()(v0, v1)
	m.AlterFieldTagFunc.appendCall(TagModifierAlterFieldTagFuncCall{v0, v1, r0})
	return r0
}

// SetDefaultHook sets function that is called when the AlterFieldTag method
// of the parent MockTagModifier instance is invoked and the hook queue is
// empty.
func (f *TagModifierAlterFieldTagFunc) SetDefaultHook(hook func(reflect.StructField, *structtag.Tags) error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// AlterFieldTag method of the parent MockTagModifier instance invokes the
// hook at the front of the queue and discards it. After the queue is empty,
// the default hook function is invoked for any future action.
func (f *TagModifierAlterFieldTagFunc) PushHook(hook func(reflect.StructField, *structtag.Tags) error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *TagModifierAlterFieldTagFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func(reflect.StructField, *structtag.Tags) error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *TagModifierAlterFieldTagFunc) PushReturn(r0 error) {
	f.PushHook(func(reflect.StructField, *structtag.Tags) error {
		return r0
	})
}

func (f *TagModifierAlterFieldTagFunc) nextHook() func(reflect.StructField, *structtag.Tags) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *TagModifierAlterFieldTagFunc) appendCall(r0 TagModifierAlterFieldTagFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of TagModifierAlterFieldTagFuncCall objects
// describing the invocations of this function.
func (f *TagModifierAlterFieldTagFunc) History() []TagModifierAlterFieldTagFuncCall {
	f.mutex.Lock()
	history := make([]TagModifierAlterFieldTagFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// TagModifierAlterFieldTagFuncCall is an object that describes an
// invocation of method AlterFieldTag on an instance of MockTagModifier.
type TagModifierAlterFieldTagFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 reflect.StructField
	// Arg1 is the value of the 2nd argument passed to this method
	// invocation.
	Arg1 *structtag.Tags
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c TagModifierAlterFieldTagFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0, c.Arg1}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c TagModifierAlterFieldTagFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
