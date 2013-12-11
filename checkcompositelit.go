package interactive

import (
	"go/ast"
)

func checkCompositeLit(ctx *Ctx, lit *ast.CompositeLit, env *Env) (aexpr *CompositeLit, errs []error) {
	aexpr = &CompositeLit{CompositeLit: lit}

	var moreErrs []error
	if aexpr.Type, moreErrs = checkTypeExpr(ctx, lit.Type, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	for i := range lit.Elts {
		if aexpr.Elts[i], moreErrs = checkExpr(ctx, lit.Elts[i], env); moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
	}
	return aexpr, errs
}
