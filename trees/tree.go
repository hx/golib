package trees

import (
	"github.com/hx/golib/ansi"
	"golang.org/x/term"
	"io"
	"os"
	"sync"
)

// Tree writes, and re-writes, a tree of information to the given io.Writer using ANSI cursor positioning.
// Indent and Prefix control how the tree is formed. Second level and below are prefixed with Prefix,
// and third level and below are prefixed with Indent * (level - 2), followed by Prefix.
//
// For example, with Indent of "--" and Prefix of "+ ", you might see:
//
//  Heading 1
//  + Heading 2
//  --+ Heading 3
//  ----+ Heading 4
//
// A zero Tree is not valid. Use New or NewWithWriter.
type Tree struct {
	// Indent is prepended to each Node the number of times required to display the tree's hierarchy.
	Indent string

	// Prefix is prepended to each of the tree's descendent Nodes, after Indent but before the Node's contents.
	Prefix string

	// TruncatedLineSuffix is appended to lines that have been truncated due to being longer than the width of the
	// terminal.
	TruncatedLineSuffix string

	*node

	writer io.Writer
	mutex  sync.Mutex
	height int
}

// New returns a new Tree that will write to os.Stdout.
func New() *Tree {
	return NewWithWriter(os.Stdout)
}

// NewWithWriter returns a new Tree that will write to writer.
func NewWithWriter(writer io.Writer) (tree *Tree) {
	tree = &Tree{
		node:                new(node),
		Indent:              "  ",
		Prefix:              "  ",
		TruncatedLineSuffix: "…",
		writer:              writer,
	}
	tree.node.tree = tree
	return
}

// render updates the tree's appearance. In summary, it:
//
//  - Moves the cursor up to the affected line
//  - Erases the line and writes the new content
//  - If the line has increased in height,
//    - Erases and re-writes all lines underneath it
//  - otherwise
//    - Moves the cursor back to the line below the last line of the tree
//
// The golang.org/x/term package is used to determine
// the terminal's width and height. Content that would be above the top
// or past the right margin of the window is not written.
func (t *Tree) render(target *node, content string) {
	var (
		// Dimensions of the terminal window.
		termWidth, termHeight, _ = term.GetSize(int(os.Stdin.Fd()))

		// From the top, number of lines skipped due to no change.
		skipped int

		// Number of lines to be hidden above the top of the window.
		// This strategy ignores hidden line updates, so resizing
		// after a render will have unpredictable results for the
		// hidden lines. Unhidden lines should remain stable.
		hidden int

		// Whether the target node has been found in the walk cycle.
		found bool

		// Whether this is the first time the node is being rendered.
		// If it is, all subsequent nodes will be re-rendered one
		// line below their former position.
		isFirstRender bool
	)
	// Once the node is non-blank, preserve its height to avoid unsupported tree shrinkage.
	if content != "" {
		isFirstRender = target.contentHeight == 0
		target.contentHeight = 1
	}
	if target.content == content {
		return
	}
	target.content = content
	t.walk(0, func(node *node, level int) (keepWalking bool) {
		keepWalking = true
		var cursorDown int
		if !found {
			if node != target {
				skipped += node.contentHeight
				return
			}
			cursorUp := t.height - skipped
			if termHeight > 0 && cursorUp+1 > termHeight {
				hidden = cursorUp + 1 - termHeight
				cursorUp -= hidden
			}
			if cursorUp > 0 {
				t.write(ansi.CursorUp(cursorUp))
			}
			if isFirstRender {
				t.height++
			} else if cursorUp > 1 && hidden == 0 {
				cursorDown = cursorUp - 1
			}
			found = true
		}
		if node.contentHeight == 0 {
			return
		}
		if hidden > 0 {
			hidden--
			return
		}
		line := t.indent(level) + node.content
		if termWidth > 0 {
			line = t.truncate(line+ansi.EraseToEOL, termWidth)
		}
		t.write(line + "\n")
		if cursorDown > 0 {
			t.write(ansi.CursorDown(cursorDown))
			keepWalking = false
		}
		return
	})
}

func (t *Tree) write(str string) {
	if _, err := t.writer.Write([]byte(str)); err != nil {
		panic(err)
	}
}

func (t *Tree) indent(level int) (indent string) {
	for i := level; i > 0; i-- {
		if i == 1 {
			indent += t.Prefix
		} else {
			indent += t.Indent
		}
	}
	return
}

func (t *Tree) truncate(str string, width int) string {
	suffixLength := len([]rune(t.TruncatedLineSuffix))
	if width > suffixLength && ansi.Len(str) > width {
		return ansi.Truncate(str, width-suffixLength) + t.TruncatedLineSuffix
	}
	return str
}
