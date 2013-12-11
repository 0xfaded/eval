package interactive

import (
	"go/ast"
)

func checkBinaryExpr(ctx *Ctx, binary *ast.BinaryExpr, env *Env) (aexpr *BinaryExpr, errs []error) {
	aexpr = &BinaryExpr{BinaryExpr: binary}

	var moreErrs []error
	if aexpr.X, moreErrs = checkExpr(ctx, binary.X, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	if aexpr.Y, moreErrs = checkExpr(ctx, binary.Y, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	return aexpr, errs
}
