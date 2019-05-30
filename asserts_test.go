package assert

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	i     interface{}
	zeros = []interface{}{
		false,
		byte(0),
		complex64(0),
		complex128(0),
		float32(0),
		float64(0),
		int(0),
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		rune(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),
		"",
		[0]interface{}{},
		[]interface{}(nil),
		struct{ x int }{},
		(*interface{})(nil),
		(func())(nil),
		nil,
		interface{}(nil),
		map[interface{}]interface{}(nil),
		(chan interface{})(nil),
		(<-chan interface{})(nil),
		(chan<- interface{})(nil),
	}
	nonzeros = []interface{}{
		true,
		byte(1),
		complex64(1),
		complex128(1),
		float32(1),
		float64(1),
		int(1),
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		rune(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintptr(1),
		"s",
		[1]interface{}{1},
		[]interface{}{},
		struct{ x int }{1},
		(*interface{})(&i),
		(func())(func() {}),
		interface{}(1),
		map[interface{}]interface{}{},
		(chan interface{})(make(chan interface{})),
		(<-chan interface{})(make(chan interface{})),
		(chan<- interface{})(make(chan interface{})),
	}
)

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct {
}

func (a *AssertionTesterConformingObject) TestMethod() {
}

// AssertionTesterUnconformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterUnconformingObject struct {
}

// bufferT implements Testing.
// Its implementation of Errorf writes the output that would be produced by
// testing.T.Errorf to an internal bytes.Buffer.
type bufferT struct {
	buf bytes.Buffer
}

func (t *bufferT) Errorf(format string, args ...interface{}) {
	// implementation of decorate is copied from testing.T
	decorate := func(s string) string {
		_, file, line, ok := runtime.Caller(3) // decorate + log + public function.
		if ok {
			// Truncate file name at last file name separator.
			if index := strings.LastIndex(file, "/"); index >= 0 {
				file = file[index+1:]
			} else if index = strings.LastIndex(file, "\\"); index >= 0 {
				file = file[index+1:]
			}
		} else {
			file = "???"
			line = 1
		}
		buf := new(bytes.Buffer)
		// Every line is indented at least one tab.
		buf.WriteByte('\t')
		fmt.Fprintf(buf, "%s:%d: ", file, line)
		lines := strings.Split(s, "\n")
		if l := len(lines); l > 1 && lines[l-1] == "" {
			lines = lines[:l-1]
		}
		for i, line := range lines {
			if i > 0 {
				// Second and subsequent lines are indented an extra tab.
				buf.WriteString("\n\t\t")
			}
			buf.WriteString(line)
		}
		buf.WriteByte('\n')
		return buf.String()
	}
	t.buf.WriteString(decorate(fmt.Sprintf(format, args...)))
}

func Test_Nil(t *testing.T) {
	mockT := new(testing.T)

	var (
		nilslice []interface{}
		nilmap   map[interface{}]interface{}
		nilchan  chan int
		nilptr   *string
	)

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, nil},
		{true, (*struct{})(nil)},
		{true, nilslice},
		{true, nilmap},
		{true, nilchan},
		{true, nilptr},
		{false, false},
		{false, 0},
		{false, ""},
		{false, "x"},
		{false, 'x'},
	}

	for _, tc := range testCases {
		if Nil(mockT, tc.actual) != tc.expected {
			t.Errorf("Nil of %#v should return %v", tc.actual, tc.expected)
		}
	}
}

func Test_NotNil(t *testing.T) {
	mockT := new(testing.T)

	var (
		nilslice []interface{}
		nilmap   map[interface{}]interface{}
		nilchan  chan int
		nilptr   *string
	)

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, nil},
		{true, (*struct{})(nil)},
		{true, nilslice},
		{true, nilmap},
		{true, nilchan},
		{true, nilptr},
		{false, false},
		{false, 0},
		{false, ""},
		{false, "x"},
		{false, 'x'},
	}

	for _, tc := range testCases {
		if NotNil(mockT, tc.actual) == tc.expected {
			t.Errorf("NotNil of %#v should return %v", tc.actual, !tc.expected)
		}
	}
}

func Test_Zero(t *testing.T) {
	mockT := new(testing.T)

	for _, v := range zeros {
		if !Zero(mockT, v) {
			t.Errorf("Expected %#v is the zero value of %v", v, reflect.TypeOf(v))
		}
	}

	for _, v := range nonzeros {
		if Zero(mockT, v) {
			t.Errorf("Expected %#v is not the zero value of %v", v, reflect.TypeOf(v))
		}
	}
}

func Test_NotZero(t *testing.T) {
	mockT := new(testing.T)

	for _, v := range zeros {
		if NotZero(mockT, v) {
			t.Errorf("Expected %#v is the zero value of %v", v, reflect.TypeOf(v))
		}
	}

	for _, v := range nonzeros {
		if !NotZero(mockT, v) {
			t.Errorf("Expected %#v is not the zero value of %v", v, reflect.TypeOf(v))
		}
	}
}

