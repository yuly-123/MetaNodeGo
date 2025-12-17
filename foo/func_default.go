//go:build !windowsx && !liunxx && !darwinx
// +build !windowsx,!liunxx,!darwinx

package foo

import "fmt"

func PlatformSpecificFunction() {
	fmt.Println("This is the Default implementation.")
}
