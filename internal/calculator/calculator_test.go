package calculator_test

import (
	"errors"
	"testing"

	"github.com/isw2-unileon/cicd/internal/calculator"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 3, 4, 7},
		{"negative numbers", -3, -4, -7},
		{"mixed signs", -3, 4, 1},
		{"zeros", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculator.Add(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Add(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 10, 3, 7},
		{"negative result", 3, 10, -7},
		{"zeros", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculator.Subtract(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Subtract(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		expected float64
	}{
		{"positive numbers", 3, 4, 12},
		{"by zero", 5, 0, 0},
		{"negative numbers", -3, -4, 12},
		{"mixed signs", -3, 4, -12},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculator.Multiply(tt.a, tt.b)
			if got != tt.expected {
				t.Errorf("Multiply(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	t.Run("valid division", func(t *testing.T) {
		got, err := calculator.Divide(10, 2)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != 5 {
			t.Errorf("Divide(10, 2) = %v, want 5", got)
		}
	})

	t.Run("division by zero", func(t *testing.T) {
		_, err := calculator.Divide(10, 0)
		if !errors.Is(err, calculator.ErrDivisionByZero) {
			t.Errorf("Divide(10, 0) error = %v, want ErrDivisionByZero", err)
		}
	})
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name      string
		operator  string
		a, b      float64
		expected  float64
		expectErr bool
	}{
		{"add", "add", 5, 3, 8, false},
		{"subtract", "subtract", 5, 3, 2, false},
		{"multiply", "multiply", 5, 3, 15, false},
		{"divide", "divide", 6, 3, 2, false},
		{"divide by zero", "divide", 6, 0, 0, true},
		{"unknown operator", "modulo", 6, 3, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := calculator.Calculate(tt.operator, tt.a, tt.b)
			if tt.expectErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("Calculate(%q, %v, %v) = %v, want %v", tt.operator, tt.a, tt.b, got, tt.expected)
			}
		})
	}
}
