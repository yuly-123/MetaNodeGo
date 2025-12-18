//go:build linux
// +build linux

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Linux implementation.")
}
