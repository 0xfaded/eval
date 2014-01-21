package eval

import (
	"reflect"
)

func evalStarExpr(ctx *Ctx, starExpr *StarExpr, env *Env) (reflect.Value, error) {
	if v, _, err := EvalExpr(ctx, starExpr, env); err != nil {
		return reflect.Value{}, err
	} else {
		return (*v)[0].Elem(), nil
	}
}
