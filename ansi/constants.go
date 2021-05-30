package ansi

const (
	// Reset is the ANSI reset sequence.
	Reset = "\033[0m"

	// EraseToEOL erases from the cursor to the end of the current line.
	EraseToEOL = "\033[K"
)