func Test_True(t *testing.T) {
	mockT := new(testing.T)

	var (
		nilslice []interface{}
		nilmap   map[interface{}]interface{}
		nilchan  chan int
		nilptr   *string
	)

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, true},
		{false, false},
		{false, (*struct{})(nil)},
		{false, nilslice},
		{false, nilmap},
		{false, nilchan},
		{false, nilptr},
		{false, 0},
		{false, ""},
		{false, "x"},
		{false, 'x'},
	}

	for _, tc := range testCases {
		if True(mockT, tc.actual) != tc.expected {
			t.Errorf("True of %#v should return %v", tc.actual, tc.expected)
		}
	}
}

func Test_False(t *testing.T) {
	mockT := new(testing.T)

	var (
		nilslice []interface{}
		nilmap   map[interface{}]interface{}
		nilchan  chan int
		nilptr   *string
	)

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, true},
		{false, false},
		{false, (*struct{})(nil)},
		{false, nilslice},
		{false, nilmap},
		{false, nilchan},
		{false, nilptr},
		{false, 0},
		{false, ""},
		{false, "x"},
		{false, 'x'},
	}

	for _, tc := range testCases {
		if False(mockT, tc.actual) == tc.expected {
			t.Errorf("False of %#v should return %v", tc.actual, !tc.expected)
		}
	}
}

func Test_IsType(t *testing.T) {
	mockT := new(testing.T)

	if !IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("IsType should return true: AssertionTesterConformingObject is the same type as AssertionTesterConformingObject")
	}
	if IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterUnconformingObject)) {
		t.Error("IsType should return false: AssertionTesterConformingObject is not the same type as AssertionTesterUnconformingObject")
	}
}

func Test_Implements(t *testing.T) {
	mockT := new(testing.T)

	if !Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}
	if Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterUnconformingObject)) {
		t.Error("Implements method should return false: AssertionTesterUnconformingObject does not implements AssertionTesterInterface")
	}
}

func Test_Equal(t *testing.T) {
	mockT := new(testing.T)

	// it should work
	testCases := []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, nil},
		{true, true},
		{false, false},
		{"", ""},
		{"Hello world", "Hello world"},
		{[]byte(""), []byte("")},
		{[]byte("Hello world"), []byte("Hello world")},
		{int(0), int(0)},
		{int8(0), int8(0)},
		{int16(0), int16(0)},
		{int32(0), int32(0)},
		{int64(0), int64(0)},
		{uint(0), uint(0)},
		{uint8(0), uint8(0)},
		{uint16(0), uint16(0)},
		{uint32(0), uint32(0)},
		{uint64(0), uint64(0)},
		{float32(0.0), float32(0.0)},
		{float64(0.0), float64(0.0)},
		{complex64(0), complex64(0)},
		{complex128(0), complex128(0)},
		{'x', 'x'},
		{struct{ a string }{}, struct{ a string }{}},
		{&struct{ a string }{}, &struct{ a string }{}},
		{map[interface{}]interface{}{}, map[interface{}]interface{}{}},
	}

	for _, tc := range testCases {
		if !Equal(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is equal to %#v", tc.actual, tc.expected)
		}
	}

	// it should not work
	testCases = []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, true},
		{nil, false},
		{nil, 0},
		{nil, ""},
		{nil, []byte("")},
		{true, false},
		{true, 0},
		{true, ""},
		{false, 0},
		{false, ""},
		{"", "Hello world"},
		{[]byte(""), []byte("Hello world")},
		{[]byte(""), ""},
		{int(0), int8(0)},
		{int(0), int16(0)},
		{int(0), int32(0)},
		{int(0), int64(0)},
		{uint(0), uint8(0)},
		{uint(0), uint16(0)},
		{uint(0), uint32(0)},
		{uint(0), uint64(0)},
		{float32(0.0), float64(0.0)},
		{complex64(0), complex128(0)},
		{'x', "x"},
		{struct{ a string }{}, struct{ b string }{}},
		{&struct{ a string }{}, &struct{ b string }{}},
		{map[interface{}]interface{}{}, map[interface{}]string{}},
	}

	for _, tc := range testCases {
		if Equal(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is NOT equal to %#v", tc.actual, tc.expected)
		}
	}
}

