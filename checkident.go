package interactive

import (
	"go/ast"
)

func checkIdent(ctx *Ctx, ident *ast.Ident, env *Env) (*Ident, []error) {
	return &Ident{Ident: ident}, nil
}
