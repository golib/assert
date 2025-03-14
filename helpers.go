package assert

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/buger/jsonparser"
	"github.com/dolab/colorize"
	"github.com/kr/pretty"
	"github.com/pmezard/go-difflib/difflib"
)

// AreEqualObjects determines if two objects are considered equal.
//
// NOTE: This func does no assertion of any kind.
func AreEqualObjects(expected, actual interface{}) bool {
	if expected == nil && actual == nil {
		return true
	}

	if expected == nil || actual == nil {
		return expected == actual
	}

	return reflect.DeepEqual(expected, actual)
}

// AreEqualValues gets whether two objects are equal, or if their
// values are equal.
func AreEqualValues(expected, actual interface{}) bool {
	if AreEqualObjects(expected, actual) {
		return true
	}

	actualValue := reflect.ValueOf(actual)
	expectedValue := reflect.ValueOf(expected)
	if !actualValue.IsValid() || !expectedValue.IsValid() {
		return false
	}

	for actualValue.Kind() == reflect.Ptr {
		actualValue = actualValue.Elem()
	}
	for expectedValue.Kind() == reflect.Ptr {
		expectedValue = expectedValue.Elem()
	}

	// Attempt comparison after type conversion
	if expectedValue.Type().ConvertibleTo(actualValue.Type()) {
		return reflect.DeepEqual(expectedValue.Convert(actualValue.Type()).Interface(), actual)
	}

	if actualValue.Type().ConvertibleTo(expectedValue.Type()) {
		return reflect.DeepEqual(actualValue.Convert(expectedValue.Type()).Interface(), actualValue)
	}

	return false
}

// StackTraces is necessary because the assert func use the testing object
// internally, causing it to print the file:line of the assert method,
// rather than where the problem actually occurred in calling code.
//
// StackTraces returns an array of strings containing the file and line number
// of each stack frame leading from the current test to the assert call that
// failed.
func StackTraces() []string {
	var (
		pc   uintptr
		file string
		name string
		line int
		ok   bool
	)

	var callers []string
	for i := 0; ; i++ {
		pc, file, line, ok = runtime.Caller(i)

		// The breaks below failed to terminate the loop, and we ran off the
		// end of the call stack.
		if !ok {
			break
		}

		// This is a huge edge case, but it will panic if this is the case.
		if file == "<autogenerated>" {
			break
		}

		f := runtime.FuncForPC(pc)
		if f == nil {
			break
		}

		name = f.Name()

		// testing.tRunner is the standard library function that calls
		// tests. Subtests are called directly by tRunner, without going through
		// the Test/Benchmark/Example function that contains the t.Run calls, so
		// with subtests we should break when we hit tRunner, without adding it
		// to the list of callers.
		if name == "testing.tRunner" {
			break
		}

		// ignore github.com/golib/assert itself
		if strings.HasPrefix(name, "github.com/golib/assert.") {
			continue
		}

		// ignore golang packages
		root, _ := os.Getwd()
		paths := strings.Split(strings.TrimPrefix(file, root), "/")
		if len(paths) < 2 {
			continue
		}

		if len(paths) > 2 {
			callers = append(callers, fmt.Sprintf("%s:%d", strings.Join(paths[len(paths)-2:], "/"), line))
		} else {
			callers = append(callers, fmt.Sprintf("%s:%d", paths[len(paths)-1], line))
		}
		// callers = append(callers, fmt.Sprintf("%s:%d", file, line))

		// Drop the package
		segments := strings.Split(name, ".")
		name = segments[len(segments)-1]
		if isTest(name, "Test") ||
			isTest(name, "Benchmark") ||
			isTest(name, "Example") {
			break
		}
	}

	if len(callers) > 1 {
		return callers[:1]
	}

	return callers
}

// FailNow fails test case and quit, or panic if Testing doesn't implement FailNow.
func FailNow(t Testing, message string, formatAndArgs ...interface{}) bool {
	Fail(t, message, formatAndArgs...)

	// We cannot extend Testing with FailNow() and
	// maintain backwards compatibility, so we fall back
	// to panicking when FailNow is not available in Testing.
	// See issue #263
	if nower, ok := t.(failNower); ok {
		nower.FailNow()
	} else {
		panic(fmt.Sprintf("test failed and %T does not implement `FailNow()`", t))
	}

	return false
}

// Fail reports a failure through
func Fail(t Testing, message string, formatAndArgs ...interface{}) bool {
	content := []labeledContent{
		{"Trace", strings.Join(StackTraces(), "\n\r\t\t\t")},
		{"Error", message},
	}

	if extras := formatExtraArgs(formatAndArgs...); len(extras) > 0 {
		content = append(content, labeledContent{"Messages", extras})
	}

	t.Errorf("\r" + getWhitespaceString() + labeledOutput(content...) + "\n")

	return false
}

