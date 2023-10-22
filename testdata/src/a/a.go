package a

import (
	"fmt"
	fmt1 "fmt"
)

func f() error { // ok because f is private.
	return nil
}

func F() error { // want "should call"
	return nil
}

func G() error { // want "should call"
	defer f()
	return nil
}

func H() error { // ok because h begins by deferring a call to fmt.Println.
	defer fmt.Println("a")
	return nil
}

func J() error { // ok because j begins by deferring a call to fmt1.Println where fmt1 is alias of fmt.
	defer fmt1.Println()
	return nil
}

//lint:ignore dwrap reason
func K() error { // ok because k is ignored by dwrap.
	fmt.Println("F")
	return nil
}

// lint:ignore otherlinter reason
func L() error { // want "should call"
	fmt.Println("F")
	return nil
}

func M() (int, error) { // ok
	defer fmt.Println("a")
	return 0, nil
}

func N() (int, error) { // want "should call"
	fmt.Println("a")
	return 0, nil
}

func P() { // ok because f does not return error.
	fmt.Println()
}

func Println() {}

func Q() error { // want "should call"
	defer Println()
	return nil
}

func R() int { // ok because f does not return error.
	return 0
}

func S() error { // want "should call"
	defer fmt.Print()
	return nil
}

func T() // ok becasue T has no body.