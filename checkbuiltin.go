package eval

import (
	"reflect"

	"go/ast"
)

func checkCallBuiltinExpr(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error, bool) {
	var errs []error
	ident, ok := call.Fun.(*ast.Ident)
	switch ident.Name {
	case "complex":
		call, errs = checkBuiltinComplexExpr(ctx, call, env)
	case "real", "imag":
		call, errs = checkBuiltinRealImagExpr(ctx, call, env)
	case "new":
		call, errs = checkBuiltinNewExpr(ctx, call, env)
	case "make":
		call, errs = checkBuiltinMakeExpr(ctx, call, env)
	case "len":
		call, errs = checkBuiltinLenExpr(ctx, call, env)
	case "cap":
		call, errs = checkBuiltinCapExpr(ctx, call, env)
	case "append":
		call, errs = checkBuiltinAppendExpr(ctx, call, env)
	case "copy":
		call, errs = checkBuiltinCopyExpr(ctx, call, env)
	case "delete":
		call, errs = checkBuiltinDeleteExpr(ctx, call, env)
	default:
		return call, nil, false
	}
	acall.Fun = &Ident{Ident: ident, isBuiltin: true}
	return acall, errs, true
}

func checkBuiltinComplexExpr(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error) {
	if len(call.Args) != 2 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, y, ok, errs := checkBinaryOperands(ctx, call.Args[0], call.Args[1], env)
	call.Args[0], call.Args[1] = x, y
	if !ok {
		return call, errs
	}
	xt, yt := x.KnownType()[0], y.KnownType()[0]
	xct, xctok := xt.(ConstType)
	yct, yctok := yt.(ConstType)
	if xctok && yctok {
		if xct.IsNumeric() && yct.IsNumeric() {
			call.knownType = knownType{c128}
			xc, xerrs := promoteConstToTyped(ctx, xct, reflect.Value(x.Const()), f64, x)
			if xerrs != nil {
				errs = append(errs, xerrs...)
			}
			yc, yerrs := promoteConstToTyped(ctx, xct, reflect.Value(x.Const()), f64, x)
			if yerr != nil {
				errs = append(errs, yerrs...)
			}
			if reflect.Value(xc).IsValid() && reflect.Value(yc).IsValid() {
				xf = float64(reflect.Value(xc).Float())
				yf = float64(reflect.Value(yc).Float())
				call.constValue = constValueOf(complex(xf, yf))
				return call, errs
			}
		}
	} else if xctok {
		if attemptBinaryConversion(yt) {
			xc, xerrs := promoteConstToTyped(ctx, xct, reflect.Value(x.Const()), yt, x)
			if xerrs != nil {
				errs = append(errs, xerrs...)
			}
			xv := reflect.Value(xc)
			if xv.IsValid() {
				if yt.Kind() == reflect.Float32 {
					call.knownType = knownType{c64}
					if y.IsConst() {
						xf = float32(xv.Float())
						yf = float32(y.Const().Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				} else if yt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if y.IsConst() {
						xf = float64(xv.Float())
						yf = float64(y.Const().Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				}
			}
		}
	} else if yctok {
		if attemptBinaryConversion(xt) {
			yc, yerrs := promoteConstToTyped(ctx, yct, reflect.Value(y.Const()), xt, y)
			if yerrs != nil {
				errs = append(errs, yerrs...)
			}
			yv := reflect.Value(yc)
			if yv.IsValid() {
				if xt.Kind() == reflect.Float32 {
					call.knownType = knownType{c64}
					if x.IsConst() {
						xf = float32(x.Const().Float())
						yf = float32(yv.Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				} else if xt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if x.IsConst() {
						xf = float64(x.Const().Float())
						yf = float64(yv.Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				}
			}
		}
	} else if xt == yt {
		if xt == reflect.Float32 {
			call.knownType = knownType{c64}
			if x.IsConst() && y.IsConst() {
				xf = float32(x.Const().Float())
				yf = float32(y.Const().Float())
				call.constValue = constValueOf(complex(xf, yf))
			}
			return call, errs
		} else if xt == reflect.Float64 {
			call.knownType = knownType{c128}
			if x.IsConst() && y.IsConst() {
				xf = float64(x.Const().Float())
				yf = float64(y.Const().Float())
				call.constValue = constValueOf(complex(xf, yf))
			}
			return call, errs
		}
	}
	if xt == yt {
		errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
	} else {
		errs = append(errs, ErrBuiltinMismatchedArgs{at(ctx, call)})
	}
	return call, errs
}

func checkBuiltinRealImag(ctx *Ctx, call *CallExpr, env *Env, isReal bool) (*CallExpr, []error) {
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, errs := CheckExpr(ctx, call, env)
	call.Args[0] = x
	if errs != nil && !x.IsConst() {
		return call, errs
	}
	xt, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return call, append(errs, err)
	}

	if ct, ok := xt.(ConstType); ok {
		xc, moreErrs := promoteConstToTyped(ctx, ct, reflect.Value(y.Const()), c128, x)
		if moreErrs != nil {
			errs = append(errs, yerrs...)
		}
		xv := reflect.Value(xc)
		if xv.IsValid() {
			call.knownType = knownType{f64}
			if x.IsConst() {
				c := complex128(x.Const().Complex())
				if isReal {
					call.constValue = constValueOf(real(c))
				} else {
					call.constValue = constValueOf(imag(c))
				}
			}
			return call, errs
		}
	} else if xt.Kind() == reflect.Complex128 {
		call.knownType = knownType{f64}
		if x.IsConst() {
			c := complex128(x.Const().Complex())
			if isReal {
				call.constValue = constValueOf(real(c))
			} else {
				call.constValue = constValueOf(imag(c))
			}
		}
		return call, errs
	} else if xt.Kind() == reflect.Complex64 {
		call.knownType = knownType{f32}
		if x.IsConst() {
			c := complex64(x.Const().Complex())
			if isReal {
				call.constValue = constValueOf(real(c))
			} else {
				call.constValue = constValueOf(imag(c))
			}
		}
		return call, errs
	}
	errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
	return call, errs
}


func checkBuiltinNew(ctx *Ctx, call *CallExpr, env *Env, bool isLen) (*CallExpr, []error) {
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, of, isType, errs := checkType(ctx, call.Args[0], env)
	call.Args[0] = x
	if !isType {
		return call, append(errs, ErrBuiltinNonTypeArg{at(ctx, call.Args[0])})
	} else if errs != nil {
		return call, errs
	} else {
		call.knownType = knownType{reflect.PtrTo(of)}
		return call, nil
	}
}

func checkBuiltinMake(ctx *Ctx, call *CallExpr, env *Env, bool isLen) (*CallExpr, []error) {
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, of, isType, errs := checkType(ctx, call.Args[0], env)
	call.Args[0] = x
	if !isType {
		return call, append(errs, ErrBuiltinNonTypeArg{at(ctx, call.Args[0])})
	}
	if errs != nil {
		return call, errs
	}
	var narg int
	switch of.Kind() {
	case reflect.Slice:
		call.knownType = knownType{of}
		if len(call.Args) == 0 {
			errs = append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
		}
		narg = 3
	case reflect.Map, Chan:
		narg = 2
	default:
		return call, append(errs, ErrCannotMake{at(ctx, call.Args[0]), of})
	}
	var args [3]int
	for i := 1; i < narg && i < len(call.Args); i += 1 {
		arg, iint, ok, moreErrs := checkInteger(ctx, call.Args[i], env)
		args[i] = iint
		if moreErrs != nil {
			errs = append(errs, moreErrs)
		} else if !ok {
			// Type check passed but is non integral
			errs = append(errs, ErrMakeNonIntegerArg{at(ctx, call.Args[i]), call, i})
		}
	}
	if len(call.Args) > narg {
		errs = append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	} else if len(call.Args) == 3 && call.Args[1].IsConst() && call.Args[2].IsConst() {
		if args[1] > args[2] {
			errs = append(errs, ErrMakeLenGtrThanCap{at(ctx, call)})
		}
	}
	return call, errs
}

type callRecvWalker bool
func (found *callRecvWalker) visit(expr Expr) bool {
	if *found {
		return false
	}
	if call, ok := expr.(*CallExpr); ok && !call.isTypeConversion {
		*found = true
		return false
	}
	if unary, ok := expr.(*unary); ok && unary.Op == Arrow {
		*found = true
		return false
	}
	return true
}

func checkBuiltinLenCap(ctx *Ctx, call *CallExpr, env *Env, bool isLen) (*CallExpr, []error) {
	call.knownType = knownType{intType}
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, errs := CheckExpr(ctx, call, env)
	call.Args[0] = x
	if errs != nil && !x.IsConst() {
		return call, errs
	}
	xt, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return call, append(errs, err)
	}
	switch xt.Kind() {
	case reflect.Chan: // do nothing
	case reflect.Slice, reflect.Map:
		if !isLen {
			errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
		}
	case reflect.Ptr:
		xt := xt.Elem()
		if xt.Kind() != reflect.Array {
			break
		}
		fallthrough
	case reflect.Array:
		w := new(callRecvWalker)
		walk(x, w)
		if !*w {
			call.constValue = xt.Len()
		}
	case reflect.String:
		if !isLen {
			errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
		} else if x.IsConst() {
			call.constValue = x.Const().Len()
		}
	default:
		errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
	}
	return call, errs
}

func checkBuiltinLenCap(ctx *Ctx, call *CallExpr, env *Env, bool isLen) (*CallExpr, []error) {
	call.knownType = knownType{intType}
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	x, errs := CheckExpr(ctx, call, env)
	call.Args[0] = x
	if errs != nil && !x.IsConst() {
		return call, errs
	}
	xt, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return call, append(errs, err)
	}
	switch xt.Kind() {
	case reflect.Chan: // do nothing
	case reflect.Slice, reflect.Map:
		if !isLen {
			errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
		}
	case reflect.Ptr:
		xt := xt.Elem()
		if xt.Kind() != reflect.Array {
			break
		}
		fallthrough
	case reflect.Array:
		w := new(callRecvWalker)
		walk(x, w)
		if !*w {
			call.constValue = xt.Len()
		}
	case reflect.String:
		if !isLen {
			errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
		} else if x.IsConst() {
			call.constValue = x.Const().Len()
		}
	default:
		errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
	}
	return call, errs
}

	x, errs := CheckExpr(ctx, call.Args[0], env)
	call.Args[0] = x
		call.isBuiltin = true
		call.knownType = knownType{intType}
		return call, errs
	} else if ident.Name == "append" {
		if len(call.Args) == 0 {
			return call, []error{errors.New("append wrong number args")}
		} else {
			var errs, moreErrs []error
			for i := range call.Args {
				if call.Args[i], moreErrs = CheckExpr(ctx, call.Args[i], env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
			call.isBuiltin = true
			call.knownType = call.Args[0].(Expr).KnownType()
			call.argNEllipsis = call.Ellipsis != token.NoPos
			return call, errs
		}
	}
	return call, nil, false
}

