package eval

import (
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/token"
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

	setKnownType(t knownType)
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
	source envSource
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

	// length of array or slice literal
	length int

	// indices specified in array or slice literal.
	// terminated by a {-1, -1} after last index.
	indices []struct{pos, index int}

	// fields of struct mapping position to struct index
	fields []int
}

type ParenExpr struct {
	*ast.ParenExpr
	knownType
	constValue
}

type SelectorExpr struct {
	*ast.SelectorExpr
	Sel *Ident
	knownType
	constValue

	// if not "", this is a package selector
	pkgName string

	// if not nil, this is a struct field
	field []int

	// if true, the method was found on the pointer
	// for this type, not the type itself.
	isPtrReceiver bool

	// the method index
	method int
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

	// Is this a type conversion
	isTypeConversion bool

	// Is this a builtin call
	isBuiltin bool

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
	dir reflect.ChanDir
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

func (*BadExpr) setKnownType(t knownType)      { panic("eval: cannot set knownType of BadExpr") }
func (*FuncLit) setKnownType(t knownType)      { panic("eval: cannot set knownType of FuncLit") }
func (*KeyValueExpr) setKnownType(t knownType) { panic("eval: cannot set knownType of KeyValueExpr") }

func (e *BasicLit) setKnownType(t knownType)       { e.knownType = t }
func (e *BinaryExpr) setKnownType(t knownType)     { e.knownType = t }
func (e *CallExpr) setKnownType(t knownType)       { e.knownType = t }
func (e *Ellipsis) setKnownType(t knownType)       { e.knownType = t }
func (e *CompositeLit) setKnownType(t knownType)   { e.knownType = t }
func (e *SelectorExpr) setKnownType(t knownType)   { e.knownType = t }
func (e *Ident) setKnownType(t knownType)          { e.knownType = t }
func (e *IndexExpr) setKnownType(t knownType)      { e.knownType = t }
func (e *ParenExpr) setKnownType(t knownType)      { e.knownType = t }
func (e *SliceExpr) setKnownType(t knownType)      { e.knownType = t }
func (e *TypeAssertExpr) setKnownType(t knownType) { e.knownType = t }
func (e *StarExpr) setKnownType(t knownType)       { e.knownType = t }
func (e *UnaryExpr) setKnownType(t knownType)      { e.knownType = t }
func (e *ArrayType) setKnownType(t knownType)      { e.knownType = t }
func (e *StructType) setKnownType(t knownType)     { e.knownType = t }
func (e *FuncType) setKnownType(t knownType)       { e.knownType = t }
func (e *InterfaceType) setKnownType(t knownType)  { e.knownType = t }
func (e *MapType) setKnownType(t knownType)        { e.knownType = t }
func (e *ChanType) setKnownType(t knownType)       { e.knownType = t }

// Does not assert that c is a valid const value type
// Should be *ConstNumber, bool, or string
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

func (ellipsis *Ellipsis) String() string {
	if ellipsis.Elt != nil {
		return fmt.Sprintf("...%v", ellipsis.Elt)
	} else {
		return "..."
	}
}

func (basicLit *BasicLit) String() string {
	if basicLit.IsConst() {
		return sprintConstValue(basicLit.KnownType()[0], basicLit.Const(), true)
	}
	return basicLit.Value
}

func (funcLit *FuncLit) String() string { return "func literal" }

func (lit *CompositeLit) String() string {
	kt := lit.KnownType()
	if kt == nil {
		// This matches gc formatting for unchecked nodes.
		// For example, complex([]int{}) produces
		// missing argument to complex - complex(composite literal, <N>)
		return "composite literal"
	}
	t := kt[0]
	if t.Name() != "" {
		return fmt.Sprintf("%s literal", t.Name())
	}

	switch t.Kind() {
	case reflect.Slice:
		return fmt.Sprintf("[]%v literal", t.Elem())
	case reflect.Map:
		return fmt.Sprintf("map[%v]%v literal", t.Key(), t.Elem())
	default:
		return "TODO composite lit"
	}
}

func (parenExpr *ParenExpr) String() string {
	if parenExpr.IsConst() {
		return sprintConstValue(parenExpr.KnownType()[0], parenExpr.Const(), true)
	}
	return fmt.Sprintf("(%v)", skipSuperfluousParens(parenExpr.X.(Expr)))
}

func (selectorExpr *SelectorExpr) String() string {
	return fmt.Sprintf("%s.%v", selectorExpr.X, selectorExpr.Sel)
}

func (index *IndexExpr) String() string {
	return fmt.Sprintf("%v[%v]", index.X, index.Index)
}

func (slice *SliceExpr) String() string {
	// TODO update for :: with go 1.2 upgrade
	var low, high string
	if slice.Low != nil {
		low = fmt.Sprint(slice.Low)
	}
	if slice.High != nil {
		high = fmt.Sprint(slice.High)
	}
	return fmt.Sprintf("%v[%v:%v]", slice.X, low, high)
}

func (assert *TypeAssertExpr) String() string {
	return fmt.Sprintf("%v.(%s)", assert.X, assert.Type)
}

func (call *CallExpr) String() string {
	if call.isTypeConversion || call.isBuiltin {
		if len(call.Args) == 0 {
			// missing argument error
			return fmt.Sprintf("%v()", call.Fun)
		} else if len(call.Args) > 1 || call.isBuiltin {
			// too many arguments error
			s := fmt.Sprintf("%v", call.Fun)
			sep := "("
			for _, arg := range call.Args {
				s += fmt.Sprintf("%v%v", sep, arg)
				sep = ", "
			}
			if call.argNEllipsis {
				s += "..."
			}
			return s + ")"
		} else {
			if call.IsConst() {
				return sprintConstValue(call.KnownType()[0], call.Const(), true)
			}
			return fmt.Sprintf("%v(%v)", call.KnownType()[0], call.Args[0])
		}
	} else {
		return fmt.Sprintf("%v()", call.Fun)
	}
}

func (star *StarExpr) String() string {
	return fmt.Sprintf("*%v", star.X)
}

func (unary *UnaryExpr) String() string {
	operand := skipSuperfluousParens(unary.X.(Expr))
	if unary.Op == token.AND {
		if _, ok := unary.X.(*StarExpr); ok {
			return fmt.Sprintf("&(%v)", operand)
		} else {
			return fmt.Sprintf("&%v", operand)
		}
	}
	return fmt.Sprintf("%v %v", unary.Op, operand)
}

func (binary *BinaryExpr) String() string {
	left := simplifyBinaryChildExpr(binary, binary.X.(Expr))
	right := simplifyBinaryChildExpr(binary, binary.Y.(Expr))

	return fmt.Sprintf("%v %v %v", left, binary.Op, right)
}

func (keyValueExpr *KeyValueExpr) String() string { return "TODO  keyValueExpr.KeyValueExpr" }

func (arrayType *ArrayType) String() string {
	if arrayType.Len != nil {
		return fmt.Sprintf("[%v]%v", arrayType.Len, arrayType.Elt)
	} else {
		return fmt.Sprintf("[]%v", arrayType.Elt)
	}
}

func (structType *StructType) String() string { return "TODO  structType.StructType" }
func (funcType *FuncType) String() string { return "TODO  funcType.FuncType" }
func (interfaceType *InterfaceType) String() string { return "TODO  interfaceType.InterfaceType" }

func (mapType *MapType) String() string {
	return fmt.Sprintf("map[%v]%v", mapType.Key, mapType.Value)
}

func (chanType *ChanType) String() string {
	value := fmt.Sprint(chanType.Value)
	switch chanType.dir {
	case reflect.SendDir:
		return "chan<- " + value
	case reflect.RecvDir:
		return "<-chan " + value
	default:
		return "chan " + value
	}
}

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
	if t == RuneType {
		// When a rune constant is larger than 32 bits, it will be stored
		// as a ConstRune, but should be printed as a Rune
		if n, ok := v.Interface().(*ConstNumber); ok {
			return fmt.Sprintf("rune(%v)", n.Value.StringShow0i(false))
		}
	}
	if isTypeDisplayed(t) {
		return fmt.Sprintf("%v(%v)", t, sprintConstUntypedValue(v, showZeroComponents))
	} else {
		return fmt.Sprintf("%v", sprintConstUntypedValue(v, showZeroComponents))
	}
}

func sprintConstUntypedValue(v reflect.Value, showZeroComponents bool) string {
	i := v.Interface()
	switch x := i.(type) {
	case *ConstNumber:
		return x.StringShow0i(false)
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

func isTypeDisplayed(t reflect.Type) bool {
	switch t {
	case ConstInt, ConstRune, ConstFloat, ConstComplex,
		ConstString, ConstNil, ConstBool,
		intType, i8, i16, i32, i64,
		uintType, u8, u16, u32, u64,
		f32, f64, c64, c128,
		boolType, stringType:
		     return false
	}
	return true
}