func Test_NotEqual(t *testing.T) {
	mockT := new(testing.T)

	// it should not work
	testCases := []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, nil},
		{true, true},
		{false, false},
		{"", ""},
		{"Hello world", "Hello world"},
		{[]byte(""), []byte("")},
		{[]byte("Hello world"), []byte("Hello world")},
		{int(0), int(0)},
		{int8(0), int8(0)},
		{int16(0), int16(0)},
		{int32(0), int32(0)},
		{int64(0), int64(0)},
		{uint(0), uint(0)},
		{uint8(0), uint8(0)},
		{uint16(0), uint16(0)},
		{uint32(0), uint32(0)},
		{uint64(0), uint64(0)},
		{float32(0.0), float32(0.0)},
		{float64(0.0), float64(0.0)},
		{complex64(0), complex64(0)},
		{complex128(0), complex128(0)},
		{'x', 'x'},
		{struct{ a string }{}, struct{ a string }{}},
		{&struct{ a string }{}, &struct{ a string }{}},
		{map[interface{}]interface{}{}, map[interface{}]interface{}{}},
	}

	for _, tc := range testCases {
		if NotEqual(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is equal to %#v", tc.actual, tc.expected)
		}
	}

	// it should work
	testCases = []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, true},
		{nil, false},
		{nil, 0},
		{nil, ""},
		{nil, []byte("")},
		{true, false},
		{true, 0},
		{true, ""},
		{false, 0},
		{false, ""},
		{"", "Hello world"},
		{[]byte(""), []byte("Hello world")},
		{[]byte(""), ""},
		{int(0), int8(0)},
		{int(0), int16(0)},
		{int(0), int32(0)},
		{int(0), int64(0)},
		{uint(0), uint8(0)},
		{uint(0), uint16(0)},
		{uint(0), uint32(0)},
		{uint(0), uint64(0)},
		{float32(0.0), float64(0.0)},
		{complex64(0), complex128(0)},
		{'x', "x"},
		{struct{ a string }{}, struct{ b string }{}},
		{&struct{ a string }{}, &struct{ b string }{}},
		{map[interface{}]interface{}{}, map[interface{}]string{}},
	}

	for _, tc := range testCases {
		if !NotEqual(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is NOT equal to %#v", tc.actual, tc.expected)
		}
	}
}

func Test_EqualFormatting(t *testing.T) {
	for i, currCase := range []struct {
		equalWant     string
		equalGot      string
		formatAndArgs []interface{}
		want          string
	}{
		{equalWant: "want", equalGot: "got", want: "\tasserts.go:167: \r                        \r\t\n\t\tError Trace:\tassert.Test_EqualFormatting:509\n\t\t\r\t\n\t\tError:      \tExpected values are NOT equal.\n\t\t\r\t             \t\n\t\t\r\t             \t--- Expected\n\t\t\r\t             \t+++ Actual\n\t\t\r\t             \t@@ -1 +1 @@\n\t\t\r\t             \t-want\n\t\t\r\t             \t+got\n\t\t\r\t             \t\n\t\t\r\t             \t\n\t\t\n"},
		{equalWant: "want", equalGot: "got", formatAndArgs: []interface{}{"hello, %v!", "world"}, want: "\tasserts.go:167: \r                        \r\t\n\t\tError Trace:\tassert.Test_EqualFormatting:509\n\t\t\r\t\n\t\tError:      \tExpected values are NOT equal.\n\t\t\r\t             \t\n\t\t\r\t             \t--- Expected\n\t\t\r\t             \t+++ Actual\n\t\t\r\t             \t@@ -1 +1 @@\n\t\t\r\t             \t-want\n\t\t\r\t             \t+got\n\t\t\r\t             \t\n\t\t\r\t             \t\n\t\t\r\tMessages:    \thello, world!\n\t\t\n"},
	} {
		mockT := &bufferT{}
		Equal(mockT, currCase.equalWant, currCase.equalGot, currCase.formatAndArgs...)
		if !strings.Contains(mockT.buf.String(), currCase.want) {
			t.Errorf("Equal(%d) output is: %#v", i, mockT.buf.String())
		}
	}
}

