package pkg3

import "fmt"

var Str1 string = "asdf"
var Str2 string = Show()

func Show() string {
	Str1 += "zxcv,"
	fmt.Println("pkg3 : ", Str1)
	return Str1
}
