//go:build !windows && !liunx && !darwin
// +build !windows,!liunx,!darwin

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Default implementation.")
}
