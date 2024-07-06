package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	fn "knative.dev/func-go/http"
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

	// Define stdout and stderr buffers
	var stdout, stderr bytes.Buffer

	//cmd := exec.Command("wasmtime", "main.wasm")
	cmd := exec.Command("go", "run", "module.go") // Run the Go module instead of the WebAssembly module for testing

	// Set the input from HTTP request body as stdin for the command
	cmd.Stdin = bytes.NewReader([]byte(input))

	// Capture stdout and stderr from the command
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute module: %v\nStderr: %s", err, stderr.String()), http.StatusInternalServerError)
		return
	}

	// Write the output of the command execution to the response
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, stdout.String())
}
