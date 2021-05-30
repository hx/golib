package ansi

import "fmt"

func CursorUp(lines int) string {
	if lines == 0 {
		return ""
	}
	if lines < 0 {
		return CursorDown(-lines)
	}
	return fmt.Sprintf("\033[%dA", lines)
}

func CursorDown(lines int) string {
	if lines == 0 {
		return ""
	}
	if lines < 0 {
		return CursorUp(-lines)
	}
	return fmt.Sprintf("\033[%dB", lines)
}

func CursorForward(columns int) string {
	if columns == 0 {
		return ""
	}
	if columns < 0 {
		return CursorBack(columns)
	}
	return fmt.Sprintf("\033[%dC", columns)
}

func CursorBack(columns int) string {
	if columns == 0 {
		return ""
	}
	if columns < 0 {
		return CursorForward(columns)
	}
	return fmt.Sprintf("\033[%dD", columns)
}

// CursorColumn moves the cursor to an absolute position, with the leftmost
// column being 1.
func CursorColumn(column int) string {
	return fmt.Sprintf("\033[%dG", column)
}
