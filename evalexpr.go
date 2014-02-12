package eval

import (
	"errors"
	"reflect"
)

// EvalExpr is the main function to call to evaluate an ast-parsed
// expression, expr. 
// Parameter env, contains an evaluation environment from
// which to get reflect.Values from. Note however that env can be
// subverted somewhat by supplying callback hooks routines which
// access variables and by supplying user-defined conversion routines.
func EvalExpr(expr Expr, env *Env) (*[]reflect.Value, bool, error) {
	switch node := expr.(type) {
	case *Ident:
		v, _, err := evalIdentExprCallback(node, env)
		if v == nil {
			return nil, false, err
		}
		ret := []reflect.Value{*v}
		return &ret, true, err
	case *Ellipsis:
	case *BasicLit:
		v, err := evalBasicLit(node)
		return &[]reflect.Value{v}, true, err
	case *FuncLit:
	case *CompositeLit:
		v, err := evalCompositeLit(node, env)
		return &[]reflect.Value{v}, true, err
	case *ParenExpr:
		return EvalExpr(node.X.(Expr), env)
	case *SelectorExpr:
		v, _, err := evalSelectorExprCallback(node, env)
		if v == nil {
			return nil, true, err
		}
		return &[]reflect.Value{*v}, true, err
	case *IndexExpr:
		vs, err := evalIndexExpr(node, env)
		return &vs, true, err
	case *SliceExpr:
		v, err := evalSliceExpr(node, env)
		return &[]reflect.Value{v}, true, err
	case *TypeAssertExpr:
		v, err := evalTypeAssertExpr(node, env)
		return &[]reflect.Value{v}, true, err
	case *CallExpr:
		vs, err := evalCallExpr(node, env)
		return &vs, true, err
	case *StarExpr:
		v, err := evalStarExpr(node, env)
		return &[]reflect.Value{v}, true, err
	case *UnaryExpr:
		vs, err := evalUnaryExpr(node, env)
		return &vs, true, err
	case *BinaryExpr:
		v, err := evalBinaryExpr(node, env)
		return &[]reflect.Value{v}, true, err
	case *KeyValueExpr:
	default:
		return nil , false, errors.New("undefined type")
	}
	return &[]reflect.Value{reflect.ValueOf("Alice")}, true, nil
}
