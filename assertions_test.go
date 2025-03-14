package assert

import (
	"errors"
	"regexp"
	"testing"
	"time"
)

func TestImplementsWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.Implements((*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}
	if it.Implements((*AssertionTesterInterface)(nil), new(AssertionTesterUnconformingObject)) {
		t.Error("Implements method should return false: AssertionTesterUnconformingObject does not implements AssertionTesterInterface")
	}
}

func TestIsTypeWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.IsType(new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("IsType should return true: AssertionTesterConformingObject is the same type as AssertionTesterConformingObject")
	}
	if it.IsType(new(AssertionTesterConformingObject), new(AssertionTesterUnconformingObject)) {
		t.Error("IsType should return false: AssertionTesterConformingObject is not the same type as AssertionTesterUnconformingObject")
	}

}

func TestEqualWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.Equal("Hello World", "Hello World") {
		t.Error("Equal should return true")
	}
	if !it.Equal(123, 123) {
		t.Error("Equal should return true")
	}
	if !it.Equal(123.5, 123.5) {
		t.Error("Equal should return true")
	}
	if !it.Equal([]byte("Hello World"), []byte("Hello World")) {
		t.Error("Equal should return true")
	}
	if !it.Equal(nil, nil) {
		t.Error("Equal should return true")
	}
}

func TestEqualValuesWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.EqualValues(uint32(10), int32(10)) {
		t.Error("EqualValues should return true")
	}
}

func TestNotNilWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.NotNil(new(AssertionTesterConformingObject)) {
		t.Error("NotNil should return true: object is not nil")
	}
	if it.NotNil(nil) {
		t.Error("NotNil should return false: object is nil")
	}

}

func TestNilWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.Nil(nil) {
		t.Error("Nil should return true: object is nil")
	}
	if it.Nil(new(AssertionTesterConformingObject)) {
		t.Error("Nil should return false: object is not nil")
	}

}

func TestTrueWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.True(true) {
		t.Error("True should return true")
	}
	if it.True(false) {
		t.Error("True should return false")
	}

}

func TestFalseWrapper(t *testing.T) {
	it := New(new(testing.T))

	if !it.False(false) {
		t.Error("False should return true")
	}
	if it.False(true) {
		t.Error("False should return false")
	}

}

func TestExactlyWrapper(t *testing.T) {
	it := New(new(testing.T))

	a := float32(1)
	b := float64(1)
	c := float32(1)
	d := float32(2)

	if it.Exactly(a, b) {
		t.Error("Exactly should return false")
	}
	if it.Exactly(a, d) {
		t.Error("Exactly should return false")
	}
	if !it.Exactly(a, c) {
		t.Error("Exactly should return true")
	}

	if it.Exactly(nil, a) {
		t.Error("Exactly should return false")
	}
	if it.Exactly(a, nil) {
		t.Error("Exactly should return false")
	}

}

func TestNotEqualWrapper(t *testing.T) {

	it := New(new(testing.T))

	if !it.NotEqual("Hello World", "Hello World!") {
		t.Error("NotEqual should return true")
	}
	if !it.NotEqual(123, 1234) {
		t.Error("NotEqual should return true")
	}
	if !it.NotEqual(123.5, 123.55) {
		t.Error("NotEqual should return true")
	}
	if !it.NotEqual([]byte("Hello World"), []byte("Hello World!")) {
		t.Error("NotEqual should return true")
	}
	if !it.NotEqual(nil, new(AssertionTesterConformingObject)) {
		t.Error("NotEqual should return true")
	}
}

func TestContainsWrapper(t *testing.T) {

	it := New(new(testing.T))
	list := []string{"Foo", "Bar"}

	if !it.Contains("Hello World", "Hello") {
		t.Error("Contains should return true: \"Hello World\" contains \"Hello\"")
	}
	if it.Contains("Hello World", "Salut") {
		t.Error("Contains should return false: \"Hello World\" does not contain \"Salut\"")
	}

	if !it.Contains(list, "Foo") {
		t.Error("Contains should return true: \"[\"Foo\", \"Bar\"]\" contains \"Foo\"")
	}
	if it.Contains(list, "Salut") {
		t.Error("Contains should return false: \"[\"Foo\", \"Bar\"]\" does not contain \"Salut\"")
	}

}

func TestNotContainsWrapper(t *testing.T) {

	it := New(new(testing.T))
	list := []string{"Foo", "Bar"}

	if !it.NotContains("Hello World", "Hello!") {
		t.Error("NotContains should return true: \"Hello World\" does not contain \"Hello!\"")
	}
	if it.NotContains("Hello World", "Hello") {
		t.Error("NotContains should return false: \"Hello World\" contains \"Hello\"")
	}

	if !it.NotContains(list, "Foo!") {
		t.Error("NotContains should return true: \"[\"Foo\", \"Bar\"]\" does not contain \"Foo!\"")
	}
	if it.NotContains(list, "Foo") {
		t.Error("NotContains should return false: \"[\"Foo\", \"Bar\"]\" contains \"Foo\"")
	}

}

