package eval

import (
	"reflect"

	"go/ast"
	"go/token"
)

func checkBinaryExpr(ctx *Ctx, binary *ast.BinaryExpr, env *Env) (*BinaryExpr, []error) {
	aexpr := &BinaryExpr{BinaryExpr: binary}
	x, y, ok, errs := checkBinaryOperands(ctx, binary.X, binary.Y, env)
	binary.X, binary.Y = x, y
	if !ok {
		return aexpr, errs
	}
	xt, yt := x.KnownType()[0], y.KnownType()[0]

	xc, xuntyped := xt.(ConstType)
	yc, yuntyped := yt.(ConstType)
	op := binary.Op
	if x.IsConst() && y.IsConst() {
		if xuntyped && yuntyped {
			yv := y.Const()
			xv := x.Const()
			var promoted ConstType
			if promoted, moreErrs = promoteConsts(ctx, xc, yc, x, y, xv, yv); moreErrs != nil {
				errs = append(errs, moreErrs...)
				errs = append(errs, ErrInvalidBinaryOperation{at(ctx, aexpr)})
			} else {
				if isBooleanOp(op) {
					aexpr.knownType = []reflect.Type{ConstBool}
				} else {
					aexpr.knownType = knownType{promoted}
				}
				aexpr.constValue, moreErrs = evalConstUntypedBinaryExpr(ctx, aexpr, promoted)
				if moreErrs != nil {
					errs = append(errs, moreErrs...)
				}
			}
		} else if yuntyped {
			z, moreErrs := evalConstTypedUntypedBinaryExpr(ctx, aexpr, x, y, true)
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			} else {
				aexpr.knownType = knownType{reflect.Value(z).Type()}
				aexpr.constValue = z
			}
		} else if xuntyped {
			z, moreErrs := evalConstTypedUntypedBinaryExpr(ctx, aexpr, y, x, false)
			if moreErrs != nil {
				errs = append(errs, moreErrs...)
			} else {
				aexpr.knownType = knownType{reflect.Value(z).Type()}
				aexpr.constValue = z
			}
		} else {
			if z, moreErrs := evalConstTypedBinaryExpr(ctx, aexpr, x, y); moreErrs != nil {
				errs = append(errs, moreErrs...)
			} else {
				aexpr.knownType = knownType{reflect.Value(z).Type()}
				aexpr.constValue = z
			}
		}
	} else {
		if yuntyped {
			// the old switcheroo
			xt, yt = yt, xt
			xc, yc = yc, xc
			x, y = y, x
			xuntyped = true
		}
		yk := yt.Kind()
		errExpr := aexpr

		// special cases for const nil
		// note that only (slice|map|func|interface|ptr) == nil are legal
		// other expressions containing ConstNil do not produces a mismatched
		// types error (ErrInvalidBinaryOperation)
		if xt == ConstNil {
			if (op == token.EQL || op == token.NEQ) &&
				(yk == reflect.Slice || yk == reflect.Map || yk == reflect.Func ||
				yk == reflect.Interface || yk == reflect.Ptr) {
				aexpr.knownType = knownType{boolType}
			} else if yk == reflect.String || yk == reflect.Slice || yk == reflect.Interface || yk == reflect.Ptr || yk == reflect.Map {
				// Except strings, they do produce mismatched types
				// instead of bad conversions
				err := ErrInvalidBinaryOperation{at(ctx, aexpr)}
				errs = append(errs, err)
			} else {
				err := ErrBadConstConversion{at(ctx, x), xt, yt, x.Const()}
				errs = append(errs, err)
			}
			// http://code.google.com/p/go/issues/detail?id=7206
			if yk == reflect.Array || yk == reflect.Uintptr {
				errs = append(errs, ErrUntypedNil{at(ctx, x)})
			}
			return aexpr, errs
		}

		xk := xt.Kind()
		var operandT reflect.Type
		// Identical types are always valid, except non comparable structs
		// and types with can only be compared to nil
                if unhackType(xt) == unhackType(yt) {
			if !comparableToNilOnly(xt) && (xk != reflect.Struct || isStructComparable(xt)) {
				operandT = xt
			}
                } else if xuntyped && attemptBinaryOpConversion(yt) {
	                c, moreErrs := promoteConstToTyped(ctx, xc, constValue(x.Const()), yt, x)
			errs = append(errs, moreErrs...)
			v := reflect.Value(c)
			if v.IsValid() {
				operandT = yt
				// Check for divide by zero. Note we only need to do this
				// if y was the untyped constant, and also note that
				// it may have been truncated
				if yuntyped && (op == token.QUO || op == token.REM) &&
					isOpDefinedOn(op, operandT) &&
					v.Interface() == reflect.Zero(yt).Interface() {

					errs = append(errs, ErrDivideByZero{at(ctx, x)})
				}
			}
		// An interface is comprable if its paired operand implements it.
		// To match gc error output, if the operator produces a boolean and
		// one operand is a type that satisfies but is not the the other
		// operand's type, wrap that node in a type cast. This will only be
		// used by errors.
		} else if yk == reflect.Interface && xt.Implements(yt) {
			operandT = yt
			if isBooleanOp(op) {
				errExpr = new(BinaryExpr)
				*errExpr = *aexpr
				errExpr.X = wrapConcreteTypeWithInterface(x, yt)
			}
		} else if xk == reflect.Interface && yt.Implements(xt) {
			operandT = xt
			if isBooleanOp(op) {
				errExpr = new(BinaryExpr)
				*errExpr = *aexpr
				errExpr.Y = wrapConcreteTypeWithInterface(y, xt)
			}
                }

		if operandT != nil {
			if !isOpDefinedOn(binary.Op, operandT) {
				errs = append(errs, ErrInvalidBinaryOperation{at(ctx, aexpr)})
			} else if isBooleanOp(binary.Op) {
				aexpr.knownType = knownType{boolType}
			} else {
				aexpr.knownType = knownType{operandT}
			}
                } else {
                        errs = append(errs, ErrInvalidBinaryOperation{at(ctx, aexpr)})
		}
        }
	return aexpr, errs
}

