## assert

> golang assertion lib copied from [testify](https://github.com/stretchr/testify)

Assertions allow you to easily write test code, and are global funcs in the `assert` package.
All assertion functions take, as the first argument, the `*testing.T` object provided by the
testing framework. This allows the assertion funcs to write the failings and other details to
the correct place.

Every assertion function also takes an optional string message as the final argument,
allowing custom error messages to be appended to the message the assertion method outputs.


### Usage

```go
import (
    "testing"

    "github.com/golib/assert"
)

func TestSomething(t *testing.T) {
    var a string = "Hello"
    var b string = "Hello"

    assert.Equal(t, a, b, "The two words should be the same.")
}

// if you assert many times, use the format below:
func TestSomething(t *testing.T) {
    assertion := assert.New(t)

    var a string = "Hello"
    var b string = "Hello"

    assertion.Equal(a, b, "The two words should be the same.")
}
```
