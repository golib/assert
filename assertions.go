package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/pmezard/go-difflib/difflib"
)

var (
	spewConfig = spew.ConfigState{
		Indent:                  " ",
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SortKeys:                true,
	}

	numericZeros = []interface{}{
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		float32(0),
		float64(0),
	}
)

// IsType asserts that the specified v are of the same type.
func IsType(t Testing, expectedType, v interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualObjects(reflect.TypeOf(v), reflect.TypeOf(expectedType)) {
		return Fail(t,
			fmt.Sprintf("Object expected to be of type %v, but was %v", reflect.TypeOf(expectedType), reflect.TypeOf(v)),
			formatAndArgs...)
	}

	return true
}

// Implements asserts that v is implemented by the specified interface.
//
//    assert.Implements(t, (*Iface)(nil), new(Obj), "Oops~")
func Implements(t Testing, iface, v interface{}, formatAndArgs ...interface{}) bool {
	ifaceType := reflect.TypeOf(iface).Elem()

	if !reflect.TypeOf(v).Implements(ifaceType) {
		return Fail(t,
			fmt.Sprintf("%T must implement %v", v, ifaceType),
			formatAndArgs...)
	}

	return true
}

// Equal asserts that two objects are equal.
//
//    assert.Equal(t, 123, 123, "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func Equal(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualObjects(expected, actual) {
		diff := diff(expected, actual)
		expected, actual = formatUnequalValues(expected, actual)

		return Fail(t,
			fmt.Sprintf("Not equal: \n"+
				"expected: %s\n"+
				"received: %s%s", expected, actual, diff),
			formatAndArgs...)
	}

	return true
}

// EqualValues asserts that two objects are equal or convertable to the same types
// and equal.
//
//    assert.EqualValues(t, uint32(123), int32(123), "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func EqualValues(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualValues(expected, actual) {
		diff := diff(expected, actual)
		expected, actual = formatUnequalValues(expected, actual)

		return Fail(t,
			fmt.Sprintf("Not equal: \n"+
				"expected: %s\n"+
				"received: %s%s", expected, actual, diff),
			formatAndArgs...)
	}

	return true
}

// formatUnequalValues takes two values of arbitrary types and returns string
// representations appropriate to be presented to the user.
//
// If the values are not of like type, the returned strings will be prefixed
// with the type name, and the value will be enclosed in parenthesis similar
// to a type conversion in the Go grammar.
func formatUnequalValues(expected, actual interface{}) (e string, a string) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return fmt.Sprintf("%T(%#v)", expected, expected), fmt.Sprintf("%T(%#v)", actual, actual)
	}

	return fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual)
}

// Exactly asserts that two objects are equal is value and type.
//
//    assert.Exactly(t, int32(123), int64(123), "123 and 123 should NOT be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func Exactly(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	aType := reflect.TypeOf(expected)
	bType := reflect.TypeOf(actual)

	if aType != bType {
		return Fail(t,
			fmt.Sprintf("Types expected to match exactly\n\r\t%v != %v", aType, bType),
			formatAndArgs...)
	}

	return Equal(t, expected, actual, formatAndArgs...)
}

// Nil asserts that the specified v is nil.
//
//    assert.Nil(t, err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if isNil(v) {
		return true
	}

	return Fail(t, fmt.Sprintf("Expected nil, but got: %#v", v), formatAndArgs...)
}

// NotNil asserts that the specified object is not nil.
//
//    assert.NotNil(t, err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func NotNil(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if !isNil(v) {
		return true
	}

	return Fail(t, "Expected value not to be nil.", formatAndArgs...)
}

// isNil checks if a specified v is nil or not, without Failing.
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

// isEmpty gets whether the specified v is considered empty or not.
func isEmpty(v interface{}) bool {
	switch v {
	case nil, "", false:
		return true
	}

	for _, num := range numericZeros {
		if num == v {
			return true
		}
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		fallthrough
	case reflect.Slice, reflect.Chan:
		{
			return (rv.Len() == 0)
		}
	case reflect.Struct:
		switch v.(type) {
		case time.Time:
			return v.(time.Time).IsZero()
		}
	case reflect.Ptr:
		{
			if rv.IsNil() {
				return true
			}

			switch v.(type) {
			case *time.Time:
				return v.(*time.Time).IsZero()
			default:
				return false
			}
		}
	}

	return false
}

// Empty asserts that the specified v is empty.  I.e. nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  assert.Empty(t, obj)
//
// Returns whether the assertion was successful (true) or not (false).
func Empty(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	pass := isEmpty(v)
	if !pass {
		Fail(t,
			fmt.Sprintf("Should be empty, but was %v", v),
			formatAndArgs...)
	}

	return pass
}

