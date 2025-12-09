package main

import (
	"MetaNodeGo/pkg3"
	"fmt"
)

const mainName string = "main"

var mainVar string = getMainVar()

func init() {
	fmt.Println("main init method invoked")
}

func main() {
	//fmt.Println("main method invoked!")

	fmt.Println("main : ", pkg3.Str2)
	fmt.Println("=====================================>")
	fmt.Println(pkg3.Show())
}

func method() (i int, s string) {
	return
}

func getMainVar() string {
	fmt.Println("main.getMainVar method invoked!")
	return mainName
}
