package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// Read input line by line
		input := scanner.Text()
		fmt.Println("====== V2: Output from WASM module ======")
		fmt.Println("WASM says", input)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
