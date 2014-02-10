package eval

import (
	"reflect"

	"go/ast"
	"go/token"
)

func checkCallBuiltinExpr(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error, bool) {
	var errs []error
	ident, ok := call.Fun.(*ast.Ident)
	if !ok {
		return call, nil, false
	}
	switch ident.Name {
	case "complex":
		call, errs = checkBuiltinComplex(ctx, call, env)
	case "real":
		call, errs = checkBuiltinRealImag(ctx, call, env, true)
	case "imag":
		call, errs = checkBuiltinRealImag(ctx, call, env, false)
	case "new":
		call, errs = checkBuiltinNew(ctx, call, env)
	case "make":
		call, errs = checkBuiltinMake(ctx, call, env)
	case "len":
		call, errs = checkBuiltinLenCap(ctx, call, env, true)
	case "cap":
		call, errs = checkBuiltinLenCap(ctx, call, env, false)
	case "append":
		call, errs = checkBuiltinAppend(ctx, call, env)
		/*
	case "copy":
		call, errs = checkBuiltinCopyExpr(ctx, call, env)
	case "delete":
		call, errs = checkBuiltinDeleteExpr(ctx, call, env)
		*/
	default:
		return call, nil, false
	}
	call.Fun = &Ident{Ident: ident}
	call.isBuiltin = true
	return call, errs, true
}

