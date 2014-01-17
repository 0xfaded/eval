package eval

import (
	"reflect"
	"go/ast"
)

func checkCompositeLit(ctx *Ctx, lit *ast.CompositeLit, env *Env) (aexpr *CompositeLit, errs []error) {
	aexpr = &CompositeLit{CompositeLit: lit}

	var moreErrs []error
        // TODO confirm the type actually exists
        aexpr.Type = lit.Type
	t, _ := evalType(ctx, aexpr.Type, env)

	if t.Kind() == reflect.Struct {
		for i := range lit.Elts {
			// Don't typecheck struct keys.
			if kv, ok := aexpr.Elts[i].(*ast.KeyValueExpr); ok {
				if kv.Value, moreErrs = CheckExpr(ctx, kv.Value, env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
				aexpr.Elts[i] = &KeyValueExpr{KeyValueExpr: kv}
			} else {
				if aexpr.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
		}
	} else {
		for i := range lit.Elts {
			if aexpr.Elts[i], moreErrs = CheckExpr(ctx, lit.Elts[i], env); moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
		}
	}
	return aexpr, errs
}