func Test_EqualValues(t *testing.T) {
	mockT := new(testing.T)

	// it should work
	testCases := []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, nil},
		{true, true},
		{false, false},
		{"", ""},
		{"Hello world", "Hello world"},
		{[]byte(""), []byte("")},
		{[]byte("Hello world"), []byte("Hello world")},
		{int(0), int(0)},
		{int8(0), int8(0)},
		{int16(0), int16(0)},
		{int32(0), int32(0)},
		{int64(0), int64(0)},
		{uint(0), uint(0)},
		{uint8(0), uint8(0)},
		{uint16(0), uint16(0)},
		{uint32(0), uint32(0)},
		{uint64(0), uint64(0)},
		{float32(0.0), float32(0.0)},
		{float64(0.0), float64(0.0)},
		{complex64(0), complex64(0)},
		{complex128(0), complex128(0)},
		{'x', 'x'},
		{struct{ a string }{}, struct{ a string }{}},
		{&struct{ a string }{}, &struct{ a string }{}},
		{map[interface{}]interface{}{}, map[interface{}]interface{}{}},
		{[]byte(""), ""},
		{'x', "x"},
		{int(0), int8(0)},
		{int(0), int16(0)},
		{int(0), int32(0)},
		{int(0), int64(0)},
		{uint(0), uint8(0)},
		{uint(0), uint16(0)},
		{uint(0), uint32(0)},
		{uint(0), uint64(0)},
		{float32(0.0), float64(0.0)},
		{complex64(0), complex128(0)},
	}

	for _, tc := range testCases {
		if !EqualValues(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is equal to %#v", tc.actual, tc.expected)
		}
	}

	// it should not work
	testCases = []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, true},
		{nil, false},
		{nil, 0},
		{nil, ""},
		{nil, []byte("")},
		{true, false},
		{true, 0},
		{true, ""},
		{false, 0},
		{false, ""},
		{"", "Hello world"},
		{[]byte(""), []byte("Hello world")},
		{struct{ a string }{}, struct{ b string }{}},
		{&struct{ a string }{}, &struct{ b string }{}},
		{map[interface{}]interface{}{}, map[interface{}]string{}},
	}

	for _, tc := range testCases {
		if EqualValues(mockT, tc.expected, tc.actual) {
			t.Errorf("Expected %#v is NOT equal to %#v", tc.actual, tc.expected)
		}
	}
}

func Test_Exactly(t *testing.T) {
	mockT := new(testing.T)

	// it should work for int

	ia := int32(0)
	ib := int64(0)
	ic := int32(0)
	id := int32(1)

	if Exactly(mockT, ia, ib) {
		t.Errorf("Exactly(%#v, %#v) should return false", ia, ib)
	}
	if Exactly(mockT, ia, id) {
		t.Errorf("Exactly(%#v, %#v) should return false", ia, id)
	}
	if !Exactly(mockT, ia, ic) {
		t.Errorf("Exactly(%#v, %#v) should return false", ia, ic)
	}

	// it should work for float
	fa := float32(0)
	fb := float64(0)
	fc := float32(0)
	fd := float32(1)

	if Exactly(mockT, fa, fb) {
		t.Errorf("Exactly(%#v, %#v) should return false", fa, fb)
	}
	if Exactly(mockT, fa, fd) {
		t.Errorf("Exactly(%#v, %#v) should return false", fa, fd)
	}
	if !Exactly(mockT, fa, fc) {
		t.Errorf("Exactly(%#v, %#v) should return true", fa, fc)
	}

	if Exactly(mockT, ia, fa) {
		t.Errorf("Exactly(%#v, %#v) should return false", ia, fa)
	}
	if Exactly(mockT, ia, nil) {
		t.Errorf("Exactly(%#v, %#v) should return false", ia, nil)
	}
	if Exactly(mockT, fa, nil) {
		t.Errorf("Exactly(%#v, %#v) should return false", fa, nil)
	}
	if Exactly(mockT, true, false) {
		t.Errorf("Exactly(%#v, %#v) should return false", true, false)
	}
	if Exactly(mockT, true, nil) {
		t.Errorf("Exactly(%#v, %#v) should return false", true, nil)
	}
	if Exactly(mockT, false, nil) {
		t.Errorf("Exactly(%#v, %#v) should return false", false, nil)
	}
}

func Test_Empty(t *testing.T) {
	mockT := new(testing.T)

	var (
		ts    time.Time
		tsptr *time.Time
		s     string
		sptr  *string
		f     os.File
		fptr  *os.File
	)

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, nil},
		{true, false},
		{true, 0},
		{true, ""},
		{true, []interface{}{}},
		{true, make(map[interface{}]interface{})},
		{true, make(chan struct{})},
		{true, ts},
		{true, tsptr},
		{true, s},
		{true, sptr},
		{true, f},
		{true, fptr},
		{false, true},
		{false, 1},
		{false, "0"},
		{false, []interface{}{0}},
		{false, map[interface{}]interface{}{0: 0}},
		{false, chWithValue},
	}

	for _, tc := range testCases {
		if Empty(mockT, tc.actual) != tc.expected {
			t.Errorf("Empty of %#v should return %v", tc.actual, tc.expected)
		}
	}
}

func Test_NotEmpty(t *testing.T) {
	mockT := new(testing.T)

	var (
		ts    time.Time
		tsptr *time.Time
		s     string
		sptr  *string
		f     os.File
		fptr  *os.File
	)

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	testCases := []struct {
		expected bool
		actual   interface{}
	}{
		{true, nil},
		{true, false},
		{true, 0},
		{true, ""},
		{true, []interface{}{}},
		{true, make(map[interface{}]interface{})},
		{true, make(chan struct{})},
		{true, ts},
		{true, tsptr},
		{true, s},
		{true, sptr},
		{true, f},
		{true, fptr},
		{false, true},
		{false, 1},
		{false, "0"},
		{false, []interface{}{0}},
		{false, map[interface{}]interface{}{0: 0}},
		{false, chWithValue},
	}

	for _, tc := range testCases {
		if NotEmpty(mockT, tc.actual) == tc.expected {
			t.Errorf("NotEmpty of %#v should return %v", tc.actual, tc.expected)
		}
	}
}

