package main

import "fmt"

type Error struct {
	errCode uint8
}

func (e *Error) Error() string {
	return "err"
}

func checkError1(err *Error) {
	if err != nil {
		panic("1")
	}
}

func checkError2(err error) {
	if err != nil {
		panic("2")
	}
}

func main() {
	var v *string                      // v 是 nil
	var i interface{}                  // 申明但未初始化，所以i是nil
	fmt.Println(i == nil)              // true
	i = v                              // 给i初始化赋值为nil，此时i的值为nil，但本身不为nil
	fmt.Println(i == nil)              // false
	fmt.Println(v == nil)              // true
	fmt.Println(interface{}(v) == nil) // false

	var e *Error
	checkError1(e) // not panic
	checkError2(e) // panic
}
