package main

import (
	"fmt"
	"os"
)

// test program to print environment variables
func main() {
	fmt.Printf("foo %s\n", os.Getenv("MY_APP_foo"))
	fmt.Printf("hello %s\n", os.Getenv("MY_APP_hello"))
}
