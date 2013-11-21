package interactive

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go/ast"
	"go/parser"
	"go/token"

	//"github.com/sbinet/go-readline/pkg/readline"
)

type pkg *env

type env struct {
	// Values
	vars map[string] reflect.Value
	consts map[string] reflect.Value
	funcs map[string] reflect.Value

	// Types
	types map[string] reflect.Type

	// Packages
	pkgs map[string] pkg
}

var (
	f32 reflect.Type = reflect.TypeOf(float32(0))
	f64 reflect.Type = reflect.TypeOf(float64(0))
	c64 reflect.Type = reflect.TypeOf(complex64(0))
	c128 reflect.Type = reflect.TypeOf(complex128(0))
)

// Builtin funcs receive Value bool pairs indicating if the corrosoponding value is typed
// They must return a Value, a bool indicating if the value is typed, and an error
// The returned Value must be valid
var builtinFuncs = map[string] reflect.Value {
	"complex": reflect.ValueOf(func(r, i reflect.Value, rt, it bool) (reflect.Value, bool, error) {
		rr, rerr := assignableValue(r, f64, rt)
		ii, ierr := assignableValue(i, f64, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex(rr.Float(), ii.Float())), true, nil
		}

		rr, rerr = assignableValue(r, f32, rt)
		ii, ierr = assignableValue(i, f32, it)
		if rerr == nil && ierr == nil {
			return reflect.ValueOf(complex64(complex(rr.Float(), ii.Float()))), true, nil
		}
		return reflect.Zero(c128), false, ErrBadComplexArguments{r, i}
	}),
	"real": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(real(zz.Complex())), true, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(real(zz.Complex()))), true, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"real", z}
		}
	}),
	"imag": reflect.ValueOf(func(z reflect.Value, zt bool) (reflect.Value, bool, error) {
		if zz, err := assignableValue(z, c128, zt); err == nil {
			return reflect.ValueOf(imag(zz.Complex())), true, nil
		} else if zz, err := assignableValue(z, c64, zt); err == nil {
			return reflect.ValueOf(float32(imag(zz.Complex()))), true, nil
		} else {
			return reflect.Zero(f64), false, ErrBadBuiltinArgument{"imag", z}
		}
	}),
}

var builtinTypes = map[string] reflect.Type{
	"int8": reflect.TypeOf(int8(0)),
	"int16": reflect.TypeOf(int16(0)),
	"int32": reflect.TypeOf(int32(0)),
	"int64": reflect.TypeOf(int64(0)),

	"uint8": reflect.TypeOf(uint8(0)),
	"uint16": reflect.TypeOf(uint16(0)),
	"uint32": reflect.TypeOf(uint32(0)),
	"uint64": reflect.TypeOf(uint64(0)),

	"float32": reflect.TypeOf(float32(0)),
	"float64": reflect.TypeOf(float64(0)),

	"complex64": reflect.TypeOf(complex64(0)),
	"complex128": reflect.TypeOf(complex128(0)),

	"bool": reflect.TypeOf(bool(false)),
	"byte": reflect.TypeOf(byte(0)),
	"rune": reflect.TypeOf(rune('☃')),
	"string": reflect.TypeOf(""),

	"error": reflect.TypeOf(errors.New("")),
}

type ErrInvalidOperands struct {
	x reflect.Value
	op token.Token
	y reflect.Value
}

type ErrBadFunArgument struct {
	fun reflect.Value
	index int
	value reflect.Value
}

type ErrBadComplexArguments struct {
	a, b reflect.Value
}

type ErrBadBuiltinArgument struct {
	fun string
	value reflect.Value
}

type ErrWrongNumberOfArgs struct {
	fun reflect.Value
	numArgs int
}

type ErrMultiInSingleContext struct {
	vals []reflect.Value
}

func Run(vars, consts, funcs map[string] reflect.Value, types map[string] reflect.Type) {
	expr, err := parser.ParseExpr("log.New(stdout, \"Awesome> \", log.Ltime).Printf(\":)\")")
	fmt.Printf("err %v\n", err)
	logPkg := &env {
		vars: vars,
		consts: consts,
		funcs: funcs,
		types: types,
	}
	env := env {
		vars: make(map[string] reflect.Value),
		consts: make(map[string] reflect.Value),
		funcs: make(map[string] reflect.Value),
		types: make(map[string] reflect.Type),
		pkgs: map[string] pkg { "log": logPkg },
	}

	vals, typed, err := evalExpr(expr, &env)
	fmt.Printf("err: %v\n", err)
	fmt.Printf("(%v) %+v\n", typed, vals[0].Interface())
}

