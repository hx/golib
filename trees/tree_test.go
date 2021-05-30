package trees_test

import (
	"bytes"
	. "github.com/hx/golib/testing"
	. "github.com/hx/golib/trees"
	"testing"
)

func TestTreeBasicOutput(t *testing.T) {
	buf := new(bytes.Buffer)
	tree := NewWithWriter(buf)
	tree.Indent = "-"
	tree.Prefix = "+"
	tree.AddChild().Update("foo")
	tree.AddChild().Update("bar")
	baz := tree.AddChild()
	baz.Update("baz")
	bazz := baz.AddChild()
	bazz.Update("bazz")
	bazzz := bazz.AddChild()
	bazzz.Update("bazzz")
	bazzz.AddChild().Update("bazzzz")
	baz.AddChild().Update("quux")
	expected := `
+foo
+bar
+baz
-+bazz
--+bazzz
---+bazzzz
-+quux
`
	Equals(t, expected[1:], buf.String())
}
