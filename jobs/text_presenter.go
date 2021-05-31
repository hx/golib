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

	// Indent is prepended to sub-jobs.
	Indent string

	job       Job
	writer    io.Writer
	formatter Formatter

	width      []int
	runningJob Job

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
		Indent:        "  ",
	}
}

// Run implements Job.
func (t *TextPresenter) Run(ctx *Context) (err error) {
	t.width = []int{0}
	t.queue = nil
	t.started = make(map[Job]struct{})
	for event := range ctx.Run(t.job) {
		text := t.formatter(event)
		if text == "" {
			continue
		}
		switch event := event.(type) {
		case *EventQueued:
			if t.runningJob != nil {
				t.increaseIndent()
				t.runningJob = nil
			}
			t.queue = append(t.queue, event)
			if width := ansi.Len(text); width > t.width[0] {
				t.width[0] = width
			}
		case *EventStarted:
			t.runningJob = event.Job()
			t.started[event.Job()] = struct{}{}
			t.indent()
			t.write(t.justify(text))
		case *EventFinished:
			if t.runningJob == event.Job() {
				if event.Error() == nil {
					t.write(t.SuffixOk + "\n")
				} else {
					t.write(t.SuffixFail + "\n")
					err = event.Error()
				}
			} else {
				t.decreaseIndent()
			}
			t.runningJob = nil
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

func (t *TextPresenter) justify(str string) string { return ansi.PadRight(str, t.width[0]) }
func (t *TextPresenter) write(str string)          { t.writer.Write([]byte(str)) }

func (t *TextPresenter) indent() {
	for i, w := range t.width {
		if i != 0 && w != 0 {
			t.write(t.Indent)
		}
	}
}

func (t *TextPresenter) increaseIndent() {
	t.write("\n")
	t.width = append([]int{0}, t.width...)
}

func (t *TextPresenter) decreaseIndent() {
	t.width = t.width[1:]
}
