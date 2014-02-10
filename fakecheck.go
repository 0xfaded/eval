package eval

import (
	"go/ast"
)

// convert an ast.Expr to an Expr without actually checking it. This
// is useful for avoiding special cases in error messages.
func fakeCheckExpr(expr ast.Expr, env *Env) Expr {
	if expr == nil {
		return nil
	}
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
			if _, t, isType, _ := checkType(&Ctx{""}, uncheckType(c.Fun), env); isType {
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

	// Types
	case *ast.ArrayType:
		arrayT := &ArrayType{ArrayType: expr}
		arrayT.Len = fakeCheckExpr(arrayT.Len, env)
		arrayT.Elt = fakeCheckExpr(arrayT.Elt, env)
		return arrayT
	case *ast.StructType:
		structT := &StructType{StructType: expr}
		return structT
	case *ast.FuncType:
		funcT := &FuncType{FuncType: expr}
		return funcT
	case *ast.InterfaceType:
		interfaceT := &InterfaceType{InterfaceType: expr}
		return interfaceT
	case *ast.MapType:
		mapT := &MapType{MapType: expr}
		mapT.Key = fakeCheckExpr(mapT.Key, env)
		mapT.Value = fakeCheckExpr(mapT.Value, env)
		return mapT
	case *ast.ChanType:
		chanT := &ChanType{ChanType: expr}
		chanT.Value = fakeCheckExpr(chanT.Value, env)
		return chanT

	default:
		// If the input is already an Expr, assume it is already type checked.
		// So far, this has been enough to produce the desired error messages.
		return expr.(Expr)
	}
}

func uncheckType(expr ast.Expr) ast.Expr {
	switch e := expr.(type) {
	case *ParenExpr:
		e.X = uncheckType(e.X)
		return e.ParenExpr
	case *Ident:
		return e.Ident
	case *StarExpr:
		e.X = uncheckType(e.X)
		return e.StarExpr
	case *StructType:
		return e.StructType
	case *ArrayType:
		return e.ArrayType
	case *InterfaceType:
		return e.InterfaceType
	case *MapType:
		return e.MapType
	case *ChanType:
		return e.ChanType
	default:
		return e
	}
}

// Remove the const value from an Expr. If the Expr is not const, do nothing. Otherwise,
// return a non const clone.
func unconstNode(expr Expr) Expr {
	// Special case for non-const builtin calls which must have their args unconsted
	if e, ok := expr.(*CallExpr); ok && (e.IsConst() || e.isBuiltin) {
		u := new(CallExpr)
		*u = *e
		u.constValue = constValue{}
		if u.isBuiltin {
			u.CallExpr = new(ast.CallExpr)
			*u.CallExpr = *e.CallExpr
			u.Args = make([]ast.Expr, len(e.Args))
			for i := range e.Args {
				if arg, ok := e.Args[i].(Expr); ok {
					u.Args[i] = unconstNode(arg)
				} else {
					u.Args[i] = e.Args[i]
				}
			}
		}
		return u
	}

	if !expr.IsConst() {
		return expr
	}
	switch e := expr.(type) {
	case *Ident:
		u := new(Ident)
		*u = *e
		u.constValue = constValue{}
		return u
	case *BasicLit:
		u := new(BasicLit)
		*u = *e
		u.constValue = constValue{}
		return u
	case *SelectorExpr:
		u := new(SelectorExpr)
		*u = *e
		u.constValue = constValue{}
		return u
	case *IndexExpr:
		u := new(IndexExpr)
		*u = *e
		u.constValue = constValue{}
		return u
	case *UnaryExpr:
		u := new(UnaryExpr)
		*u = *e
		u.constValue = constValue{}
		return u
	case *BinaryExpr:
		u := new(BinaryExpr)
		*u = *e
		u.constValue = constValue{}
		return u
	default:
		panic("eval: impossible. non-const node IsConst() returned true")
	}
}

// shorthand for unconstNode
func uc(expr Expr) Expr {
	return unconstNode(expr)
}
