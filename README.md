## assert

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/golib/assert/tree/master.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/golib/assert/tree/master)

> golang assert helpers modified from [testify](https://github.com/stretchr/testify)

Assertions allow you to easily write testing codes, and are global funcs in the `assert` package.
All assertion funcs take, as the first argument, the `*testing.T` object provided by the
testing framework. This allows the assertion funcs to write the failings and other details to
the correct place.

Every assertion func also takes an optional string message as the final argument,
allowing custom error messages to be appended to the message the assertion method outputs.

### Basic Usage

```go
import (
    "testing"

    "github.com/golib/assert"
)

func TestSomething(t *testing.T) {
    var (
        a string = "Hello"
        b string = "Hello"
    )

    assert.Equal(t, a, b, "The two words should be the same.")
}
```

### Advanced Usage
```go
import (
    "testing"
    "net/http"

    "github.com/golib/assert"
)

// if you assert many times, use the format below:
func TestSomething(t *testing.T) {
    it := assert.New(t)

    req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
    if it.Nil(err) {
        resp, err := http.DefaultClient.Do(req)
        if it.Nil(err) {
            it.Equal("HIT", resp.Header().Get("x-cache"))
        }
    }
}
```
