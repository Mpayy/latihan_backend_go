package main

import (
	"fmt"
	"testing"
)

// Soal 10 — Validate Balanced Parentheses
// Ini pattern klasik lain yang worth dikuasai: penggunaan stack.

// Diberikan sebuah string yang hanya berisi karakter (, ), {, }, [, ] — tentukan apakah string tersebut punya kurung yang valid/seimbang (setiap kurung buka punya pasangan tutup yang sesuai, dengan urutan yang benar).

// func IsValid(s string) bool

// Contoh:
// "()"      -> true
// "()[]{}"  -> true
// "(]"      -> false
// "([)]"    -> false
// "{[]}"    -> true

func IsValid(s string) bool {
	var stack []rune

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else {
			if len(stack) == 0 {
				return false
			}

			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if char == ')' && top != '(' ||
				(char == ']' && top != '[') ||
				(char == '}' && top != '{') {
				return false
			}
		}
	}

	return len(stack) == 0
}

func TestIsValid(t *testing.T) {
	s := "{[]}"
	isValid := IsValid(s)
	fmt.Println(isValid)
}
