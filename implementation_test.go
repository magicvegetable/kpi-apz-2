package lab2

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrefixToInfix(t *testing.T) {
	type Case struct {
		name   string
		input  string
		result string
	}

	SimpleCases := []Case{
		{"empty", "", ""},
		{"sum 2", "+ 2 2", "2 + 2"},
		{"sum 3", "+ + 19.95 2 2", "19.95 + 2 + 2"},
		{"difference 2", "- 1.2 3", "1.2 - 3"},
		{"difference 3", "- - 189.2 3 33", "189.2 - 3 - 33"},
		{"unary minus", "-33", "-33"},
		{"product 2", "* 13 37", "13 * 37"},
		{"product 3", "* * 15 13 -37", "15 * 13 * (-37)"},
		{"quotient 2", "/ (-99.6) 37", "(-99.6) / 37"},
		{"quotient 3", "/ / (-99.6) 37 10", "(-99.6) / 37 / 10"},
	}

	for _, c := range SimpleCases {
		t.Run(c.name, func(t *testing.T) {
			res, err := PrefixToInfix(c.input)
			if assert.Nil(t, err) {
				assert.Equal(t, c.result, res)
			}
		})
	}

	ComplexCases := []Case{
		{"3 unary & sum 2 & difference 2", "/(-(-(+ 2 -2))) 1.0", "(-(-(2 + (-2)))) / 1.0"},
		{"sum 3 & difference 3 & product 2", "+ + (+ + 19.95 2 2) (- - 189.2 3 33) (* 13 37)", "(19.95 + 2 + 2) + (189.2 - 3 - 33) + (13 * 37)"},
		{"product 7", "* * * * * * (-1) 7987243.343 67 11 15 13 -37", "(-1) * 7987243.343 * 67 * 11 * 15 * 13 * (-37)"},
		{"quotient 9", "/ / / / / / / / (-99.6) 37 10 (-99.6) 37 10 (-99.6) 37 10", "(-99.6) / 37 / 10 / (-99.6) / 37 / 10 / (-99.6) / 37 / 10"},
	}

	for _, c := range ComplexCases {
		t.Run(c.name, func(t *testing.T) {
			res, err := PrefixToInfix(c.input)

			if assert.Nil(t, err) {
				assert.Equal(t, c.result, res)
			}
		})
	}

	type ErrorCase struct {
		name  string
		input string
		err   error
	}

	ErrorCases := []ErrorCase{
		{"Invalid operand", "abc", fmt.Errorf("Invalid operand %v", "abc")},
		{"No applicable unary operator", "+9", fmt.Errorf("No applicable unary operator %v", "+")},
		{"Lack of operators", "9 9", fmt.Errorf("Lack of operators %v", "1")},
		{"Not found matching opening parantheses", "9)", fmt.Errorf("Not found matching opening parantheses %v", "9)")},
	}

	for _, c := range ErrorCases {
		t.Run(c.name, func(t *testing.T) {
			_, err := PrefixToInfix(c.input)

			if assert.NotNil(t, err) {
				assert.Equal(t, c.err, err)
			}
		})
	}
}

func ExamplePrefixToInfix() {
	res, _ := PrefixToInfix("+ 2 2")
	fmt.Println(res)

	// Output:
	// 2 + 2
}
