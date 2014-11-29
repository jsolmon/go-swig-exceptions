package main

import (
	"testing"

	"github.com/jsolmon/go-swig-exceptions/demolib"
)

var (
	demo = demolib.NewDemoLib()
)

func TestThrowsNegativeDoesNotThrow(t *testing.T) {
	n, err := demo.NegativeThrows(2)
	if err != nil {
		t.Error("Expected nil error, got ", err.Error())
	}
	if n != 2 {
		t.Error("Expected 2, got ", n)
	}
}

func TestThrowsNegativeThrows(t *testing.T) {
	expectedErr := "NegativeThrows threw exception"
	_, err := demo.NegativeThrows(-1)

	if err == nil {
		t.Fatal("Expected an error.")
	}
	if err.Error() != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, err.Error())
	}
}

func TestDivideByTwo(t *testing.T) {
	n, err := demo.DivideBy(1)

	if err != nil {
		t.Error("Expected nil error, got ", err.Error())
	}
	if n != 1.0 {
		t.Errorf("Expected 1.0 but got %v", n)
	}
}

func TestDivideByZero(t *testing.T) {
	expectedErr := "Cannot divide by zero"
	_, err := demo.DivideBy(0)

	if err == nil {
		t.Fatal("Expected an error when dividing by zero.")
	}
	if err.Error() != expectedErr {
		t.Errorf("Expected error message %v but got %v", expectedErr, err.Error())
	}
}

func TestNeverThrowsReturnsInput(t *testing.T) {
	n := demo.NeverThrows(-1)

	if n != -1 {
		t.Errorf("Expected -1 but got %v", n)
	}
}