func evalExpr(expr ast.Expr, env *env) ([]reflect.Value, bool, error) {
	fmt.Printf("%T %+v\n", expr, expr)
	switch node := expr.(type) {
	case *ast.Ident:
		v, typed, err := evalIdentExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.Ellipsis:
	case *ast.BasicLit:
		v, typed, err := evalBasicLit(node)
		return []reflect.Value{v}, typed, err
	case *ast.FuncLit:
	case *ast.CompositeLit:
		v, typed, err := evalCompositeLit(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.ParenExpr:
		return evalExpr(node.X, env)
	case *ast.SelectorExpr:
		v, typed, err := evalSelectorExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.IndexExpr:
	case *ast.SliceExpr:
	case *ast.TypeAssertExpr:
	case *ast.CallExpr:
		return evalCallExpr(node, env)
	case *ast.StarExpr:
	case *ast.UnaryExpr:
	case *ast.BinaryExpr:
		v, typed, err := evalBinaryExpr(node, env)
		return []reflect.Value{v}, typed, err
	case *ast.KeyValueExpr:
	}
	return []reflect.Value{reflect.ValueOf("Alice")}, true, nil
}

func evalType(expr ast.Expr, env *env) (reflect.Type, error) {
	fmt.Printf("%T %+v\n", expr, expr)
	switch node := expr.(type) {
	case *ast.Ident:
		if t, ok := env.types[node.Name]; ok {
			return t, nil
		} else if t, ok := builtinTypes[node.Name]; ok {
			return t, nil
		} else {
			return t, errors.New("undefined type: " + node.Name)
		}
	case *ast.ArrayType:
		return nil, errors.New("array types not implemented")
	case *ast.StructType:
		return nil, errors.New("struct types not implemented")
	case *ast.FuncType:
		return nil, errors.New("func types not implemented")
	case *ast.InterfaceType:
		return nil, errors.New("interface types not implemented")
	case *ast.MapType:
		return nil, errors.New("map types not implemented")
	case *ast.ChanType:
		return nil, errors.New("chan types not implemented")
	default:
		return nil, errors.New(fmt.Sprintf("Type: Bad type (%+v)", node))
	}
}

func evalIdentExpr(ident *ast.Ident, env *env) (reflect.Value, bool, error) {
	name := ident.Name
	if v, ok := env.vars[name]; ok {
		return v.Elem(), true, nil
	} else if v, ok := env.consts[name]; ok {
		return v, true, nil
	} else if v, ok := env.funcs[name]; ok {
		return v, true, nil
	} else if v, ok := builtinFuncs[name]; ok {
		return v, false, nil
	} else if p, ok := env.pkgs[name]; ok {
		return reflect.ValueOf(p), true, nil
	} else {
		return reflect.Value{}, false, errors.New(fmt.Sprintf("%s undefined", name))
	}
}

func evalBasicLit(lit *ast.BasicLit) (reflect.Value, bool, error) {
	switch lit.Kind {
	case token.STRING:
		return reflect.ValueOf(lit.Value), false, nil
	case token.INT:
		i, err := strconv.ParseInt(lit.Value, 0, 0)
		return reflect.ValueOf(i), false, err
	case token.FLOAT:
		f, err := strconv.ParseFloat(lit.Value, 64)
		return reflect.ValueOf(f), false, err
	case token.IMAG:
		f, err := strconv.ParseFloat(lit.Value[:len(lit.Value)-1], 64)
		return reflect.ValueOf(complex(0, f)), false, err
	default:
		return reflect.Value{}, false, errors.New(fmt.Sprintf("BasicLit: Bad token type (%+v)", lit))
	}
}

func evalCompositeLit(lit *ast.CompositeLit, env *env) (reflect.Value, bool, error) {
	t, err := evalType(lit.Type, env)
	if err != nil {
		return reflect.Value{}, true, err
	}

	switch t.Kind() {
	//case reflect.Array:
	//case reflect.Map:
	//case reflect.Slice:
	case reflect.Struct:
		vp := reflect.New(t)
		v := vp.Elem()

		if len(lit.Elts) == 0 {
			return v, true, nil
		}

		_, pairs := lit.Elts[0].(*ast.KeyValueExpr)
		for i, elt := range lit.Elts {
			var field, value reflect.Value
			var typed bool
			var fname string
			if kv, ok := elt.(*ast.KeyValueExpr); ok != pairs {
				return v, true, errors.New("Elements are either all key value pairs or not")
			} else if pairs {
				if k, ok := kv.Key.(*ast.Ident); !ok {
					return v, true, errors.New(fmt.Sprintf("Invalid key node %v %T", kv.Key, kv.Key))
				} else if f := v.FieldByName(k.Name); !f.IsValid() {
					return v, true, errors.New(t.Name() + " has no field " + k.Name)
				} else if fv, ft, err := evalExpr(kv.Value, env); err != nil {
					return v, true, err
				} else {
					field, value = f, fv[0]
					typed = ft
					tfield, _ := v.Type().FieldByName(k.Name)
					fname = tfield.Name
				}
			} else {
				if i >= v.NumField() {
					return v, true, errors.New("Too many elements for struct " + t.Name())
				} else if _, ok := elt.(*ast.KeyValueExpr); ok {
					return v, true, errors.New("Elements are either all key value pairs or not")
				} else if fv, ft, err := evalExpr(elt, env); err != nil {
					return v, true, err
				} else {
					field, value = v.Field(i), fv[0]
					typed = ft
					fname = v.Type().Field(i).Name
				}
			}

			if err := setTypedValue(field, value, typed); err != nil {
				return v, true, errors.New(fmt.Sprintf("cannot set %s.%s : %v",
					t.Name(), fname, value.String()))
			}
		}
		return v, true, nil

	default:
		return reflect.Value{}, true, errors.New(fmt.Sprintf("invalid type for composite literal %s", t.Name()))

	}
}

func evalCallExpr(call *ast.CallExpr, env *env) ([]reflect.Value, bool, error) {
	fmt.Printf("==%+v %+v %T %T\n", call.Fun, call.Args, call.Fun, call.Args)
	if t, err := evalType(call.Fun, env); err == nil {
		if v, typed, err := evalCallTypeExpr(t, call, env); err != nil {
			return []reflect.Value{}, false, err
		} else {
			return []reflect.Value{v}, typed, nil
		}
	} else if fun, _, err := evalExpr(call.Fun, env); err == nil {
		return evalCallFunExpr(fun[0], call, env)
	} else {
		return []reflect.Value{}, false, err
	}
}

func evalCallTypeExpr(t reflect.Type, call *ast.CallExpr, env *env) (reflect.Value, bool, error) {
	var r reflect.Value
	if call.Args == nil {
		return r, false, errors.New(fmt.Sprintf("missing argument to conversion to %v", t))
	} else if len(call.Args) > 1 {
		return r, false, errors.New(fmt.Sprintf("too many arguments to conversion to %v", t))
	} else if arg, typed, err := evalExpr(call.Args[0], env); err != nil {
		return r, false, err
	} else if cast, err := assignableValue(arg[0], t, typed); err != nil {
		return r, false, err
	} else {
		return cast, true, nil
	}
}

func evalCallFunExpr(fun reflect.Value, call *ast.CallExpr, env *env) ([]reflect.Value, bool, error) {
	fmt.Printf("==%+v %+v %T %T\n", call.Fun, call.Args, call.Fun, call.Args)

	var err error
	var v, out []reflect.Value
	var typed bool
	if v, typed, err = evalExpr(call.Fun, env); err != nil {
		return []reflect.Value{}, false, err
	}
	if v[0].Kind() != reflect.Func {
		return []reflect.Value{}, false, errors.New(fmt.Sprintf("Cannot call %v", v[0]))
	}
	builtin := !typed

	// Special case handling doesn't play well with nil Args
	ftype := v[0].Type()
	if call.Args == nil {
		if ftype.NumIn() == 0 {
			out = v[0].Call([]reflect.Value{})
			return out, true, nil
		} else {
			return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], 0}
		}
	}

	args := make([][]reflect.Value, len(call.Args))
	atyped := make([]bool, len(call.Args))

	for i := range call.Args {
		if args[i], atyped[i], err = evalExpr(call.Args[i], env); err != nil {
			return []reflect.Value{}, false, err
		}
	}

	_, firstArgIsFun := call.Args[0].(*ast.CallExpr)
	// Special case for f(g()), where g may return multiple values
	if len(call.Args) == 1 && firstArgIsFun {
		splat := make([][]reflect.Value, len(args[0]))
		atyped = make([]bool, len(args[0]))
		for i := range args[0] {
			splat[i] = []reflect.Value{args[0][i]}
			atyped[i] = true
		}
		args = splat
	}

	// Parse args into a slice suitable for calling the function
	actualNumIn := ftype.NumIn()
	if builtin {
		// See builtinFuncs comment 
		actualNumIn /= 2
	}

	in := make([]reflect.Value, actualNumIn)
	intyped := make([]bool, actualNumIn)

	if !ftype.IsVariadic() && len(call.Args) == actualNumIn {
		for i := range in {
			if len(args[i]) > 1 {
				return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
			}
			in[i] = args[i][0]
			intyped[i] = atyped[i]
		}
	} else if ftype.IsVariadic() && actualNumIn-1 <= len(call.Args) {
		var i int
		for i = 0; i < len(in)-1; i += 1 {
			if len(args[i]) > 1 {
				return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
			}
			in[i] = args[i][0]
			intyped[i] = atyped[i]
		}
		fmt.Printf("%d %d %v\n", i, len(call.Args), call.Ellipsis)
		if i == len(call.Args)-1 && call.Ellipsis != token.NoPos {
			// Assert this indeed is the ellipsis
			_ = call.Args[i].(*ast.Ellipsis)
			in[i], err = makeSliceWithValues(args[i], ftype.In(i))
			intyped[i] = true
			if err != nil {
				return []reflect.Value{}, false, ErrBadFunArgument{v[0], i, in[i]}
			}
		} else if i <= len(call.Args) && call.Ellipsis == token.NoPos {
			remainingArgs := len(call.Args) - actualNumIn + 1
			in[i] = reflect.MakeSlice(ftype.In(i), remainingArgs, remainingArgs)

			intyped[i] = true
			etype := in[i].Type().Elem()
			for j := i; j < len(call.Args); j += 1 {
				if len(args[j]) > 1 {
					return []reflect.Value{}, false, ErrMultiInSingleContext{args[i]}
				} else if arg, err := assignableValue(args[j][0], etype, atyped[j]); err != nil {
					return []reflect.Value{}, false, ErrBadFunArgument{v[0], j, args[j][0]}
				} else {
					in[i].Index(j-i).Set(arg)
				}
			}
		} else {
			return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], len(call.Args)}
		}
	} else {
		return []reflect.Value{}, false, ErrWrongNumberOfArgs{v[0], len(call.Args)}
	}

	if builtin {
		// Builtin functions take and return raw values as well as typing information
		bin := make([]reflect.Value, len(in) * 2)
		for i := range in {
			bin[i] = reflect.ValueOf(in[i])
			bin[i+len(in)] = reflect.ValueOf(intyped[i])
		}
		in = bin
	} else {
		// Check argument types
		for i := range in {
			if in[i], err = assignableValue(in[i], ftype.In(i), intyped[i]); err != nil {
				return []reflect.Value{}, false, ErrBadFunArgument{v[0], i, in[i]}
			}
		}
	}

	if ftype.IsVariadic() {
		out = v[0].CallSlice(in)
	} else {
		out = v[0].Call(in)
	}

	if builtin {
		otyped := out[1].Bool()
		var err error = nil
		if !out[2].IsNil() {
			err = out[2].Interface().(error)
		}
		// Unwrap the Value of a Value
		out = []reflect.Value{out[0].Interface().(reflect.Value)}
		return out, otyped, err
	} else {
		return out, true, nil
	}
}

