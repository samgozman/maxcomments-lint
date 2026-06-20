package maxcomments

import (
	"fmt"
	"go/ast"
)

// funcScope is one independently-checkable function unit (a FuncDecl or an
// anonymous FuncLit). Each comment is attributed to the innermost scope
// containing it, so a closure's comments are never folded into its encloser.
type funcScope struct {
	node     ast.Node // *ast.FuncDecl or *ast.FuncLit
	name     string
	doc      *ast.CommentGroup // doc comment (FuncDecl only; nil otherwise)
	comments []*ast.CommentGroup
}

// collectFuncScopes returns every FuncDecl and FuncLit as a funcScope, with
// each comment group attributed to its innermost scope.
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
// whose source range contains it. Doc comments are skipped (handled separately
// via funcScope.doc).
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

// scopeSpan is the byte span of a scope's node — smaller means more deeply nested.
func scopeSpan(s *funcScope) int {
	return int(s.node.End() - s.node.Pos())
}
