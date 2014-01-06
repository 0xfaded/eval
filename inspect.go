package eval

import (
	"fmt"
	"reflect"
	"strconv"
)

func InspectPtr(val reflect.Value) string {
	// fall back to %v when we panic here.
	defer func() string {
		recover()
		return fmt.Sprintf("%v", val)
	}()
	return "&" + Inspect(val.Elem())
}


// Inspect prints a reflect.Value the way you would enter it.
// Some like this should really be part of the reflect package.
func Inspect(val reflect.Value) string {

	if val.CanInterface() && val.Interface() == nil {
		return "nil"
	}
	switch val.Kind() {
	case reflect.String:
		return strconv.QuoteToASCII(val.String())
	case reflect.Bool, reflect.Int,	reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		if val.CanInterface() {
			fmt.Sprintf("%v", val.Interface())
		} else {
			fmt.Sprintf("<%v>", val.Type())
		}
	case reflect.Slice, reflect.Array:
		sep := "{"
		str := ""
		for i:=0; i < val.Len(); i++ {
			str += sep
			sep = ", "
			str += Inspect(val.Index(i))
		}
		str += "}"
		return str

	case reflect.Struct:
		str := "{"
		n  := val.NumField()
		for i := 0; i < n; i++ {
			if i > 0 { str += " " }
			name  := val.Type().Field(i).Name
			field := val.Field(i)
			str += fmt.Sprintf("%s: %s,", name, Inspect(field))
		}
		str += "}"
		return str

	case reflect.Ptr:
		if val.IsNil() {
			return "nil"
		} else {
			return "&" + Inspect(reflect.Indirect(val))
		}
	}


	// FIXME: add more Kinds as folks are annoyed with the output of
	// the below:
	if val.CanInterface() {
		return fmt.Sprintf("%v", val.Interface())
	} else {
		return fmt.Sprintf("<%v>", val.Type())
	}
}
