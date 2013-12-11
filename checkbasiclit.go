package interactive

import (
	"go/ast"
)

func checkBasicLit(ctx *Ctx, lit *ast.BasicLit, env *Env) (*BasicLit, []error) {
	return &BasicLit{BasicLit: lit}, nil
}
