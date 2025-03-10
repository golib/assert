package testing

import "fmt"

type Testing struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	Addresses []string `json:"addresses"`
}

func Hello() string {
	return fmt.Sprintf("Hello, %s!", "world")
}
