package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/token"
)

func evalCallExpr(ctx *Ctx, call *CallExpr, env *Env) (*[]reflect.Value, bool, error) {
	if t, err := evalType(ctx, call.Fun.(Expr), env); err == nil {
		if v, typed, err := evalCallTypeExpr(ctx, t, call, env); err != nil {
			return nil, false, err
		} else {
			ret := []reflect.Value{v}
			return &ret, typed, nil
		}
	} else if fun, _, err := EvalExpr(ctx, call.Fun.(Expr), env); err == nil {
		return evalCallFunExpr(ctx, (*fun)[0], call, env)
	} else {
		return nil, false, err
	}
}

func evalCallTypeExpr(ctx *Ctx, t reflect.Type, call *CallExpr, env *Env) (reflect.Value, bool, error) {
	var r reflect.Value
	if call.Args == nil {
		return r, false, errors.New(fmt.Sprintf("missing argument to conversion to %v", t))
	} else if len(call.Args) > 1 {
		return r, false, errors.New(fmt.Sprintf("too many arguments to conversion to %v", t))
	} else if arg, typed, err := EvalExpr(ctx, call.Args[0].(Expr), env); err != nil {
		return r, false, err
	} else if cast, err := assignableValue((*arg)[0], t, typed); err != nil {
		return r, false, err
	} else {
		return cast, true, nil
	}
}

func evalCallFunExpr(ctx *Ctx, fun reflect.Value, call *CallExpr, env *Env) (*[]reflect.Value, bool, error) {
	var err error
	var v *[]reflect.Value
	var typed bool
	if v, typed, err = EvalExpr(ctx, call.Fun.(Expr), env); v == nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	if (*v)[0].Kind() != reflect.Func {
		return nil, false, errors.New(fmt.Sprintf("Cannot call %v", (*v)[0]))
	}
	builtin := !typed

	// Special case handling doesn't play well with nil Args
	ftype := (*v)[0].Type()
	if call.Args == nil {
		if ftype.NumIn() == 0 {
			out := (*v)[0].Call([]reflect.Value{})
			return &out, true, nil
		} else {
			return nil, false, ErrWrongNumberOfArgs{(*v)[0], 0}
		}
	}

	args := make([]*[]reflect.Value, len(call.Args))
	atyped := make([]bool, len(call.Args))

	// Evaluate each arg
	for i := range call.Args {
		var err error
		args[i], atyped[i], err = EvalExpr(ctx, call.Args[i].(Expr), env)
		if err != nil {
			return nil, false, err
		}
	}

	_, firstArgIsFun := call.Args[0].(*CallExpr)
	// Special case for f(g()), where g may return multiple values
	wasSplat := false
	if len(call.Args) == 1 && firstArgIsFun {
		arg := *(args[0])

		// g := func() {}; h := func() {}; _ = g(h()) is illegal
		if len(arg) == 0 {
			return nil, false, ErrMissingValue{at(ctx, call.Args[0])}
		}

		splat := make([]*[]reflect.Value, len(arg))
		atyped = make([]bool, len(arg))
		for i := range arg {
			splat[i] = &[]reflect.Value{arg[i]}
			atyped[i] = true
		}
		args = splat
		wasSplat = true
	}

	// Parse args into a slice suitable for calling the function
	actualNumIn := ftype.NumIn()
	if builtin {
		// See builtinFuncs comment
		actualNumIn /= 2
	}

	in := make([]reflect.Value, actualNumIn)
	intyped := make([]bool, actualNumIn)

	if !ftype.IsVariadic() && len(args) == actualNumIn {
		// Standard call
		for i := range in {
			var arg reflect.Value;
			var err error

			// In the case of a splat, we cannot possibly be dealing with multi values here
			if wasSplat {
				arg = (*args[i])[0]
			} else if arg, err = expectSingleValue(ctx, *(args[i]), call.Args[i]); err != nil {
				return nil, false, err
			}
			in[i] = arg
			intyped[i] = atyped[i]
		}
	} else if ftype.IsVariadic() && actualNumIn-1 <= len(args) {
		// Varadic call
		var i int
		for i = 0; i < len(in)-1; i += 1 {
			var arg reflect.Value;
			var err error
			if wasSplat {
				arg = (*args[i])[0]
			} else if arg, err = expectSingleValue(ctx, *(args[i]), call.Args[i]); err != nil {
				return nil, false, err
			}
			in[i] = arg
			intyped[i] = atyped[i]
		}
		if i == len(args)-1 && call.Ellipsis != token.NoPos {
			// Call of form f(first, second, ...others)
			arg := *(args[i])
			// Assert this indeed is the ellipsis
			_ = call.Args[i].(*Ellipsis)
			in[i], err = makeSliceWithValues(arg, ftype.In(i))
			intyped[i] = true
			if err != nil {
				return nil, false, ErrBadFunArgument{(*v)[0], i, in[i]}
			}
		} else if i <= len(args) && call.Ellipsis == token.NoPos {
			// Call of form f(first, second, third, fourth and so on)
			remainingArgs := len(args) - actualNumIn + 1
			in[i] = reflect.MakeSlice(ftype.In(i), remainingArgs, remainingArgs)

			intyped[i] = true
			etype := in[i].Type().Elem()
			for j := i; j < len(args); j += 1 {
				if arg, err := expectSingleValue(ctx, *(args[j]), call.Args[j]); err != nil {
					return nil, false, err
				} else if arg, err := assignableValue(arg, etype, atyped[j]); err != nil {
					return nil, false, ErrBadFunArgument{(*v)[0], j, arg}
				} else {
					in[i].Index(j-i).Set(arg)
				}
			}
		} else {
			return nil, false, ErrWrongNumberOfArgs{(*v)[0], len(call.Args)}
		}
	} else {
		return nil, false, ErrWrongNumberOfArgs{(*v)[0], len(call.Args)}
	}

	if builtin {
		// Builtin functions take and return raw values as well as typing information
		bin := make([]reflect.Value, len(in) * 2)
		for i := range in {
			bin[i] = reflect.ValueOf(in[i])
			bin[i+len(in)] = reflect.ValueOf(intyped[i])
		}
		in = bin
	} else {
		// Check argument types
		for i := range in {
			var checked reflect.Value
			if checked, err = assignableValue(in[i], ftype.In(i), intyped[i]); err != nil {
				return nil, false, ErrBadFunArgument{(*v)[0], i, in[i]}
			} else {
				in[i] = checked
			}
		}
	}

	var ret []reflect.Value
	if ftype.IsVariadic() {
		ret = (*v)[0].CallSlice(in)
	} else {
		ret = (*v)[0].Call(in)
	}
	out := &ret

	if builtin {
		otyped := ret[1].Bool()
		var err error = nil
		if !ret[2].IsNil() {
			err = ret[2].Interface().(error)
		}
		// Unwrap the Value of a Value
		out = &[]reflect.Value{ret[0].Interface().(reflect.Value)}
		return out, otyped, err
	} else {
		return out, true, nil
	}
}
