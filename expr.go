package eval

import (
	"errors"
	"fmt"
	"reflect"
)

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
		v, typed, err := evalCompositeLit(ctx, node, env)
		return &[]reflect.Value{*v}, typed, err
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
	case *TypeAssertExpr:
	case *CallExpr:
		return evalCallExpr(ctx, node, env)
	case *StarExpr:
	case *UnaryExpr:
		v, typed, err := evalUnaryExpr(ctx, node, env)
		return &[]reflect.Value{v}, typed, err
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

func evalType(ctx *Ctx, expr Expr, env *Env) (reflect.Type, error) {
	switch node := expr.(type) {
	case *Ident:
		if t, ok := env.Types[node.Name]; ok {
			return t, nil
		} else if t, ok := builtinTypes[node.Name]; ok {
			return t, nil
		} else {
			return t, errors.New("undefined type: " + node.Name)
		}
	case *ArrayType:
		return nil, errors.New("array types not implemented")
	case *StructType:
		return nil, errors.New("struct types not implemented")
	case *FuncType:
		return nil, errors.New("func types not implemented")
	case *InterfaceType:
		return nil, errors.New("interface types not implemented")
	case *MapType:
		return nil, errors.New("map types not implemented")
	case *ChanType:
		return nil, errors.New("chan types not implemented")
	default:
		return nil, errors.New(fmt.Sprintf("Type: Bad type (%+v)", node))
	}
}
