package eval

import (
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
)

// Annotated ast.Expr nodes
type Expr interface {
	ast.Expr

	// The type of this expression if known. Certain expr types have special interpretations
	// Constant expr: a ConstType will be returned if constant is untyped
	// Ellipsis expr: a single reflect.Type represents the type of all unpacked values
	KnownType() []reflect.Type

	// Returns true if this expression evaluates to a constant. In this
	// case, Value() will return the evalutated constant. Nodes below
	// this expression can be ignored.
	IsConst() bool

	// Returns the const value, if known.
	Const() reflect.Value

	// String() matches the print format of expressions in go errors
	String() string
}

type knownType []reflect.Type
type constValue reflect.Value

type BadExpr struct {
	*ast.BadExpr
}

type Ident struct {
	*ast.Ident
	knownType
	constValue
}

type Ellipsis struct {
	*ast.Ellipsis
	knownType
}

type BasicLit struct {
	*ast.BasicLit
	knownType
	constValue
}

type FuncLit struct {
	*ast.FuncLit
}

type CompositeLit struct {
	*ast.CompositeLit
	knownType
}

type ParenExpr struct {
	*ast.ParenExpr
	knownType
	constValue
}

type SelectorExpr struct {
	*ast.SelectorExpr
	knownType
}

type IndexExpr struct {
	*ast.IndexExpr
	knownType

	// Const value only relevant for strings.
	// "abc"[2] is a const expr
	// []int{1}[2] or [1]int{1}[2] are not
	constValue
}

type SliceExpr struct {
	*ast.SliceExpr
	knownType
}

type TypeAssertExpr struct {
	*ast.TypeAssertExpr
	knownType
}

type CallExpr struct {
	*ast.CallExpr

	// Is this a type conversion. If true, 't' will be non-nil
	isTypeConversion bool
	knownType
	constValue

	// Does this function a take single, multivalued expession that is unpacked as its arguments
	arg0MultiValued bool

	// Is an ellipsis expression used to unpack variadic arguments
	argNEllipsis bool
}

type StarExpr struct {
	*ast.StarExpr
	knownType
}

type UnaryExpr struct {
	*ast.UnaryExpr
	knownType
	constValue
}

type BinaryExpr struct {
	*ast.BinaryExpr
	knownType
	constValue
}

type KeyValueExpr struct {
	*ast.KeyValueExpr
}

type ArrayType struct {
	*ast.ArrayType
	knownType
}

type StructType struct {
	*ast.StructType
	knownType
}

type FuncType struct {
	*ast.FuncType
	knownType
}

type InterfaceType struct {
	*ast.InterfaceType
	knownType
}

type MapType struct {
	*ast.MapType
	knownType
}

type ChanType struct {
	*ast.ChanType
	knownType
}

func (t knownType) KnownType() []reflect.Type {
	return []reflect.Type(t)
}

func (c constValue) IsConst() bool {
	return reflect.Value(c).IsValid()
}

func (c constValue) Const() reflect.Value {
	return reflect.Value(c)
}

func (*BadExpr) KnownType() []reflect.Type      { return nil }
func (*FuncLit) KnownType() []reflect.Type      { return nil }
func (*KeyValueExpr) KnownType() []reflect.Type { return nil }

func (*BadExpr) IsConst() bool        { return false }
func (*Ellipsis) IsConst() bool       { return false }
func (*FuncLit) IsConst() bool        { return false }
func (*CompositeLit) IsConst() bool   { return false }
func (*SelectorExpr) IsConst() bool   { return false }
func (*IndexExpr) IsConst() bool      { return false }
func (*SliceExpr) IsConst() bool      { return false }
func (*TypeAssertExpr) IsConst() bool { return false }
func (*StarExpr) IsConst() bool       { return false }
func (*KeyValueExpr) IsConst() bool   { return false }
func (*ArrayType) IsConst() bool      { return false }
func (*StructType) IsConst() bool     { return false }
func (*FuncType) IsConst() bool       { return false }
func (*InterfaceType) IsConst() bool  { return false }
func (*MapType) IsConst() bool        { return false }
func (*ChanType) IsConst() bool       { return false }

func (*BadExpr) Const() reflect.Value        { return reflect.Value{} }
func (*Ellipsis) Const() reflect.Value       { return reflect.Value{} }
func (*FuncLit) Const() reflect.Value        { return reflect.Value{} }
func (*CompositeLit) Const() reflect.Value   { return reflect.Value{} }
func (*SelectorExpr) Const() reflect.Value   { return reflect.Value{} }
func (*IndexExpr) Const() reflect.Value      { return reflect.Value{} }
func (*SliceExpr) Const() reflect.Value      { return reflect.Value{} }
func (*TypeAssertExpr) Const() reflect.Value { return reflect.Value{} }
func (*StarExpr) Const() reflect.Value       { return reflect.Value{} }
func (*KeyValueExpr) Const() reflect.Value   { return reflect.Value{} }
func (*ArrayType) Const() reflect.Value      { return reflect.Value{} }
func (*StructType) Const() reflect.Value     { return reflect.Value{} }
func (*FuncType) Const() reflect.Value       { return reflect.Value{} }
func (*InterfaceType) Const() reflect.Value  { return reflect.Value{} }
func (*MapType) Const() reflect.Value        { return reflect.Value{} }
func (*ChanType) Const() reflect.Value       { return reflect.Value{} }

