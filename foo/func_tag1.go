//go:build windows
// +build windows

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Windows implementation.")
}
