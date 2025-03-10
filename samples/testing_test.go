package testing

import (
	"math/rand/v2"
	"testing"

	"github.com/golib/assert"
)

func Test_Testing_Assert(t *testing.T) {
	assert.Empty(t, Hello())
}

func Test_Testing_AssertWithDiff(t *testing.T) {
	expect := &Testing{
		Name: "expect",
		Age:  rand.IntN(100),
		Addresses: []string{
			"test1 street",
			"test1@mail.com",
		},
	}

	actual := &Testing{}
	assert.Equal(t, actual, expect)
}

func Test_Testing_Assertion(t *testing.T) {
	t.Run("it should work", func(t *testing.T) {
		it := assert.New(t)

		it.Empty(Hello(), "it should be empty")
		it.Equal("expected interface{}", "actual interface{}")
	})

	t.Run("it should work with json", func(t *testing.T) {
		it := assert.New(t)

		jsonStr := `{"hello": "world", "foo": ["foo", "bar"]}`
		it.NotEmptyJSON(jsonStr, "world")
	})

	t.Run("it should work with struct", func(t *testing.T) {
		it := assert.New(t)

		actual := &Testing{}
		expect := "age"

		it.NotContains(actual, expect)
	})
}
