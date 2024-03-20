package lab2

import "io"
import "fmt"

type ComputeHandler struct {
	Input  io.Reader
	Output io.Writer
}

func (ch *ComputeHandler) Compute() error {
	if ch.Input == nil {
		return fmt.Errorf("No Input in compute handler")
	}

	if ch.Output == nil {
		return fmt.Errorf("No Output in compute handler")
	}

	inputb, err := io.ReadAll(ch.Input)
	if err != nil {
		return err
	}

	var output string
	output, err = PrefixToInfix(string(inputb))

	if err != nil {
		return err
	}

	if _, err = io.WriteString(ch.Output, output); err != nil {
		return err
	}

	return nil
}