// NotEmpty asserts that the specified v is NOT empty.  I.e. not nil, "", false, 0 or either
// a slice or a channel with len == 0.
//
//  if assert.NotEmpty(t, obj) {
//    assert.Equal(t, "two", obj[1])
//  }
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmpty(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	pass := !isEmpty(v)
	if !pass {
		Fail(t, fmt.Sprintf("Should NOT be empty, but was %v", v), formatAndArgs...)
	}

	return pass
}

// getLen try to get length of v.
// return (false, 0) if impossible.
func getLen(v interface{}) (n int, ok bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
		}
	}()

	return reflect.ValueOf(v).Len(), true
}

// Len asserts that the specified v has specific length.
// Len also fails if the v has a type that len() not accept.
//
//    assert.Len(t, mySlice, 3, "The size of slice is not 3")
//
// Returns whether the assertion was successful (true) or not (false).
func Len(t Testing, v interface{}, length int, formatAndArgs ...interface{}) bool {
	n, ok := getLen(v)
	if !ok {
		return Fail(t,
			fmt.Sprintf("\"%s\" could not be applied builtin len()", v),
			formatAndArgs...)
	}

	if n != length {
		return Fail(t,
			fmt.Sprintf("\"%s\" should have %d item(s), but has %d", v, length, n),
			formatAndArgs...)
	}

	return true
}

// True asserts that the specified value is true.
//
//    assert.True(t, myBool, "myBool should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func True(t Testing, value bool, formatAndArgs ...interface{}) bool {
	if value != true {
		return Fail(t, "Should be true", formatAndArgs...)
	}

	return true
}

// False asserts that the specified value is false.
//
//    assert.False(t, myBool, "myBool should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func False(t Testing, value bool, formatAndArgs ...interface{}) bool {
	if value != false {
		return Fail(t, "Should be false", formatAndArgs...)
	}

	return true
}

// NotEqual asserts that the specified values are NOT equal.
//
//    assert.NotEqual(t, obj1, obj2, "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
//
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
func NotEqual(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if AreEqualObjects(expected, actual) {
		return Fail(t,
			fmt.Sprintf("Should not be: %#v\n", actual),
			formatAndArgs...)
	}

	return true
}

// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
func includeElement(list, element interface{}) (ok, found bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			found = false
		}
	}()

	listValue := reflect.ValueOf(list)
	elementValue := reflect.ValueOf(element)

	if reflect.TypeOf(list).Kind() == reflect.String {
		return true, strings.Contains(listValue.String(), elementValue.String())
	}

	if reflect.TypeOf(list).Kind() == reflect.Map {
		mapKeys := listValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if AreEqualObjects(mapKeys[i].Interface(), element) {
				return true, true
			}
		}

		return true, false
	}

	for i := 0; i < listValue.Len(); i++ {
		if AreEqualObjects(listValue.Index(i).Interface(), element) {
			return true, true
		}
	}

	return true, false
}

// Contains asserts that the specified string, list(array, slice...) or map contains the
// specified substring or element.
//
//    assert.Contains(t, "Hello World", "World", "But 'Hello World' does contain 'World'")
//    assert.Contains(t, ["Hello", "World"], "World", "But ["Hello", "World"] does contain 'World'")
//    assert.Contains(t, {"Hello": "World"}, "Hello", "But {'Hello': 'World'} does contain 'Hello'")
//
// Returns whether the assertion was successful (true) or not (false).
func Contains(t Testing, list, v interface{}, formatAndArgs ...interface{}) bool {
	ok, found := includeElement(list, v)
	if !ok {
		return Fail(t,
			fmt.Sprintf("\"%s\" could not be applied builtin len()", list),
			formatAndArgs...)
	}

	if !found {
		return Fail(t,
			fmt.Sprintf("\"%s\" does not contain \"%s\"", list, v),
			formatAndArgs...)
	}

	return true
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//    assert.NotContains(t, "Hello World", "Earth", "But 'Hello World' does NOT contain 'Earth'")
//    assert.NotContains(t, ["Hello", "World"], "Earth", "But ['Hello', 'World'] does NOT contain 'Earth'")
//    assert.NotContains(t, {"Hello": "World"}, "Earth", "But {'Hello': 'World'} does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(t Testing, list, v interface{}, formatAndArgs ...interface{}) bool {
	ok, found := includeElement(list, v)
	if !ok {
		return Fail(t,
			fmt.Sprintf("\"%s\" could not be applied builtin len()", list),
			formatAndArgs...)
	}

	if found {
		return Fail(t,
			fmt.Sprintf("\"%s\" should not contain \"%s\"", list, v),
			formatAndArgs...)
	}

	return true
}

