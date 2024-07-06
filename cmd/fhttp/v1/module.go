package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: provide an input string as an argument")
		return
	}
	input := os.Args[1]
	fmt.Println("WASM says '", input, "'")
}