func TestConditionWrapper(t *testing.T) {

	it := New(new(testing.T))

	if !it.Condition(func() bool { return true }, "Truth") {
		t.Error("Condition should return true")
	}

	if it.Condition(func() bool { return false }, "Lie") {
		t.Error("Condition should return false")
	}

}

func TestPanicsWrapper(t *testing.T) {

	it := New(new(testing.T))

	if !it.Panics(func() {
		panic("Panic!")
	}) {
		t.Error("Panics should return true")
	}

	if it.Panics(func() {
	}) {
		t.Error("Panics should return false")
	}

}

func TestNotPanicsWrapper(t *testing.T) {

	it := New(new(testing.T))

	if !it.NotPanics(func() {
	}) {
		t.Error("NotPanics should return true")
	}

	if it.NotPanics(func() {
		panic("Panic!")
	}) {
		t.Error("NotPanics should return false")
	}

}

func TestNotErrorWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error

	it.True(mockAssert.NotError(err), "NotError should return True for nil arg")

	// now set an error
	err = errors.New("Some error")

	it.False(mockAssert.NotError(err), "NotError with error should return False")

}

func TestErrorWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error

	it.False(mockAssert.IsError(err), "IsError should return False for nil arg")

	// now set an error
	err = errors.New("Some error")

	it.True(mockAssert.IsError(err), "IsError with error should return True")

}

func TestEqualErrorWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error
	it.False(mockAssert.EqualError(err, ""),
		"EqualError should return false for nil arg")

	// now set an error
	err = errors.New("some error")
	it.False(mockAssert.EqualError(err, "Not some error"),
		"EqualError should return false for different error string")
	it.True(mockAssert.EqualError(err, "some error"),
		"EqualError should return true")
}

func TestEqualErrorsWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	// start with a nil error
	var err error
	it.False(mockAssert.EqualErrors(err, nil),
		"EqualError should return false for nil arg")

	// now set an error
	err = errors.New("some error")
	it.False(mockAssert.EqualErrors(err, errors.New("Not some error")),
		"EqualError should return false for different error string")
	it.True(mockAssert.EqualErrors(err, errors.New("some error")),
		"EqualError should return true")
}

func TestEmptyWrappePr(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	it.True(mockAssert.Empty(""), "Empty string is empty")
	it.True(mockAssert.Empty(nil), "Nil is empty")
	it.True(mockAssert.Empty([]string{}), "Empty string array is empty")
	it.True(mockAssert.Empty(0), "Zero int value is empty")
	it.True(mockAssert.Empty(false), "False value is empty")

	it.False(mockAssert.Empty("something"), "Non Empty string is not empty")
	it.False(mockAssert.Empty(errors.New("something")), "Non nil object is not empty")
	it.False(mockAssert.Empty([]string{"something"}), "Non empty string array is not empty")
	it.False(mockAssert.Empty(1), "Non-zero int value is not empty")
	it.False(mockAssert.Empty(true), "True value is not empty")

}

func TestNotEmptyWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	it.False(mockAssert.NotEmpty(""), "Empty string is empty")
	it.False(mockAssert.NotEmpty(nil), "Nil is empty")
	it.False(mockAssert.NotEmpty([]string{}), "Empty string array is empty")
	it.False(mockAssert.NotEmpty(0), "Zero int value is empty")
	it.False(mockAssert.NotEmpty(false), "False value is empty")

	it.True(mockAssert.NotEmpty("something"), "Non Empty string is not empty")
	it.True(mockAssert.NotEmpty(errors.New("something")), "Non nil object is not empty")
	it.True(mockAssert.NotEmpty([]string{"something"}), "Non empty string array is not empty")
	it.True(mockAssert.NotEmpty(1), "Non-zero int value is not empty")
	it.True(mockAssert.NotEmpty(true), "True value is not empty")

}

func TestLenWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	it.False(mockAssert.Len(nil, 0), "nil does not have length")
	it.False(mockAssert.Len(0, 0), "int does not have length")
	it.False(mockAssert.Len(true, 0), "true does not have length")
	it.False(mockAssert.Len(false, 0), "false does not have length")
	it.False(mockAssert.Len('A', 0), "Rune does not have length")
	it.False(mockAssert.Len(struct{}{}, 0), "Struct does not have length")

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3

	cases := []struct {
		v interface{}
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range cases {
		it.True(mockAssert.Len(c.v, c.l), "%#v have %d items", c.v, c.l)
	}
}

func TestWithinDurationWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))
	a := time.Now()
	b := a.Add(10 * time.Second)

	it.True(mockAssert.WithinDuration(a, b, 10*time.Second), "A 10s difference is within a 10s time difference")
	it.True(mockAssert.WithinDuration(b, a, 10*time.Second), "A 10s difference is within a 10s time difference")

	it.False(mockAssert.WithinDuration(a, b, 9*time.Second), "A 10s difference is not within a 9s time difference")
	it.False(mockAssert.WithinDuration(b, a, 9*time.Second), "A 10s difference is not within a 9s time difference")

	it.False(mockAssert.WithinDuration(a, b, -9*time.Second), "A 10s difference is not within a 9s time difference")
	it.False(mockAssert.WithinDuration(b, a, -9*time.Second), "A 10s difference is not within a 9s time difference")

	it.False(mockAssert.WithinDuration(a, b, -11*time.Second), "A 10s difference is not within a 9s time difference")
	it.False(mockAssert.WithinDuration(b, a, -11*time.Second), "A 10s difference is not within a 9s time difference")
}