// Evaluates a const binary Expr. May return a sensical constValue
// even if ErrTruncatedConst errors are present
func evalConstUntypedBinaryExpr(ctx *Ctx, constExpr *BinaryExpr, promotedType ConstType) (constValue, []error) {
	x := constExpr.X.(Expr).Const()
	y := constExpr.Y.(Expr).Const()
	switch promotedType.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		xx := x.Interface().(*ConstNumber)
		yy := y.Interface().(*ConstNumber)
		return evalConstBinaryNumericExpr(ctx, constExpr, xx, yy)
	case ConstStringType:
		xx := x.String()
		yy := y.String()
		return evalConstBinaryStringExpr(ctx, constExpr, xx, yy)
	case ConstBoolType:
		xx := x.Bool()
		yy := y.Bool()
		return evalConstBinaryBoolExpr(ctx, constExpr, xx, yy)
	default:
		// It is possible that both x and y are ConstNil, however no operator is defined, not even ==
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}

}

func evalConstBinaryNumericExpr(ctx *Ctx, constExpr *BinaryExpr, x, y *ConstNumber) (constValue, []error) {
	var errs []error

	switch constExpr.Op {
	case token.ADD:
		return constValueOf(new(ConstNumber).Add(x, y)), nil
	case token.SUB:
		return constValueOf(new(ConstNumber).Sub(x, y)), nil
	case token.MUL:
		return constValueOf(new(ConstNumber).Mul(x, y)), nil
	case token.QUO:
		if y.Value.IsZero() {
			return constValue{}, []error{ErrDivideByZero{at(ctx, constExpr.Y)}}
		}
		return constValueOf(new(ConstNumber).Quo(x, y)), nil
	case token.REM:
		if y.Value.IsZero() {
			return constValue{}, []error{ErrDivideByZero{at(ctx, constExpr.Y)}}
		} else if !(x.Type.IsIntegral() && y.Type.IsIntegral()) {
			return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
		} else {
			return constValueOf(new(ConstNumber).Rem(x, y)), nil
		}
	case token.AND, token.OR, token.XOR, token.AND_NOT:
		if !(x.Type.IsIntegral() && y.Type.IsIntegral()) {
			return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
		}

		switch constExpr.Op {
		case token.AND:
			return constValueOf(new(ConstNumber).And(x, y)), nil
		case token.OR:
			return constValueOf(new(ConstNumber).Or(x, y)), nil
		case token.XOR:
			return constValueOf(new(ConstNumber).Xor(x, y)), nil
		case token.AND_NOT:
			return constValueOf(new(ConstNumber).AndNot(x, y)), nil
		default:
			panic("go-interactive: impossible")
		}

	case token.EQL:
		return constValueOf(x.Value.Equals(&y.Value)), nil
	case token.NEQ:
		return constValueOf(!x.Value.Equals(&y.Value)), nil

	case token.LEQ, token.GEQ, token.LSS, token.GTR:
		var b bool
		if !(x.Type.IsReal() && y.Type.IsReal()) {
			return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
		}
		cmp := x.Value.Re.Cmp(&y.Value.Re)
		switch constExpr.Op {
		case token.NEQ:
			b = cmp != 0
		case token.LEQ:
			b = cmp <= 0
		case token.GEQ:
			b = cmp >= 0
		case token.LSS:
			b = cmp < 0
		case token.GTR:
			b = cmp > 0
		}
		return constValueOf(b), errs
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

func evalConstBinaryStringExpr(ctx *Ctx, constExpr *BinaryExpr, x, y string) (constValue, []error) {
	switch constExpr.Op {
	case token.ADD:
		return constValueOf(x + y), nil
	case token.EQL:
		return constValueOf(x == y), nil
	case token.NEQ:
		return constValueOf(x != y), nil
	case token.LEQ:
		return constValueOf(x <= y), nil
	case token.GEQ:
		return constValueOf(x >= y), nil
	case token.LSS:
		return constValueOf(x < y), nil
	case token.GTR:
		return constValueOf(x > y), nil
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

func evalConstBinaryBoolExpr(ctx *Ctx, constExpr *BinaryExpr, x, y bool) (constValue, []error) {
	switch constExpr.Op {
	case token.EQL:
		return constValueOf(x == y), nil
	case token.NEQ:
		return constValueOf(x != y), nil
	case token. LAND:
		return constValueOf(x && y), nil
	case token.LOR:
		return constValueOf(x || y), nil
	default:
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, constExpr)}}
	}
}

// Evaluate x op y
func evalConstTypedUntypedBinaryExpr(ctx *Ctx, binary *BinaryExpr, typedExpr, untypedExpr Expr, reversed bool) (
	constValue, []error) {

	xt := untypedExpr.KnownType()[0].(ConstType)
	yt := typedExpr.KnownType()[0]

	// x must be convertible to target type
	xUntyped := untypedExpr.Const()
	x, xConvErrs := promoteConstToTyped(ctx, xt, constValue(xUntyped), yt, untypedExpr)

        // If the untyped operand is nil, gc simply says it could not convert the nil type
        if xt == ConstNil {
                // ... unless, its a string. In which case, report mismatched types
                if yt.Kind() == reflect.String {
		        return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, binary)}}
                } else {
		        return constValue{}, xConvErrs
                }
	} else if !isOpDefinedOn(binary.Op, yt) {
		return constValue{}, append(xConvErrs, ErrInvalidBinaryOperation{at(ctx, binary)})
	}

        // Check for an impossible conversion. This occurs when the types
        // are incompatible, such as string(1.5). For other errors, such as
        // integer overflows, the type check should continue as if the conversion
        // succeeded
        if !reflect.Value(x).IsValid() {
		errs := append(xConvErrs, ErrInvalidBinaryOperation{at(ctx, binary)})
		return constValue{}, errs
        }

	switch xt.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		xx, xok := convertTypedToConstNumber(reflect.Value(x))
		yy, yok := convertTypedToConstNumber(typedExpr.Const())

		// If a child node errored, then it is possible that, typedExpr.Const() is
		// actually a *ConstNumber to avoid loss of precision in error messages.
		if !yok {
			yy, yok = reflect.Value(typedExpr.Const()).Interface().(*ConstNumber)
		}

		if !xok || !yok {
			// This is a non numeric expression. Return the errors encountered so far
			return constValue{}, append(xConvErrs, ErrInvalidBinaryOperation{at(ctx, binary)})
		}

		if reversed {
			xx, yy = yy, xx
		}

		z, errs := evalConstBinaryNumericExpr(ctx, binary, xx, yy)
		if errs != nil {
			return constValueOf(z), append(xConvErrs, errs...)
		}
		errs = append(xConvErrs, errs...)

		var zt ConstType
		var rt reflect.Type
		if isBooleanOp(binary.Op) {
			zt = ConstBool
			rt = reflect.TypeOf(false)
		} else {
			zt = reflect.Value(z).Interface().(*ConstNumber).Type
			rt = yt
		}

		r, moreErrs := promoteConstToTyped(ctx, zt, z, rt, binary)
		return constValue(r), append(errs, moreErrs...)

	case ConstStringType:
		if yt.Kind() == reflect.String {
			xstring := reflect.Value(x).String()
			ystring := typedExpr.Const().String()

			if reversed {
				xstring, ystring = ystring, xstring
			}

			z, errs := evalConstBinaryStringExpr(ctx, binary, xstring, ystring)
                        if errs != nil {
                                return constValue{}, errs
                        }

                        var zt ConstType
                        var rt reflect.Type
                        if isBooleanOp(binary.Op) {
                                zt = ConstBool
                                rt = reflect.TypeOf(false)
                        } else {
                                zt = ConstString
                                rt = reflect.TypeOf("")
                        }

			r, errs := promoteConstToTyped(ctx, zt, z, rt, binary)
			return constValue(r), errs
		}

	case ConstBoolType:
		if yt.Kind() == reflect.Bool {
			xbool := reflect.Value(x).Bool()
			ybool := typedExpr.Const().Bool()

			if reversed {
				xbool, ybool = ybool, xbool
			}

			z, errs := evalConstBinaryBoolExpr(ctx, binary, xbool, ybool)
			return constValue(z), errs
		}
	}
	return constValue{}, append(xConvErrs, ErrInvalidBinaryOperation{at(ctx, binary)})
}

