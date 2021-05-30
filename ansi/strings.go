package ansi

import (
	"regexp"
	"strings"
)

// EscapeSequencePattern is the pattern used for matching (and removing) ANSI escape sequences.
// From https://stackoverflow.com/questions/14693701/how-can-i-remove-the-ansi-escape-sequences-from-a-string-in-python
var EscapeSequencePattern = regexp.MustCompile(`(\x9B|\x1B\[)[0-?]*[ -/]*[@-~]`)

// Len returns the character (rune) length of str with all ANSI escape sequences removed.
func Len(str string) int {
	return len([]rune(Strip(str)))
}

// Strip returns str with all ANSI escape sequences removed.
func Strip(str string) string {
	return EscapeSequencePattern.ReplaceAllString(str, "")
}

// PadRightRune pads the right side of the given string using rune, up to length, as measured by Len.
func PadRightRune(str string, length int, padding rune) string {
	length -= Len(str)
	if length > 0 {
		return str + strings.Repeat(string(padding), length)
	}
	return str
}

// PadRight pads the right side of the given string using space characters, up to length, as measured by Len.
func PadRight(str string, length int) string {
	return PadRightRune(str, length, ' ')
}

// Truncate a string containing ANSI escape sequences, to maxLength characters (runes), excluding
// escape sequences. All escape sequences are preserved.
func Truncate(str string, maxLength int) string {
	if maxLength == 0 {
		return ""
	}
	if maxLength < 0 {
		return str
	}
	length := Len(str)
	if length <= maxLength {
		return str
	}
	escapes := EscapeSequencePattern.FindAllStringIndex(str, -1)
	if len(escapes) == 0 {
		return string([]rune(str)[0:maxLength])
	}
	escapeIndex := len(escapes)
	escapes = append([][]int{{0, 0}}, append(escapes, []int{len(str)})...)
	for {
		var (
			a         = escapes[escapeIndex][1]
			b         = escapes[escapeIndex+1][0]
			runes     = []rune(str[a:b])
			newLength = length - len(runes)
		)
		if newLength <= maxLength {
			return str[:a] + string(runes[:maxLength-newLength]) + str[b:]
		}
		str = str[:a] + str[b:]
		length -= len(runes)
		escapeIndex--
	}
}