func formatExtraArgs(formatAndArgs ...interface{}) string {
	if len(formatAndArgs) == 0 || formatAndArgs == nil {
		return ""
	}

	if len(formatAndArgs) == 1 {
		switch t := formatAndArgs[0].(type) {
		case string:
			return t
		case []byte:
			return string(t)
		default:
			return fmt.Sprintf("%v", t)
		}
	}

	if len(formatAndArgs) > 1 {
		switch t := formatAndArgs[0].(type) {
		case string:
			return fmt.Sprintf(t, formatAndArgs[1:]...)
		default:
			return fmt.Sprintf("%v", formatAndArgs)
		}
	}

	return ""
}

// prettifyValues takes two values of arbitrary types and returns string
// representations appropriate to be presented to the user.
//
// If the values are not of like type, the returned strings will be prefixed
// with the type name, and the value will be enclosed in parentheses similar
// to a type conversion in the Go grammar.
func prettifyValues(expected, actual interface{}) (es, as string) {
	if extype, ok := expected.(reflect.Type); ok {
		es = extype.Name()
	} else {
		es = pretty.Sprintf("%#v", expected)
	}

	if actype, ok := actual.(reflect.Type); ok {
		as = actype.Name()
	} else {
		as = pretty.Sprintf("%#v", actual)
	}

	return
}

type labeledContent struct {
	label   string
	content string
}

// labeledOutput returns a string consisting of the provided labeledContent.
// Each labeled output is appended in the following manner:
//
//	\r\t{{label}}:{{align_spaces}}\t{{content}}\n
//
// The initial carriage return is required to undo/erase any padding added by testing.T.Errorf. The "\t{{label}}:" is for the label.
// If a label is shorter than the longest label provided, padding spaces are added to make all the labels match in length. Once this
// alignment is achieved, "\t{{content}}\n" is added for the output.
//
// If the content of the labeledOutput contains line breaks, the subsequent lines are aligned so that they start at the same location as the first line.
func labeledOutput(content ...labeledContent) string {
	longestLabel := 0
	for _, v := range content {
		if len(v.label) > longestLabel {
			longestLabel = len(v.label)
		}
	}

	var output string
	for _, v := range content {
		output += fmt.Sprintf("\r\t%s:%s\t%s\n",
			v.label,
			strings.Repeat(" ", longestLabel-len(v.label)),
			paddingLines(v.content, longestLabel),
		)
	}

	return output
}

// getWhitespaceString returns a string that is long enough to overwrite the default
// output from the go testing framework.
func getWhitespaceString() string {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return ""
	}

	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return strings.Repeat(" ", len(fmt.Sprintf("%s:%d:    ", file, line)))
}

// Aligns the provided message so that all lines after the first line start at the same location as the first line.
// Assumes that the first line starts at the correct location (after carriage return, tab, label, spacer and tab).
// The longestLabelLen parameter specifies the length of the longest label in the output (required becaues this is the
// basis on which the alignment occurs).
func paddingLines(message string, longestLabelLen int) string {
	out := new(bytes.Buffer)

	for i, scanner := 0, bufio.NewScanner(strings.NewReader(message)); scanner.Scan(); i++ {
		// no need to align first line because it starts at the correct location (after the label)
		if i != 0 {
			// append alignLen+1 spaces to align with "{{longestLabel}}:" before adding tab
			out.WriteString("\n\r\t" + strings.Repeat(" ", longestLabelLen+1) + "\t")
		}

		out.WriteString(scanner.Text())
	}

	return out.String()
}

// Stolen from the `go test` tool.
// isTest tells whether name looks like a test (or benchmark, according to prefix).
// It is a Test (say) if there is a character after Test that is not a lower-case letter.
// We don't want TesticularCancer.
func isTest(name, prefix string) bool {
	if !strings.HasPrefix(name, prefix) {
		return false
	}

	if len(name) == len(prefix) { // "Test" is ok
		return true
	}

	r, _ := utf8.DecodeRuneInString(name[len(prefix):])

	return !unicode.IsLower(r)
}

// isNil checks if a specified v is nil or not, without Failing.
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}

	value := reflect.ValueOf(v)
	switch value.Kind() {
	case reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return value.IsNil()
	default:
		return false
	}
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

