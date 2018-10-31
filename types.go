package assert

// Testing is an interface wrapper around *testing.T
type Testing interface {
	Errorf(format string, args ...interface{})
}

// Comparison a custom func that returns true on success and false on failure
type Comparison func() (ok bool)
