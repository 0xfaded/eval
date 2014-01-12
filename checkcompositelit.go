package eval

import (
	"go/ast"
)

func checkCompositeLit(ctx *Ctx, lit *ast.CompositeLit, env *Env) (aexpr *CompositeLit, errs []error) {
	aexpr = &CompositeLit{CompositeLit: lit}

	var moreErrs []error
        // TODO confirm the type actually exists
        aexpr.Type = lit.Type

	for i := range lit.Elts {
		if aexpr.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env); moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
	}
	return aexpr, errs
}
