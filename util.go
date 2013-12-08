package interactive

import (
	"errors"
	"fmt"
	"reflect"
)

// Wrapper around reflect.New that must be used when creating possibly anonymous types.
// Returned value is a settable value (not Ptr) of type t
func makeAnonValue(t reflect.Type) reflect.Value {
	// It is impossible to define an anonymous composite type where an anonymous type is 
	// an element of a named type. Therefore we only need to dig until we find a named type
	if arrt, ok := t.(arrayType); ok {
		// Create a slice of the underlying elem
		v := reflect.New(reflect.SliceOf(arrt.Type.Elem())).Elem()
		v.Set(reflect.MakeSlice(reflect.SliceOf(arrt.Type.Elem()), t.Len(), t.Len()))

		for i := 0; i < t.Len(); i += 1 {
			v.Index(i).Set(makeAnonValue(arrt.Type.Elem()))
		}
		return reflect.ValueOf(arrayValue{arrt, v})
	} else {
		return reflect.New(t).Elem()
	}
}

// Wrapper around reflect.SliceOf that must be used when creating slices of possibly anonymous elts.
func anonSliceOf(elt reflect.Type) reflect.Type {
	switch unwrapped := elt.(type) {
	case arrayType:
		return reflect.SliceOf(unwrapped.Type)
	default:
		return reflect.SliceOf(elt)
	}
}

// Produces a value assignable to a value of type 'to'
func assignableValue(x reflect.Value, to reflect.Type, xTyped bool) (reflect.Value, error) {
	if xx, err := unwrappedAssignableValue(x, to); err == nil {
		x = xx
	}
	if xTyped {
		if x.Type().AssignableTo(to) {
			return x, nil
		}
	} else {
		var err error
		if x, err = promoteUntypedNumeral(x, to); err == nil {
			return x, nil
		}
	}
	return x, errors.New(fmt.Sprintf("Cannot convert %v to type %v", x, to))
}

// Similar to assignableValue, but only performs unwrapping
func unwrappedAssignableValue(x reflect.Value, to reflect.Type) (reflect.Value, error) {
	if arr, _ := recoverArray(x); arr.IsValid() {
		return promoteSliceToArray(arr, to)
	} else if x.Type() == to {
		return x, nil
	}
	return reflect.Value{}, errors.New(fmt.Sprintf("Cannot convert %v to type %v", x, to))
}

func setTypedValue(dst, src reflect.Value, srcTyped bool) error {
	fmt.Printf("set %v %v\n", dst, src)
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

func promoteSliceToArray(slice reflect.Value, to reflect.Type) (reflect.Value, error) {
	fmt.Printf(">> %v %v\n", slice, to)
	if to.Kind() != reflect.Array || slice.Len() != to.Len() {
		fmt.Printf(">> %v %v %v\n", to.Kind(), slice, to)
		return reflect.Value{}, errors.New(fmt.Sprintf("cannot convert %v to %v", slice, to))
	}

	v := reflect.New(to).Elem()
	for i := 0; i < to.Len(); i += 1 {
		if elem, err := unwrappedAssignableValue(slice.Index(i), to.Elem());
			err != nil || elem.Type() != to.Elem() {
			fmt.Printf(">>> %v %v %v\n", err, elem.Type(), to.Elem())
			return reflect.Value{}, errors.New(fmt.Sprintf("cannot convert %v to %v", slice, to))
		} else {
			v.Index(i).Set(elem)
		}
	}
	fmt.Printf(">>>> %v\n", v)
	return v, nil
}
