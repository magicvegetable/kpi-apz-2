package lab2

import "fmt"
import "slices"
import "unicode"
import "strings"
import "regexp"

func CheckSingleOperand(operand string) bool {
	re := regexp.MustCompile("(?m)^[0-9]+(\\.[0-9])?[0-9]*$")
	return re.MatchString(operand)
}

func IsOperator(operator rune) bool {
	operators := []rune{'+', '-', '*', '/'}
	return slices.Contains(operators, operator)
}

func ParseEnclosingOperand(str string, offset int) (string, error) {
	openings := 0
	enclosings := 0

	for i := offset; i >= 0; i-- {
		ch := str[i]

		if ch == ')' {
			enclosings++
			continue
		}

		if ch == '(' {
			openings++
			if enclosings == openings {
				return str[i : offset+1], nil
			}
		}
	}

	return "", fmt.Errorf("Not found matching opening parantheses %v", str[:offset+1])
}

func ParseOperand(str string, offset int) (string, error) {
	if str[offset] == ')' {
		return ParseEnclosingOperand(str, offset)
	}

	for i := offset; i >= 0; i-- {
		ch := rune(str[i])

		if unicode.IsSpace(ch) || IsOperator(ch) {
			return str[i+1 : offset+1], nil
		}
	}

	operand := str[0 : offset+1]
	if !CheckSingleOperand(operand) {
		return operand, fmt.Errorf("Invalid operand %v", operand)
	}

	return operand, nil
}

const separator = " "

func OperandAppliedOperatorsFunc(operand string, fn func(map[rune]bool) bool) bool {
	operators := map[rune]bool{
		'+': false,
		'-': false,
		'*': false,
		'/': false,
	}

	for i := len(operand) - 1; i >= 0; i-- {
		ch := rune(operand[i])

		if unicode.IsSpace(ch) {
			continue
		}

		if IsOperator(ch) {
			operators[ch] = true

			if fn(operators) {
				return true
			}

			continue
		}

		suboperand, _ := ParseOperand(operand, i)

		offset := len(suboperand) - 1
		i -= offset
	}

	return false
}

func IsHaveLowerPriorityOperator(operators map[rune]bool) bool {
	return operators['+'] || operators['-']
}

func IsHaveAnyOperator(operators map[rune]bool) bool {
	return operators['+'] || operators['-'] || operators['*'] || operators['/']
}

func ApplyOperator(operator rune, operands []string) ([]string, error) {
	HighPriorityOperators := []rune{'*', '/'}

	LenOperands := len(operands)

	var operand []string
	if LenOperands >= 2 {
		LeftOperand := []string{operands[LenOperands-1]}
		RightOperand := []string{operands[LenOperands-2]}

		if slices.Contains(HighPriorityOperators, operator) {
			HaveLowerPriorityOperator := OperandAppliedOperatorsFunc(operands[LenOperands-2], IsHaveLowerPriorityOperator)
			if HaveLowerPriorityOperator {
				RightOperand = []string{"(", RightOperand[0], ")"}
			}

			HaveLowerPriorityOperator = OperandAppliedOperatorsFunc(operands[LenOperands-1], IsHaveLowerPriorityOperator)
			if HaveLowerPriorityOperator {
				LeftOperand = []string{"(", LeftOperand[0], ")"}
			}
		}

		if len(RightOperand) == 1 && RightOperand[0][0] == '-' {
			RightOperand = []string{"(", RightOperand[0], ")"}
		}

		operand = []string{
			strings.Join(LeftOperand, ""),
			separator,
			string(operator),
			separator,
			strings.Join(RightOperand, ""),
		}

		operands = operands[:LenOperands-2]
		operands = append(operands, strings.Join(operand, ""))

		return operands, nil
	}

	if LenOperands == 1 && operator == '-' {
		HaveAnyOperator := OperandAppliedOperatorsFunc(operands[LenOperands-1], IsHaveAnyOperator)

		if HaveAnyOperator {
			operand = []string{
				string(operator),
				"(",
				operands[LenOperands-1],
				")",
			}
		} else {
			operand = []string{
				string(operator),
				operands[LenOperands-1],
			}
		}

		operands = operands[:LenOperands-1]
		operands = append(operands, strings.Join(operand, ""))

		return operands, nil
	}

	if LenOperands == 0 {
		return operands, fmt.Errorf("No operands has given")
	}

	return operands, fmt.Errorf("No applicable unary operator %v", string(operator))
}

func PrefixToInfix(prefix string) (string, error) {
	var operands []string

	for i := len(prefix) - 1; i >= 0; i-- {
		ch := rune(prefix[i])

		if unicode.IsSpace(ch) {
			continue
		}

		var err error
		if IsOperator(ch) {
			operands, err = ApplyOperator(ch, operands)

			if err == nil {
				continue
			}

			return prefix, err
		}

		var operand string
		operand, err = ParseOperand(prefix, i)

		if err != nil {
			return prefix, err
		}

		offset := len(operand) - 1
		if operand[0] == '(' {
			operand, err = PrefixToInfix(operand[1:offset])

			if err != nil {
				return operand, err
			}

			operand = fmt.Sprintf("(%v)", operand)
		}

		operands = append(operands, operand)
		i -= offset
	}

	if len(operands) > 1 {
		return strings.Join(operands, ""), fmt.Errorf("Lack of operators %v", len(operands)-1)
	}

	return strings.Join(operands, ""), nil
}
