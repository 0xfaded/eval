package eval

import (
	"reflect"
)

func evalCallExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	if call.IsConst() {
		return []reflect.Value{call.Const()}, nil
	} else if call.isBuiltin {
		return evalCallBuiltinExpr(ctx, call, env)
	} else if call.isTypeConversion {
		return evalCallTypeExpr(ctx, call, env)
	} else {
		return evalCallFunExpr(ctx, call, env)
	}
}

func evalCallTypeExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	// Arg0 cannot be const, otherwise the entire expression would be const and
	// evalCallExpr will have already returned.
	if arg, _, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
		return nil, nil
	} else {
		cast, _ := assignableValue((*arg)[0], call.KnownType()[0], true)
		return []reflect.Value{cast}, nil
	}
}

func evalCallFunExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	v, _, err := EvalExpr(ctx, call.Fun.(Expr), env)
	if err != nil {
		return nil, err
	}

	fun := (*v)[0]
	ft := fun.Type()
	numIn := ft.NumIn()

	// Evaluate arguments
	args := make([]reflect.Value, len(call.Args))
	if call.arg0MultiValued {
		// TODO clean this up once EvalExpr type changes
		if argp, _, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
			return nil, err
		} else {
			args = *argp
		}
	} else if len(args) != 0 {
		var i int
		for i = 0; i < numIn - 1; i += 1 {
			arg := call.Args[i].(Expr)
			argType := knownType{ft.In(i)}
			if argV, err := evalTypedExpr(ctx, arg, argType, env); err != nil {
				return nil, err
			} else {
				args[i] = argV[0]
			}
		}
		argNT := ft.In(i)
		if ft.IsVariadic() && !call.argNEllipsis {
			argNT = argNT.Elem()
		}
		argNKnownType := knownType{argNT}
		for ; i < len(call.Args); i += 1 {
			arg := call.Args[i].(Expr)
			if argV, err := evalTypedExpr(ctx, arg, argNKnownType, env); err != nil {
				return nil, err
			} else {
				args[i] = argV[0]
			}
		}
	}

	var out []reflect.Value
	if call.argNEllipsis {
		out = fun.CallSlice(args)
	} else {
		out = fun.Call(args)
	}
	return out, nil
}
