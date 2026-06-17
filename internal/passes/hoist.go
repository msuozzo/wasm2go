package passes

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// HoistVars lifts wasm2go's per-block temporaries out of the anonymous block
// scopes it emits to mirror wasm block structure. This unlocks more UnnestBlocks
// candidates (Go forbids a goto jumping over an in-scope declaration) which can
// flatten deeply-nested block constructs.
//
// A temp is lifted only as far as the innermost enclosing *natural* Go scope --
// a for/if/switch body or a switch case clause -- rather than all the way to
// function scope. Anonymous blocks exist solely to model wasm and are the ones
// UnnestBlocks flattens, so hoisting past them is what enables flattening;
// for/if/switch blocks are idiomatic Go we want to keep, so a temp already
// declared directly inside one is left in place. This still empties anonymous
// blocks of declarations (the natural ancestor always encloses them), so they
// remain unnestable, while keeping declarations closer to their use.
//
// Codegen emits temps long-form (`var tN T = e`); HoistVars splits each into a
// hoisted `var tN T` plus an in-place `tN = e`, then UnnestBlocks collapses the
// cascade. Safe because temps are single-assignment and uniquely named, and the
// value-stack discipline keeps every use within the subtree of the temp's block
// -- so the innermost natural ancestor always encloses every use. The type rides
// in on the decl (no inference). Runs before UnnestBlocks/InlineSingleGoto.
func HoistVars(fn *ast.FuncDecl) {
	if fn.Body != nil {
		hoistFuncVars(fn)
	}
}

type hoistedDecl struct {
	name string
	typ  ast.Expr
}

func hoistFuncVars(fn *ast.FuncDecl) {
	// Count declaration sites per name across the whole function.
	// A name declared more than once will not be hoisted: sibling blocks may
	// legally declare the same temp name as different types and, to simplify
	// reasoning about safety, we detect and skip these hoists.
	declCount := map[string]int{}
	ast.Inspect(fn.Body, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.GenDecl: // long form: var x T  /  var x T = e
			if n.Tok == token.VAR {
				for _, spec := range n.Specs {
					if vs, ok := spec.(*ast.ValueSpec); ok {
						for _, id := range vs.Names {
							declCount[id.Name]++
						}
					}
				}
			}
		case *ast.AssignStmt: // short form: x := e
			if n.Tok == token.DEFINE {
				for _, lhs := range n.Lhs {
					if id, ok := lhs.(*ast.Ident); ok {
						declCount[id.Name]++
					}
				}
			}
		}
		return true
	})
	hoistable := func(name string) bool { return name != "_" && declCount[name] == 1 }

	// Hoisted decls accumulate per target scope -- the innermost enclosing
	// natural (for/if/switch) scope -- keyed by the node owning its statement
	// list (a block or case clause), then prepended after the walk.
	hoistedByTarget := map[ast.Node][]hoistedDecl{}

	// Stack of natural scopes currently open; the top is the active hoist
	// target, defaulting to fn.Body when no natural scope encloses the decl.
	var scopes []ast.Node
	target := func() ast.Node {
		if len(scopes) == 0 {
			return fn.Body
		}
		return scopes[len(scopes)-1]
	}
	astutil.Apply(fn.Body, func(c *astutil.Cursor) bool {
		if isNaturalScope(c) {
			scopes = append(scopes, c.Node())
		}
		return true
	}, func(c *astutil.Cursor) bool {
		if isNaturalScope(c) {
			scopes = scopes[:len(scopes)-1]
			return true
		}
		ds, ok := c.Node().(*ast.DeclStmt)
		if !ok {
			return true
		} else if gd, ok := ds.Decl.(*ast.GenDecl); !ok || gd.Tok != token.VAR || len(gd.Specs) != 1 {
			return true
		} else if vs, ok := gd.Specs[0].(*ast.ValueSpec); !ok || vs.Type == nil {
			return true
		} else {
			tgt := target()
			// A decl already directly inside the natural scope is left in place:
			// it is exactly where we would hoist it to.
			directlyInScope := scopeOwner(c.Parent()) == tgt
			if len(vs.Values) == 0 { // Decl form `var ... T`
				if directlyInScope {
					return true
				}
				// Nested: condition temps / block results. Hoist only when every
				// name is uniquely declared; otherwise leave the decl in place.
				for _, id := range vs.Names {
					if id.Name != "_" && !hoistable(id.Name) {
						return true
					}
				}
				// Hoist the bare `var x T` but leave a `x = 0` value reset.
				var lhs, rhs []ast.Expr
				for _, id := range vs.Names {
					if id.Name == "_" {
						continue
					}
					hoistedByTarget[tgt] = append(hoistedByTarget[tgt], hoistedDecl{id.Name, vs.Type})
					lhs = append(lhs, id)
					rhs = append(rhs, zeroValue(vs.Type))
				}
				if len(lhs) == 0 {
					c.Delete() // e.g. `var _ T` -- nothing to hoist
				} else {
					c.Replace(&ast.AssignStmt{Tok: token.ASSIGN, Lhs: lhs, Rhs: rhs})
				}
			} else if len(vs.Names) == 1 { // Value form `var x T = e`
				name := vs.Names[0].Name
				switch {
				case name == "_":
					// Blanked by RemoveUnusedLocals: no declaration needed, but keep
					// rhs's side effects and simplify assign so block can flatten.
					c.Replace(valueAssign(vs))
				case directlyInScope:
					// Already directly inside a control-flow block: leave in place.
				case hoistable(name):
					// Hoist `var x T` and leave `x = e` in place.
					hoistedByTarget[tgt] = append(hoistedByTarget[tgt], hoistedDecl{name, vs.Type})
					c.Replace(valueAssign(vs))
				default:
					// Declared more than once: leave the block-scoped declaration.
				}
			}
		}
		return true
	})

	for tgt, hoisted := range hoistedByTarget {
		decls := make([]ast.Stmt, len(hoisted))
		for i, h := range hoisted {
			decls[i] = &ast.DeclStmt{Decl: &ast.GenDecl{Tok: token.VAR, Specs: []ast.Spec{
				&ast.ValueSpec{Names: []*ast.Ident{ast.NewIdent(h.name)}, Type: h.typ}}}}
		}
		prependStmts(tgt, decls)
	}
}