func checkBuiltinComplex(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error) {
	var errs []error
	if call.argNEllipsis = call.Ellipsis != token.NoPos; call.argNEllipsis {
		errs = append(errs, ErrBuiltinInvalidEllipsis{at(ctx, call)})
	}
	if len(call.Args) != 2 {
		fakeCheckRemainingArgs(call, 0, env)
		return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	}
	x, y, ok, moreErrs := checkBinaryOperands(ctx, call.Args[0], call.Args[1], env)
	call.Args[0], call.Args[1] = x, y
	if moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	if !ok {
		return call, errs
	}
	xt, yt := x.KnownType()[0], y.KnownType()[0]
	xct, xctok := xt.(ConstType)
	yct, yctok := yt.(ConstType)
	if xctok && yctok {
		if xct.IsNumeric() && yct.IsNumeric() {
			call.knownType = knownType{c128}
			xc, xerrs := promoteConstToTyped(ctx, xct, constValue(x.Const()), f64, x)
			if xerrs != nil {
				errs = append(errs, xerrs...)
			}
			yc, yerrs := promoteConstToTyped(ctx, yct, constValue(y.Const()), f64, y)
			if yerrs != nil {
				errs = append(errs, yerrs...)
			}
			if reflect.Value(xc).IsValid() && reflect.Value(yc).IsValid() {
				xf := float64(reflect.Value(xc).Float())
				yf := float64(reflect.Value(yc).Float())
				call.constValue = constValueOf(complex(xf, yf))
				return call, errs
			}
		}
	} else if xctok {
		if attemptBinaryOpConversion(yt) {
			xc, xerrs := promoteConstToTyped(ctx, xct, constValue(x.Const()), yt, x)
			if xerrs != nil {
				errs = append(errs, xerrs...)
				if xt == ConstNil {
					// No MismatchedTypes error for nils
					return call, errs
				}
			}
			xv := reflect.Value(xc)
			if xv.IsValid() {
				if yt.Kind() == reflect.Float32 {
					call.knownType = knownType{c64}
					if y.IsConst() {
						xf := float32(xv.Float())
						yf := float32(y.Const().Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				} else if yt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if y.IsConst() {
						xf := float64(xv.Float())
						yf := float64(y.Const().Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				}
			}
		} else {
			if xt == ConstNil && isNillable(yt) {
				errs = append(errs, ErrBuiltinWrongArgType{at(ctx, y), call})
				return call, errs
			}
		}
	} else if yctok {
		if attemptBinaryOpConversion(xt) {
			yc, yerrs := promoteConstToTyped(ctx, yct, constValue(y.Const()), xt, y)
			if yerrs != nil {
				errs = append(errs, yerrs...)
				if yt == ConstNil {
					// No MismatchedTypes error for nils
					return call, errs
				}
			} else if yt == ConstNil {
				errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
				return call, errs
			}
			yv := reflect.Value(yc)
			if yv.IsValid() {
				if xt.Kind() == reflect.Float32 {
					call.knownType = knownType{c64}
					if x.IsConst() {
						xf := float32(x.Const().Float())
						yf := float32(yv.Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				} else if xt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if x.IsConst() {
						xf := float64(x.Const().Float())
						yf := float64(yv.Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs
				}
			}
		} else {
			if yt == ConstNil && isNillable(xt) {
				errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
				return call, errs
			}
		}
	} else if xt == yt {
		if xt.Kind() == reflect.Float32 {
			call.knownType = knownType{c64}
			if x.IsConst() && y.IsConst() {
				xf := float32(x.Const().Float())
				yf := float32(y.Const().Float())
				call.constValue = constValueOf(complex(xf, yf))
			}
			return call, errs
		} else if xt.Kind() == reflect.Float64 {
			call.knownType = knownType{c128}
			if x.IsConst() && y.IsConst() {
				xf := float64(x.Const().Float())
				yf := float64(y.Const().Float())
				call.constValue = constValueOf(complex(xf, yf))
			}
			return call, errs
		}
	}
	if unhackType(xt) == unhackType(yt) {
		errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
	} else {
		errs = append(errs, ErrBuiltinMismatchedArgs{at(ctx, call), xt, yt})
	}
	return call, errs
}

func checkBuiltinRealImag(ctx *Ctx, call *CallExpr, env *Env, isReal bool) (*CallExpr, []error) {
	var errs []error
	if call.argNEllipsis = call.Ellipsis != token.NoPos; call.argNEllipsis {
		errs = append(errs, ErrBuiltinInvalidEllipsis{at(ctx, call)})
	}
	if len(call.Args) != 1 {
		fakeCheckRemainingArgs(call, 0, env)
		return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	}
	x, moreErrs := CheckExpr(ctx, call.Args[0], env)
	if moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	call.Args[0] = x
	if moreErrs != nil && !x.IsConst() {
		return call, errs
	}
	xt, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return call, append(errs, err)
	}
	if xt == ConstNil {
		return call, append(errs, ErrUntypedNil{at(ctx, x)})
	}

	if ct, ok := xt.(ConstType); ok {
		if ct == ConstComplex {
			xc, moreErrs := promoteConstToTyped(ctx, ct, constValue(x.Const()), c128, x)
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
			xv := reflect.Value(xc)
			if xv.IsValid() {
				call.knownType = knownType{f64}
				c := complex128(reflect.Value(xc).Complex())
				if isReal {
					call.constValue = constValueOf(real(c))
				} else {
					call.constValue = constValueOf(imag(c))
				}
				return call, errs
			}
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
	errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
	return call, errs
}

func checkBuiltinNew(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error) {
	var errs []error
	if call.argNEllipsis = call.Ellipsis != token.NoPos; call.argNEllipsis {
		errs = append(errs, ErrBuiltinInvalidEllipsis{at(ctx, call)})
	}
	if len(call.Args) == 0 {
		return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	}
	x, of, isType, moreErrs := checkType(ctx, call.Args[0], env)
	if !isType {
		x, moreErrs = CheckExpr(ctx, call.Args[0], env)
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
		call.Args[0] = x
		fakeCheckRemainingArgs(call, 1, env)
		if moreErrs == nil {
			errs = append(errs, ErrBuiltinNonTypeArg{at(ctx, call.Args[0])})
		}
		return call, errs
	} else if len(call.Args) != 1 {
		fakeCheckRemainingArgs(call, 0, env)
		return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	} else if moreErrs != nil {
		return call, append(errs, moreErrs...)
	} else {
		call.Args[0] = x
		call.knownType = knownType{reflect.PtrTo(of)}
		return call, nil
	}
}

func checkBuiltinMake(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error) {
	if len(call.Args) == 0 {
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
	case reflect.Map, reflect.Chan:
		narg = 2
	default:
		return call, append(errs, ErrMakeBadType{at(ctx, call.Args[0]), of})
	}
	var args [3]int
	for i := 1; i < narg && i < len(call.Args); i += 1 {
		arg, iint, ok, moreErrs := checkInteger(ctx, call.Args[i], env)
		call.Args[i] = arg
		args[i] = iint
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		} else if !ok {
			// Type check passed but is non integral
			errs = append(errs, ErrMakeNonIntegerArg{at(ctx, call.Args[i]), call, i})
		}
	}
	if len(call.Args) > narg {
		errs = append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	} else if len(call.Args) == 3 && call.Args[1].(Expr).IsConst() && call.Args[2].(Expr).IsConst() {
		if args[1] > args[2] {
			errs = append(errs, ErrMakeLenGtrThanCap{at(ctx, call), args[1], args[2]})
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
	if unary, ok := expr.(*UnaryExpr); ok && unary.Op == token.ARROW {
		*found = true
		return false
	}
	return true
}

func checkBuiltinLenCap(ctx *Ctx, call *CallExpr, env *Env, isLen bool) (*CallExpr, []error) {
	call.knownType = knownType{intType}
	var errs []error
	if call.argNEllipsis = call.Ellipsis != token.NoPos; call.argNEllipsis {
		errs = append(errs, ErrBuiltinInvalidEllipsis{at(ctx, call)})
	}
	if len(call.Args) != 1 {
		fakeCheckRemainingArgs(call, 0, env)
		return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
	}

	x, moreErrs := CheckExpr(ctx, call.Args[0], env)
	if moreErrs != nil {
		errs = append(errs, moreErrs...)
	}
	call.Args[0] = x
	if errs != nil && !x.IsConst() {
		return call, errs
	}
	xt, err := expectSingleType(ctx, x.KnownType(), x)
	if err != nil {
		return call, append(errs, err)
	}
	if xt == ConstNil {
		return call, append(errs, ErrUntypedNil{at(ctx, x)})
	}
	switch xt.Kind() {
	case reflect.Chan, reflect.Slice: // do nothing
	case reflect.Map:
		if !isLen {
			errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
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
			call.constValue = constValueOf(xt.Len())
		}
	case reflect.String:
		if !isLen {
			errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
		} else if x.IsConst() {
			call.constValue = constValueOf(x.Const().Len())
		}
	default:
		errs = append(errs, ErrBuiltinWrongArgType{at(ctx, x), call})
	}
	return call, errs
}

func checkBuiltinAppend(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error) {
	if len(call.Args) < 1 {
		fakeCheckRemainingArgs(call, 0, env)
		return call, []error{ErrBuiltinWrongNumberOfArgs{at(ctx, call)}}
	}
	slice, errs := CheckExpr(ctx, call.Args[0], env)
	call.Args[0] = slice
	var sliceT reflect.Type
	var isSlice bool
	if errs == nil || slice.IsConst() {
		var err error
		sliceT, err = expectSingleType(ctx, slice.KnownType(), slice)
		if err != nil {
			fakeCheckRemainingArgs(call, 1, env)
			return call, append(errs, err)
		}
		if sliceT != ConstNil {
			isSlice = sliceT.Kind() == reflect.Slice
			call.knownType = knownType{sliceT}
		}
	}
	if call.Ellipsis != token.NoPos {
		call.argNEllipsis = true
		if len(call.Args) == 1 {
			return call, append(errs, ErrAppendFirstArgNotVariadic{at(ctx, call.Args[0])})
		} else if len(call.Args) != 2 {
			return call, append(errs, ErrBuiltinWrongNumberOfArgs{at(ctx, call)})
		} else {
			arg1, moreErrs := CheckExpr(ctx, call.Args[1], env)
			call.Args[1] = arg1
			if moreErrs != nil && !slice.IsConst() {
				return call, append(errs, moreErrs...)
			}
			arg1T, err := expectSingleType(ctx, arg1.KnownType(), arg1)
			if err != nil {
				errs = append(errs, err)
			} else if isSlice {
				if arg1T != sliceT {
					errs = append(errs, ErrBuiltinWrongArgType{at(ctx, arg1), call})
				}
			} else if sliceT != nil {
				errs = append(errs, ErrAppendFirstArgNotSlice{at(ctx, call.Args[0])})
			}
		}
	} else {
		skipTypeCheck := make([]bool, len(call.Args))
		for i := 1; i < len(call.Args); i += 1 {
			argI, moreErrs := CheckExpr(ctx, call.Args[i], env)
			call.Args[i] = argI
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			}
			if moreErrs == nil || argI.IsConst() {
				if _, err := expectSingleType(ctx, argI.KnownType(), argI); err != nil {
					skipTypeCheck[i] = true
					errs = append(errs, err)
				}
			} else {
				skipTypeCheck[i] = true
			}
		}
		if isSlice {
			eltT := sliceT.Elem()
			for i := 1; i < len(call.Args); i += 1 {
				if skipTypeCheck[i] {
					continue
				}
				argI := call.Args[i].(Expr)
				ok := false
				if argI.IsConst() {
					var ct ConstType
					if ct, ok = argI.KnownType()[0].(ConstType); ok {
						x, moreErrs := promoteConstToTyped(ctx, ct, constValue(argI.Const()), eltT, argI)
						if !reflect.Value(x).IsValid() {
							errs = append(errs, ErrBuiltinWrongArgType{at(ctx, argI), call})
						} else if moreErrs != nil {
							errs = append(errs, moreErrs...)
						}
					}
				}
				if !ok && unhackType(argI.KnownType()[0]) != unhackType(eltT) {
					errs = append(errs, ErrBuiltinWrongArgType{at(ctx, argI), call})
				}
			}
		} else if sliceT != nil {
			errs = append(errs, ErrAppendFirstArgNotSlice{at(ctx, call.Args[0])})
		}
	}
	return call, errs
}

func fakeCheckRemainingArgs(call *CallExpr, from int, env *Env) {
	for i := from; i < len(call.Args); i += 1 {
		call.Args[i] = fakeCheckExpr(call.Args[i], env)
	}
}