type astruct struct {
	Name, Value string
}

func Test_Contains(t *testing.T) {
	mockT := new(testing.T)

	var (
		list  = []string{"Foo", "Bar"}
		xlist = []*astruct{
			{"b", "c"},
			{"d", "e"},
			{"g", "h"},
			{"j", "k"},
		}
		amap = map[interface{}]interface{}{"Foo": "Bar"}
	)

	testCases := []struct {
		expected bool
		list     interface{}
		actual   interface{}
	}{
		{true, "Hello World", "Hello"},
		{true, "Hello World", "World"},
		{true, "Hello World", ""},
		{false, "Hello World", "Salt"},
		{true, list, "Foo"},
		{true, list, "Bar"},
		{false, list, ""},
		{false, list, "Salt"},
		{true, xlist, &astruct{"b", "c"}},
		{true, xlist, &astruct{"g", "h"}},
		{false, xlist, &astruct{"a", "b"}},
		{false, xlist, &astruct{}},
		{true, amap, "Foo"},
		{false, amap, ""},
		{false, amap, "Bar"},
	}

	for _, tc := range testCases {
		if Contains(mockT, tc.list, tc.actual) != tc.expected {
			t.Errorf("%#v contains %#v should return %v", tc.list, tc.actual, tc.expected)
		}
	}
}

func Test_NotContains(t *testing.T) {
	mockT := new(testing.T)

	var (
		list  = []string{"Foo", "Bar"}
		xlist = []*astruct{
			{"b", "c"},
			{"d", "e"},
			{"g", "h"},
			{"j", "k"},
		}
		amap = map[interface{}]interface{}{"Foo": "Bar"}
	)

	testCases := []struct {
		expected bool
		list     interface{}
		actual   interface{}
	}{
		{true, "Hello World", "Hello"},
		{true, "Hello World", "World"},
		{true, "Hello World", ""},
		{false, "Hello World", "Salt"},
		{true, list, "Foo"},
		{true, list, "Bar"},
		{false, list, ""},
		{false, list, "Salt"},
		{true, xlist, &astruct{"b", "c"}},
		{true, xlist, &astruct{"g", "h"}},
		{false, xlist, &astruct{"a", "b"}},
		{false, xlist, &astruct{}},
		{true, amap, "Foo"},
		{false, amap, ""},
		{false, amap, "Bar"},
	}

	for _, tc := range testCases {
		if NotContains(mockT, tc.list, tc.actual) == tc.expected {
			t.Errorf("%#v contains %#v should return %v", tc.list, tc.actual, !tc.expected)
		}
	}
}

func Test_Match(t *testing.T) {
	mockT := new(testing.T)

	testCases := []struct {
		rx, str string
		ok      bool
	}{
		{"", "Hello, world", true},
		{"^start", "start of the line", true},
		{"end$", "in the end", true},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12.34", true},
		{"Hello, world", "", false},
		{"^asdfastart", "Not the start of the line", false},
		{"end$", "in the end.", false},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12a.34", false},
	}

	for _, tc := range testCases {
		if Match(mockT, tc.rx, tc.str) != tc.ok {
			t.Errorf("Expected string(%s) to match regexp(%s)", tc.str, tc.rx)
		}

		if Match(mockT, regexp.MustCompile(tc.rx), tc.str) != tc.ok {
			t.Errorf("Expected string(%s) to match regexp(%s)", tc.str, tc.rx)
		}
	}
}

func Test_NotMatch(t *testing.T) {
	mockT := new(testing.T)

	testCases := []struct {
		rx, str string
		ok      bool
	}{
		{"", "Hello, world", true},
		{"^start", "start of the line", true},
		{"end$", "in the end", true},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12.34", true},
		{"Hello, world", "", false},
		{"^asdfastart", "Not the start of the line", false},
		{"end$", "in the end.", false},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12a.34", false},
	}

	for _, tc := range testCases {
		if NotMatch(mockT, tc.rx, tc.str) == tc.ok {
			t.Errorf("Expected string(%s) NOT to match regexp(%s)", tc.str, tc.rx)
		}

		if NotMatch(mockT, regexp.MustCompile(tc.rx), tc.str) == tc.ok {
			t.Errorf("Expected string(%s) NOT to match regexp(%s)", tc.str, tc.rx)
		}
	}
}

func Test_Condition(t *testing.T) {
	mockT := new(testing.T)

	if !Condition(mockT, func() bool { return true }, "Truth") {
		t.Error("Condition should return true")
	}

	if Condition(mockT, func() bool { return false }, "Lie") {
		t.Error("Condition should return false")
	}
}

