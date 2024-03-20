package lab2

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeHandlerCompute(t *testing.T) {
	type Case struct {
		name   string
		input  *bytes.Buffer
		output *bytes.Buffer
	}

	cases := []Case{
		{"sum 2", bytes.NewBufferString("+ 2 2"), bytes.NewBufferString("")},
		{"sum 3", bytes.NewBufferString("+ + 19.95 2 2"), bytes.NewBufferString("")},
		{"difference 2", bytes.NewBufferString("- 1.2 3"), bytes.NewBufferString("")},
		{"quotient 3", bytes.NewBufferString("/ / (-99.6) 37 10"), bytes.NewBufferString("")},
	}

	var ch ComputeHandler
	for _, c := range cases {
		ch.Input = c.input
		ch.Output = c.output

		t.Run(c.name, func(t *testing.T) {
			err := ch.Compute()

			if assert.Nil(t, err) {
				var outputb []byte
				outputb, err = io.ReadAll(c.output)

				if assert.Nil(t, err) {
					assert.NotEqual(t, string(outputb), "")
				}
			}
		})
	}

	type ErrorCase struct {
		name   string
		input  io.Reader
		output io.Writer
		err    error
	}

	errorCases := []ErrorCase{
		{"No input", nil, bytes.NewBufferString(""), fmt.Errorf("No Input in compute handler")},
		{"No output", bytes.NewBufferString("+ + 19.95 2 2"), nil, fmt.Errorf("No Output in compute handler")},
		{"No applicable unary operator", bytes.NewBufferString("+9"), bytes.NewBufferString(""), fmt.Errorf("No applicable unary operator %v", "+")},
	}

	for _, c := range errorCases {
		ch.Input = c.input
		ch.Output = c.output

		t.Run(c.name, func(t *testing.T) {
			err := ch.Compute()

			if assert.NotNil(t, err) {
				assert.Equal(t, c.err, err)
			}
		})
	}
}