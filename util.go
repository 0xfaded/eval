package interactive

import (
	"errors"
	"fmt"
	"reflect"
)

func assignableValue(x reflect.Value, to reflect.Type, xTyped bool) (reflect.Value, error) {
	var err error
	if xTyped {
		if x.Type().AssignableTo(to) {
			return x, nil
		}
	} else {
		if x, err = promoteUntypedNumeral(x, to); err == nil {
			return x, nil
		}
	}
	return x, errors.New(fmt.Sprintf("Cannot convert %v to type %v", x, to))
}

func setTypedValue(dst, src reflect.Value, srcTyped bool) error {
	if assignable, err := assignableValue(src, dst.Type(), srcTyped); err != nil {
		return errors.New(fmt.Sprintf("Cannot assign %v = %v", dst, src))
	} else {
		dst.Set(assignable)
		return nil
	}
}

func makeSliceWithValues(elts []reflect.Value, sliceType reflect.Type) (reflect.Value, error) {
	slice := reflect.MakeSlice(sliceType, len(elts), len(elts))
	for i := 0; i < slice.Len(); i += 1 {
		if err := setTypedValue(slice.Index(i), elts[i], true); err != nil {
			return reflect.Value{}, nil
		}
	}
	return slice, nil
}


// Only considers untyped kinds produced by our runtime. Assumes input type is unnamed
func isUntypedNumeral(x reflect.Value) bool {
	switch x.Kind() {
	case reflect.Int64, reflect.Float64, reflect.Complex128:
		return true
	default:
		return false
	}
}

func promoteUntypedNumeral(untyped reflect.Value, to reflect.Type) (reflect.Value, error) {
	// The only valid promotion that cannot be directly converted is int|float -> complex
	if untyped.Type().ConvertibleTo(to) {
		return untyped.Convert(to), nil
	} else if to.Kind() == reflect.Complex64 || to.Kind() == reflect.Complex128 {
		floatType := reflect.TypeOf(float64(0))
		if untyped.Type().ConvertibleTo(floatType) {
			return reflect.ValueOf(complex(untyped.Convert(floatType).Float(), 0)), nil
		}
	}
	return reflect.Value{}, errors.New(fmt.Sprintf("cannot convert %v to %v", untyped, to))
}

// Only considers untyped kinds produced by our runtime. Assumes input type is unnamed
func promoteUntypedNumerals(x, y reflect.Value) (reflect.Value, reflect.Value) {
	switch x.Kind() {
	case reflect.Int64:
		switch y.Kind() {
		case reflect.Int64:
			return x, y
		case reflect.Float64:
			return x.Convert(y.Type()), y
		case reflect.Complex128:
			return reflect.ValueOf(complex(float64(x.Int()), 0)), y
		}
	case reflect.Float64:
		switch y.Kind() {
		case reflect.Int64:
			return x, y.Convert(x.Type())
		case reflect.Float64:
			return x, y
		case reflect.Complex128:
			return reflect.ValueOf(complex(x.Float(), 0)), y
		}
	case reflect.Complex128:
		switch y.Kind() {
		case reflect.Int64:
			return x, reflect.ValueOf(complex(float64(y.Int()), 0))
		case reflect.Float64:
			return x, reflect.ValueOf(complex(y.Float(), 0))
		case reflect.Complex128:
			return x, y
		}
	}
	panic(fmt.Sprintf("runtime: bad untyped numeras %v and %v", x, y))
}