// Condition uses a Comparison to assert a complex condition.
func Condition(t Testing, comp Comparison, formatAndArgs ...interface{}) bool {
	result := comp()
	if !result {
		Fail(t, "Condition test failed!", formatAndArgs...)
	}

	return result
}

// PanicTestFunc defines a func that should be passed to the assert.Panics and assert.NotPanics
// methods, and represents a simple func that takes no arguments, and returns nothing.
type PanicTestFunc func()

// didPanic returns true if the function passed to it panics. Otherwise, it returns false.
func didPanic(f PanicTestFunc) (bool, interface{}) {
	didPanic := false

	var message interface{}
	func() {
		defer func() {
			if message = recover(); message != nil {
				didPanic = true
			}
		}()

		// call the target function
		f()
	}()

	return didPanic, message
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//   assert.Panics(t, func(){
//     GoCrazy()
//   }, "Calling GoCrazy() should panic")
//
// Returns whether the assertion was successful (true) or not (false).
func Panics(t Testing, f PanicTestFunc, formatAndArgs ...interface{}) bool {
	if funcDidPanic, panicValue := didPanic(f); !funcDidPanic {
		return Fail(t,
			fmt.Sprintf("Func %T should panic\n\r\tPanic value:\t%v", f, panicValue),
			formatAndArgs...)
	}

	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//   assert.NotPanics(t, func(){
//     RemainCalm()
//   }, "Calling RemainCalm() should NOT panic")
//
// Returns whether the assertion was successful (true) or not (false).
func NotPanics(t Testing, f PanicTestFunc, formatAndArgs ...interface{}) bool {
	if funcDidPanic, panicValue := didPanic(f); funcDidPanic {
		return Fail(t,
			fmt.Sprintf("Func %T should not panic\n\r\tPanic value:\t%v", f, panicValue),
			formatAndArgs...)
	}

	return true
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//   assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second, "The difference should not be more than 10s")
//
// Returns whether the assertion was successful (true) or not (false).
func WithinDuration(t Testing, expected, actual time.Time, delta time.Duration, formatAndArgs ...interface{}) bool {
	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		return Fail(t,
			fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt),
			formatAndArgs...)
	}

	return true
}

func toFloat(x interface{}) (float64, bool) {
	var xf float64
	xok := true

	switch xn := x.(type) {
	case uint8:
		xf = float64(xn)
	case uint16:
		xf = float64(xn)
	case uint32:
		xf = float64(xn)
	case uint64:
		xf = float64(xn)
	case int:
		xf = float64(xn)
	case int8:
		xf = float64(xn)
	case int16:
		xf = float64(xn)
	case int32:
		xf = float64(xn)
	case int64:
		xf = float64(xn)
	case float32:
		xf = float64(xn)
	case float64:
		xf = float64(xn)
	default:
		xok = false
	}

	return xf, xok
}

// InDelta asserts that the two numerals are within delta of each other.
//
// 	 assert.InDelta(t, math.Pi, (22 / 7.0), 0.01)
//
// Returns whether the assertion was successful (true) or not (false).
func InDelta(t Testing, expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	af, aok := toFloat(expected)
	bf, bok := toFloat(actual)

	if !aok || !bok {
		return Fail(t, fmt.Sprintf("Parameters must be numerical"), formatAndArgs...)
	}

	if math.IsNaN(af) {
		return Fail(t, fmt.Sprintf("Actual must not be NaN"), formatAndArgs...)
	}

	if math.IsNaN(bf) {
		return Fail(t, fmt.Sprintf("Expected %v with delta %v, but was NaN", expected, delta), formatAndArgs...)
	}

	dt := af - bf
	if dt < -delta || dt > delta {
		return Fail(t, fmt.Sprintf("Max difference between %v and %v allowed is %v, but difference was %v", expected, actual, delta, dt), formatAndArgs...)
	}

	return true
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func InDeltaSlice(t Testing, expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Slice ||
		reflect.TypeOf(expected).Kind() != reflect.Slice {
		return Fail(t, fmt.Sprintf("Parameters must be slice"), formatAndArgs...)
	}

	actualSlice := reflect.ValueOf(actual)
	expectedSlice := reflect.ValueOf(expected)

	for i := 0; i < actualSlice.Len(); i++ {
		result := InDelta(t, actualSlice.Index(i).Interface(), expectedSlice.Index(i).Interface(), delta)
		if !result {
			return result
		}
	}

	return true
}

func calcRelativeError(expected, actual interface{}) (float64, error) {
	af, aok := toFloat(expected)
	if !aok {
		return 0, fmt.Errorf("expected value %q cannot be converted to float", expected)
	}
	if af == 0 {
		return 0, fmt.Errorf("expected value must have a value other than zero to calculate the relative error")
	}
	bf, bok := toFloat(actual)
	if !bok {
		return 0, fmt.Errorf("expected value %q cannot be converted to float", actual)
	}

	return math.Abs(af-bf) / math.Abs(af), nil
}

// InEpsilon asserts that expected and actual have a relative error less than epsilon
//
// Returns whether the assertion was successful (true) or not (false).
func InEpsilon(t Testing, expected, actual interface{}, epsilon float64, formatAndArgs ...interface{}) bool {
	actualEpsilon, err := calcRelativeError(expected, actual)
	if err != nil {
		return Fail(t, err.Error(), formatAndArgs...)
	}
	if actualEpsilon > epsilon {
		return Fail(t, fmt.Sprintf("Relative error is too high: %#v (expected)\n"+
			"        < %#v (actual)", actualEpsilon, epsilon), formatAndArgs...)
	}

	return true
}

// InEpsilonSlice is the same as InEpsilon, except it compares each value from two slices.
func InEpsilonSlice(t Testing, expected, actual interface{}, epsilon float64, formatAndArgs ...interface{}) bool {
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Slice ||
		reflect.TypeOf(expected).Kind() != reflect.Slice {
		return Fail(t, fmt.Sprintf("Parameters must be slice"), formatAndArgs...)
	}

	actualSlice := reflect.ValueOf(actual)
	expectedSlice := reflect.ValueOf(expected)

	for i := 0; i < actualSlice.Len(); i++ {
		result := InEpsilon(t, actualSlice.Index(i).Interface(), expectedSlice.Index(i).Interface(), epsilon)
		if !result {
			return result
		}
	}

	return true
}

/*
	Errors
*/

// NoError asserts that a function returned no error (i.e. `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.NoError(t, err) {
//	   assert.Equal(t, actualObj, expectedObj)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func NoError(t Testing, err error, formatAndArgs ...interface{}) bool {
	if err != nil {
		return Fail(t,
			fmt.Sprintf("Received unexpected error:\n%+v", err),
			formatAndArgs...)
	}

	return true
}

// Error asserts that a function returned an error (i.e. not `nil`).
//
//   actualObj, err := SomeFunction()
//   if assert.Error(t, err, "An error was expected") {
//	   assert.Equal(t, err, expectedError)
//   }
//
// Returns whether the assertion was successful (true) or not (false).
func Error(t Testing, err error, formatAndArgs ...interface{}) bool {
	if err == nil {
		return Fail(t, "An error is expected but got nil.", formatAndArgs...)
	}

	return true
}

// EqualError asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualError(t, err,  expectedErrorString, "An error was expected")
//
// Returns whether the assertion was successful (true) or not (false).
func EqualError(t Testing, theError error, errString string, formatAndArgs ...interface{}) bool {
	if !Error(t, theError, formatAndArgs...) {
		return false
	}

	expected := errString
	actual := theError.Error()

	// don't need to use deep equals here, we know they are both strings
	if expected != actual {
		return Fail(t,
			fmt.Sprintf("Error message not equal:\n"+
				"expected: %q\n"+
				"received: %q", expected, actual),
			formatAndArgs...)
	}

	return true
}

// EqualErrors asserts that a function returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//   actualObj, err := SomeFunction()
//   assert.EqualErrors(t, err,  expectedError, "An error was expected")
//
// Returns whether the assertion was successful (true) or not (false).
func EqualErrors(t Testing, actualError, expectedError error, formatAndArgs ...interface{}) bool {
	if !Error(t, actualError, formatAndArgs...) || !Error(t, expectedError, formatAndArgs...) {
		return false
	}

	expected := expectedError.Error()
	actual := actualError.Error()

	// don't need to use deep equals here, we know they are both strings
	if expected != actual {
		return Fail(t,
			fmt.Sprintf("Error message not equal:\n"+
				"expected: %q\n"+
				"received: %q", expected, actual),
			formatAndArgs...)
	}

	return true
}

// regexpMatch return true if a specified regexp matches a string.
func regexpMatch(reg, str interface{}) bool {
	var r *regexp.Regexp
	if rr, ok := reg.(*regexp.Regexp); ok {
		r = rr
	} else {
		r = regexp.MustCompile(fmt.Sprint(reg))
	}

	return (r.FindStringIndex(fmt.Sprint(str)) != nil)
}

// Match asserts that a specified regexp matches a string.
//
//  assert.Match(t, regexp.MustCompile("start"), "it's starting")
//  assert.Match(t, "start...$", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func Match(t Testing, reg, str interface{}, formatAndArgs ...interface{}) bool {
	match := regexpMatch(reg, str)
	if !match {
		Fail(t,
			fmt.Sprintf("Expect \"%v\" to match \"%v\"", str, reg),
			formatAndArgs...)
	}

	return match
}

