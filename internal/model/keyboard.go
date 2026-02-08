package model

import (
	"unicode"
)

// LayoutRows represents a Finnish/Swedish QWERTY letter-only layout.
var LayoutRows = []string{
	"qwertyuiopå",
	"asdfghjklöä",
	"zxcvbnm",
}

var leftHandSet = map[rune]struct{}{}
var rightHandSet = map[rune]struct{}{}

func init() {
	for _, r := range "qwertasdfgzxcvb" {
		leftHandSet[r] = struct{}{}
	}
	for _, r := range "yuiopåhjklöänm" {
		rightHandSet[r] = struct{}{}
	}
}

func IsLetterAllowed(r rune) bool {
	r = unicode.ToLower(r)
	for _, row := range LayoutRows {
		for _, c := range row {
			if r == c {
				return true
			}
		}
	}
	return false
}

func IsLeftHandLetter(r rune) bool {
	_, ok := leftHandSet[unicode.ToLower(r)]
	return ok
}

func IsRightHandLetter(r rune) bool {
	_, ok := rightHandSet[unicode.ToLower(r)]
	return ok
}

func BuildColumnMap(columnWidth int) map[rune]int {
	columns := make(map[rune]int)
	for rowIdx, row := range LayoutRows {
		for colIdx, ch := range row {
			offset := 2 + rowIdx // stagger rows slightly for readability
			columns[ch] = (colIdx+offset)*columnWidth + 2
		}
	}
	return columns
}
