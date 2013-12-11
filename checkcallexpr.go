package interactive

import (
	"go/ast"
)

func checkCallExpr(ctx *Ctx, callExpr *ast.CallExpr, env *Env) (aexpr *CallExpr, errs []error) {
	aexpr = &CallExpr{CallExpr: callExpr}

	var moreErrs []error
	if aexpr.Fun, moreErrs = checkExpr(ctx, callExpr.Fun, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	for i := range callExpr.Args {
		if aexpr.Args[i], moreErrs = checkExpr(ctx, callExpr.Args[i], env); moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
	}
	return aexpr, errs
}
