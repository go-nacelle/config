// Code generated by go-mockgen 0.1.0; DO NOT EDIT.

package config

import "sync"

// MockSourcer is a mock implementation of the Sourcer interface (from the
// package github.com/go-nacelle/config) used for unit testing.
type MockSourcer struct {
	// AssetsFunc is an instance of a mock function object controlling the
	// behavior of the method Assets.
	AssetsFunc *SourcerAssetsFunc
	// DumpFunc is an instance of a mock function object controlling the
	// behavior of the method Dump.
	DumpFunc *SourcerDumpFunc
	// GetFunc is an instance of a mock function object controlling the
	// behavior of the method Get.
	GetFunc *SourcerGetFunc
	// InitFunc is an instance of a mock function object controlling the
	// behavior of the method Init.
	InitFunc *SourcerInitFunc
	// TagsFunc is an instance of a mock function object controlling the
	// behavior of the method Tags.
	TagsFunc *SourcerTagsFunc
}

// NewMockSourcer creates a new mock of the Sourcer interface. All methods
// return zero values for all results, unless overwritten.
func NewMockSourcer() *MockSourcer {
	return &MockSourcer{
		AssetsFunc: &SourcerAssetsFunc{
			defaultHook: func() []string {
				return nil
			},
		},
		DumpFunc: &SourcerDumpFunc{
			defaultHook: func() map[string]string {
				return nil
			},
		},
		GetFunc: &SourcerGetFunc{
			defaultHook: func([]string) (string, SourcerFlag, error) {
				return "", 0, nil
			},
		},
		InitFunc: &SourcerInitFunc{
			defaultHook: func() error {
				return nil
			},
		},
		TagsFunc: &SourcerTagsFunc{
			defaultHook: func() []string {
				return nil
			},
		},
	}
}

// NewMockSourcerFrom creates a new mock of the MockSourcer interface. All
// methods delegate to the given implementation, unless overwritten.
func NewMockSourcerFrom(i Sourcer) *MockSourcer {
	return &MockSourcer{
		AssetsFunc: &SourcerAssetsFunc{
			defaultHook: i.Assets,
		},
		DumpFunc: &SourcerDumpFunc{
			defaultHook: i.Dump,
		},
		GetFunc: &SourcerGetFunc{
			defaultHook: i.Get,
		},
		InitFunc: &SourcerInitFunc{
			defaultHook: i.Init,
		},
		TagsFunc: &SourcerTagsFunc{
			defaultHook: i.Tags,
		},
	}
}

// SourcerAssetsFunc describes the behavior when the Assets method of the
// parent MockSourcer instance is invoked.
type SourcerAssetsFunc struct {
	defaultHook func() []string
	hooks       []func() []string
	history     []SourcerAssetsFuncCall
	mutex       sync.Mutex
}

