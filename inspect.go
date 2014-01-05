package eval

import (
	"fmt"
	"reflect"
	"strconv"
)

// Inspect prints a reflect.Value the way you would enter it.
// Some like this should really be part of the reflect package.
func Inspect(val reflect.Value) string {
	if val.Interface() == nil {
		return "nil"
	}
	switch val.Kind() {
	case reflect.String:
		return strconv.QuoteToASCII(val.String())
	case reflect.Bool, reflect.Int,	reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		fmt.Sprintf("%v", val.Interface())
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
	}
	// FIXME: add more Kinds as folks are annoyed with the output of
	// the below:
	return fmt.Sprintf("%v", val.Interface())
}
