package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
	"go/token"
)

func evalCallExpr(call *ast.CallExpr, env *Env) (*[]reflect.Value, bool, error) {
	if t, err := evalType(call.Fun, env); err == nil {
		if v, typed, err := evalCallTypeExpr(t, call, env); err != nil {
			return nil, false, err
		} else {
			ret := []reflect.Value{v}
			return &ret, typed, nil
		}
	} else if fun, _, err := EvalExpr(call.Fun, env); err == nil {
		return evalCallFunExpr((*fun)[0], call, env)
	} else {
		return nil, false, err
	}
}

func evalCallTypeExpr(t reflect.Type, call *ast.CallExpr, env *Env) (reflect.Value, bool, error) {
	var r reflect.Value
	if call.Args == nil {
		return r, false, errors.New(fmt.Sprintf("missing argument to conversion to %v", t))
	} else if len(call.Args) > 1 {
		return r, false, errors.New(fmt.Sprintf("too many arguments to conversion to %v", t))
	} else if arg, typed, err := EvalExpr(call.Args[0], env); err != nil {
		return r, false, err
	} else if cast, err := assignableValue((*arg)[0], t, typed); err != nil {
		return r, false, err
	} else {
		return cast, true, nil
	}
}

func evalCallFunExpr(fun reflect.Value, call *ast.CallExpr, env *Env) (*[]reflect.Value, bool, error) {
	var err error
	var v *[]reflect.Value
	var typed bool
	if v, typed, err = EvalExpr(call.Fun, env); v == nil {
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
		args[i], atyped[i], err = EvalExpr(call.Args[i], env)
		if err != nil {
			return nil, false, err
		}
	}

	_, firstArgIsFun := call.Args[0].(*ast.CallExpr)
	// Special case for f(g()), where g may return multiple values
	if len(call.Args) == 1 && firstArgIsFun {
		arg   := *(args[0])
		splat := make([]*[]reflect.Value, len(arg))
		atyped = make([]bool, len(arg))
		for i := range arg {
			splat[i] = &[]reflect.Value{arg[i]}
			atyped[i] = true
		}
		args = splat
	}

	// Parse args into a slice suitable for calling the function
	actualNumIn := ftype.NumIn()
	if builtin {
		// See builtinFuncs comment
		actualNumIn /= 2
	}

	in := make([]reflect.Value, actualNumIn)
	intyped := make([]bool, actualNumIn)

	if !ftype.IsVariadic() && len(call.Args) == actualNumIn {
		for i := range in {
			arg := *(args[i])
			if len(arg) > 1 {
				return nil, false, ErrMultiInSingleContext{arg}
			}
			in[i] = arg[0]
			intyped[i] = atyped[i]
		}
	} else if ftype.IsVariadic() && actualNumIn-1 <= len(call.Args) {
		var i int
		for i = 0; i < len(in)-1; i += 1 {
			arg := *(args[i])
			if len(arg) > 1 {
				return nil, false, ErrMultiInSingleContext{arg}
			}
			in[i] = arg[0]
			intyped[i] = atyped[i]
		}
		if i == len(call.Args)-1 && call.Ellipsis != token.NoPos {
			arg := *(args[i])
			// Assert this indeed is the ellipsis
			_ = call.Args[i].(*ast.Ellipsis)
			in[i], err = makeSliceWithValues(arg, ftype.In(i))
			intyped[i] = true
			if err != nil {
				return nil, false, ErrBadFunArgument{(*v)[0], i, in[i]}
			}
		} else if i <= len(call.Args) && call.Ellipsis == token.NoPos {
			remainingArgs := len(call.Args) - actualNumIn + 1
			in[i] = reflect.MakeSlice(ftype.In(i), remainingArgs, remainingArgs)

			intyped[i] = true
			etype := in[i].Type().Elem()
			for j := i; j < len(call.Args); j += 1 {
				arg := *(args[j])
				if len(arg) > 1 {
					return nil, false, ErrMultiInSingleContext{arg}
				} else if arg, err := assignableValue(arg[0], etype, atyped[j]); err != nil {
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
