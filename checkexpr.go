package eval

import (
	"errors"
	"reflect"
	"go/ast"
)

// Type check an ast.Expr to produce an Expr. Errors are accumulated and
// returned as a single slice. When evaluating constant expressions,
// non fatal truncation/overflow errors may be raised but type checking
// will continue. A common pattern to detect errors is
//
//  if expr, errs := CheckExpr(...); errs != nil && !expr.IsConst() {
//    fatal
//  }
//
// if expr.IsConst() is true, then the resulting Expr has been successfully
// checked, regardless of if errors are present.
func CheckExpr(ctx *Ctx, expr ast.Expr, env *Env) (Expr, []error) {
	if t, _, isType, _ := checkType(ctx, expr, env); isType {
		return t, []error{ErrTypeUsedAsExpression{at(ctx, t)}}
	}

	switch expr := expr.(type) {
	case *ast.BadExpr:
		return &BadExpr{BadExpr: expr}, nil
	case *ast.Ident:
		return checkIdent(ctx, expr, env)
	case *ast.Ellipsis:
		return &Ellipsis{Ellipsis: expr}, nil
	case *ast.BasicLit:
		return checkBasicLit(ctx, expr, env)
	case *ast.FuncLit:
		return &FuncLit{FuncLit: expr}, nil
	case *ast.CompositeLit:
		return checkCompositeLit(ctx, expr, env)
	case *ast.ParenExpr:
		return checkParenExpr(ctx, expr, env)
	case *ast.SelectorExpr:
		return checkSelectorExpr(ctx, expr, env)
	case *ast.IndexExpr:
		return checkIndexExpr(ctx, expr, env)
	case *ast.SliceExpr:
		return checkSliceExpr(ctx, expr, env)
	case *ast.TypeAssertExpr:
		return checkTypeAssertExpr(ctx, expr, env)
	case *ast.CallExpr:
		return checkCallExpr(ctx, expr, env)
	case *ast.StarExpr:
		return checkStarExpr(ctx, expr, env)
	case *ast.UnaryExpr:
		return checkUnaryExpr(ctx, expr, env)
	case *ast.BinaryExpr:
		return checkBinaryExpr(ctx, expr, env)
	case *ast.KeyValueExpr:
		panic("eval: KeyValueExpr checked")
	default:
		panic("eval: Bad expr")
	}
}

func checkType(ctx *Ctx, expr ast.Expr, env *Env) (Expr, reflect.Type, bool, []error) {
	for parens, ok := expr.(*ast.ParenExpr); ok; parens, ok = expr.(*ast.ParenExpr) {
		expr = parens.X
	}
	switch node := expr.(type) {
	case *ast.Ident:
		ident := &Ident{Ident: node}
		if t, ok := env.Types[node.Name]; ok {
			return ident, t, true, nil
		} else if t, ok := builtinTypes[node.Name]; ok {
			return ident, t, true, nil
		} else {
			return ident, nil, false, []error{ErrUndefined{at(ctx, ident)}}
		}
	case *ast.StarExpr:
		star := &StarExpr{StarExpr: node}
		elem, elemT, isType, errs := checkType(ctx, node.X, env)
		if isType {
			// Only set X if X is a type, as * can be part of an expression or type
			star.X = elem
		}
		if errs != nil {
			return star, nil, isType, errs
		} else {
			return star, reflect.PtrTo(elemT), isType, nil
		}
	case *ast.ArrayType:
		arrayT := &ArrayType{ArrayType: node}
		if node.Len != nil {
			return arrayT, nil, true, []error{errors.New("array types not implemented")}
		} else {
			elt, eltT, _, errs := checkType(ctx, node.Elt, env);
			arrayT.Elt = elt
			if errs != nil {
				return arrayT, nil, true, errs
			} else {
				return arrayT, reflect.SliceOf(unhackType(eltT)), true, nil
			}
		}
	case *ast.StructType:
		structT := &StructType{StructType: node}
		return structT, nil, true, []error{errors.New("struct types not implemented")}
	case *ast.FuncType:
		funcT := &FuncType{FuncType: node}
		return funcT, nil, true, []error{errors.New("func types not implemented")}
	case *ast.InterfaceType:
		interfaceT := &InterfaceType{InterfaceType: node}
		// Allow interface{}'s
		if node.Methods.List == nil {
			return interfaceT, emptyInterface, true, nil
		}
		return interfaceT, nil, true, []error{errors.New("interface types not implemented")}
	case *ast.MapType:
		mapT := &MapType{MapType: node}
		keyT, k, _, errs := checkType(ctx, mapT.Key, env)
		mapT.Key = keyT
		if k != nil && !isStaticTypeComparable(k) {
			errs = append(errs, ErrUncomparableMapKey{at(ctx, node), k})
		}
		valueT, v, _, moreErrs := checkType(ctx, mapT.Value, env)
		mapT.Value = valueT
		if moreErrs != nil {
			errs = append(errs, moreErrs...)
		}
		if errs == nil {
			return mapT, reflect.MapOf(unhackType(k), unhackType(v)), true, nil
		}
		return mapT, nil, true, errs
	case *ast.ChanType:
		chanT := &ChanType{ChanType: node}
		value, valueT, _, errs := checkType(ctx, node.Value, env);
		chanT.Value = value
		if errs != nil {
			return chanT, nil, true, errs
		} else {
			if node.Dir == ast.SEND {
				chanT.dir = reflect.SendDir
			} else if node.Dir == ast.RECV {
				chanT.dir = reflect.RecvDir
			} else {
				chanT.dir = reflect.BothDir
			}
			return chanT, reflect.ChanOf(chanT.dir, unhackType(valueT)), true, nil
		}
	}
	// Note this error should never be shown to the user. It is used to detect
	// when a CallExpr is a type conversion
	return nil, nil, false, []error{errors.New("Bad type")}
}
