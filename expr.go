package eval

import (
	"errors"
	"fmt"
	"reflect"

        "go/ast"
)

// EvalExpr is the main function to call to evaluate an ast-parsed
// expression, expr.  Parameter ctx contains a string representation
// of expr. Parameter env, contains an evaluation environment from
// which to get reflect.Values from. Note however that env can be
// subverted somewhat by supplying callback hooks routines which
// access variables and by supplying user-defined conversion routines.
func EvalExpr(ctx *Ctx, expr Expr, env *Env) (*[]reflect.Value, bool, error) {
	switch node := expr.(type) {
	case *Ident:
		v, typed, err := evalIdentExprCallback(ctx, node, env)
		if v == nil {
			return nil, false, err
		}
		ret := []reflect.Value{*v}
		return &ret, typed, err
	case *Ellipsis:
	case *BasicLit:
		v, typed, err := evalBasicLit(ctx, node)
		return &[]reflect.Value{v}, typed, err
	case *FuncLit:
	case *CompositeLit:
		v, err := evalCompositeLit(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *ParenExpr:
		return EvalExpr(ctx, node.X.(Expr), env)
	case *SelectorExpr:
		v, typed, err := evalSelectorExprCallback(ctx, node, env)
		if v == nil {
			return nil, typed, err
		}
		return &[]reflect.Value{*v}, typed, err
	case *IndexExpr:
		v, typed, err := evalIndexExpr(ctx, node, env)
		if v == nil {
			return nil, typed, err
		}
		return &[]reflect.Value{*v}, typed, err
	case *SliceExpr:
		v, typed, err := evalSliceExpr(ctx, node, env)
		if v == nil {
			return nil, typed, err
		}
		return &[]reflect.Value{*v}, typed, err
	case *TypeAssertExpr:
	case *CallExpr:
		vs, err := evalCallExpr(ctx, node, env)
		return &vs, true, err
	case *StarExpr:
		v, err := evalStarExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *UnaryExpr:
		v, err := evalUnaryExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *BinaryExpr:
		v, typed, err := evalBinaryExpr(ctx, node, env)
		return &[]reflect.Value{v}, typed, err
	case *KeyValueExpr:
	default:
		panic(node)
		return nil , false, errors.New("undefined type")
	}
	return &[]reflect.Value{reflect.ValueOf("Alice")}, true, nil
}

func evalType(ctx *Ctx, expr ast.Expr, env *Env) (reflect.Type, error) {
	switch node := expr.(type) {
	case *ast.Ident:
		if t, ok := env.Types[node.Name]; ok {
			return t, nil
		} else if t, ok := builtinTypes[node.Name]; ok {
			return t, nil
		} else {
			return t, errors.New("undefined type: " + node.Name)
		}
	case *ast.ArrayType:
		if node.Len != nil {
			return nil, errors.New("array types not implemented")
		} else if eltT, err := evalType(ctx, node.Elt, env); err != nil {
			return nil, err
		} else {
			return reflect.SliceOf(eltT), nil
		}
	case *ast.StructType:
		return nil, errors.New("struct types not implemented")
	case *ast.FuncType:
		return nil, errors.New("func types not implemented")
	case *ast.InterfaceType:
		return nil, errors.New("interface types not implemented")
	case *ast.MapType:
		return nil, errors.New("map types not implemented")
	case *ast.ChanType:
		return nil, errors.New("chan types not implemented")
	default:
		return nil, errors.New(fmt.Sprintf("Type: Bad type (%+v)", node))
	}
}
