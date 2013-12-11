package interactive

import (
	"reflect"

	"go/ast"
)

// Annotated ast.Expr nodes
type Expr interface {
	ast.Expr

	// The type of this expression if known. Certain expr types have special interpretations
	// Constant expr: a value of nil implies the constant is untyped
	// Ellipsis expr: a single reflect.Type represents the type of all unpacked values
	KnownType() []reflect.Type

	// Returns true if this expression evaluates to a constant. In this
	// case, Value() will return the evalutated constant. Nodes below
	// this expression can be ignored.
	IsConst() bool

	// Returns the const value, if known.
	Const() reflect.Value
}

type knownType []reflect.Type
type constValue reflect.Value

type BadExpr struct {
	*ast.BadExpr
}

type Ident struct {
	*ast.Ident
	knownType
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

	// Does this function take  single, multivalued expession that is unpacked as its arguments
	isSplat bool

	// Is the function type variadic
	isVariadic bool

	// Is an ellipsis expression used to unpack variadic arguments
	hasEllipsis bool
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
func (*Ident) IsConst() bool          { return false }
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
func (*Ident) Const() reflect.Value          { return reflect.Value{} }
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
