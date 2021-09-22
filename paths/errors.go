package paths

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrNonDirectory Error = "not a directory"
)