func getJsonValue(jsonStr, jsonKey string) ([]byte, error) {
	var (
		buf  = []byte(jsonStr)
		data []byte
		err  error
	)

	for {
		// first, try with the raw key
		data, _, _, err = jsonparser.Get(buf, jsonKey)
		if err == nil {
			buf = data
			break
		}

		// second, pop first key if dot existed
		parts := strings.SplitN(jsonKey, ".", 2)

		yek := parts[0]

		data, _, _, err = jsonparser.Get(buf, yek)
		if err == nil {
			buf = data
			if len(parts) != 2 {
				break
			}

			jsonKey = parts[1]

			continue
		}

		// is the yek an subscript?
		n, e := strconv.ParseInt(yek, 10, 32)
		if e != nil {
			break
		}

		var i int64
		_, err = jsonparser.ArrayEach(buf, func(arrBuf []byte, arrType jsonparser.ValueType, arrOffset int, arrErr error) {
			if i == n {
				data = arrBuf
				buf = data
				err = arrErr
			}

			i++
		})
		if err != nil {
			break
		}
		if len(parts) != 2 {
			break
		}

		jsonKey = parts[1]
	}
	if err != nil {
		return nil, err
	}

	return data, nil
}

func isJsonEqualObject(data string, obj interface{}) bool {
	// first, reform json string
	var value interface{}
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return false
	}

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false
	}

	// second, marshal obj with json
	objValue, err := json.Marshal(obj)
	if err != nil {
		return false
	}

	return bytes.Equal(jsonValue, objValue)
}

// containsElement try loop over the list check if the list includes the element.
// return (false, false) if impossible.
// return (true, false) if element was not found.
// return (true, true) if element was found.
func containsElement(actual, expect interface{}) (ok, found bool) {
	defer func() {
		if err := recover(); err != nil {
			ok = false
			found = false
		}
	}()

	actualValue := reflect.ValueOf(actual)
	for actualValue.Kind() == reflect.Ptr {
		actualValue = actualValue.Elem()
	}

	expectValue := reflect.ValueOf(expect)
	for expectValue.Kind() == reflect.Ptr {
		expectValue = expectValue.Elem()
	}

	switch actualValue.Kind() {
	case reflect.String:
		return true, strings.Contains(actualValue.String(), expectValue.String())

	case reflect.Map:
		mapKeys := actualValue.MapKeys()
		for i := 0; i < len(mapKeys); i++ {
			if AreEqualObjects(mapKeys[i].Interface(), expect) {
				return true, true
			}
		}

		return true, false

	case reflect.Struct:
		for i := 0; i < actualValue.NumField(); i++ {
			field := actualValue.Type().Field(i)
			if !field.IsExported() {
				continue
			}

			if tagName := field.Tag.Get("json"); tagName != "" {
				if AreEqualObjects(tagName, expect) {
					return true, true
				}
			} else {
				if AreEqualObjects(field.Name, expect) {
					return true, true
				}
			}
		}

		return true, false

	default:
		for i := 0; i < actualValue.Len(); i++ {
			if !actualValue.Index(i).CanInterface() {
				continue
			}

			if AreEqualObjects(actualValue.Index(i).Interface(), expect) {
				return true, true
			}
		}
	}

	return true, false
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

// diffValues returns a diff of both values as long as both are of the same type and
// are a struct, map, slice or array. Otherwise, it returns an empty string.
func diffValues(expected, actual interface{}) string {
	expectStr, actualStr := prettifyValues(expected, actual)

	diffs, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(expectStr),
		B:        difflib.SplitLines(actualStr),
		FromFile: "Expected",
		FromDate: "",
		ToFile:   "Actual",
		ToDate:   "",
		Context:  1,
	})

	if err != nil || len(diffs) == 0 {
		diffs := pretty.Diff(expected, actual)
		if len(diffs) == 0 {
			return ""
		}

		return fmt.Sprintf("\n\n%v\n", diffs)
	}

	return fmt.Sprintf("\n\n%s\n", diffColorize(diffs))
}

func diffColorize(diffs string) string {
	paint := colorize.New("yellow")

	lines := strings.Split(diffs, "\n")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "+"):
			paint.SetFgColor(colorize.ColorBlue)
		case strings.HasPrefix(line, "-"):
			paint.SetFgColor(colorize.ColorRed)
		default:
			paint.SetFgColor(colorize.ColorGray)
		}

		lines[i] = paint.Paint(line)
	}

	return strings.Join(lines, "\n")
}

// panicRecovery returns true if the function passed to it panics. Otherwise, it returns false.
func panicRecovery(f PanicTestFunc) (bool, interface{}) {
	isRecovered := false

	var message interface{}
	func() {
		defer func() {
			if message = recover(); message != nil {
				isRecovered = true
			}
		}()

		// call the target function
		f()
	}()

	return isRecovered, message
}

// tryMatch returns true if a specified regexp matches a string.
func tryMatch(reg, str interface{}) (ok bool) {
	defer func() {
		if perr := recover(); perr != nil {
			ok = false
		}
	}()

	var r *regexp.Regexp
	if tmpr, ok := reg.(*regexp.Regexp); ok {
		r = tmpr
	} else {
		r = regexp.MustCompile(fmt.Sprint(reg))
	}

	return len(r.FindStringIndex(fmt.Sprint(str))) > 0
}
