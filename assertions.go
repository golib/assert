package assert

import (
	"errors"
	"io"
	"time"
)

// Assertions provides asserts around the
// Testing interface.
type Assertions struct {
	t Testing
}

// New creates a new *Assertions for the Testing specified.
func New(t Testing) *Assertions {
	return &Assertions{
		t: t,
	}
}

// Fail reports a failure through
func (it *Assertions) Fail(message string, formatAndArgs ...interface{}) bool {
	return Fail(it.t, message, formatAndArgs...)
}

// FailNow fails test
func (it *Assertions) FailNow(message string, formatAndArgs ...interface{}) bool {
	return FailNow(it.t, message, formatAndArgs...)
}

// IsType asserts that the v is of the same type.
//
//	it.IsType(int, 123)
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) IsType(expectedType, v interface{}, formatAndArgs ...interface{}) bool {
	return IsType(it.t, expectedType, v, formatAndArgs...)
}

// Implements asserts that the v is implemented by the interface.
//
//	it.Implements((*Iface)(nil), new(v))
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Implements(iface, v interface{}, formatAndArgs ...interface{}) bool {
	return Implements(it.t, iface, v, formatAndArgs...)
}

// Contains asserts that the list(string, array, slice...) or map contains the
// sub string or element.
//
//	it.Contains("Hello World", "World", "'Hello World' does contain 'World'")
//	it.Contains([]string{"Hello", "World"}, "World", "["Hello", "World"] does contain 'World'")
//	it.Contains(map[string]string{"Hello": "World"}, "Hello", "{'Hello': 'World'} does contain 'Hello'")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Contains(list, contains interface{}, formatAndArgs ...interface{}) bool {
	return Contains(it.t, list, contains, formatAndArgs...)
}

// NotContains asserts that the list(string, array, slice...) or map does NOT contain the
// sub string or element.
//
//	it.NotContains("Hello World", "Earth", "'Hello World' does NOT contain 'Earth'")
//	it.NotContains([]string{"Hello", "World", "Earth", "['Hello', 'World'] does NOT contain 'Earth'")
//	it.NotContains(map[string]string{"Hello": "World"}, "Earth", "{'Hello': 'World'} does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotContains(list, contains interface{}, formatAndArgs ...interface{}) bool {
	return NotContains(it.t, list, contains, formatAndArgs...)
}

// Match asserts that the regexp matches a string.
//
//	it.Match(regexp.MustCompile("start"), "it's starting")
//	it.Match("start...$", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Match(reg, str interface{}, formatAndArgs ...interface{}) bool {
	return Match(it.t, reg, str, formatAndArgs...)
}

// NotMatch asserts that the regexp does not match a string.
//
//	it.NotMatch(regexp.MustCompile("starts"), "it's starting")
//	it.NotMatch("^start", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotMatch(reg, str interface{}, formatAndArgs ...interface{}) bool {
	return NotMatch(it.t, reg, str, formatAndArgs...)
}

// Equal asserts that two objects are equal.
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
//
//	it.Equal(123, 123, "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Equal(expected, actual interface{}, formatAndArgs ...interface{}) bool {
	return Equal(it.t, expected, actual, formatAndArgs...)
}

// NotEqual asserts that the two objects are NOT equal.
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
//
//	it.NotEqual(123, "123", "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotEqual(expected, actual interface{}, formatAndArgs ...interface{}) bool {
	return NotEqual(it.t, expected, actual, formatAndArgs...)
}

// EqualValues asserts that two objects are equal
// or convertable to the same types and equal.
//
//	it.EqualValues(uint32(123), int32(123), "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) EqualValues(expected, actual interface{}, formatAndArgs ...interface{}) bool {
	return EqualValues(it.t, expected, actual, formatAndArgs...)
}

// Exactly asserts that two objects are equal in both values and types.
//
//	it.Exactly(int32(123), int64(123), "int32 and int64 should NOT be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Exactly(expected, actual interface{}, formatAndArgs ...interface{}) bool {
	return Exactly(it.t, expected, actual, formatAndArgs...)
}

// Condition uses a custom Comparison to assert a complex condition.
func (it *Assertions) Condition(comp Comparison, formatAndArgs ...interface{}) bool {
	return Condition(it.t, comp, formatAndArgs...)
}

// Empty asserts that the v is empty.  I.e. nil, "", false, 0,
// or list(slice, map, channel) with len == 0.
//
//	it.Empty(v)
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Empty(v interface{}, formatAndArgs ...interface{}) bool {
	return Empty(it.t, v, formatAndArgs...)
}

// NotEmpty asserts that the v is NOT empty.  I.e. not nil, "", false, 0,
// or list(slice, map, channel) with len == 0.
//
//	if it.NotEmpty(v) {
//	  assert.Equal(t, "two", v[1])
//	}
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotEmpty(v interface{}, formatAndArgs ...interface{}) bool {
	return NotEmpty(it.t, v, formatAndArgs...)
}

// True asserts that the specified value is true.
//
//	it.True(myBool, "myBool should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) True(value bool, formatAndArgs ...interface{}) bool {
	return True(it.t, value, formatAndArgs...)
}

// False asserts that the specified value is false.
//
//	it.False(myBool, "myBool should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) False(value bool, formatAndArgs ...interface{}) bool {
	return False(it.t, value, formatAndArgs...)
}