func Test_Len(t *testing.T) {
	mockT := new(testing.T)

	// for invalid types
	invalidCases := []struct {
		expected int
		actual   interface{}
		ok       bool
	}{
		{0, nil, false},
		{0, 0, false},
		{0, true, false},
		{0, false, false},
		{0, 'a', false},
		{0, struct{}{}, false},
		{0, func() {}, false},
	}

	for _, tc := range invalidCases {
		if Len(mockT, tc.actual, tc.expected) != tc.ok {
			t.Errorf("Len(`%#v`, %d) should return %v", tc.actual, tc.expected, tc.ok)
		}
	}

	// for valid types
	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3

	validCases := []struct {
		expected int
		actual   interface{}
		ok       bool
	}{
		{0, "", true},
		{1, "", false},
		{3, "ABC", true},
		{0, "ABC", false},
		{4, "ABC", false},
		{0, []int(nil), true},
		{1, []int(nil), false},
		{0, []interface{}{}, true},
		{1, []interface{}{}, false},
		{3, []interface{}{1, 2, 3}, true},
		{4, []interface{}{1, 2, 3}, false},
		{3, [...]interface{}{1, 2, 3}, true},
		{4, [...]interface{}{1, 2, 3}, false},
		{0, map[int]int(nil), true},
		{1, map[int]int(nil), false},
		{0, map[interface{}]interface{}{}, true},
		{1, map[interface{}]interface{}{}, false},
		{3, map[interface{}]interface{}{1: 2, 2: 4, 3: 6}, true},
		{4, map[interface{}]interface{}{1: 2, 2: 4, 3: 6}, false},
		{0, (chan int)(nil), true},
		{1, (chan int)(nil), false},
		{0, make(chan int), true},
		{1, make(chan int), false},
		{3, ch, true},
		{4, ch, false},
	}

	for _, tc := range validCases {
		if Len(mockT, tc.actual, tc.expected) != tc.ok {
			t.Errorf("Expected `%#v` have %d item(s) return %v", tc.actual, tc.expected, tc.ok)
		}
	}
}

type customError struct{}

func (*customError) Error() string { return "fail" }

func Test_Error(t *testing.T) {
	mockT := new(testing.T)

	// start with a nil error
	var err error

	if Error(mockT, err) {
		t.Errorf("Error should return false for `%#v`", err)
	}

	// now set an error
	err = errors.New("some error")

	if !Error(mockT, err) {
		t.Errorf("Error should return true for `%#v`", err)
	}

	// returning an empty error interface
	var tmperr *customError

	if !Error(mockT, tmperr) {
		t.Errorf("Error should return true with empty error interface for `%#v`", err)
	}
}

func Test_NotError(t *testing.T) {
	mockT := new(testing.T)

	// start with a nil error
	var err error

	if !NotError(mockT, err) {
		t.Errorf("NotError should return true for `%#v`", err)
	}

	// now set an error
	err = errors.New("some error")

	if NotError(mockT, err) {
		t.Errorf("NotError should return false for `%#v`", err)
	}

	// returning an empty error interface
	var tmperr *customError

	if NotError(mockT, tmperr) {
		t.Errorf("NotError should return false with empty error interface for `%#v`", err)
	}
}

func Test_EqualErrors(t *testing.T) {
	mockT := new(testing.T)

	// start with a nil error
	var (
		err    error
		tmperr *customError
	)

	if !EqualErrors(mockT, tmperr, tmperr) {
		t.Errorf("EqualErrors should return true for the same %#v", tmperr)
	}

	if EqualErrors(mockT, err, tmperr) {
		t.Error("EqualErrors should return false for error of nil")
	}

	if EqualErrors(mockT, tmperr, err) {
		t.Error("EqualErrors should return false for error of nil")
	}

	// now set an error
	newerr := errors.New("some error")

	if !EqualErrors(mockT, newerr, newerr) {
		t.Errorf("EqualErrors should return true for error %#v", newerr)
	}

	if EqualErrors(mockT, tmperr, newerr) {
		t.Errorf("EqualErrors should return false for different %#v and %#v", err, tmperr)
	}
}

func Test_Panics(t *testing.T) {
	mockT := new(testing.T)

	if !Panics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("Panics should return true")
	}

	if Panics(mockT, func() {}) {
		t.Error("Panics should return false")
	}
}

func Test_NotPanics(t *testing.T) {
	mockT := new(testing.T)

	if !NotPanics(mockT, func() {}) {
		t.Error("NotPanics should return true")
	}

	if NotPanics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("NotPanics should return false")
	}
}

