package passes

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// ReconstructLoops rewrites wasm2go's goto-based loops into native Go for
// statements. Codegen emits a wasm `loop` as a block labeled at its start, each
// back-edge `br` to it a `goto` to that label; a wasm `block` (a forward break
// target) is instead labeled at its end and left in goto form, since Go has no
// break for a bare block. So a label wrapping a block and targeted from within
// uniquely marks a loop:
//
//	l1: {
//		if cond { goto l0 }   // forward exit (to an enclosing block end)
//		...
//		goto l1               // unconditional back-edge
//	}
//	l0: ;
//
// becomes
//
//	for {
//		if cond { goto l0 }   // forward exits keep their gotos
//		...
//	}                             // trailing back-edge dropped; for{} reiterates
//	l0: ;
//
// Back-edges turn into continue (labeled when nested inside another for, where an
// unlabeled continue would target the inner loop). A trailing unconditional
// back-edge is dropped, since a Go for{} re-iterates on its own; conversely, a
// body that can fall off its end gets an explicit break, since a for{} would
// otherwise loop forever. Forward exits keep their gotos (handled by later
// passes). Runs after InlineSwitchTargets, whose walk predates for statements,
// and before HoistVars, so hoisted temps land inside the reconstructed loop.
func ReconstructLoops(fn *ast.FuncDecl) {
	if fn.Body == nil {
		return
	}
	// Post-order: inner loops become for statements before their enclosing loop,
	// so back-edge nesting (and thus continue labeling) is computed correctly.
	astutil.Apply(fn.Body, nil, func(c *astutil.Cursor) bool {
		ls, ok := c.Node().(*ast.LabeledStmt)
		if !ok {
			return true
		}
		body, ok := ls.Stmt.(*ast.BlockStmt)
		if !ok || !hasBackEdge(body, ls.Label.Name) {
			return true // a block end-label, not a loop start
		}
		label := ls.Label.Name

		// Drop a trailing unconditional back-edge: the for{} reiterates for us,
		// and the body then provably never falls off its end.
		reiterates := false
		if n := len(body.List); n > 0 {
			if br, ok := body.List[n-1].(*ast.BranchStmt); ok && br.Tok == token.GOTO && br.Label.Name == label {
				body.List = body.List[:n-1]
				reiterates = true
			}
		}

		// Remaining back-edges become continue.
		labeled := convertBackEdges(body, label)

		// Without a trailing back-edge, falling off the body's end exits the loop
		// in wasm; a Go for{} would re-iterate, so break out explicitly.
		if !reiterates && canFallThrough(body.List) {
			body.List = append(body.List, &ast.BranchStmt{Tok: token.BREAK})
		}

		forStmt := &ast.ForStmt{Body: body}
		if labeled {
			ls.Stmt = forStmt // a continue references the label; keep it
		} else {
			c.Replace(forStmt) // label is now unused; drop it
		}
		return true
	})
}

// hasBackEdge reports whether body contains a `goto label`, i.e. label marks a
// loop start (branched to from within) rather than a block end.
func hasBackEdge(body *ast.BlockStmt, label string) bool {
	found := false
	ast.Inspect(body, func(n ast.Node) bool {
		if br, ok := n.(*ast.BranchStmt); ok && br.Tok == token.GOTO && br.Label.Name == label {
			found = true
		}
		return !found
	})
	return found
}

// convertBackEdges turns every `goto label` in body into a continue, returning
// whether any needed a label because it sits inside a nested for statement (an
// unlabeled continue there would target the inner loop, not this one).
func convertBackEdges(body *ast.BlockStmt, label string) (labeled bool) {
	depth := 0
	astutil.Apply(body, func(c *astutil.Cursor) bool {
		if _, ok := c.Node().(*ast.ForStmt); ok {
			depth++
		}
		return true
	}, func(c *astutil.Cursor) bool {
		if _, ok := c.Node().(*ast.ForStmt); ok {
			depth--
			return true
		}
		if br, ok := c.Node().(*ast.BranchStmt); ok && br.Tok == token.GOTO && br.Label.Name == label {
			if depth > 0 {
				c.Replace(&ast.BranchStmt{Tok: token.CONTINUE, Label: ast.NewIdent(label)})
				labeled = true
			} else {
				c.Replace(&ast.BranchStmt{Tok: token.CONTINUE})
			}
		}
		return true
	})
	return labeled
}

// canFallThrough reports whether control can reach the end of a loop body, so a
// Go for{} needs an explicit break to stop re-iterating.
func canFallThrough(list []ast.Stmt) bool {
	if len(list) == 0 {
		return true
	}
	// continue/break/goto/fallthrough all leave the bottom of the body; terminates
	// covers goto/fallthrough plus return/panic and terminating if/switch.
	if _, ok := list[len(list)-1].(*ast.BranchStmt); ok {
		return false
	}
	return !terminates(list)
}
