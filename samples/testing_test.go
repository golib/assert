package testing

import (
	"testing"

	"github.com/golib/assert"
)

func Test_Testing_Assert(t *testing.T) {
	assert.Empty(t, Hello())
}

func Test_Testing_Assertion(t *testing.T) {
	it := assert.New(t)

	it.Empty(Hello(), "it should be empty")
	it.Equal("expected interface{}", "actual interface{}")
}