func TestWithinDuration(t *testing.T) {

	mockT := new(testing.T)
	a := time.Now()
	b := a.Add(10 * time.Second)

	True(t, WithinDuration(mockT, a, b, 10*time.Second), "A 10s difference is within a 10s time difference")
	True(t, WithinDuration(mockT, b, a, 10*time.Second), "A 10s difference is within a 10s time difference")

	False(t, WithinDuration(mockT, a, b, 9*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, 9*time.Second), "A 10s difference is not within a 9s time difference")

	False(t, WithinDuration(mockT, a, b, -9*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, -9*time.Second), "A 10s difference is not within a 9s time difference")

	False(t, WithinDuration(mockT, a, b, -11*time.Second), "A 10s difference is not within a 9s time difference")
	False(t, WithinDuration(mockT, b, a, -11*time.Second), "A 10s difference is not within a 9s time difference")
}

func TestInDelta(t *testing.T) {
	mockT := new(testing.T)

	True(t, InDelta(mockT, 1.001, 1, 0.01), "|1.001 - 1| <= 0.01")
	True(t, InDelta(mockT, 1, 1.001, 0.01), "|1 - 1.001| <= 0.01")
	True(t, InDelta(mockT, 1, 2, 1), "|1 - 2| <= 1")
	False(t, InDelta(mockT, 1, 2, 0.5), "Expected |1 - 2| <= 0.5 to fail")
	False(t, InDelta(mockT, 2, 1, 0.5), "Expected |2 - 1| <= 0.5 to fail")
	False(t, InDelta(mockT, "", nil, 1), "Expected non numerals to fail")
	False(t, InDelta(mockT, 42, math.NaN(), 0.01), "Expected NaN for actual to fail")
	False(t, InDelta(mockT, math.NaN(), 42, 0.01), "Expected NaN for expected to fail")

	cases := []struct {
		a, b  interface{}
		delta float64
	}{
		{uint8(2), uint8(1), 1},
		{uint16(2), uint16(1), 1},
		{uint32(2), uint32(1), 1},
		{uint64(2), uint64(1), 1},

		{int(2), int(1), 1},
		{int8(2), int8(1), 1},
		{int16(2), int16(1), 1},
		{int32(2), int32(1), 1},
		{int64(2), int64(1), 1},

		{float32(2), float32(1), 1},
		{float64(2), float64(1), 1},
	}

	for _, tc := range cases {
		True(t, InDelta(mockT, tc.a, tc.b, tc.delta), "Expected |%V - %V| <= %v", tc.a, tc.b, tc.delta)
	}
}

func TestInDeltaSlice(t *testing.T) {
	mockT := new(testing.T)

	True(t, InDeltaSlice(mockT,
		[]float64{1.001, 0.999},
		[]float64{1, 1},
		0.1), "{1.001, 0.009} is element-wise close to {1, 1} in delta=0.1")

	True(t, InDeltaSlice(mockT,
		[]float64{1, 2},
		[]float64{0, 3},
		1), "{1, 2} is element-wise close to {0, 3} in delta=1")

	False(t, InDeltaSlice(mockT,
		[]float64{1, 2},
		[]float64{0, 3},
		0.1), "{1, 2} is not element-wise close to {0, 3} in delta=0.1")

	False(t, InDeltaSlice(mockT, "", nil, 1), "Expected non numeral slices to fail")
}

func testAutogeneratedFunction() {
	defer func() {
		if err := recover(); err == nil {
			panic("did not panic")
		}
		StackTraces()
	}()
	t := struct {
		io.Closer
	}{}
	var c io.Closer
	c = t
	c.Close()
}

func TestStackTracesWithAutogeneratedFunctions(t *testing.T) {
	NotPanics(t, func() {
		testAutogeneratedFunction()
	})
}

func TestEqualJSON_EqualSONString(t *testing.T) {
	mockT := new(testing.T)
	True(t, EqualJSON(mockT, `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`))
}

func TestEqualJSON_EquivalentButNotEqual(t *testing.T) {
	mockT := new(testing.T)
	True(t, EqualJSON(mockT, `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`))
}

func TestEqualJSON_HashOfArraysAndHashes(t *testing.T) {
	mockT := new(testing.T)
	True(t, EqualJSON(mockT, "{\r\n\t\"numeric\": 1.5,\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]],\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\"\r\n}",
		"{\r\n\t\"numeric\": 1.5,\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\",\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]]\r\n}"))
}

func TestEqualJSON_Array(t *testing.T) {
	mockT := new(testing.T)
	True(t, EqualJSON(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`))
}

func TestEqualJSON_HashAndArrayNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`))
}