// Zero asserts that v is the zero value for its type and returns the truth.
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Zero(v interface{}, formatAndArgs ...interface{}) bool {
	return Zero(it.t, v, formatAndArgs...)
}

// NotZero asserts that the v is not the zero value.
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotZero(v interface{}, formatAndArgs ...interface{}) bool {
	return NotZero(it.t, v, formatAndArgs...)
}

// Len asserts that the a v has specific length.
// Len also fails if the v has a type that len() not accept.
//
//	it.Len(mySlice, 3, "The size of slice is not 3")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Len(v interface{}, length int, formatAndArgs ...interface{}) bool {
	return Len(it.t, v, length, formatAndArgs...)
}

// Nil asserts that the v is nil.
//
//	it.Nil(err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Nil(v interface{}, formatAndArgs ...interface{}) bool {
	return Nil(it.t, v, formatAndArgs...)
}

// NotNil asserts that the v is not nil.
//
//	it.NotNil(err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotNil(v interface{}, formatAndArgs ...interface{}) bool {
	return NotNil(it.t, v, formatAndArgs...)
}

// Error asserts that a func returned an error (i.e. not `nil`).
//
//	  actual, err := SomeFunc()
//	  if it.Error(err, "An error was expected") {
//		   assert.Equal(t, err, ErrNotFound)
//	  }
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Error(err error, formatAndArgs ...interface{}) bool {
	return Error(it.t, err, formatAndArgs...)
}

// NotError asserts that a func returned not an error (i.e. `nil`).
//
//	  actual, err := SomeFunction()
//	  if it.NotError(err) {
//		   assert.Equal(t, actual, expected)
//	  }
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotError(err error, formatAndArgs ...interface{}) bool {
	return NotError(it.t, err, formatAndArgs...)
}

// EqualError asserts that an error.Error() (i.e. not `nil`) is equal to expected string.
//
//	_, err := SomeFunc()
//	it.EqualError(err,  errString, "An error was expected")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) EqualError(err error, str string, formatAndArgs ...interface{}) bool {
	return EqualErrors(it.t, err, errors.New(str), formatAndArgs...)
}

// EqualErrors asserts that two errors (i.e. not `nil`) are equal.
//
//	_, err := SomeFunc()
//	it.EqualErrors(err,  ErrNotFound, "An not found error was expected")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) EqualErrors(expectedErr, actualErr error, formatAndArgs ...interface{}) bool {
	return EqualErrors(it.t, actualErr, expectedErr, formatAndArgs...)
}

// InDelta asserts that the two numerals are within delta of each other.
//
//	it.InDelta(math.Pi, (22 / 7.0), 0.01)
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) InDelta(expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	return InDelta(it.t, expected, actual, delta, formatAndArgs...)
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) InDeltaSlice(expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	return InDeltaSlice(it.t, expected, actual, delta, formatAndArgs...)
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//	it.WithinDuration(time.Now(), time.Now(), 10*time.Second, "The difference should not be more than 10s")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) WithinDuration(expected time.Time, actual time.Time, delta time.Duration, formatAndArgs ...interface{}) bool {
	return WithinDuration(it.t, expected, actual, delta, formatAndArgs...)
}

// ReaderContains asserts that io.Reader contains the specified sub string or element.
//
//	reader := bytes.NewBuffer([]byte("Hello, world!"))
//	it.ReaderContains(reader, "world")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) ReaderContains(reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	return ReaderContains(it.t, reader, contains, formatAndArgs...)
}

// ReaderNotContains asserts that reader does NOT contain the specified substring or element.
//
//	reader := bytes.NewBuffer([]byte("Hello, world!"))
//	it.ReaderNotContains(reader, "test")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) ReaderNotContains(reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	return ReaderNotContains(it.t, reader, contains, formatAndArgs...)
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//	it.Panics(func(){
//	  GoCrazy()
//	}, "Calling GoCrazy() should panic")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) Panics(f PanicTestFunc, formatAndArgs ...interface{}) bool {
	return Panics(it.t, f, formatAndArgs...)
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//	it.NotPanics(func(){
//	  RemainCalm()
//	}, "Calling RemainCalm() should NOT panic")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotPanics(f PanicTestFunc, formatAndArgs ...interface{}) bool {
	return NotPanics(it.t, f, formatAndArgs...)
}

// EqualJSON asserts that two JSON strings are equivalent.
//
//	it.EqualJSON(`{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) EqualJSON(expected string, actual string, formatAndArgs ...interface{}) bool {
	return EqualJSON(it.t, expected, actual, formatAndArgs...)
}

// ContainsJSON asserts that JSON string contains value of the key.
//
//	it.ContainsJSON(`{"hello": "world", "foo": "bar"}`, "hello", "world")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) ContainsJSON(actual, key string, v interface{}) bool {
	return ContainsJSON(it.t, actual, key, v)
}

// NotContainsJSON asserts that JSON string does not contain attribute of the key.
//
//	it.NotContainsJSON(`{"hello": "world", "foo": "bar"}`, "world")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotContainsJSON(actual, key string) bool {
	return NotContainsJSON(it.t, actual, key)
}

// NotEmptyJSON asserts that JSON string contains attribute of the key with not empty value.
//
//	it.NotEmptyJSON(`{"hello": "world", "foo": "bar"}`, "foo")
//
// Returns whether the assertion was successful (true) or not (false).
func (it *Assertions) NotEmptyJSON(actual, key string) bool {
	return NotEmptyJSON(it.t, actual, key)
}
