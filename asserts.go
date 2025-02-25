package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/dolab/types"
	"github.com/kr/pretty"
)

// Nil asserts that the v is nil.
//
//	assert.Nil(t, err, "err should be nothing")
//
// Returns whether the assertion was successful (true) or not (false).
func Nil(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if isNil(v) {
		return true
	}

	return Fail(t, pretty.Sprintf("Expected to be nil, but got: %# v", v), formatAndArgs...)
}

// NotNil asserts that the v is not nil.
//
//	assert.NotNil(t, err, "err should be something")
//
// Returns whether the assertion was successful (true) or not (false).
func NotNil(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if !isNil(v) {
		return true
	}

	return Fail(t, "Expected NOT to be nil.", formatAndArgs...)
}

// Zero asserts that v is the zero value for its type and returns the truth.
func Zero(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v != nil && !reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return Fail(t,
			pretty.Sprintf("Should be zero, but got: %# v", v),
			formatAndArgs...)
	}

	return true
}

// NotZero asserts that v is not the zero value for its type and returns the truth.
func NotZero(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v == nil || reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface()) {
		return Fail(t,
			pretty.Sprintf("Should NOT be zero, but got: %# v", v),
			formatAndArgs...)
	}

	return true
}

// True asserts that the value is true.
//
//	assert.True(t, ok, "ok should be true")
//
// Returns whether the assertion was successful (true) or not (false).
func True(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	var tv bool
	switch v.(type) {
	case bool:
		tv = v.(bool)
	}

	if tv != true {
		return Fail(t, pretty.Sprintf("Expected %# v to be true", v), formatAndArgs...)
	}

	return true
}

// False asserts that the value is false.
//
//	assert.False(t, ko, "ko should be false")
//
// Returns whether the assertion was successful (true) or not (false).
func False(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	var fv bool
	switch v.(type) {
	case bool:
		fv = v.(bool)
	}

	if fv != false {
		return Fail(t, pretty.Sprintf("Expected %# v to be false", v), formatAndArgs...)
	}

	return true
}

// IsType asserts that the v is of the same type with expected type.
//
//	assert.IsType(t, int, 123)
//
// Returns whether the assertion was successful (true) or not (false).
func IsType(t Testing, expectedType, v interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualObjects(reflect.TypeOf(v), reflect.TypeOf(expectedType)) {
		return Fail(t,
			fmt.Sprintf(
				"Expect type of values are NOT the same.%s",
				diffValues(reflect.TypeOf(expectedType), reflect.TypeOf(v)),
			),
			formatAndArgs...)
	}

	return true
}

// Implements asserts that v implements the expected interface.
//
//	assert.Implements(t, (*Iface)(nil), new(v))
//
// Returns whether the assertion was successful (true) or not (false).
func Implements(t Testing, iface, v interface{}, formatAndArgs ...interface{}) bool {
	ifaceType := reflect.TypeOf(iface).Elem()

	if !reflect.TypeOf(v).Implements(ifaceType) {
		return Fail(t,
			fmt.Sprintf("Expect %T to implement %v", v, ifaceType),
			formatAndArgs...)
	}

	return true
}

// Equal asserts that two objects are equal.
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
//
//	assert.Equal(t, 123, 123)
//
// Returns whether the assertion was successful (true) or not (false).
func Equal(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualObjects(expected, actual) {
		return Fail(t,
			fmt.Sprintf(
				"Expected values are NOT equal.%s",
				diffValues(expected, actual),
			),
			formatAndArgs...)
	}

	return true
}

// NotEqual asserts that the values are NOT equal.
// Pointer variable equality is determined based on the equality of the
// referenced values (as opposed to the memory addresses).
//
//	assert.NotEqual(t, obj1, obj2, "two objects shouldn't be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func NotEqual(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if AreEqualObjects(expected, actual) {
		expected, actual = prettifyValues(expected, actual)

		return Fail(t, fmt.Sprintf(
			"Expected values are NOT equal in value.%s",
			diffValues(expected, actual),
		), formatAndArgs...)
	}

	return true
}

// EqualValues asserts that two objects are equal in value.
//
//	assert.EqualValues(t, uint32(123), int32(123), "123 and 123 should be equal")
//
// Returns whether the assertion was successful (true) or not (false).
func EqualValues(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if !AreEqualValues(expected, actual) {
		return Fail(t,
			fmt.Sprintf(
				"Expected values are NOT equal in value.%s",
				diffValues(expected, actual),
			),
			formatAndArgs...)
	}

	return true
}