func evalConstTypedBinaryExpr(ctx *Ctx, binary *BinaryExpr, xexpr, yexpr Expr) (constValue, []error) {

	// These are known not to be ConstTypes
	xt := xexpr.KnownType()[0]
	yt := yexpr.KnownType()[0]

        // Check that the types are compatible, handling the special alias type for runes
        // For the sake of error messages, for expressions involving int32 and rune, the
        // resulting type is that of the left operand
	var zt reflect.Type
        if xt == yt {
		zt = xt
        } else if xt == RuneType && yt == RuneType.Type {
                zt = RuneType
        } else if yt == RuneType && xt == RuneType.Type {
                zt = xt
        } else {
		return constValue{}, []error{ErrInvalidBinaryOperation{at(ctx, binary)}}
	}

	x, xok := convertTypedToConstNumber(xexpr.Const())
	y, yok := convertTypedToConstNumber(yexpr.Const())

	if xok && yok {
		z, errs := evalConstBinaryNumericExpr(ctx, binary, x, y)
                if isBooleanOp(binary.Op) {
                        return constValue(z), errs
                }
                if errs != nil {
                        if _, ok := errs[0].(ErrInvalidBinaryOperation); ok {
                                // This happens if the operator is not defined on x and y
		                return constValue(z), errs
                        }
                }
		from := reflect.Value(z).Interface().(*ConstNumber).Type
		r, moreErrs := promoteConstToTyped(ctx, from, z, zt, binary)
		return constValue(r), append(errs, moreErrs...)
	} else if !xok && !yok {
		switch zt.Kind() {
		case reflect.String:
			xstring := xexpr.Const().String()
			ystring := yexpr.Const().String()
			z, errs := evalConstBinaryStringExpr(ctx, binary, xstring, ystring)
                        if isBooleanOp(binary.Op) {
                                return constValue(z), errs
                        }
			r, moreErrs := promoteConstToTyped(ctx, ConstString, z, zt, binary)
			return constValue(r), append(errs, moreErrs...)

		case reflect.Bool:
			xbool := xexpr.Const().Bool()
			ybool := yexpr.Const().Bool()
			z, errs := evalConstBinaryBoolExpr(ctx, binary, xbool, ybool)
			return constValue(z), errs
		}
	}
	panic("go-interactive: impossible")
}

