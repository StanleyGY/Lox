package main

import (
	"fmt"
	"os"
)

func main() {
	filename := os.Args[1]

	// Read program
	buf, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	code := string(buf)

	// Tokenize
	scanner := &ScannerImpl{}
	tokens, err := scanner.Scan(code)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Build ast tree
	parser := &RDParser{}
	stmts, err := parser.Parse(tokens)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Evaluate
	interpreter := MakeInterpreter()
	err = interpreter.Evaluate(stmts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
