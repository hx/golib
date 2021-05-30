package jobs

import (
	"github.com/hx/golib/ansi"
	"io"
)

// TextPresenter wraps a Job, and uses a Formatter to display each of Job's child jobs on a new line.
type TextPresenter struct {
	// When true, jobs that were queued but did not start will be included in TextPresenter's output.
	ShowSkipped bool

	// SuffixOk is appended to lines of jobs that succeed.
	SuffixOk string

	// SuffixFail is appended to lines of jobs that fail.
	SuffixFail string

	// SuffixSkipped is appended to lines of jobs that are queued but do not run.
	SuffixSkipped string

	job       Job
	writer    io.Writer
	formatter Formatter

	width int

	queue   []Event
	started map[Job]struct{}
}

// NewTextPresenter creates a TextPresenter .
func NewTextPresenter(job Job, writer io.Writer, formatter Formatter) *TextPresenter {
	return &TextPresenter{
		job:           job,
		writer:        writer,
		formatter:     formatter,
		ShowSkipped:   true,
		SuffixOk:      "  ok",
		SuffixFail:    "  FAIL",
		SuffixSkipped: "  -",
	}
}

// Run implements Job.
func (t *TextPresenter) Run(ctx *Context) (err error) {
	t.width = 0
	t.queue = nil
	t.started = make(map[Job]struct{})
	for event := range ctx.Run(t.job) {
		text := t.formatter(event)
		if text == "" {
			continue
		}
		switch event := event.(type) {
		case *EventQueued:
			t.queue = append(t.queue, event)
			if width := len(text); width > t.width {
				t.width = width
			}
		case *EventStarted:
			t.started[event.Job()] = struct{}{}
			t.write(t.justify(text))
		case *EventFinished:
			if event.Error() == nil {
				t.write(t.SuffixOk + "\n")
			} else {
				t.write(t.SuffixFail + "\n")
				err = event.Error()
			}
		}
	}
	if t.ShowSkipped {
		for _, event := range t.queue {
			if _, ok := t.started[event.Job()]; ok {
				continue
			}
			if text := t.formatter(event); text != "" {
				t.write(t.justify(text) + t.SuffixSkipped + "\n")
			}
		}
	}
	return
}

func (t *TextPresenter) justify(str string) string { return ansi.PadRight(str, t.width) }
func (t *TextPresenter) write(str string)          { t.writer.Write([]byte(str)) }
