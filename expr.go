package interactive

import (
	"errors"
	"fmt"
	"reflect"

	"go/ast"
)

func EvalExpr(expr ast.Expr, env *Env) ([]reflect.Value, bool, error) {
	switch node := expr.(type) {
	case *ast.Ident:
		v, typed, err := evalIdentExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.Ellipsis:
	case *ast.BasicLit:
		v, typed, err := evalBasicLit(node)
		return []reflect.Value{v}, typed, err
	case *ast.FuncLit:
	case *ast.CompositeLit:
		v, typed, err := evalCompositeLit(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.ParenExpr:
		return EvalExpr(node.X, env)
	case *ast.SelectorExpr:
		v, typed, err := evalSelectorExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
	case *ast.CallExpr:
		return evalCallExpr(node, env)
	case *ast.StarExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
		v, typed, err := evalBinaryExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.KeyValueExpr:
	}
	return []reflect.Value{reflect.ValueOf("Alice")}, true, nil
}

func evalType(expr ast.Expr, env *Env) (reflect.Type, error) {
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
		return nil, errors.New("array types not implemented")
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
