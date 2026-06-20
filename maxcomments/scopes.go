package maxcomments

import (
	"fmt"
	"go/ast"
)

// funcScope is one independently-checkable function unit: a FuncDecl (named
// function or method) or an anonymous FuncLit. Each comment in a file is
// attributed to the innermost scope that contains it, so a closure's comments
// are checked against the closure and never folded into the function that
// encloses it.
type funcScope struct {
	node     ast.Node            // *ast.FuncDecl or *ast.FuncLit
	name     string              // human label used in diagnostics
	doc      *ast.CommentGroup   // doc comment (FuncDecl only; nil otherwise)
	comments []*ast.CommentGroup // body comments attributed to this scope
}

// collectFuncScopes walks the file and returns every FuncDecl and FuncLit as a
// funcScope, then attributes each comment group to its innermost scope.
func collectFuncScopes(file *ast.File) []*funcScope {
	var scopes []*funcScope

	ast.Inspect(file, func(n ast.Node) bool {
		switch fn := n.(type) {
		case *ast.FuncDecl:
			scopes = append(scopes, &funcScope{
				node: fn,
				name: fmt.Sprintf("function %q", fn.Name.Name),
				doc:  fn.Doc,
			})
		case *ast.FuncLit:
			scopes = append(scopes, &funcScope{
				node: fn,
				name: "function literal",
			})
		}
		return true
	})

	attributeComments(scopes, file.Comments)
	return scopes
}

// attributeComments assigns each non-doc comment group to the innermost scope
// whose source range contains it. Doc comments are skipped because they sit
// outside the function body range and are accounted for via funcScope.doc.
func attributeComments(scopes []*funcScope, comments []*ast.CommentGroup) {
	docs := make(map[*ast.CommentGroup]struct{}, len(scopes))
	for _, s := range scopes {
		if s.doc != nil {
			docs[s.doc] = struct{}{}
		}
	}

	for _, g := range comments {
		if _, isDoc := docs[g]; isDoc {
			continue
		}

		var best *funcScope
		for _, s := range scopes {
			if g.Pos() < s.node.Pos() || g.End() > s.node.End() {
				continue
			}
			if best == nil || scopeSpan(s) < scopeSpan(best) {
				best = s
			}
		}

		if best != nil {
			best.comments = append(best.comments, g)
		}
	}
}

// scopeSpan is the byte span of a scope's node, used to pick the innermost
// (smallest) of several nested scopes that all contain a comment.
func scopeSpan(s *funcScope) int {
	return int(s.node.End() - s.node.Pos())
}
