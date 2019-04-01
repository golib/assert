package assert

import "testing"

func Test_AreEqualObjects(t *testing.T) {
	// it should work
	testCases := []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, nil},
		{true, true},
		{false, false},
		{123, 123},
		{123.45, 123.45},
		{complex64(1), complex64(1)},
		{"Hello world", "Hello world"},
		{[]byte("Hello world"), []byte("Hello world")},
		{map[interface{}]interface{}{"foo": 123}, map[interface{}]interface{}{"foo": 123}},
	}
	for _, tc := range testCases {
		if !AreEqualObjects(tc.expected, tc.actual) {
			t.Errorf("Expect %#v is equal to %#v", tc.actual, tc.expected)
		}
	}

	// it should not work
	testCases = []struct {
		expected interface{}
		actual   interface{}
	}{
		{nil, ""},
		{"", nil},
		{nil, 0},
		{0, nil},
		{nil, false},
		{false, nil},
		{true, false},
		{false, true},
		{0, 0.123},
		{0.123, 0},
		{int32(123), int64(123)},
		{int64(123), int32(123)},
		{uint32(123), int32(123)},
		{int32(123), uint32(123)},
		{complex64(0), complex64(1)},
		{complex64(1), complex64(0)},
		{"Hello world", "hello world"},
		{"hello world", "Hello world"},
		{[]byte("Hello world"), []byte("hello world")},
		{[]byte("hello world"), []byte("Hello world")},
		{'x', "x"},
		{"x", 'x'},
	}
	for _, tc := range testCases {
		if AreEqualObjects(tc.expected, tc.actual) {
			t.Errorf("Expect %#v is not equal to %#v", tc.actual, tc.expected)
		}
	}
}

func Test_includeElement(t *testing.T) {

	list1 := []string{"Foo", "Bar"}
	list2 := []int{1, 2}
	simpleMap := map[interface{}]interface{}{"Foo": "Bar"}

	ok, found := includeElement("Hello World", "World")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Foo")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Bar")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list2, 1)
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list2, 2)
	True(t, ok)
	True(t, found)

	ok, found = includeElement(list1, "Foo!")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(list2, 3)
	True(t, ok)
	False(t, found)

	ok, found = includeElement(list2, "1")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(simpleMap, "Foo")
	True(t, ok)
	True(t, found)

	ok, found = includeElement(simpleMap, "Bar")
	True(t, ok)
	False(t, found)

	ok, found = includeElement(1433, "1")
	False(t, ok)
	False(t, found)
}

func Test_getLen(t *testing.T) {
	falseCases := []interface{}{
		nil,
		0,
		true,
		false,
		'A',
		struct{}{},
	}
	for _, v := range falseCases {
		n, ok := getLen(v)
		Equal(t, 0, n, "getLen should return 0 for %+v", v)
		False(t, ok, "Expected getLen fail to get length of %+v", v)
	}

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	trueCases := []struct {
		v interface{}
		n int
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

	for _, c := range trueCases {
		n, ok := getLen(c.v)
		Equal(t, c.n, n)
		True(t, ok, "Expected getLen success to get length of %+v", c.v)
	}
}

func Test_panicRecovery(t *testing.T) {
	if isRecovered, _ := panicRecovery(func() {
		panic("Panic!")
	}); !isRecovered {
		t.Error("panicRecovery should return true for paniced calling")
	}

	if isRecovered, _ := panicRecovery(func() {}); isRecovered {
		t.Error("panicRecovery should return false for non paniced calling")
	}
}
