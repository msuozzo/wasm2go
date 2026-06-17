package passes

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"testing"
)

func TestReconstructLoops(t *testing.T) {
	for _, tc := range []struct {
		name          string
		input, expect string
	}{{
		name:   "trailing back-edge becomes a bare for",
		input:  `l1:{a();goto l1}`,
		expect: `for {a()}`,
	}, {
		name:   "conditional back-edge becomes continue with a fall-through break",
		input:  `l1:{a();if c() {goto l1}}`,
		expect: `for {a();if c() {continue};break}`,
	}, {
		name:   "forward exit keeps its goto",
		input:  `l1:{if c() {goto l0};a();goto l1};l0:;b()`,
		expect: `for {if c() {goto l0};a()};l0:;b()`,
	}, {
		name:   "no back-edge is left untouched (a block end-label)",
		input:  `{if c() {goto l0};a()};l0:;b()`,
		expect: `{if c() {goto l0};a()};l0:;b()`,
	}, {
		name:   "nested loops use a labeled continue for the outer back-edge",
		input:  `l2:{l1:{if c() {goto l2};goto l1}}`,
		expect: `l2:for {for {if c() {continue l2}};break}`,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got := rewriteLoops(t, tc.input, true)
			want := rewriteLoops(t, tc.expect, false)
			if got != want {
				t.Errorf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

// rewriteLoops parses src and conditionally applies ReconstructLoops.
func rewriteLoops(t *testing.T, src string, transform bool) string {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", "package p\nfunc f() {\n"+src+"\n}", parser.SkipObjectResolution)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	fn := file.Decls[0].(*ast.FuncDecl)
	if transform {
		ReconstructLoops(fn)
	}
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, fn); err != nil {
		t.Fatal(err)
	}
	return buf.String()
}
