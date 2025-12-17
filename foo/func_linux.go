//go:build linuxx
// +build linuxx

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Linux implementation.")
}
