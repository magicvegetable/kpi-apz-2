package main

import (
	"flag"
	"fmt"
	lab2 "github.com/maaagicvegetable/kpi-apz-2"
	"io"
	"os"
	"strings"
)

var (
	expressionInput = flag.String("e", "", "Expression to compute")
	expressionFile  = flag.String("f", "", "File containing the expression to compute")
	outputFile      = flag.String("o", "", "File to place output")
)

func main() {
	flag.Parse()

	var Input io.Reader
	Output := os.Stdout

	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			fmt.Fprintln(os.Stderr, "Fallback to stdout as output")
		} else {
			Output = file
		}
	}

	if *expressionFile != "" {
		file, err := os.Open(*expressionFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else {
			Input = file
		}
	}

	if *expressionInput != "" {
		if Input != nil {
			fmt.Fprintln(os.Stderr, "The expression is already set")
			fmt.Fprintln(os.Stderr, "Ignoring:", *expressionInput)
		} else {
			Input = strings.NewReader(*expressionInput)
		}
	}

	ch := lab2.ComputeHandler{Input, Output}
	err := ch.Compute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
