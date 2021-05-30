package ansi

import (
	"strconv"
	"strings"
)

// EscapeSequence represents a basic ANSI escape sequence, in a way that allows bitwise combination.
// For example:
//
//  fmt.Print((Yellow | BlueBG | Bold).String())
//  # => \033[1;33;44m
type EscapeSequence uint

const (
	// Black foreground
	Black EscapeSequence = iota

	// Red foreground
	Red

	// Green foreground
	Green

	// Yellow foreground
	Yellow

	// Blue foreground
	Blue

	// Magenta foreground
	Magenta

	// Cyan foreground
	Cyan

	// White foreground
	White
)

const (
	// BlackBG background text
	BlackBG EscapeSequence = iota << 3

	// RedBG background text
	RedBG

	// GreenBG background text
	GreenBG

	// YellowBG background text
	YellowBG

	// BlueBG background text
	BlueBG

	// MagentaBG background text
	MagentaBG

	// CyanBG background text
	CyanBG

	// WhiteBG background text
	WhiteBG
)

const (
	// Bright text
	Bright EscapeSequence = 1 << (iota + 6)

	// BrightBG background
	BrightBG

	// Bold text
	Bold

	// Underline text
	Underline
)

func (s EscapeSequence) IsBright() bool    { return s&Bright != 0 }
func (s EscapeSequence) IsBrightBG() bool  { return s&BrightBG != 0 }
func (s EscapeSequence) IsBold() bool      { return s&Bold != 0 }
func (s EscapeSequence) IsUnderline() bool { return s&Underline != 0 }

func (s EscapeSequence) Bright() EscapeSequence    { return s | Bright }
func (s EscapeSequence) BrightBG() EscapeSequence  { return s | BrightBG }
func (s EscapeSequence) Bold() EscapeSequence      { return s | Bold }
func (s EscapeSequence) Underline() EscapeSequence { return s | Underline }

func (s EscapeSequence) NotBright() EscapeSequence    { return s & ^Bright }
func (s EscapeSequence) NotBrightBG() EscapeSequence  { return s & ^BrightBG }
func (s EscapeSequence) NotBold() EscapeSequence      { return s & ^Bold }
func (s EscapeSequence) NotUnderline() EscapeSequence { return s & ^Underline }

// String converts the EscapeSequence to a string, for writing to a terminal.
func (s EscapeSequence) String() (str string) {
	var nums []int
	if s == 0 {
		nums = []int{0}
	} else {
		if s.IsBold() {
			nums = []int{1}
		}
		if s.IsUnderline() {
			nums = append(nums, 4)
		}
		fg := s&0b111 + 30
		if s.IsBright() {
			fg += 60
		}
		if fg != 30 {
			nums = append(nums, int(fg))
		}
		bg := s&0b111000>>3 + 40
		if s.IsBrightBG() {
			bg += 60
		}
		if bg != 40 {
			nums = append(nums, int(bg))
		}
	}
	parts := make([]string, len(nums))
	for i, num := range nums {
		parts[i] = strconv.Itoa(num)
	}
	return "\033[" + strings.Join(parts, ";") + "m"
}

// Bytes converts the EscapeSequence to a byte array, for writing to a terminal.
func (s EscapeSequence) Bytes() (bytes []byte) { return []byte(s.String()) }

// Wrap the given string with the escape sequence as a prefix, and Reset as a suffix.
func (s EscapeSequence) Wrap(stringToWrap string) string {
	return s.String() + stringToWrap + Reset
}
