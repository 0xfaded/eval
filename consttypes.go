package interactive

import "reflect"

// Constant types, should not be used with the reflect package, but
// can annotate information needed for evaluating const expressions
type ConstType interface {
	reflect.Type
}

type ConstIntType struct { reflect.Type }
type ConstRuneType struct { reflect.Type }
type ConstFloatType struct { reflect.Type }
type ConstComplexType struct { reflect.Type }
type ConstStringType struct { reflect.Type }
type ConstNilType struct { reflect.Type }
type ConstBoolType struct { reflect.Type }

var (
	ConstInt = ConstIntType { reflect.TypeOf(0) }
	ConstRune = ConstRuneType { reflect.TypeOf('\000') }
	ConstFloat = ConstFloatType { reflect.TypeOf(0.0) }
	ConstComplex = ConstComplexType { reflect.TypeOf(0i) }
	ConstString = ConstStringType { reflect.TypeOf("") }
	ConstNil = ConstNilType { nil }
	ConstBool = ConstBoolType { reflect.TypeOf(false) }
)

func (ConstIntType) String() string { return "int" }
func (ConstRuneType) String() string { return "rune" }
func (ConstFloatType) String() string { return "float64" }
func (ConstComplexType) String() string { return "complex128" }
func (ConstStringType) String() string { return "string" }
func (ConstNilType) String() string { return "<T>" }
func (ConstBoolType) String() string { return "bool" }

// Returns the ConstType of a binary, non-boolean, expression invalving const types of
// x and y.
func promoteConsts(ctx *Ctx, x, y ConstType, yexpr Expr, yval reflect.Value) (ConstType, error) {
	switch x.(type) {
	case ConstIntType:
		switch y.(type) {
		case ConstIntType:
			return x, nil
		case ConstRuneType, ConstFloatType, ConstComplexType:
			return y, nil
		}
	case ConstRuneType:
		switch y.(type) {
		case ConstIntType, ConstRuneType:
			return x, nil
		case ConstFloatType, ConstComplexType:
			return y, nil
		}
	case ConstFloatType:
		switch y.(type) {
		case ConstIntType, ConstRuneType, ConstFloatType:
			return x, nil
		case ConstComplexType:
			return y, nil
		}
	case ConstComplexType:
		switch y.(type) {
		case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
			return x, nil
		}
	case ConstStringType:
		if _, ok := y.(ConstStringType); ok {
			return x, nil
		}
	case ConstNilType:
		if _, ok := y.(ConstNilType); ok {
			return x, nil
		}
	case ConstBoolType:
		if _, ok := y.(ConstBoolType); ok {
			return x, nil
		}
	}
	return nil, ErrBadConversion{at(ctx, yexpr), y, x, yval}
}

func convertConstToTyped(ctx *Ctx, from ConstType, c constValue, to reflect.Type, expr Expr) (
	v reflect.Value, errs []error) {

	v = reflect.New(to).Elem()

	switch from.(type) {
	case ConstIntType, ConstRuneType, ConstFloatType, ConstComplexType:
		underlying := reflect.Value(c).Interface().(*BigComplex)
		switch to.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			var errs []error
			i, truncation, overflow := underlying.Int(to.Bits())
			if truncation {
				errs = append(errs, ErrTruncatedConstant{at(ctx, expr), ConstInt, underlying})
			}
			if overflow {
				integer, _ := underlying.Integer()
				errs = append(errs, ErrOverflowedConstant{at(ctx, expr), from, to, integer})
			}
			v.SetInt(i)
			return v, errs

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			u, truncation, overflow := underlying.Uint(to.Bits())
			if truncation {
				errs = append(errs, ErrTruncatedConstant{at(ctx, expr), ConstInt, underlying})
			}
			if overflow {
				integer, _ := underlying.Integer()
				errs = append(errs, ErrOverflowedConstant{at(ctx, expr), from, to, integer})
			}
			v.SetUint(u)
			return v, errs

		case reflect.Float32, reflect.Float64:
			f, truncation, _ := underlying.Float64()
			if truncation {
				errs = append(errs, ErrTruncatedConstant{at(ctx, expr), ConstFloat, underlying})
			}
			v.SetFloat(f)
			return v, errs

		case reflect.Complex64, reflect.Complex128:
			cmplx, _ := underlying.Complex128()
			v.SetComplex(cmplx)
			return v, errs
		}
	case ConstStringType:
		// Check Kind == String ourselves, as reflect.Value.String() doesn't panic
		// on non string values
		if v.Type().Kind() != reflect.String {
			panic("go-interactive: string constant has wrong underlying type")
		}
		v.SetString(reflect.Value(c).String())

	case ConstBoolType:
		v.SetBool(reflect.Value(c).Bool())

	case ConstNilType:
		// Unfortunately there is no reflect.Type.CanNil()
		switch to.Kind() {
		case reflect.Chan, reflect.Func, reflect.Interface,
			reflect.Map, reflect.Ptr, reflect.Slice:

			// v is already nil
			return v, nil
		}
	}
	return v, []error{ErrBadConversion{at(ctx, expr), from, to, reflect.Value(c)}}
}

