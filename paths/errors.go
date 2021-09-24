package paths

type Error string

func (e Error) Error() string { return string(e) }

const (
	ErrDirectory    Error = "file is a directory"
	ErrNonDirectory Error = "not a directory"
	ErrNonLink      Error = "not a link"
	ErrPathNotFound Error = "path not found"
	ErrFileExists   Error = "path already exists"
	ErrBrokenLink   Error = "broken link"
	ErrNotWritable  Error = "not writable"
	ErrInvalid      Error = "invalid"
)
