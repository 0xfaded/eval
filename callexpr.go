package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
	"go/token"
)

func evalCallExpr(call *ast.CallExpr, env *Env) ([]reflect.Value, bool, error) {
	if t, err := evalType(call.Fun, env); err == nil {
		if v, typed, err := evalCallTypeExpr(t, call, env); err != nil {
			return []reflect.Value{}, false, err
		} else {
			return []reflect.Value{v}, typed, nil
		}
	} else if fun, _, err := EvalExpr(call.Fun, env); err == nil {
		return evalCallFunExpr(fun[0], call, env)
	} else {
		return []reflect.Value{}, false, err
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
	} else if cast, err := assignableValue(arg[0], t, typed); err != nil {
		return r, false, err
	} else {
		return cast, true, nil
	}
}

func evalCallFunExpr(fun reflect.Value, call *ast.CallExpr, env *Env) ([]reflect.Value, bool, error) {
	var err error
	var v, out []reflect.Value
	var typed bool
	if v, typed, err = EvalExpr(call.Fun, env); err != nil {
		return []reflect.Value{}, false, err
	}
	if v[0].Kind() != reflect.Func {
		return []reflect.Value{}, false, errors.New(fmt.Sprintf("Cannot call %v", v[0]))
	}
	builtin := !typed

	// Special case handling doesn't play well with nil Args
	ftype := v[0].Type()
	if call.Args == nil {
		if ftype.NumIn() == 0 {
			out = v[0].Call([]reflect.Value{})
			return out, true, nil
		} else {
			return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], 0}
		}
	}

	args := make([][]reflect.Value, len(call.Args))
	atyped := make([]bool, len(call.Args))

	// Evaluate each arg
	for i := range call.Args {
		if args[i], atyped[i], err = EvalExpr(call.Args[i], env); err != nil {
			return []reflect.Value{}, false, err
		}
	}

	_, firstArgIsFun := call.Args[0].(*ast.CallExpr)
	// Special case for f(g()), where g may return multiple values
	if len(call.Args) == 1 && firstArgIsFun {
		splat := make([][]reflect.Value, len(args[0]))
		atyped = make([]bool, len(args[0]))
		for i := range args[0] {
			splat[i] = []reflect.Value{args[0][i]}
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
			if len(args[i]) > 1 {
				return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
			}
			in[i] = args[i][0]
			intyped[i] = atyped[i]
		}
	} else if ftype.IsVariadic() && actualNumIn-1 <= len(call.Args) {
		var i int
		for i = 0; i < len(in)-1; i += 1 {
			if len(args[i]) > 1 {
				return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
			}
			in[i] = args[i][0]
			intyped[i] = atyped[i]
		}
		if i == len(call.Args)-1 && call.Ellipsis != token.NoPos {
			// Assert this indeed is the ellipsis
			_ = call.Args[i].(*ast.Ellipsis)
			in[i], err = makeSliceWithValues(args[i], ftype.In(i))
			intyped[i] = true
			if err != nil {
				return []reflect.Value{}, false, ErrBadFunArgument{v[0], i, in[i]}
			}
		} else if i <= len(call.Args) && call.Ellipsis == token.NoPos {
			remainingArgs := len(call.Args) - actualNumIn + 1
			in[i] = reflect.MakeSlice(ftype.In(i), remainingArgs, remainingArgs)

			intyped[i] = true
			etype := in[i].Type().Elem()
			for j := i; j < len(call.Args); j += 1 {
				if len(args[j]) > 1 {
					return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
				} else if arg, err := assignableValue(args[j][0], etype, atyped[j]); err != nil {
					return []reflect.Value{}, false, ErrBadFunArgument{v[0], j, args[j][0]}
				} else {
					in[i].Index(j-i).Set(arg)
				}
			}
		} else {
			return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], len(call.Args)}
		}
	} else {
		return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], len(call.Args)}
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
				return []reflect.Value{}, false, ErrBadFunArgument{v[0], i, in[i]}
			} else {
				in[i] = checked
			}
		}
	}

	if ftype.IsVariadic() {
		out = v[0].CallSlice(in)
	} else {
		out = v[0].Call(in)
	}

	if builtin {
		otyped := out[1].Bool()
		var err error = nil
		if !out[2].IsNil() {
			err = out[2].Interface().(error)
		}
		// Unwrap the Value of a Value
		out = []reflect.Value{out[0].Interface().(reflect.Value)}
		return out, otyped, err
	} else {
		return out, true, nil
	}
}