func checkBinaryOperands(ctx, yexpr, xexpr ast.Expr, env *Env) (Expr, Expr, bool, errs) {
	var xok, yok bool
	var xt, yt reflect.Type
	var err error

	x, errs := CheckExpr(ctx, binary.X, env)
	if errs == nil || x.IsConst() {
		if xt, err = expectSingleType(ctx, x.KnownType(), x); err != nil {
			errs = append(errs, err)
		} else {
			xok = true
		}
	}

	y, moreErrs := CheckExpr(ctx, binary.Y, env)
	if moreErrs == nil || y.IsConst() {
		if yt, err = expectSingleType(ctx, y.KnownType(), y); err != nil {
			errs = append(moreErrs, err)
		} else {
			yok = true
		}
	}
	errs = append(errs, moreErrs...)
	return x, y, xok && yok, errs
}


func wrapConcreteTypeWithInterface(operand Expr, interfaceT reflect.Type) Expr {
	// Rig the token positions to such that typeConv.(Len|Pos) match operand
	typeConv := &CallExpr{CallExpr: new(ast.CallExpr)}
	typeConv.Fun = &Ident{Ident: &ast.Ident{Name: "", NamePos: operand.Pos()}}
	typeConv.Lparen = operand.Pos()
	typeConv.Rparen = operand.End() - 1
	typeConv.knownType = knownType{interfaceT}
	typeConv.Args = []ast.Expr{operand}
	typeConv.isTypeConversion = true
	return typeConv;
}
