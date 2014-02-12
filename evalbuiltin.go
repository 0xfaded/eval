package eval

import (
	"reflect"
)

func evalCallBuiltinExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	ident := call.Fun.(*Ident)
	switch ident.Name {
	case "complex":
		return evalBuiltinComplexExpr(ctx, call, env)
	case "real":
		return evalBuiltinRealExpr(ctx, call, env)
	case "imag":
		return evalBuiltinImagExpr(ctx, call, env)
	case "new":
		return evalBuiltinNewExpr(ctx, call, env)
	case "make":
		return evalBuiltinMakeExpr(ctx, call, env)
	case "len":
		return evalBuiltinLenExpr(ctx, call, env)
	case "cap":
		return evalBuiltinCapExpr(ctx, call, env)
	case "append":
		return evalBuiltinAppendExpr(ctx, call, env)
	case "copy":
		return evalBuiltinCopyExpr(ctx, call, env)
	case "delete":
		return evalBuiltinDeleteExpr(ctx, call, env)
	case "panic":
		return evalBuiltinPanicExpr(ctx, call, env)
	default:
		panic("eval: unimplemented builtin " + ident.Name)
	}
}

func evalBuiltinComplexExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	var err error

	resT := call.KnownType()[0]
	argT := knownType{comprisingFloatType(resT)}

	var re, im []reflect.Value
	if re, err = evalTypedExpr(ctx, call.Args[0].(Expr), argT, env); err != nil {
		return nil, err
	} else if im, err = evalTypedExpr(ctx, call.Args[1].(Expr), argT, env); err != nil {
		return nil, err
	}
	cplx := builtinComplex(re[0], im[0])
	return []reflect.Value{cplx}, nil
}

func evalBuiltinRealExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	var err error

	resT := call.KnownType()[0]
	argT := knownType{comprisingFloatType(resT)}

	var cplx []reflect.Value
	if cplx, err = evalTypedExpr(ctx, call.Args[0].(Expr), argT, env); err != nil {
		return nil, err
	}
	re := builtinReal(cplx[0])
	return []reflect.Value{re}, nil
}

func evalBuiltinImagExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	var err error

	resT := call.KnownType()[0]
	argT := knownType{comprisingFloatType(resT)}

	var cplx []reflect.Value
	if cplx, err = evalTypedExpr(ctx, call.Args[0].(Expr), argT, env); err != nil {
		return nil, err
	}
	im := builtinImag(cplx[0])
	return []reflect.Value{im}, nil
}

func evalBuiltinNewExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	resT := call.KnownType()[0]
	ptr := builtinNew(resT.Elem())
	return []reflect.Value{ptr}, nil
}

func evalBuiltinMakeExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	resT := call.KnownType()[0]
	length, capacity := 0, 0
	var err error
	if len(call.Args) > 1 {
		if length, err = evalInteger(ctx, call.Args[1].(Expr), env); err != nil {
			return nil, err
		}
	}
	if len(call.Args) > 2 {
		if capacity, err = evalInteger(ctx, call.Args[2].(Expr), env); err != nil {
			return nil, err
		}
	}
	var res reflect.Value
	switch resT.Kind() {
	case reflect.Slice:
		res = reflect.MakeSlice(resT, length, capacity)
	case reflect.Map:
		res = reflect.MakeMap(resT)
	case reflect.Chan:
		res = reflect.MakeChan(resT, length)
	default:
		panic(dytc("make(bad type)"))
	}
	return []reflect.Value{res}, nil
}

func evalBuiltinLenExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	if arg0, _, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
		return nil, err
	} else {
		length := builtinLen((*arg0)[0])
		return []reflect.Value{length}, nil
	}
}

func evalBuiltinCapExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	if arg0, _, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
		return nil, err
	} else {
		capacity := builtinCap((*arg0)[0])
		return []reflect.Value{capacity}, nil
	}
}

func evalBuiltinAppendExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	sliceT := call.KnownType()
	head, err := evalTypedExpr(ctx, call.Args[0].(Expr), sliceT, env)
	if err != nil {
		return nil, err
	}
	var tail reflect.Value
	if call.argNEllipsis {
		xs, _, err := EvalExpr(ctx, call.Args[1].(Expr), env)
		if err != nil {
			return nil, err
		}
		tail = (*xs)[0]
	} else {
		numXs := len(call.Args) - 1
		tail = reflect.MakeSlice(sliceT[0], numXs, numXs)
		xT := knownType{sliceT[0].Elem()}
		for i := 1; i < len(call.Args); i += 1 {
			if x, err := evalTypedExpr(ctx, call.Args[i].(Expr), xT, env); err != nil {
				return nil, err
			} else {
				tail.Index(i-1).Set(x[0])
			}
		}
	}

	res := builtinAppend(head[0], tail)
	return []reflect.Value{res}, nil
}

func evalBuiltinCopyExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	if x, _, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
		return nil, err
	} else if y, _, err := EvalExpr(ctx, call.Args[1].(Expr), env); err != nil {
		return nil, err
	} else {
		n := builtinCopy((*x)[0], (*y)[0])
		return []reflect.Value{n}, nil
	}
}

func evalBuiltinDeleteExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	m := call.Args[0].(Expr)
	mT := m.KnownType()[0]
	if x, _, err := EvalExpr(ctx, m, env); err != nil {
		return nil, err
	} else if y, err := evalTypedExpr(ctx, call.Args[1].(Expr), knownType{mT.Key()}, env); err != nil {
		return nil, err
	} else {
		builtinDelete((*x)[0], y[0])
		return []reflect.Value{}, nil
	}
}

func evalBuiltinPanicExpr(ctx *Ctx, call *CallExpr, env *Env) ([]reflect.Value, error) {
	if arg0, err := evalTypedExpr(ctx, call.Args[0].(Expr), knownType{emptyInterface}, env); err != nil {
		return nil, err
	} else {
		err := builtinPanic(arg0[0])
		return []reflect.Value{}, err
	}
}