// isNaturalScope reports whether the cursor's node introduces a natural Go scope
// -- a for/if/switch body or a switch/select case clause -- as opposed to an
// anonymous block that wasm2go emits purely to mirror wasm block structure.
func isNaturalScope(c *astutil.Cursor) bool {
	switch c.Node().(type) {
	case *ast.CaseClause, *ast.CommClause:
		return true
	case *ast.BlockStmt:
		switch c.Parent().(type) {
		case *ast.IfStmt, *ast.ForStmt, *ast.RangeStmt:
			return true
		}
	}
	return false
}

// scopeOwner returns the node owning the statement list a decl's parent belongs
// to (a block or case clause), or nil if the parent owns no statement list.
func scopeOwner(parent ast.Node) ast.Node {
	switch parent.(type) {
	case *ast.BlockStmt, *ast.CaseClause, *ast.CommClause:
		return parent
	}
	return nil
}

// prependStmts inserts decls at the front of the target scope's statement list.
func prependStmts(target ast.Node, decls []ast.Stmt) {
	switch t := target.(type) {
	case *ast.BlockStmt:
		t.List = append(decls, t.List...)
	case *ast.CaseClause:
		t.Body = append(decls, t.Body...)
	case *ast.CommClause:
		t.Body = append(decls, t.Body...)
	}
}

// valueAssign turns a `var x T = e` spec into the assignment `x = e`.
func valueAssign(vs *ast.ValueSpec) *ast.AssignStmt {
	return &ast.AssignStmt{Tok: token.ASSIGN, Lhs: []ast.Expr{vs.Names[0]}, Rhs: []ast.Expr{vs.Values[0]}}
}

// zeroValue is the zero value used to re-initialize a hoisted no-initializer
func zeroValue(typ ast.Expr) ast.Expr {
	if id, ok := typ.(*ast.Ident); ok && id.Name == "any" {
		return ast.NewIdent("nil")
	}
	return &ast.BasicLit{Kind: token.INT, Value: "0"}
}
