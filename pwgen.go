package pwgen

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/nbutton23/zxcvbn-go"
)

const (
	factorDigit = 1 << iota
	factorAlphabetLarge
	factorAlphabetSmall
	factorUnderscore
	factorSpecialChars
)

const (
	codeAlphabetLargeStart = 0x41
	codeAlphabetLargeEnd   = 0x5a
	codeAlphabetSmallStart = 0x61
	codeAlphabetSmallEnd   = 0x7a
)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
)

var (
	codesetDigit         []byte
	codesetAlphabetLarge []byte
	codesetAlphabetSmall []byte
	codesetSpecialChars  []byte
)

type option struct {
	len     int
	factors int
}

func init() {
	codesetDigit = getCodesetDigit()
	codesetAlphabetLarge = getCodesetAlphabetLarge()
	codesetAlphabetSmall = getCodesetAlphabetSmall()
	codesetSpecialChars = getCodesetSpecialChars()

	rand.Seed(time.Now().UnixNano())
}

func getCodesetDigit() []byte {
	return []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
}

func getCodesetAlphabetLarge() []byte {
	codeset := []byte{}

	for i := codeAlphabetLargeStart; i <= codeAlphabetLargeEnd; i++ {
		codeset = append(codeset, byte(i))
	}

	return codeset
}

func getCodesetAlphabetSmall() []byte {
	codeset := []byte{}

	for i := codeAlphabetSmallStart; i <= codeAlphabetSmallEnd; i++ {
		codeset = append(codeset, byte(i))
	}

	return codeset
}

func getCodesetSpecialChars() []byte {
	return []byte{'!', '"', '#', '$', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/', ':', ';', '<', '=', '>', '?', '@', '[', ']', '^', '`', '{', '|', '}', '~'}
}

func randomByte(codeset []byte) byte {
	if len(codeset) <= 0 {
		return ' '
	}
	index := rand.Int31n(int32(len(codeset)))

	return codeset[index]
}

func expandCodeset(factors int) []byte {
	codeset := []byte{}

	for factors > 0 {
		switch {
		case factors&factorDigit != 0:
			codeset = append(codeset, codesetDigit...)
			factors -= factorDigit

		case factors&factorAlphabetLarge != 0:
			codeset = append(codeset, codesetAlphabetLarge...)
			factors -= factorAlphabetLarge

		case factors&factorAlphabetSmall != 0:
			codeset = append(codeset, codesetAlphabetSmall...)
			factors -= factorAlphabetSmall

		case factors&factorUnderscore != 0:
			codeset = append(codeset, '_')
			factors -= factorUnderscore

		case factors&factorSpecialChars != 0:
			codeset = append(codeset, codesetSpecialChars...)
			factors -= factorSpecialChars
		}
	}

	return codeset
}

func generatePassword(opt option) string {
	chars := []byte{}
	codeset := expandCodeset(opt.factors)
	if len(codeset) <= 0 {
		return ""
	}

	for i := 0; i < opt.len; i++ {
		chars = append(chars, randomByte(codeset))
	}

	return string(chars)
}

// Pwgen is ...
func Pwgen(w io.Writer, length int, count int, isDigit bool, isLarge bool, isSmall bool, isUnderscore bool, isSpecial bool) error {
	if length < minPasswordLength || maxPasswordLength < length {
		return fmt.Errorf("length range error %d\n", length)
	}
	if count <= 0 {
		return fmt.Errorf("count range error %d\n", count)
	}

	factors := 0
	if isDigit {
		factors |= factorDigit
	}
	if isLarge {
		factors |= factorAlphabetLarge
	}
	if isSmall {
		factors |= factorAlphabetSmall
	}
	if isUnderscore {
		factors |= factorUnderscore
	}
	if isSpecial {
		factors |= factorSpecialChars
	}

	opt := option{
		len:     length,
		factors: factors,
	}
	passwords := []string{}

	for i := 0; i < count; {
		password := generatePassword(opt)
		strength := zxcvbn.PasswordStrength(password, nil)

		if strength.Score == 4 { // Score is [0,1,2,3,4] (from risky to strong)
			fmt.Fprintln(w, password)
			passwords = append(passwords, password)
			i++
		}
	}

	return nil
}
