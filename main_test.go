package main

import (
	"testing"
)

func TestPrecedence(t *testing.T) {
	if precedence('+') != 1 {
		t.Error("Expected precedence of '+' to be 1")
	}
	if precedence('-') != 1 {
		t.Error("Expected precedence of '-' to be 1")
	}
	if precedence('*') != 2 {
		t.Error("Expected precedence of '*' to be 2")
	}
	if precedence('/') != 2 {
		t.Error("Expected precedence of '/' to be 2")
	}
	if precedence('a') != 0 {
		t.Error("Expected precedence of 'a' to be 0")
	}
}

func TestApplyOperation(t *testing.T) {
	if applyOperation(1, 2, '+') != 3 {
		t.Error("Expected result of 1 + 2 to be 3")
	}
	if applyOperation(1, 2, '-') != -1 {
		t.Error("Expected result of 1 - 2 to be -1")
	}
	if applyOperation(1, 2, '*') != 2 {
		t.Error("Expected result of 1 * 2 to be 2")
	}
	if applyOperation(1, 2, '/') != 0 {
		t.Error("Expected result of 1 / 2 to be 0")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for division by zero")
		}
	}()
	applyOperation(1, 0, '/')
}
