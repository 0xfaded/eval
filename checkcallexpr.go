package interactive

import (
	"reflect"
	"go/ast"
)

func checkCallExpr(ctx *Ctx, callExpr *ast.CallExpr, env *Env) (acall *CallExpr, errs []error) {
	acall = &CallExpr{CallExpr: callExpr}

	var moreErrs []error
	if acall.Fun, moreErrs = CheckExpr(ctx, callExpr.Fun, env); moreErrs != nil {
		errs = append(errs, moreErrs...)
	}

	for i := range callExpr.Args {
		if acall.Args[i], moreErrs = CheckExpr(ctx, callExpr.Args[i], env); moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
	}

	if errs != nil {
		return acall, errs
	}

	fun := acall.Fun.(Expr)
	if to, err := evalType(ctx, acall.Fun.(Expr), env); err == nil {
		return checkCallTypeExpr(ctx, to, acall, env)
	} else if fun.IsConst() && fun.KnownType()[0] == ConstNil {
		return acall, []error{ErrUntypedNil{at(ctx, fun)}}
	}
	// TODO eval function calls
	return acall, errs
}

func checkCallTypeExpr(ctx *Ctx, to reflect.Type, call *CallExpr, env *Env) (*CallExpr, []error) {
	call.knownType = []reflect.Type{to}
	call.isTypeConversion = true

	if len(call.Args) != 1 {
		return call, []error{ErrWrongNumberOfArgs{at(ctx, call)}}
	}

	arg := call.Args[0].(Expr)
	from, err := expectSingleType(ctx, arg.KnownType(), arg)
	if err != nil {
		return call, []error{err}
	}

	if arg.IsConst() {
		// For bad constant conversions, gc produces two error messages. E.g. string to uint64
		// cannot convert "abc" to type uint64
		// cannot convert "abc" (type string) to type uint64
		//
		// I've separated these into ErrBadConstConversiond and ErrBadConversion
		// The exception is if the conversion is from nil
		v, errs := convertConstToTyped(ctx, from.(ConstType), constValue(arg.Const()), to, arg)
		if errs != nil {
			if b, ok := errs[0].(ErrBadConstConversion); ok && b.from != ConstNil {
				err := ErrBadConversion{b.ErrorContext, b.from, b.to, b.v}
				errs = append(errs, err)
			}
			return call, errs
		} else {
			call.constValue = v
			return call, nil
		}
	} else {
		if from.ConvertibleTo(to) {
			return call, nil
		} else {
			return call, []error{ErrBadConstConversion{at(ctx, call), from, to, reflect.Value{}}}
		}
	}
}