// NotMatch asserts that a specified regexp does not match a string.
//
//  assert.NotMatch(t, regexp.MustCompile("starts"), "it's starting")
//  assert.NotMatch(t, "^start", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func NotMatch(t Testing, reg, str interface{}, formatAndArgs ...interface{}) bool {
	match := regexpMatch(reg, str)
	if match {
		Fail(t,
			fmt.Sprintf("Expect \"%v\" to NOT match \"%v\"", str, reg),
			formatAndArgs...)
	}

	return !match
}

// Zero asserts that v is the zero value for its type and returns the truth.
func Zero(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v != nil && !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return Fail(t,
			fmt.Sprintf("Should be zero, but was %v", v),
			formatAndArgs...)
	}

	return true
}

// NotZero asserts that v is not the zero value for its type and returns the truth.
func NotZero(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v == nil || reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return Fail(t,
			fmt.Sprintf("Should NOT be zero, but was %v", v),
			formatAndArgs...)
	}

	return true
}

// ReaderContains asserts that the specified io.Reader contains the specified sub string or element.
//
//    assert.ReaderContains(t, http.Response.Body, "Earth", "But 'http.Response.Body' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func ReaderContains(t Testing, reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Fail(t,
			fmt.Sprintf("Error read from \"%T\" of \"%s\"", reader, err.Error()),
			formatAndArgs...)
	}

	// try to close reader if it's io.Closer and reset reader
	if ioc, ok := reader.(io.Closer); ok {
		ioc.Close()
	}
	reader = ioutil.NopCloser(bytes.NewReader(data))

	return Contains(t, string(data), contains, formatAndArgs...)
}

