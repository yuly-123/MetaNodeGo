//go:build darwin
// +build darwin

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Darwin implementation.")
}