func TestEqualJSON_HashesNotEquivalent(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`))
}

func TestEqualJSON_ActualIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, `{"foo": "bar"}`, "Not JSON"))
}

func TestEqualJSON_ExpectedIsNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, "Not JSON", `{"foo": "bar", "hello": "world"}`))
}

func TestEqualJSON_ExpectedAndActualNotJSON(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, "Not JSON", "Not JSON"))
}

func TestEqualJSON_ArraysOfDifferentOrder(t *testing.T) {
	mockT := new(testing.T)
	False(t, EqualJSON(mockT, `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`))
}

func TestDiff(t *testing.T) {
	expected := `

--- Expected
+++ Actual
@@ -1 +1 @@
-struct { foo string }{foo:"hello"}
+struct { foo string }{foo:"bar"}


`
	actual := diffValues(
		struct{ foo string }{"hello"},
		struct{ foo string }{"bar"},
	)
	Equal(t, expected, actual)

	expected = `

--- Expected
+++ Actual
@@ -1 +1 @@
-[]int{1, 2, 3, 4}
+[]int{1, 3, 5, 7}


`
	actual = diffValues(
		[]int{1, 2, 3, 4},
		[]int{1, 3, 5, 7},
	)
	Equal(t, expected, actual)

	expected = `

--- Expected
+++ Actual
@@ -1 +1 @@
-[]int{1, 2, 3}
+[]int{1, 3, 5}


`
	actual = diffValues(
		[]int{1, 2, 3, 4}[0:3],
		[]int{1, 3, 5, 7}[0:3],
	)
	Equal(t, expected, actual)

	//	// NOTE: map is unsorted!
	//	expected = `
	//
	//--- Expected
	//+++ Actual
	//@@ -1 +1 @@
	//-map[string]int{"one":1, "two":2, "three":3, "four":4}
	//+map[string]int{"one":1, "three":3, "five":5, "seven":7}
	//
	//
	//`
	//
	//	actual = diffValues(
	//		map[string]int{"one": 1, "two": 2, "three": 3, "four": 4},
	//		map[string]int{"one": 1, "three": 3, "five": 5, "seven": 7},
	//	)
	//	Equal(t, expected, actual)
}

func TestDiffEmptyCases(t *testing.T) {
	Equal(t, "", diffValues(nil, nil))
	Equal(t, "", diffValues("", ""))
}

// Ensure there are no data races
func TestDiffRace(t *testing.T) {
	t.Parallel()

	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}

	actual := map[string]string{
		"d": "D",
		"e": "E",
		"f": "F",
	}

	// run diffs in parallel simulating tests with t.Parallel()
	numRoutines := 10
	rChans := make([]chan string, numRoutines)
	for idx := range rChans {
		rChans[idx] = make(chan string)
		go func(ch chan string) {
			defer close(ch)
			ch <- diffValues(expected, actual)
		}(rChans[idx])
	}

	for _, ch := range rChans {
		for msg := range ch {
			NotZero(t, msg) // dummy assert
		}
	}
}

type mockTesting struct {
}

func (m *mockTesting) Errorf(format string, args ...interface{}) {}

func TestFailNowWithPlainTesting(t *testing.T) {
	mockT := &mockTesting{}

	Panics(t, func() {
		FailNow(mockT, "failed")
	}, "should panic since mockT is missing FailNow()")
}

type mockFailNowTesting struct {
}

func (m *mockFailNowTesting) Errorf(format string, args ...interface{}) {}

func (m *mockFailNowTesting) FailNow() {}

func TestFailNowWithFullTesting(t *testing.T) {
	mockT := &mockFailNowTesting{}

	NotPanics(t, func() {
		FailNow(mockT, "failed")
	}, "should call mockT.FailNow() rather than panicking")
}

func TestReaderContains(t *testing.T) {

	mockT := new(testing.T)
	reader := strings.NewReader("Hello, World")

	if !ReaderContains(mockT, reader, "Hello") {
		t.Error("Contains should return true: \"Hello World\" contains \"Hello\"")
	}
	if ReaderContains(mockT, reader, "Salut") {
		t.Error("Contains should return false: \"Hello World\" does not contain \"Salut\"")
	}
}

func TestReaderNotContains(t *testing.T) {

	mockT := new(testing.T)
	reader := strings.NewReader("Hello, World")

	if ReaderNotContains(mockT, reader, "Hello") {
		t.Error("Contains should return true: \"Hello World\" contains \"Hello\"")
	}
	if !ReaderNotContains(mockT, reader, "Salut") {
		t.Error("Contains should return false: \"Hello World\" does not contain \"Salut\"")
	}
}

func TestContainsJSON(t *testing.T) {
	mockT := new(testing.T)

	jsonstr := `{"name":"testing","items":["one", 2],"status":true}`

	if !ContainsJSON(mockT, jsonstr, "name", "testing") {
		t.Error("ContainsJSON should return true")
	}

	if !ContainsJSON(mockT, jsonstr, "items.0", "one") {
		t.Error("ContainsJSON should return true")
	}

	if !ContainsJSON(mockT, jsonstr, "items.1", 2) {
		t.Error("ContainsJSON should return true")
	}

	if !ContainsJSON(mockT, jsonstr, "status", true) {
		t.Error("ContainsJSON should return true")
	}
}
