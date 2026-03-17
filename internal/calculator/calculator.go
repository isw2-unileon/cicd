// Package calculator provides basic arithmetic operations.
// This is the core business logic of our application — kept intentionally
// simple so students can focus on the CI/CD workflows, not the code.
package calculator

import "errors"

// ErrDivisionByZero is returned when attempting to divide by zero.
var ErrDivisionByZero = errors.New("division by zero")

// ErrUnknownOperation is returned when an unsupported operation is requested.
var ErrUnknownOperation = errors.New("unknown operation")

// Add returns the sum of a and b.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract returns the difference of a and b.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply returns the product of a and b.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide returns a divided by b.
// Returns ErrDivisionByZero if b is zero.
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return a / b, nil
}

// Calculate dispatches to the correct operation based on the operator string.
// Supported operators: "add", "subtract", "multiply", "divide".
func Calculate(operator string, a, b float64) (float64, error) {
	switch operator {
	case "add":
		return Add(a, b), nil
	case "subtract":
		return Subtract(a, b), nil
	case "multiply":
		return Multiply(a, b), nil
	case "divide":
		return Divide(a, b)
	default:
		return 0, ErrUnknownOperation
	}
}