// Does not assert that c is a valid const value type
// Should be *BigComplex, bool, or string
func constValueOf(i interface{}) constValue {
	return constValue(reflect.ValueOf(i))
}

func (badExpr *BadExpr) String() string {
	return "BadExpr"
}

func (ident *Ident) String() string {
	if ident.IsConst() {
		return sprintConstValue(ident.KnownType()[0], ident.Const(), true)
	}
	return ident.Ident.String()
}

func (ellipsis *Ellipsis) String() string { return "TODO  ellipsis.Ellipsis" }

func (basicLit *BasicLit) String() string {
	if basicLit.IsConst() {
		return sprintConstValue(basicLit.KnownType()[0], basicLit.Const(), true)
	}
	return basicLit.Value
}

func (funcLit *FuncLit) String() string { return "TODO  funcLit.FuncLit" }
func (compositeLit *CompositeLit) String() string { return "TODO  compositeLit.CompositeLit" }

func (parenExpr *ParenExpr) String() string {
	if parenExpr.IsConst() {
		return sprintConstValue(parenExpr.KnownType()[0], parenExpr.Const(), true)
	}
	return skipSuperfluousParens(parenExpr).String()
}

func (selectorExpr *SelectorExpr) String() string { return "TODO  selectorExpr.SelectorExpr" }
func (indexExpr *IndexExpr) String() string { return "TODO  indexExpr.IndexExpr" }
func (sliceExpr *SliceExpr) String() string { return "TODO  sliceExpr.SliceExpr" }
func (typeAssertExpr *TypeAssertExpr) String() string { return "TODO  typeAssertExpr.TypeAssertExpr" }

func (callExpr *CallExpr) String() string {
	if callExpr.IsConst() {
		return sprintConstValue(callExpr.KnownType()[0], callExpr.Const(), true)
	}
	return "TODO  callExpr.CallExpr"
}

func (starExpr *StarExpr) String() string { return "TODO  starExpr.StarExpr" }

func (unary *UnaryExpr) String() string {
	operand := skipSuperfluousParens(unary.X.(Expr))
	return fmt.Sprintf("%v %v", unary.Op, operand)
}

func (binary *BinaryExpr) String() string {
	left := simplifyBinaryChildExpr(binary, binary.X.(Expr))
	right := simplifyBinaryChildExpr(binary, binary.Y.(Expr))

	return fmt.Sprintf("%v %v %v", left, binary.Op, right)
}

func (keyValueExpr *KeyValueExpr) String() string { return "TODO  keyValueExpr.KeyValueExpr" }
func (arrayType *ArrayType) String() string { return "TODO  arrayType.ArrayType" }
func (structType *StructType) String() string { return "TODO  structType.StructType" }
func (funcType *FuncType) String() string { return "TODO  funcType.FuncType" }
func (interfaceType *InterfaceType) String() string { return "TODO  interfaceType.InterfaceType" }
func (mapType *MapType) String() string { return "TODO  mapType.MapType" }
func (chanType *ChanType) String() string { return "TODO  chanType.ChanType" }

// Returns a printable interface{} which replaces constant expressions with their constants
func simplifyBinaryChildExpr(parent *BinaryExpr, expr Expr) interface{} {
        if expr.IsConst() {
		return sprintConstValue(expr.KnownType()[0], expr.Const(), true)
	}
	expr = skipSuperfluousParens(expr)
	if p, ok := expr.(*ParenExpr); ok {
		// Remove parens all together from 1 + (2 * 3)
		if b, ok := p.X.(*BinaryExpr); ok && b.Op.Precedence() > parent.Op.Precedence() {
			return p.X
		}
	}
	return expr
}

func sprintConstValue(t reflect.Type, v reflect.Value, showZeroComponents bool) string {
	i := v.Interface()
	switch x := i.(type) {
	case *ConstNumber:
                // This is hear to print overflowing typed runes, which
                // are represented by a *ConstNumber and a RuneType
                if t == RuneType {
                        r := *x
                        r.Type = ConstInt
		        return "rune(" + r.StringShow0i(false) + ")"
                }
		return x.StringShow0i(false)
        case rune:
                if t == RuneType {
		        return fmt.Sprintf("rune(%v)", x)
                }
	case float32, float64:
		return fmt.Sprintf("%.6g", x)
	case complex64:
		re := real(x)
		im := imag(x)
		if re == 0 && !showZeroComponents {
			return fmt.Sprintf("%.6gi", im)
		} else if im == 0 && !showZeroComponents {
			return fmt.Sprintf("%.6g", re)
		} else {
			return fmt.Sprintf("(%.6g+%.6gi)", re, im)
		}
	case complex128:
		re := real(x)
		im := imag(x)
		if re == 0 && !showZeroComponents {
			return fmt.Sprintf("%.6gi", im)
		} else if im == 0 && !showZeroComponents {
			return fmt.Sprintf("%.6g", re)
		} else {
			return fmt.Sprintf("(%.6g+%.6gi)", re, im)
		}
	}

	return fmt.Sprint(quoteString(i))
}

func quoteString(i interface{}) interface{} {
	if s, ok := i.(string); ok {
		return strconv.Quote(s)
	} else {
		return i
	}
}

func (c constValue) String() string {
        return reflect.Value(c).String()
}
