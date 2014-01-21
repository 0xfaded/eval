package eval

import (
	"errors"
	"fmt"
	//"reflect"

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
		return &SliceExpr{SliceExpr: expr}, nil
	case *ast.TypeAssertExpr:
		return &TypeAssertExpr{TypeAssertExpr: expr}, nil
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
		return nil, []error{errors.New(fmt.Sprintf("Type: Bad expr (%+v)", expr))}
	}

}

func checkTypeExpr(ctx *Ctx, expr ast.Expr, env *Env) (Expr, []error) {
	switch expr := expr.(type) {
	case *ast.Ident:
		return &Ident{Ident: expr}, nil
	case *ast.ArrayType:
		return &ArrayType{ArrayType: expr}, nil
	case *ast.StructType:
		return &StructType{StructType: expr}, nil
	case *ast.FuncType:
		return &FuncType{FuncType: expr}, nil
	case *ast.InterfaceType:
		return &InterfaceType{InterfaceType: expr}, nil
	case *ast.MapType:
		return &MapType{MapType: expr}, nil
	case *ast.ChanType:
		return &ChanType{ChanType: expr}, nil
	default:
		return nil, []error{errors.New(fmt.Sprintf("Type: Bad type (%+v)", expr))}
	}
}
