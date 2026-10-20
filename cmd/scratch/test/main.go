package main

import (
	"fmt"
	"runtime"
)

func main() {
	os := runtime.GOOS
	fmt.Println("Operating system:", os)
}