// ReaderNotContains asserts that the specified io.Reader does not contain the specified substring or element.
//
//    assert.ReaderNotContains(t, http.Response.Body, "Earth", "But 'http.Response.Body' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func ReaderNotContains(t Testing, reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Fail(t,
			fmt.Sprintf("Error read from \"%T\" of \"%s\"", reader, err.Error()),
			formatAndArgs...)
	}

	// try to close reader if it's io.Closer and reset reader
	if ioc, ok := reader.(io.Closer); ok {
		ioc.Close()
	}
	reader = ioutil.NopCloser(bytes.NewReader(data))

	return NotContains(t, string(data), contains, formatAndArgs...)
}

// JSONEqual asserts that two JSON strings are equivalent.
//
//  assert.JSONEqual(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
//
// Returns whether the assertion was successful (true) or not (false).
func JSONEqual(t Testing, expected, actual string, formatAndArgs ...interface{}) bool {
	var expectedJSONAsInterface, actualJSONAsInterface interface{}

	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		return Fail(t,
			fmt.Sprintf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error()),
			formatAndArgs...)
	}

	if err := json.Unmarshal([]byte(actual), &actualJSONAsInterface); err != nil {
		return Fail(t,
			fmt.Sprintf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", actual, err.Error()),
			formatAndArgs...)
	}

	return Equal(t, expectedJSONAsInterface, actualJSONAsInterface, formatAndArgs...)
}

func typeAndKind(v interface{}) (reflect.Type, reflect.Kind) {
	t := reflect.TypeOf(v)
	k := t.Kind()

	if k == reflect.Ptr {
		t = t.Elem()
		k = t.Kind()
	}

	return t, k
}

// diff returns a diff of both values as long as both are of the same type and
// are a struct, map, slice or array. Otherwise it returns an empty string.
func diff(expected interface{}, actual interface{}) string {
	if expected == nil || actual == nil {
		return ""
	}

	et, ek := typeAndKind(expected)
	at, _ := typeAndKind(actual)

	if et != at {
		return ""
	}

	if ek != reflect.Struct && ek != reflect.Map && ek != reflect.Slice && ek != reflect.Array {
		return ""
	}

	e := spewConfig.Sdump(expected)
	a := spewConfig.Sdump(actual)

	diff, _ := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(e),
		B:        difflib.SplitLines(a),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})

	return "\n\nDiff:\n" + diff
}
