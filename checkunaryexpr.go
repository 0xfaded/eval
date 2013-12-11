package interactive

import (
	"go/ast"
)

func checkUnaryExpr(ctx *Ctx, unary *ast.UnaryExpr, env *Env) (aexpr *UnaryExpr, errs []error) {
	aexpr = &UnaryExpr{UnaryExpr: unary}

	var moreErrs []error
	if aexpr.X, moreErrs = checkExpr(ctx, unary.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	return aexpr, errs
}
