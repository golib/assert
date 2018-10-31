package assert

// Assertions provides assertion methods around the
// Testing interface.
type Assertions struct {
	t Testing
}

// New makes a new Assertions object for the specified Testing.
func New(t Testing) *Assertions {
	return &Assertions{
		t: t,
	}
}

//go:generate go run ../_codegen/main.go -output-package=assert -template=assertion_forward.go.tmpl