// Exactly asserts that two objects are equal in both values and types.
//
//	assert.Exactly(t, int32(123), int64(123))
//
// Returns whether the assertion was successful (true) or not (false).
func Exactly(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	expectedType := reflect.TypeOf(expected)
	actualType := reflect.TypeOf(actual)

	if expectedType != actualType {
		return Fail(t,
			fmt.Sprintf(
				"Expected values are NOT equal in type.%s",
				diffValues(expectedType, actualType),
			),
			formatAndArgs...)
	}

	return Equal(t, expected, actual, formatAndArgs...)
}

// Empty asserts that the v is empty.  I.e. nil, "", false, 0 or either
// a list(slice, map, channel) with len == 0.
//
//	assert.Empty(t, v)
//
// Returns whether the assertion was successful (true) or not (false).
func Empty(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v == nil {
		return true
	}

	if !types.IsEmpty(v) {
		return Fail(t,
			pretty.Sprintf("Expected to be empty, but got: %# v", v),
			formatAndArgs...)
	}

	return true
}

// NotEmpty asserts that the v is NOT empty.  I.e. not nil, "", false, 0 or either
// a list(slice, map, channel) with len == 0.
//
//	if assert.NotEmpty(t, vs) {
//	  assert.Equal(t, "two", vs[0])
//	}
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmpty(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	if v == nil || types.IsEmpty(v) {
		return Fail(t,
			pretty.Sprintf("Expected not to be empty, but got: %# v", v),
			formatAndArgs...)
	}

	return true
}

// Contains asserts that the list(string, array, slice...) or map contains the
// specific sub string or element.
//
//	assert.Contains(t, "Hello World", "World", `"Hello World" does contain "World"`)
//	assert.Contains(t, []string{"Hello", "World"}, "World", `["Hello", "World"] does contain "World"`)
//	assert.Contains(t, map[string]string{"Hello": "World"}, "Hello", `{"Hello":"World"} does contain "Hello"`)
//
// Returns whether the assertion was successful (true) or not (false).
func Contains(t Testing, list, v interface{}, formatAndArgs ...interface{}) bool {
	ok, found := includeElement(list, v)
	if !ok {
		return Fail(t,
			pretty.Sprintf("Could not apply len() with %# v", v),
			formatAndArgs...)
	}

	if !found {
		return Fail(t,
			pretty.Sprintf("%# v does not contain `%v`", list, v),
			formatAndArgs...)
	}

	return true
}

// NotContains asserts that the specified string, list(array, slice...) or map does NOT contain the
// specified substring or element.
//
//	assert.NotContains(t, "Hello World", "Earth", `"Hello World" does NOT contain "Earth"`)
//	assert.NotContains(t, ["Hello", "World"], "Earth", `["Hello", "World"] does NOT contain "Earth"`)
//	assert.NotContains(t, {"Hello": "World"}, "Earth", `{"Hello": "World"} does NOT contain "Earth"`)
//
// Returns whether the assertion was successful (true) or not (false).
func NotContains(t Testing, list, v interface{}, formatAndArgs ...interface{}) bool {
	ok, found := includeElement(list, v)
	if !ok {
		return Fail(t,
			pretty.Sprintf("Could not apply len() with %# v", v),
			formatAndArgs...)
	}

	if found {
		return Fail(t,
			pretty.Sprintf("%# v contains `%v`", list, v),
			formatAndArgs...)
	}

	return true
}

// Match asserts that a specified regexp matches a string.
//
//	assert.Match(t, regexp.MustCompile("start"), "it's starting")
//	assert.Match(t, "start...$", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func Match(t Testing, reg, str interface{}, formatAndArgs ...interface{}) bool {
	if !tryMatch(reg, str) {
		return Fail(t,
			fmt.Sprintf("Expect string(%s) to match regexp(%s)", fmt.Sprint(str), fmt.Sprint(reg)),
			formatAndArgs...)
	}

	return true
}

// NotMatch asserts that a specified regexp does not match a string.
//
//	assert.NotMatch(t, regexp.MustCompile("starts"), "it's starting")
//	assert.NotMatch(t, "^starting", "it's not starting")
//
// Returns whether the assertion was successful (true) or not (false).
func NotMatch(t Testing, reg, str interface{}, formatAndArgs ...interface{}) bool {
	if tryMatch(reg, str) {
		return Fail(t,
			fmt.Sprintf("Expect string(%s) to NOT match regexp(%s)", fmt.Sprint(str), fmt.Sprint(reg)),
			formatAndArgs...)
	}

	return true
}