func TestInDeltaWrapper(t *testing.T) {
	it := New(new(testing.T))

	True(t, it.InDelta(1.001, 1, 0.01), "|1.001 - 1| <= 0.01")
	True(t, it.InDelta(1, 1.001, 0.01), "|1 - 1.001| <= 0.01")
	True(t, it.InDelta(1, 2, 1), "|1 - 2| <= 1")
	False(t, it.InDelta(1, 2, 0.5), "Expected |1 - 2| <= 0.5 to fail")
	False(t, it.InDelta(2, 1, 0.5), "Expected |2 - 1| <= 0.5 to fail")
	False(t, it.InDelta("", nil, 1), "Expected non numerals to fail")

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
		True(t, it.InDelta(tc.a, tc.b, tc.delta), "Expected |%V - %V| <= %v", tc.a, tc.b, tc.delta)
	}
}

func TestRegexpWrapper(t *testing.T) {

	it := New(new(testing.T))

	cases := []struct {
		rx, str string
	}{
		{"^start", "start of the line"},
		{"end$", "in the end"},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12.34"},
	}

	for _, tc := range cases {
		True(t, it.Match(tc.rx, tc.str))
		True(t, it.Match(regexp.MustCompile(tc.rx), tc.str))
		False(t, it.NotMatch(tc.rx, tc.str))
		False(t, it.NotMatch(regexp.MustCompile(tc.rx), tc.str))
	}

	cases = []struct {
		rx, str string
	}{
		{"^asdfastart", "Not the start of the line"},
		{"end$", "in the end."},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12a.34"},
	}

	for _, tc := range cases {
		False(t, it.Match(tc.rx, tc.str), "Expected \"%s\" to not match \"%s\"", tc.rx, tc.str)
		False(t, it.Match(regexp.MustCompile(tc.rx), tc.str))
		True(t, it.NotMatch(tc.rx, tc.str))
		True(t, it.NotMatch(regexp.MustCompile(tc.rx), tc.str))
	}
}

func TestZeroWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	for _, test := range zeros {
		it.True(mockAssert.Zero(test), "Zero should return true for %v", test)
	}

	for _, test := range nonzeros {
		it.False(mockAssert.Zero(test), "Zero should return false for %v", test)
	}
}

func TestNotZeroWrapper(t *testing.T) {
	it := New(t)
	mockAssert := New(new(testing.T))

	for _, test := range zeros {
		it.False(mockAssert.NotZero(test), "Zero should return true for %v", test)
	}

	for _, test := range nonzeros {
		it.True(mockAssert.NotZero(test), "Zero should return false for %v", test)
	}
}

func TestJSONEqWrapper_EqualSONString(t *testing.T) {
	it := New(new(testing.T))
	if !it.EqualJSON(`{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`) {
		t.Error("JSONEq should return true")
	}

}

func TestJSONEqWrapper_EquivalentButNotEqual(t *testing.T) {
	it := New(new(testing.T))
	if !it.EqualJSON(`{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`) {
		t.Error("JSONEq should return true")
	}

}

func TestJSONEqWrapper_HashOfArraysAndHashes(t *testing.T) {
	it := New(new(testing.T))
	if !it.EqualJSON("{\r\n\t\"numeric\": 1.5,\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]],\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\"\r\n}",
		"{\r\n\t\"numeric\": 1.5,\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\",\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]]\r\n}") {
		t.Error("JSONEq should return true")
	}
}

func TestJSONEqWrapper_Array(t *testing.T) {
	it := New(new(testing.T))
	if !it.EqualJSON(`["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`) {
		t.Error("JSONEq should return true")
	}

}

func TestJSONEqWrapper_HashAndArrayNotEquivalent(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON(`["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`) {
		t.Error("JSONEq should return false")
	}
}

func TestJSONEqWrapper_HashesNotEquivalent(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON(`{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`) {
		t.Error("JSONEq should return false")
	}
}

func TestJSONEqWrapper_ActualIsNotJSON(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON(`{"foo": "bar"}`, "Not JSON") {
		t.Error("JSONEq should return false")
	}
}

func TestJSONEqWrapper_ExpectedIsNotJSON(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON("Not JSON", `{"foo": "bar", "hello": "world"}`) {
		t.Error("JSONEq should return false")
	}
}

func TestJSONEqWrapper_ExpectedAndActualNotJSON(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON("Not JSON", "Not JSON") {
		t.Error("JSONEq should return false")
	}
}

func TestJSONEqWrapper_ArraysOfDifferentOrder(t *testing.T) {
	it := New(new(testing.T))
	if it.EqualJSON(`["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`) {
		t.Error("JSONEq should return false")
	}
}
