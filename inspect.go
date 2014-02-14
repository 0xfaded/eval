package eval

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode"
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
	case reflect.Slice, reflect.Array:
		if val.Len() == 0 {
			return "[]"
		}
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
		if val.Type() == untypedNilType {
			return "nil"
		}
		n  := val.NumField()
		if n == 0 {
			return "{}"
		}
		str := val.Type().String()
		sep := " {\n\t"
		unexported := false
		for i := 0; i < n; i++ {
			name  := val.Type().Field(i).Name
			if unicode.IsLower([]rune(name)[0]) {
				unexported = true
				continue
			}
			field := val.Field(i)
			str += fmt.Sprintf("%s%s: %s", sep, name, Inspect(field))
			sep = ",\n\t"
		}
		if unexported {
			str += fmt.Sprintf("\n\t// unexported fields")
		}
		str += "\n}"
		return str

	case reflect.Ptr:
		// Internal const numbers
		i := val.Interface()
		if cn, ok := i.(*ConstNumber); ok {
			return fmt.Sprint(cn)
		} else if val.IsNil() {
			return "nil"
		} else {
			return "&" + Inspect(reflect.Indirect(val))
		}
	default:
		// FIXME: add more Kinds as folks are annoyed with the output of
		// the below:
		if val.CanInterface() {
			return fmt.Sprintf("%v", val.Interface())
		} else {
			return fmt.Sprintf("<%v>", val.Type())
		}
	}
}