// Condition uses a Comparison to assert a complex condition.
//
//	assert.Condition(t, func()bool{return true;}, "It should return true")
//
// Returns whether the assertion was successful (true) or not (false).
func Condition(t Testing, comp Comparison, formatAndArgs ...interface{}) bool {
	if !comp() {
		return Fail(t, "Condition is failed!", formatAndArgs...)
	}

	return true
}

// Len asserts that the v has specific length.
// It fails if the v has a type that len() not accept.
//
//	assert.Len(t, aslice, 3, "The size of slice is not 3")
//
// Returns whether the assertion was successful (true) or not (false).
func Len(t Testing, v interface{}, length int, formatAndArgs ...interface{}) bool {
	n, ok := getLen(v)
	if !ok {
		return Fail(t,
			pretty.Sprintf("Could not apply len() with %# v", v),
			formatAndArgs...)
	}

	if n != length {
		return Fail(t,
			pretty.Sprintf("Expected %# v should have %d item(s), but got: %d item(s)", v, length, n),
			formatAndArgs...)
	}

	return true
}

// Error asserts that a func returned an error (i.e. not `nil`).
//
//	  v, err := SomeFunc()
//	  if assert.Error(t, err) {
//		   assert.EqualErrors(t, err, ErrNotFound)
//	  }
//
// Returns whether the assertion was successful (true) or not (false).
func Error(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	err, ok := v.(error)
	if !ok || err == nil {
		return Fail(t,
			pretty.Sprintf("Expected value is an error, but got: %# v", v),
			formatAndArgs...)
	}

	return true
}

// NotError asserts that a func returned no error (i.e. `nil`).
//
//	  v, err := SomeFunc()
//	  if assert.NotError(t, err) {
//		   assert.Equal(t, v, "OK")
//	  }
//
// Returns whether the assertion was successful (true) or not (false).
func NotError(t Testing, v interface{}, formatAndArgs ...interface{}) bool {
	err, ok := v.(error)
	if ok && err != nil {
		return Fail(t,
			pretty.Sprintf("Expected valus is NOT an error, but got: %# v", err),
			formatAndArgs...)
	}

	return true
}

// EqualErrors asserts that a func returned an error (i.e. not `nil`)
// and that it is equal to the provided error.
//
//	v, err := SomeFunc()
//	assert.EqualErrors(t, err,  ErrNotFound, "Error shoule be not found")
//
// Returns whether the assertion was successful (true) or not (false).
func EqualErrors(t Testing, expected, actual interface{}, formatAndArgs ...interface{}) bool {
	if !Error(t, expected, formatAndArgs...) {
		return false
	}
	if !Error(t, actual, formatAndArgs...) {
		return false
	}

	return Equal(t, expected.(error), actual.(error), formatAndArgs...)
}

// Panics asserts that the code inside the specified PanicTestFunc panics.
//
//	assert.Panics(t, func(){
//	  panic("Oops~")
//	}, "Calling should panic")
//
// Returns whether the assertion was successful (true) or not (false).
func Panics(t Testing, f PanicTestFunc, formatAndArgs ...interface{}) bool {
	if isRecovered, _ := panicRecovery(f); !isRecovered {
		return Fail(t,
			fmt.Sprintf("Expected Func(%T) should panic.", f),
			formatAndArgs...)
	}

	return true
}

// NotPanics asserts that the code inside the specified PanicTestFunc does NOT panic.
//
//	assert.NotPanics(t, func(){
//	  RemainCalm()
//	}, "Calling should NOT panic")
//
// Returns whether the assertion was successful (true) or not (false).
func NotPanics(t Testing, f PanicTestFunc, formatAndArgs ...interface{}) bool {
	if isRecovered, panicValue := panicRecovery(f); isRecovered {
		return Fail(t,
			fmt.Sprintf("Expected Func(%T) should not panic, but paniced with: %v", f, panicValue),
			formatAndArgs...)
	}

	return true
}

// WithinDuration asserts that the two times are within duration delta of each other.
//
//	assert.WithinDuration(t, time.Now(), time.Now(), 10*time.Second, "The difference should not be more than 10s")
//
// Returns whether the assertion was successful (true) or not (false).
func WithinDuration(t Testing, expected, actual time.Time, delta time.Duration, formatAndArgs ...interface{}) bool {
	dt := expected.Sub(actual)
	if dt < -delta || dt > delta {
		return Fail(t,
			fmt.Sprintf("Expected max difference between %v and %v allowed is %v, but got: %v", expected, actual, delta, dt),
			formatAndArgs...)
	}

	return true
}

