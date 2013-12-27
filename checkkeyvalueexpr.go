package eval

import (
	"go/ast"
)

func checkKeyValueExpr(ctx *Ctx, keyValue *ast.KeyValueExpr, env *Env) (aexpr *KeyValueExpr, errs []error) {
	aexpr = &KeyValueExpr{KeyValueExpr: keyValue}

	var moreErrs []error
	if aexpr.Key, moreErrs = CheckExpr(ctx, keyValue.Key, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	if aexpr.Value, moreErrs = CheckExpr(ctx, keyValue.Value, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	return aexpr, errs
}