func evalSelectorExpr(selector *ast.SelectorExpr, env *env) (reflect.Value, bool, error) {
	var err error
	var x []reflect.Value
	if x, _, err = evalExpr(selector.X, env); err != nil {
		return reflect.Value{}, true, err
	}
	sel := selector.Sel.Name
	xname := x[0].Type().Name()

	if x[0].Kind() == reflect.Ptr {
		// Special case for handling packages
		if x[0].Type() == reflect.TypeOf(pkg(nil)) {
			return evalIdentExpr(selector.Sel, x[0].Interface().(pkg))
		} else if !x[0].IsNil() && x[0].Elem().Kind() == reflect.Struct {
			x[0] = x[0].Elem()
		}
	}

	switch x[0].Type().Kind() {
	case reflect.Struct:
		if v := x[0].FieldByName(sel); v.IsValid() {
			return v, true, nil
		} else if x[0].CanAddr() {
			if v := x[0].Addr().MethodByName(sel); v.IsValid() {
				return v, true, nil
			}
		}
		return reflect.Value{}, true, errors.New(fmt.Sprintf("%s has no field or method %s", xname, sel))
	case reflect.Interface:
		if v := x[0].MethodByName(sel); !v.IsValid() {
			return v, true, errors.New(fmt.Sprintf("%s has no method %s", xname, sel))
		} else {
			return v, true, nil
		}
	default:
		err = errors.New(fmt.Sprintf("%s.%s undefined (%s has no field or method %s)",
			xname, sel, xname, sel))
		return reflect.Value{}, true, err
	}
}

