package eval

import (
	"go/ast"
)

func checkSelectorExpr(ctx *Ctx, selector *ast.SelectorExpr, env *Env) (aexpr *SelectorExpr, errs []error) {
	aexpr = &SelectorExpr{SelectorExpr: selector}

	var moreErrs []error
	if aexpr.X, moreErrs = CheckExpr(ctx, selector.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	return aexpr, errs
}
