package jobs

// Formatter transforms an event into a description, to be displayed by a presenter.
type Formatter func(event Event) (text string)