func evalBinaryExpr(b *ast.BinaryExpr, env *env) (r reflect.Value, rtyped bool, err error) {
	var xx, yy []reflect.Value
	var xtyped, ytyped bool
	if xx, xtyped, err = evalExpr(b.X, env); err != nil {
		return reflect.Value{}, false, err
	}
	if yy, ytyped, err = evalExpr(b.Y, env); err != nil {
		return reflect.Value{}, false, err
	}
	rtyped = xtyped || ytyped
	x, y := xx[0], yy[0]

	// Rearrange x and y such that y is assignable to x, if possible
	if xtyped && ytyped {
		if x.Type().AssignableTo(y.Type()) {
			x = x.Convert(y.Type())
		} else if !y.Type().AssignableTo(x.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
	} else if xtyped {
		if !y.Type().ConvertibleTo(x.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
		y = y.Convert(x.Type())
	} else if ytyped {
		if !x.Type().ConvertibleTo(y.Type()) {
			return r, rtyped, ErrInvalidOperands{x, b.Op, y}
		}
		x = x.Convert(y.Type())
	} else if isUntypedNumeral(x) && isUntypedNumeral(y) {
		x, y = promoteUntypedNumerals(x, y)
	} else {
		return r, rtyped, ErrInvalidOperands{x, b.Op, y}
	}

	switch x.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r, err = evalBinaryIntExpr(x, b.Op, y)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r, err = evalBinaryUintExpr(x, b.Op, y)
	case reflect.Float32, reflect.Float64:
		r, err = evalBinaryFloatExpr(x, b.Op, y)
	case reflect.Complex64, reflect.Complex128:
		r, err = evalBinaryComplexExpr(x, b.Op, y)
	case reflect.String:
		r, err = evalBinaryStringExpr(x, b.Op, y)
	default:
		err = ErrInvalidOperands{x, b.Op, y}
	}
	return
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryIntExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var r int64
	var err error

	xx, yy := x.Int(), y.Int()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	case token.REM: r = xx % yy
	case token.AND: r = xx & yy
	case token.OR:  r = xx | yy
	case token.XOR: r = xx ^ yy
	case token.AND_NOT: r = xx &^ yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryUintExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r uint64

	xx, yy := x.Uint(), y.Uint()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	case token.REM: r = xx % yy
	case token.AND: r = xx & yy
	case token.OR:  r = xx | yy
	case token.XOR: r = xx ^ yy
	case token.AND_NOT: r = xx &^ yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryFloatExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r float64

	xx, yy := x.Float(), y.Float()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryComplexExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r complex128

	xx, yy := x.Complex(), y.Complex()
	switch op {
	case token.ADD: r = xx + yy
	case token.SUB: r = xx - yy
	case token.MUL: r = xx * yy
	case token.QUO: r = xx / yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

// Assumes y is assignable to x, panics otherwise
func evalBinaryStringExpr(x reflect.Value, op token.Token, y reflect.Value) (reflect.Value, error) {
	var err error
	var r string

	xx, yy := x.String(), y.String()
	switch op {
	case token.ADD: r = xx + yy
	default: err = ErrInvalidOperands{x, op, y}
	}
	return reflect.ValueOf(r).Convert(x.Type()), err
}

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


func (err ErrInvalidOperands) Error() string {
	return fmt.Sprintf("invalid binary operation %v %v %v", err.x, err.op, err.y)
}

func (err ErrBadFunArgument) Error() string {
	return fmt.Sprintf("invalid type (%v) for argument %d of %v", err.value.Type(), err.index, err.fun)
}

func (err ErrBadComplexArguments) Error() string {
	return fmt.Sprintf("invalid operation: complex(%v, %v)", err.a, err.b)
}

func (err ErrBadBuiltinArgument) Error() string {
	return fmt.Sprintf("invalid operation: %s(%v)", err.fun, err.value)
}

func (err ErrWrongNumberOfArgs) Error() string {
	expected := err.fun.Type().NumIn()
	if err.numArgs < expected {
		return fmt.Sprintf("not enouch args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	} else {
		return fmt.Sprintf("too many args (%d) to call %v (%d)", err.numArgs, err.fun, expected)
	}
}

func (err ErrMultiInSingleContext) Error() string {
	strvals := ""
	for i, val := range err.vals {
		if i == 0 {
			strvals += fmt.Sprintf("%v", val)
		} else {
			strvals += fmt.Sprintf(", %v", val)
		}
	}
	return fmt.Sprintf("multiple-value (%s) in single value context", strvals)
}