// InDelta asserts that the two numerals are within delta of each other.
//
//	assert.InDelta(t, math.Pi, (22 / 7.0), 0.01)
//
// Returns whether the assertion was successful (true) or not (false).
func InDelta(t Testing, expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	af, aok := toFloat(expected)
	bf, bok := toFloat(actual)

	if !aok || !bok {
		return Fail(t, fmt.Sprintf("Parameters must be numerical"), formatAndArgs...)
	}

	if math.IsNaN(af) {
		return Fail(t,
			fmt.Sprintf("Actual must not be NaN"),
			formatAndArgs...)
	}

	if math.IsNaN(bf) {
		return Fail(t,
			fmt.Sprintf("Expected %v with delta %v, but got: NaN", expected, delta),
			formatAndArgs...)
	}

	dt := af - bf
	if dt < -delta || dt > delta {
		return Fail(t,
			fmt.Sprintf("Expected max difference between %v and %v allowed is %v, but got: %v", expected, actual, delta, dt),
			formatAndArgs...)
	}

	return true
}

// InDeltaSlice is the same as InDelta, except it compares two slices.
func InDeltaSlice(t Testing, expected, actual interface{}, delta float64, formatAndArgs ...interface{}) bool {
	if expected == nil || actual == nil ||
		reflect.TypeOf(actual).Kind() != reflect.Slice ||
		reflect.TypeOf(expected).Kind() != reflect.Slice {
		return Fail(t,
			fmt.Sprintf("Parameters must be slice"),
			formatAndArgs...)
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

// ReaderContains asserts that the specified io.Reader contains the specified sub string or element.
//
//	assert.ReaderContains(t, http.Response.Body, "Earth", "But 'http.Response.Body' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func ReaderContains(t Testing, reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	data, err := io.ReadAll(reader)
	if err != nil {
		return Fail(t,
			fmt.Sprintf("Error read from \"%T\" of \"%s\"", reader, err.Error()),
			formatAndArgs...)
	}

	// try to close reader if it's io.Closer and reset reader
	if ioc, ok := reader.(io.Closer); ok {
		ioc.Close()
	}

	reader = io.NopCloser(bytes.NewReader(data))

	return Contains(t, string(data), contains, formatAndArgs...)
}

// ReaderNotContains asserts that the specified io.Reader does not contain the specified substring or element.
//
//	assert.ReaderNotContains(t, http.Response.Body, "Earth", "But 'http.Response.Body' does NOT contain 'Earth'")
//
// Returns whether the assertion was successful (true) or not (false).
func ReaderNotContains(t Testing, reader io.Reader, contains interface{}, formatAndArgs ...interface{}) bool {
	data, err := io.ReadAll(reader)
	if err != nil {
		return Fail(t,
			fmt.Sprintf("Error read from \"%T\" of \"%s\"", reader, err.Error()),
			formatAndArgs...)
	}

	// try to close reader if it's io.Closer and reset reader
	if ioc, ok := reader.(io.Closer); ok {
		ioc.Close()
	}
	reader = io.NopCloser(bytes.NewReader(data))

	return NotContains(t, string(data), contains, formatAndArgs...)
}

// EqualJSON asserts that two JSON strings are equivalent.
//
//	assert.EqualJSON(t, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`)
//
// Returns whether the assertion was successful (true) or not (false).
func EqualJSON(t Testing, expected, actual string, formatAndArgs ...interface{}) bool {
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

// ContainsJSON asserts that the js string contains JSON value of the key.
//
//	assert.ContainsJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "hello", "world")
//	assert.ContainsJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "foo.1", "bar")
//
// Returns whether the assertion was successful (true) or not (false).
func ContainsJSON(t Testing, actual, key string, value interface{}) bool {
	data, err := getJsonValue(actual, key)
	if err != nil {
		t.Errorf("Expected contains actual key %s of value %s, but got: %+v", key, value, err)

		return false
	}

	keyValue := string(data)

	switch expected := value.(type) {
	case []byte:
		return EqualValues(t, expected, data,
			"Expected contains actual key %q of byte: %s, but got: %s", key, expected, data)

	case string:
		return EqualValues(t, expected, keyValue,
			"Expected contains actual key %q of string: %s, but got: %s", key, expected, data)

	case int8:
		actualValue, _ := strconv.Atoi(keyValue)

		return EqualValues(t, expected, int8(actualValue),
			"Expected contains actual key %q of int8: %v, but got: %s", key, expected, data)

	case int:
		actualValue, _ := strconv.Atoi(keyValue)

		return EqualValues(t, expected, int(actualValue),
			"Expected contains actual key %q of int: %v, but got: %s", key, expected, data)

	case int16:
		actualValue, _ := strconv.ParseInt(keyValue, 10, 16)

		return EqualValues(t, expected, int16(actualValue),
			"Expected contains actual key %q of int16: %v, but got %s", key, expected, data)

	case int32:
		actualValue, _ := strconv.ParseInt(keyValue, 10, 32)

		return EqualValues(t, expected, int32(actualValue),
			"Expected contains actual key %q of int32: %v, but got: %s", key, expected, data)

	case int64:
		actualValue, _ := strconv.ParseInt(keyValue, 10, 64)

		return EqualValues(t, expected, actualValue,
			"Expected contains actual key %q of int64: %v, but got: %s", key, expected, data)

	case float32:
		actualValue, _ := strconv.ParseFloat(keyValue, 32)

		return EqualValues(t, expected, float32(actualValue),
			"Expected contains actual key %q of float32: %v, but got: %v", key, expected, data)

	case float64:
		actualValue, _ := strconv.ParseFloat(keyValue, 64)

		return EqualValues(t, expected, actualValue,
			"Expected contains actual key %q of float64: %v, but got: %v", key, expected, data)

	case bool:
		switch strings.ToLower(keyValue) {
		case "true", "1", "on", "yes":
			return True(t, expected,
				"Expected contains actual key %q of [true|1|on], but got: %s", data)

		default:
			return False(t, expected,
				"Expected contains actual key %q of [false|0|off], but got: %s", data)
		}

	default:
		expectType := reflect.TypeOf(value)
		switch expectType.Kind() {
		case reflect.Ptr:
			if !isJsonEqualObject(keyValue, value) {
				t.Errorf("Expected contains actual key %s of value %s, but got: %s", key, value, keyValue)

				return false
			}

		case reflect.Array:
			fallthrough
		case reflect.Slice:
			// first, try with reflect
			actualValue := reflect.New(expectType)
			if err := json.Unmarshal(data, actualValue.Interface()); err == nil {
				expectedValue := reflect.ValueOf(expected)
				return EqualValues(t, expectedValue.Interface(), actualValue.Elem().Interface(),
					"Expected contains actual key %q of slice: %+v, but got: %+v", key, expectedValue.Interface(), actualValue.Elem().Interface())
			}

			// second, try with json string
			if !isJsonEqualObject(keyValue, value) {
				t.Errorf("Expected contains actual key %s of value %s, but got: %s", key, value, keyValue)

				return false
			}

		case reflect.Struct:
			// first, try with reflect
			actualValue := reflect.New(expectType)
			if err := json.Unmarshal(data, actualValue.Interface()); err == nil {
				expectedValue := reflect.ValueOf(expected)
				return EqualValues(t, expectedValue.Interface(), actualValue.Elem().Interface(),
					"Expected contains actual key %q of slice: %+v, but got: %+v", key, expectedValue.Interface(), actualValue.Elem().Interface())
			}

			// second, try with json string
			if !isJsonEqualObject(keyValue, value) {
				t.Errorf("Expected contains actual key %s of value %s, but got: %s", key, value, keyValue)

				return false
			}

		case reflect.Func:
			if !isJsonEqualObject(keyValue, value) {
				t.Errorf("Expected contains actual key %s of value %s, but got: %s", key, value, keyValue)

				return false
			}

		}
	}

	t.Errorf("Expected contains actual key %s of value %s, but got: %s", key, value, keyValue)

	return false
}

// NotContainsJSON asserts that the actual does not contain JSON key.
//
//	assert.NotContainsJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "world")
//	assert.NotContainsJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "foo.3")
//
// Returns whether the assertion was successful (true) or not (false).
func NotContainsJSON(t Testing, actual, key string) bool {
	if data, err := getJsonValue(actual, key); err == nil {
		t.Errorf("Expected does not contain json key %q, but got: %s", key, data)

		return false
	}

	return true
}

// NotEmptyJSON asserts that the actual contains JSON key, and the value is not empty.
//
//	assert.NotEmptyJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "world")
//	assert.NotEmptyJSON(t, `{"hello": "world", "foo": ["foo", "bar"]}`, "foo.3")
//
// Returns whether the assertion was successful (true) or not (false).
func NotEmptyJSON(t Testing, actual, key string) bool {
	data, err := getJsonValue(actual, key)
	if err != nil {
		t.Errorf("Expected contains json key %q, but got: %+v", key, err)

		return false
	}
	if len(data) == 0 {
		t.Errorf("Expected contains json key %q, but got: <empty>", key)

		return false
	}

	return true
}
