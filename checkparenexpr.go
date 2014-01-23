package eval

import (
	"go/ast"
)

func checkParenExpr(ctx *Ctx, paren *ast.ParenExpr, env *Env) (*ParenExpr, []error) {
	aexpr := &ParenExpr{ParenExpr: paren}
	x, errs := CheckExpr(ctx, paren.X, env)

	aexpr.X = x
	aexpr.knownType = knownType(x.KnownType())
	aexpr.constValue = constValue(x.Const())
	return aexpr, errs
}
