package eval

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	intType reflect.Type = reflect.TypeOf(int(0))
	i8 reflect.Type = reflect.TypeOf(int8(0))
	i16 reflect.Type = reflect.TypeOf(int16(0))
	i32 reflect.Type = reflect.TypeOf(int32(0))
	i64 reflect.Type = reflect.TypeOf(int64(0))

	uintType reflect.Type = reflect.TypeOf(uint(0))
	u8 reflect.Type = reflect.TypeOf(uint8(0))
	u16 reflect.Type = reflect.TypeOf(uint16(0))
	u32 reflect.Type = reflect.TypeOf(uint32(0))
	u64 reflect.Type = reflect.TypeOf(uint64(0))

	f32 reflect.Type = reflect.TypeOf(float32(0))
	f64 reflect.Type = reflect.TypeOf(float64(0))
	c64 reflect.Type = reflect.TypeOf(complex64(0))
	c128 reflect.Type = reflect.TypeOf(complex128(0))

	boolType reflect.Type = reflect.TypeOf(bool(false))
	stringType reflect.Type = reflect.TypeOf(string(""))
)

// For each parameter in a builtin function, a bool parameter is passed
// to indicate if the corresponding value is typed. The boolean(s) appear
// after entire builtin parameter list.
//
// Builtin functions must return the builtin function reflect.Value, a
// bool indicating if the return value is typed, and an error if there was one.
// The returned Value must be valid
var builtinFuncs = map[string] reflect.Value {
	"complex": reflect.ValueOf(func(r, i reflect.Value, rt, it bool) (reflect.Value, bool, error) {
		rr, rerr := assignableValue(r, f64, rt)
		ii, ierr := assignableValue(i, f64, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex(rr.Float(), ii.Float())), rt || it, nil
		}

		rr, rerr = assignableValue(r, f32, rt)
		ii, ierr = assignableValue(i, f32, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex64(complex(rr.Float(), ii.Float()))), rt || it, nil
		}
		return reflect.Zero(c128), false, ErrBadComplexArguments{r, i}
	}),
	"real": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(real(zz.Complex())), zt, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(real(zz.Complex()))), zt, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"real", z}
		}
	}),
	"imag": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(imag(zz.Complex())), zt, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(imag(zz.Complex()))), zt, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"imag", z}
		}
	}),
	"append": reflect.ValueOf(builtinAppend),
	"cap"   : reflect.ValueOf(builtinCap),
	"len"   : reflect.ValueOf(builtinLen),
	"new"   : reflect.ValueOf(builtinNew),
	"panic" : reflect.ValueOf(builtinPanic),
}

var builtinTypes = map[string] reflect.Type{
	"int": intType,
	"int8": i8,
	"int16": i16,
	"int32": i32,
	"int64": i64,

	"uint": uintType,
	"uint8": u8,
	"uint16": u16,
	"uint32": u32,
	"uint64": u64,

	"float32": f32,
	"float64": f64,

	"complex64": c64,
	"complex128": c128,

	"bool": boolType,
	"byte": ByteType,
	"rune": RuneType,
	"string": stringType,

	"error": reflect.TypeOf(errors.New("")),
}

// FIXME: the real append is variadic. We can only handle one arg.

func builtinAppend(s, t reflect.Value, st, tt bool) (reflect.Value, bool, error) {
	if s.Kind() != reflect.Slice {
		return reflect.ValueOf(nil), true,
		errors.New(fmt.Sprintf("first argument to append must be a slice; " +
			"have %v", s.Type()))
	}
	stype, ttype := s.Type().Elem(), t.Type()
	if !ttype.AssignableTo(stype) {
		return reflect.ValueOf(nil), false,
		errors.New(fmt.Sprintf("cannot use type %v as type %v in append",
			ttype, stype))
	}
	return reflect.Append(s, t), true, nil
}

func builtinCap(v reflect.Value, vt bool) (reflect.Value, bool, error) {
	switch v.Kind() {
	case reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(v.Cap()), true, nil
	default:
		return reflect.Zero(intType), false,
		errors.New(fmt.Sprintf("invalid argument %v (type %v) for cap",
			v.Interface(), v.Type()))
	}
}

func builtinLen(z reflect.Value, zt bool) (reflect.Value, bool, error) {
	switch z.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return reflect.ValueOf(z.Len()), true, nil
	default:
		return reflect.ValueOf(nil), false, ErrBadBuiltinArgument{"len", z}
	}
}

func builtinNew(rtyp reflect.Value, bt bool) (reflect.Value, bool, error) {
	if typ, ok := rtyp.Interface().(reflect.Type); ok {
		return reflect.New(typ), true, nil
	} else {
		return reflect.ValueOf(nil), false, errors.New("new parameter is not a type")
	}
}

func builtinPanic(z reflect.Value, zt bool) (reflect.Value, bool, error) {
	// FIXME: we want results relative to the evaluated environment rather
	// than a panic inside the evaluator. We might use error, but panic's
	// parameter isn't the same as error's?
	panic(z.Interface())
	return reflect.ValueOf(nil), false, errors.New("Panic")
}
