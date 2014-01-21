package eval

import (
	"go/ast"
)

// convert an ast.Expr to an Expr without actually checking it. This
// is useful for avoiding special cases in error messages.
func fakeCheckExpr(expr ast.Expr, env *Env) Expr {
	switch expr := expr.(type) {
	case *ast.BadExpr:
		return &BadExpr{BadExpr: expr}
	case *ast.Ident:
		return &Ident{Ident: expr}
	case *ast.Ellipsis:
		return &Ellipsis{Ellipsis: expr}
	case *ast.BasicLit:
		return &BasicLit{BasicLit: expr}
	case *ast.FuncLit:
		return &FuncLit{FuncLit: expr}
	case *ast.CompositeLit:
		c := &CompositeLit{CompositeLit: expr}
		for i := range c.Elts {
			c.Elts[i] = fakeCheckExpr(c.Elts[i], env)
		}
		return c
	case *ast.ParenExpr:
		p := &ParenExpr{ParenExpr: expr}
		p.X = fakeCheckExpr(p.X, env)
		return p
	case *ast.SelectorExpr:
		s := &SelectorExpr{SelectorExpr: expr}
		s.X = fakeCheckExpr(s.X, env)
		return s
	case *ast.IndexExpr:
		i := &IndexExpr{IndexExpr: expr}
		i.X = fakeCheckExpr(i.X, env)
		i.Index = fakeCheckExpr(i.Index, env)
		return i
	case *ast.SliceExpr:
		s := &SliceExpr{SliceExpr: expr}
		s.X = fakeCheckExpr(s.X, env)
		if s.Low != nil {
			s.Low = fakeCheckExpr(s.Low, env)
		}
		if s.High != nil {
			s.High = fakeCheckExpr(s.High, env)
		}
		// TODO[crc] go 1.2 introduces the [::] notation. Add after upgrade
		//if s.Max != nil {
			//s.Max = fakeCheckExpr(s.Max, env)
		//}
		return s
	case *ast.TypeAssertExpr:
		a := &TypeAssertExpr{TypeAssertExpr: expr}
		a.X = fakeCheckExpr(a.X, env)
		return a
	case *ast.CallExpr:
		c := &CallExpr{CallExpr: expr}
		if ident, ok := c.Fun.(*ast.Ident); ok {
			if _, ok := builtinFuncs[ident.Name]; ok {
				c.isBuiltin = true
			}
		}
		if !c.isBuiltin {
			if t, err := evalType(&Ctx{""}, c.Fun, env); err == nil {
				c.isTypeConversion = true
				c.knownType = knownType{t}
			}
		}
		c.Fun = fakeCheckExpr(c.Fun, env)
		for i := range c.Args {
			c.Args[i] = fakeCheckExpr(c.Args[i], env)
		}
		return c
	case *ast.StarExpr:
		s := &StarExpr{StarExpr: expr}
		s.X = fakeCheckExpr(s.X, env)
		return s
	case *ast.UnaryExpr:
		u := &UnaryExpr{UnaryExpr: expr}
		u.X = fakeCheckExpr(u.X, env)
		return u
	case *ast.BinaryExpr:
		b := &BinaryExpr{BinaryExpr: expr}
		b.X = fakeCheckExpr(b.X, env)
		b.Y = fakeCheckExpr(b.Y, env)
		return b
	case *ast.KeyValueExpr:
		kv := &KeyValueExpr{KeyValueExpr: expr}
		kv.Key = fakeCheckExpr(kv.Key, env)
		kv.Value = fakeCheckExpr(kv.Value, env)
		return kv
	default:
		return &BadExpr{}
	}
}