// Assets delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSourcer) Assets() []string {
	r0 := m.AssetsFunc.nextHook()()
	m.AssetsFunc.appendCall(SourcerAssetsFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Assets method of the
// parent MockSourcer instance is invoked and the hook queue is empty.
func (f *SourcerAssetsFunc) SetDefaultHook(hook func() []string) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Assets method of the parent MockSourcer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SourcerAssetsFunc) PushHook(hook func() []string) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SourcerAssetsFunc) SetDefaultReturn(r0 []string) {
	f.SetDefaultHook(func() []string {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SourcerAssetsFunc) PushReturn(r0 []string) {
	f.PushHook(func() []string {
		return r0
	})
}

func (f *SourcerAssetsFunc) nextHook() func() []string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SourcerAssetsFunc) appendCall(r0 SourcerAssetsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SourcerAssetsFuncCall objects describing
// the invocations of this function.
func (f *SourcerAssetsFunc) History() []SourcerAssetsFuncCall {
	f.mutex.Lock()
	history := make([]SourcerAssetsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SourcerAssetsFuncCall is an object that describes an invocation of method
// Assets on an instance of MockSourcer.
type SourcerAssetsFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []string
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SourcerAssetsFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SourcerAssetsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// SourcerDumpFunc describes the behavior when the Dump method of the parent
// MockSourcer instance is invoked.
type SourcerDumpFunc struct {
	defaultHook func() map[string]string
	hooks       []func() map[string]string
	history     []SourcerDumpFuncCall
	mutex       sync.Mutex
}

// Dump delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSourcer) Dump() map[string]string {
	r0 := m.DumpFunc.nextHook()()
	m.DumpFunc.appendCall(SourcerDumpFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Dump method of the
// parent MockSourcer instance is invoked and the hook queue is empty.
func (f *SourcerDumpFunc) SetDefaultHook(hook func() map[string]string) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Dump method of the parent MockSourcer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SourcerDumpFunc) PushHook(hook func() map[string]string) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SourcerDumpFunc) SetDefaultReturn(r0 map[string]string) {
	f.SetDefaultHook(func() map[string]string {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SourcerDumpFunc) PushReturn(r0 map[string]string) {
	f.PushHook(func() map[string]string {
		return r0
	})
}

func (f *SourcerDumpFunc) nextHook() func() map[string]string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SourcerDumpFunc) appendCall(r0 SourcerDumpFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SourcerDumpFuncCall objects describing the
// invocations of this function.
func (f *SourcerDumpFunc) History() []SourcerDumpFuncCall {
	f.mutex.Lock()
	history := make([]SourcerDumpFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SourcerDumpFuncCall is an object that describes an invocation of method
// Dump on an instance of MockSourcer.
type SourcerDumpFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 map[string]string
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SourcerDumpFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SourcerDumpFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// SourcerGetFunc describes the behavior when the Get method of the parent
// MockSourcer instance is invoked.
type SourcerGetFunc struct {
	defaultHook func([]string) (string, SourcerFlag, error)
	hooks       []func([]string) (string, SourcerFlag, error)
	history     []SourcerGetFuncCall
	mutex       sync.Mutex
}

// Get delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSourcer) Get(v0 []string) (string, SourcerFlag, error) {
	r0, r1, r2 := m.GetFunc.nextHook()(v0)
	m.GetFunc.appendCall(SourcerGetFuncCall{v0, r0, r1, r2})
	return r0, r1, r2
}

// SetDefaultHook sets function that is called when the Get method of the
// parent MockSourcer instance is invoked and the hook queue is empty.
func (f *SourcerGetFunc) SetDefaultHook(hook func([]string) (string, SourcerFlag, error)) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Get method of the parent MockSourcer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SourcerGetFunc) PushHook(hook func([]string) (string, SourcerFlag, error)) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SourcerGetFunc) SetDefaultReturn(r0 string, r1 SourcerFlag, r2 error) {
	f.SetDefaultHook(func([]string) (string, SourcerFlag, error) {
		return r0, r1, r2
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SourcerGetFunc) PushReturn(r0 string, r1 SourcerFlag, r2 error) {
	f.PushHook(func([]string) (string, SourcerFlag, error) {
		return r0, r1, r2
	})
}

func (f *SourcerGetFunc) nextHook() func([]string) (string, SourcerFlag, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SourcerGetFunc) appendCall(r0 SourcerGetFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SourcerGetFuncCall objects describing the
// invocations of this function.
func (f *SourcerGetFunc) History() []SourcerGetFuncCall {
	f.mutex.Lock()
	history := make([]SourcerGetFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SourcerGetFuncCall is an object that describes an invocation of method
// Get on an instance of MockSourcer.
type SourcerGetFuncCall struct {
	// Arg0 is the value of the 1st argument passed to this method
	// invocation.
	Arg0 []string
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 string
	// Result1 is the value of the 2nd result returned from this method
	// invocation.
	Result1 SourcerFlag
	// Result2 is the value of the 3rd result returned from this method
	// invocation.
	Result2 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SourcerGetFuncCall) Args() []interface{} {
	return []interface{}{c.Arg0}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SourcerGetFuncCall) Results() []interface{} {
	return []interface{}{c.Result0, c.Result1, c.Result2}
}

// SourcerInitFunc describes the behavior when the Init method of the parent
// MockSourcer instance is invoked.
type SourcerInitFunc struct {
	defaultHook func() error
	hooks       []func() error
	history     []SourcerInitFuncCall
	mutex       sync.Mutex
}

// Init delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSourcer) Init() error {
	r0 := m.InitFunc.nextHook()()
	m.InitFunc.appendCall(SourcerInitFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Init method of the
// parent MockSourcer instance is invoked and the hook queue is empty.
func (f *SourcerInitFunc) SetDefaultHook(hook func() error) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Init method of the parent MockSourcer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SourcerInitFunc) PushHook(hook func() error) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SourcerInitFunc) SetDefaultReturn(r0 error) {
	f.SetDefaultHook(func() error {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SourcerInitFunc) PushReturn(r0 error) {
	f.PushHook(func() error {
		return r0
	})
}

func (f *SourcerInitFunc) nextHook() func() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SourcerInitFunc) appendCall(r0 SourcerInitFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SourcerInitFuncCall objects describing the
// invocations of this function.
func (f *SourcerInitFunc) History() []SourcerInitFuncCall {
	f.mutex.Lock()
	history := make([]SourcerInitFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SourcerInitFuncCall is an object that describes an invocation of method
// Init on an instance of MockSourcer.
type SourcerInitFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 error
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SourcerInitFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SourcerInitFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}

// SourcerTagsFunc describes the behavior when the Tags method of the parent
// MockSourcer instance is invoked.
type SourcerTagsFunc struct {
	defaultHook func() []string
	hooks       []func() []string
	history     []SourcerTagsFuncCall
	mutex       sync.Mutex
}

// Tags delegates to the next hook function in the queue and stores the
// parameter and result values of this invocation.
func (m *MockSourcer) Tags() []string {
	r0 := m.TagsFunc.nextHook()()
	m.TagsFunc.appendCall(SourcerTagsFuncCall{r0})
	return r0
}

// SetDefaultHook sets function that is called when the Tags method of the
// parent MockSourcer instance is invoked and the hook queue is empty.
func (f *SourcerTagsFunc) SetDefaultHook(hook func() []string) {
	f.defaultHook = hook
}

// PushHook adds a function to the end of hook queue. Each invocation of the
// Tags method of the parent MockSourcer instance invokes the hook at the
// front of the queue and discards it. After the queue is empty, the default
// hook function is invoked for any future action.
func (f *SourcerTagsFunc) PushHook(hook func() []string) {
	f.mutex.Lock()
	f.hooks = append(f.hooks, hook)
	f.mutex.Unlock()
}

// SetDefaultReturn calls SetDefaultDefaultHook with a function that returns
// the given values.
func (f *SourcerTagsFunc) SetDefaultReturn(r0 []string) {
	f.SetDefaultHook(func() []string {
		return r0
	})
}

// PushReturn calls PushDefaultHook with a function that returns the given
// values.
func (f *SourcerTagsFunc) PushReturn(r0 []string) {
	f.PushHook(func() []string {
		return r0
	})
}

func (f *SourcerTagsFunc) nextHook() func() []string {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	if len(f.hooks) == 0 {
		return f.defaultHook
	}

	hook := f.hooks[0]
	f.hooks = f.hooks[1:]
	return hook
}

func (f *SourcerTagsFunc) appendCall(r0 SourcerTagsFuncCall) {
	f.mutex.Lock()
	f.history = append(f.history, r0)
	f.mutex.Unlock()
}

// History returns a sequence of SourcerTagsFuncCall objects describing the
// invocations of this function.
func (f *SourcerTagsFunc) History() []SourcerTagsFuncCall {
	f.mutex.Lock()
	history := make([]SourcerTagsFuncCall, len(f.history))
	copy(history, f.history)
	f.mutex.Unlock()

	return history
}

// SourcerTagsFuncCall is an object that describes an invocation of method
// Tags on an instance of MockSourcer.
type SourcerTagsFuncCall struct {
	// Result0 is the value of the 1st result returned from this method
	// invocation.
	Result0 []string
}

// Args returns an interface slice containing the arguments of this
// invocation.
func (c SourcerTagsFuncCall) Args() []interface{} {
	return []interface{}{}
}

// Results returns an interface slice containing the results of this
// invocation.
func (c SourcerTagsFuncCall) Results() []interface{} {
	return []interface{}{c.Result0}
}
