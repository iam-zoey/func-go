package main

import (
	"fmt"
	"io"
	fn "knative.dev/func-go/http"
	"net/http"
	"os"
	"os/exec"
)

// Main illustrates how scaffolding works to wrap a user's function.
func main() {
	// Instanced example (in scaffolding, 'New()' will be in module 'f')
	if err := fn.Start(New()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	// Static example (in scaffolding 'Handle' will be in module f
	// if err := fn.Start(fn.DefaultHandler{Handle}); err != nil {
	// 	fmt.Fprintln(os.Stderr, err.Error())
	// 	os.Exit(1)
	// }
}

// Example Static HTTP Handler implementation.
func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTTP handler invoked")
	fmt.Fprintln(w, "HTTP Handler invoked")
}

// MyFunction is an example instanced HTTP function implementation.
type MyFunction struct{}

func New() *MyFunction {
	fmt.Println("New function instance created")
	return &MyFunction{}
}

func (f *MyFunction) Handle(w http.ResponseWriter, r *http.Request) {
	runWasm(w, r)
}

func runWasm(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("RunWasm funciton is called")
	var input string

	switch r.Method {
	case http.MethodGet:
		// Read input from query parameter for GET requests
		input = r.URL.Query().Get("input")
	case http.MethodPost:
		// Read input from request body for POST requests
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()
		input = string(body)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if input == "" {
		http.Error(w, "Error: 'input' parameter is required", http.StatusBadRequest)
		return
	}
	// Run Wasmtime command with the input string as an argument.
	cmd := exec.Command("wasmtime", "main.wasm", input)

	// Execute the Wasmtime command and read the output.
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(w, "Error running Wasmtime: %v", err)
		return
	}

	// Return the output to the client.
	fmt.Fprintf(w, "==== V1: Output from wasm module:\n%s", output)

}
