package eval

import (
	"reflect"

	"go/ast"
)

func checkCallBuiltinExpr(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error, bool) {
	ident, ok := call.Fun.(*ast.Ident)
	switch ident.Name {
	case "complex":
		return checkBuiltinComplexExpr(ctx, call, env)
	case "real", "imag":
		return checkBuiltinRealImagExpr(ctx, call, env)
	case "new":
		return checkBuiltinNewExpr(ctx, call, env)
	case "make":
		return checkBuiltinMakeExpr(ctx, call, env)
	case "len":
		return checkBuiltinLenExpr(ctx, call, env)
	case "cap":
		return checkBuiltinCapExpr(ctx, call, env)
	case "append":
		return checkBuiltinAppendExpr(ctx, call, env)
	case "copy":
		return checkBuiltinCopyExpr(ctx, call, env)
	case "delete":
		return checkBuiltinDeleteExpr(ctx, call, env)
	default:
		return call, nil, false
	}
}

func checkBuiltinComplexExpr(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error, bool) {
	call.Fun = &Ident{Ident: ident}
	if len(call.Args) != 2 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at{ctx, call}}}, true
	}
	x, y, ok, errs := checkBinaryOperands(ctx, call.Args[0], call.Args[1], env)
	call.Args[0], call.Args[1] = x, y
	if !ok {
		return call, errs, true
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
				return call, errs, true
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
					return call, errs, true
				} else if yt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if y.IsConst() {
						xf = float64(xv.Float())
						yf = float64(y.Const().Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs, true
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
					return call, errs, true
				} else if xt.Kind() == reflect.Float64 {
					call.knownType = knownType{c128}
					if x.IsConst() {
						xf = float64(x.Const().Float())
						yf = float64(yv.Float())
						call.constValue = constValueOf(complex(xf, yf))
					}
					return call, errs, true
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
			return call, errs, true
		} else if xt == reflect.Float64 {
			call.knownType = knownType{c128}
			if x.IsConst() && y.IsConst() {
				xf = float64(x.Const().Float())
				yf = float64(y.Const().Float())
				call.constValue = constValueOf(complex(xf, yf))
			}
			return call, errs, true
		}
	}
	if xt == yt {
		errs = append(errs, ErrBuiltinMistypedArgs{at(ctx, call)})
	} else {
		errs = append(errs, ErrBuiltinMismatchedArgs{at(ctx, call)})
	}
	return call, errs, true
}

func checkBuiltinRealImag(ctx *Ctx, call *CallExpr, env *Env) (*CallExpr, []error, bool) {
	if len(call.Args) != 1 {
		return call, []error{ErrBuiltinWrongNumberOfArgs{at{ctx, call}}}, true
	}
		var errs []error
		call.Args[0], errs = CheckExpr(ctx, call.Args[0], env)
		call.Fun = &Ident{Ident: ident}
		call.isBuiltin = true
		call.knownType = knownType{f64}
		return call, errs, true
	} else if ident.Name == "len" || ident.Name == "cap" {
		if len(call.Args) != 1 {
			return call, []error{errors.New(ident.Name + " wrong number args")}, true
		}
		var errs []error
		call.Args[0], errs = CheckExpr(ctx, call.Args[0], env)
		call.Fun = &Ident{Ident: ident}
		call.isBuiltin = true
		call.knownType = knownType{intType}
		return call, errs, true
	} else if ident.Name == "new" {
		if len(call.Args) != 1 {
			return call, []error{errors.New("new wrong number args")}, true
		} else if _, of, _, errs := checkType(ctx, call.Args[0], env); errs != nil {
			return call, append(errs, errors.New("new bad type")), true
		} else {
			call.Fun = &Ident{Ident: ident}
			call.isBuiltin = true
			call.knownType = knownType{reflect.PtrTo(of)}
			return call, nil, true
		}
	} else if ident.Name == "append" {
		if len(call.Args) == 0 {
			return call, []error{errors.New("append wrong number args")}, true
		} else {
			var errs, moreErrs []error
			for i := range call.Args {
				if call.Args[i], moreErrs = CheckExpr(ctx, call.Args[i], env); moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
			call.Fun = &Ident{Ident: ident}
			call.isBuiltin = true
			call.knownType = call.Args[0].(Expr).KnownType()
			call.argNEllipsis = call.Ellipsis != token.NoPos
			return call, errs, true
		}
	}
	return call, nil, false
}

