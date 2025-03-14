package assert

type (
	// Testing is an interface wrapper around *testing.T
	Testing interface {
		Errorf(format string, args ...interface{})
	}

	failNower interface {
		FailNow()
	}
)

type (
	// Comparison a custom func that returns true on success and false on failure
	Comparison func() (ok bool)

	// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
	// methods, and represents a simple func that takes no arguments, and returns nothing.
	PanicTestFunc func()
)
