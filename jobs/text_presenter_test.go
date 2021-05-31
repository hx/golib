package jobs_test

import (
	"bytes"
	. "github.com/hx/golib/jobs"
	. "github.com/hx/golib/testing"
	"testing"
)

type textJob struct {
	name string
	work []Job
}

func (j textJob) Run(ctx *Context) (err error) {
	if j.work != nil {
		err = ctx.Run(j.work...).Wait()
	}
	return
}

func TestTextPresenter_Run(t *testing.T) {
	sequence := &Sequence{
		&textJob{"foo", nil},
		&textJob{"bar", []Job{
			&textJob{"b1", nil},
			&textJob{"b2", []Job{
				&textJob{"b2a", nil},
			}},
			&textJob{"b3", nil},
		}},
		&textJob{"bazzz", nil},
	}

	formatter := func(topName string) func(Event) string {
		return func(event Event) string {
			if job, ok := event.Job().(*textJob); ok {
				return job.name
			}
			return topName
		}
	}

	t.Run("without top level heading", func(t *testing.T) {
		expected := `
foo    ok
bar  
  b1  ok
  b2
    b2a  ok
  b3  ok
bazzz  ok
`[1:]

		buf := new(bytes.Buffer)
		Run(NewTextPresenter(sequence, buf, formatter(""))).Wait()
		Equals(t, expected, buf.String())

	})

	t.Run("with top level heading", func(t *testing.T) {
		expected := `
TOP
  foo    ok
  bar  
    b1  ok
    b2
      b2a  ok
    b3  ok
  bazzz  ok
`[1:]

		buf := new(bytes.Buffer)
		Run(NewTextPresenter(sequence, buf, formatter("TOP"))).Wait()
		Equals(t, expected, buf.String())
	})
}
