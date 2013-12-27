package eval

import (
	"go/ast"
)

func checkIndexExpr(ctx *Ctx, index *ast.IndexExpr, env *Env) (aexpr *IndexExpr, errs []error) {
	aexpr = &IndexExpr{IndexExpr: index}

	var moreErrs []error
	if aexpr.X, moreErrs = CheckExpr(ctx, index.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	if aexpr.Index, moreErrs = CheckExpr(ctx, index.Index, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	return aexpr, errs
}
