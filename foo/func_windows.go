//go:build windowsx
// +build windowsx

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Windows implementation.")
}
