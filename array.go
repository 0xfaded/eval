package interactive

import (
	"reflect"
	"fmt"
)

// Array Types are internally represented by arrayType which satisfies the
// reflect.Type interface.
// reflect.Value is not an interface. To keep the runtime homogeneous,
// arrays are represented by reflect.ValueOf(arrayType(reflect.Value))
// Use recoverArray() to avoid the specifics

type arrayType struct {
	elt reflect.Type
	size int
	reflect.Type
}

type arrayValue struct {
	Type arrayType
	reflect.Value
}

func arrayOf(elt reflect.Type, size int) reflect.Type {
	return arrayType {
		elt: elt,
		size: size,
		Type: anonSliceOf(elt),
	}
}

func (arr arrayType) Elem() reflect.Type {
	return arr.elt
}

func (arr arrayType) Len() int {
	return arr.size
}

func (arr arrayType) Kind() reflect.Kind {
	return reflect.Array
}

func (arr arrayType) String() string {
	return fmt.Sprintf("[%d]%v", arr.size, arr.elt)
}

func (arr arrayValue) String() string {
	return fmt.Sprintf("%v", arr.Value.Interface())
}

// Returns the array wrapped by the runtime, along with its type. If the value
// was not an array, the Zero Value is returned
func recoverArray(wrapped reflect.Value) (reflect.Value, arrayType) {
	if av, ok := wrapped.Interface().(arrayValue); !ok {
		return reflect.Value{}, arrayType{}
	} else {
		return av.Value, av.Type
	}
}
