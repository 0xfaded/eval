package eval

import (
	"errors"
	"reflect"
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
		v, _, err := evalIdentExprCallback(ctx, node, env)
		if v == nil {
			return nil, false, err
		}
		ret := []reflect.Value{*v}
		return &ret, true, err
	case *Ellipsis:
	case *BasicLit:
		v, err := evalBasicLit(ctx, node)
		return &[]reflect.Value{v}, true, err
	case *FuncLit:
	case *CompositeLit:
		v, err := evalCompositeLit(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *ParenExpr:
		return EvalExpr(ctx, node.X.(Expr), env)
	case *SelectorExpr:
		v, _, err := evalSelectorExprCallback(ctx, node, env)
		if v == nil {
			return nil, true, err
		}
		return &[]reflect.Value{*v}, true, err
	case *IndexExpr:
		v, err := evalIndexExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *SliceExpr:
		v, err := evalSliceExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *TypeAssertExpr:
		v, err := evalTypeAssertExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
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
		v, err := evalBinaryExpr(ctx, node, env)
		return &[]reflect.Value{v}, true, err
	case *KeyValueExpr:
	default:
		return nil , false, errors.New("undefined type")
	}
	return &[]reflect.Value{reflect.ValueOf("Alice")}, true, nil
}
