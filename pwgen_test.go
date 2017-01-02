package pwgen

import (
	"bytes"
	"strings"
	"testing"
)

var randomByteTests = []struct {
	in []byte
}{
	{[]byte{}},
	{[]byte{'a'}},
	{[]byte{'a', 'b', 'c'}},
}

func TestRandomByte(t *testing.T) {
	for _, tt := range randomByteTests {
		s := randomByte(tt.in)
		if len(tt.in) == 0 && s != ' ' {
			t.Errorf("got '%#v'\nwant '%#v'", s, ' ')
		}
		if len(tt.in) > 0 && bytes.IndexByte(tt.in, s) == -1 {
			t.Errorf("got '%#v'\nnot in '%#v'", s, tt.in)
		}
	}
}

var expandCodesetTests = []struct {
	in  int
	out []byte
}{
	{factorDigit, []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}},
	{factorAlphabetLarge, []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}},
	{factorAlphabetSmall, []byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}},
	{factorUnderscore, []byte{'_'}},
	{factorSpecialChars, []byte{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', ']', '^', '`', '{', '|', '}', '~'}},
	{factorDigit | factorDigit, []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}},
	{factorDigit | factorUnderscore, []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '_'}},
	{factorDigit & factorUnderscore, []byte{}},
}

func TestExpandCodeset(t *testing.T) {
	for _, tt := range expandCodesetTests {
		s := expandCodeset(tt.in)
		if !bytes.Equal(s, tt.out) {
			t.Errorf("got '%#v'\nwant '%#v'", s, tt.out)
		}
	}
}

var generatePasswordTests = []struct {
	in option
}{
	{option{len: 0, factors: 0}},
	{option{len: 8, factors: 0}},
	{option{len: 0, factors: factorDigit}},
	{option{len: 10, factors: factorDigit}},
	{option{len: 10, factors: factorDigit | factorAlphabetLarge}},
}

func TestGeneratePassword(t *testing.T) {
	for _, tt := range generatePasswordTests {
		s := generatePassword(tt.in)

		// length
		if tt.in.factors == 0 && len(s) != 0 {
			t.Errorf("got '%#v'\nwant '%#v'", len(s), tt.in.len)
		}
		if tt.in.factors != 0 && len(s) != tt.in.len {
			t.Errorf("got '%#v'\nwant '%#v'", len(s), tt.in.len)
		}

		// contains
		s2 := s
		for _, code := range expandCodeset(tt.in.factors) {
			s2 = strings.Replace(s2, string(code), " ", -1)
		}
		if s2 != strings.Repeat(" ", len(s2)) {
			t.Errorf("got '%#v'\nwant '%#v'", s2, strings.Repeat(" ", len(s2)))
		}
	}
}
