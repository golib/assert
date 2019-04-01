package assert

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Errorf(format string, args ...interface{})
}

// Comparison a custom func that returns true on success and false on failure
type Comparison func() (ok bool)

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()
