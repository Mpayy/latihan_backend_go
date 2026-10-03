package main

import (
	"fmt"
	"testing"
	"unicode"
)

// Soal 15 — Valid Palindrome (Two Pointer, pattern baru)
// Diberikan sebuah string, tentukan apakah string itu adalah palindrome — hanya mempertimbangkan karakter alfanumerik (huruf dan angka), dan mengabaikan perbedaan huruf besar/kecil.

// go
// func IsPalindrome(s string) bool

// Contoh:
// Input:  "A man, a plan, a canal: Panama"
// Output: true   // kalau dibersihkan jadi "amanaplanacanalpanama", itu palindrome

// Input:  "race a car"
// Output: false

func isAlnum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

func IsPalindrome(s string) bool {
	left := 0
	right := len(s) - 1

	for left < right {
		for left < right && !isAlnum(rune(s[left])) {
			left++
		}

		for left < right && !isAlnum(rune(s[right])) {
			right--
		}

		if unicode.ToLower(rune(s[left])) != unicode.ToLower(rune(s[right])) {
			return false
		}

		left++
		right--
	}

	return true
}

func TestIsPalindrome(t *testing.T) {
	input := "A man, a plan, a canal: Panama"

	result := IsPalindrome(input)
	fmt.Println(result)
}
