package main

import (
	"fmt"
	"io"
        "strings"
	"text/template"
	"go/token"
	"github.com/0xfaded/go-testgen"
)

type Test struct{}

var comment = template.Must(template.New("Comment").Parse(
`// Test {{ .Lhs.Name }} {{ .Op.Value }} {{ .Rhs.Name }}
`))

var body = template.Must(template.New("Body").Parse(
`	env := makeCheckBinaryNonConstExprEnv()
	{{ .DefA }}; env.Vars["a"] = reflect.ValueOf(&a)
	{{ if .DefB }}{{ .DefB }}; env.Vars["b"] = reflect.ValueOf(&b){{ end }}{{ if .Errors }}
	expectCheckError(t, `+"`{{ .Expr }}`"+`, env,{{ range .Errors }}
		`+"`{{ . }}`"+`,{{ end }}
	)
{{ else }}
	expectType(t, `+"`{{ .Expr }}`"+`, env, reflect.TypeOf({{ .Expr }})){{ end }}
`))

var typeDefs = `
type interfaceX interface { x() }
type interfaceY interface { x() }
type interfaceZ interface { z() }
type XinterfaceX int
func (XinterfaceX) x() {}
type structT struct {
	a int
	_ []int
}
type structUncompT struct {
	a int
	b []int
}
`
var globals = typeDefs + `
func makeCheckBinaryNonConstExprEnv() *Env {
	env := makeEnv()
	env.Types["interfaceX"] = reflect.TypeOf(new(interfaceX)).Elem()
	env.Types["interfaceY"] = reflect.TypeOf(new(interfaceY)).Elem()
	env.Types["interfaceZ"] = reflect.TypeOf(new(interfaceZ)).Elem()
	env.Types["XinterfaceX"] = reflect.TypeOf(XinterfaceX(0))
	env.Types["structT"] = reflect.TypeOf(structT{})
	env.Types["structUncompT"] = reflect.TypeOf(structUncompT{})
	return env
}
`

func (*Test) Package() string {
	return "eval"
}

func (*Test) Prefix() string {
	return "CheckBinaryNonConstExpr"
}

func (*Test) Imports() map[string]string {
	return map[string]string { "reflect": "" }
}

func (*Test) Globals(w io.Writer) error {
	_, err := w.Write([]byte(globals))
	return err
}

func (*Test) Dimensions() []testgen.Dimension {
	rhs := []testgen.Element{
		{"ConstInt", "4"},
		{"ConstRune", "'@'"},
		{"ConstFloat", "2.0"},
		{"ConstComplex", "8.0i"},
		{"ConstBool", "true"},
		{"ConstString", `"abc"`},
		{"ConstNil", "nil"},
		{"Int", "int(1)"},
		{"Float32", "float32(1.5)"},
		{"Complex128", "complex128(1i)"},
		{"Rune", "rune('a')"},
		{"String", `string("abc")`},
		{"BoolT", "bool(true)"},
		{"Slice", "sliceT(nil)"},
		{"Array", "arrayT{}"},
		{"XinterfaceX", "XinterfaceX(0)"},
		{"InterfaceX", "interfaceX(nil)"},
		{"InterfaceY", "interfaceY(nil)"},
		{"InterfaceZ", "interfaceZ(nil)"},
		{"Ptr", "(*int)(nil)"},
		{"Struct", "structT(true)"},
		{"StructUncomp", "structUncompT(true)"},
	}
	ops := []testgen.Element{
		{"Add", token.ADD},
		{"And", token.AND},
		{"Rem", token.REM},
		{"Eql", token.EQL},
		{"Gtr", token.GTR},
	}
	// exclude const types
	lhs := rhs[7:]
	return []testgen.Dimension{
		lhs[:3],
		ops,
		rhs[:5],
	}
}

func (*Test) Comment(w io.Writer, elts ...testgen.Element) error {
	vars := map[string] interface{} {
		"Lhs": elts[0],
		"Op": elts[1],
		"Rhs": elts[2],
	}

	return comment.Execute(w, vars)
}

func (*Test) Body(w io.Writer, elts ...testgen.Element) error {
	op := elts[1].Value.(token.Token)
	a := "a"
	adef := "a := " + elts[0].Value.(string)
	b := "b"
	bdef := "b := " + elts[2].Value.(string)
	if strings.HasPrefix(elts[2].Name, "Const") {
		b = elts[2].Value.(string)
		bdef = ""
	}

	expr := fmt.Sprintf("%v %v %v", a, op, b)
	defs := adef + "\n" + bdef

	compileErrs, err := compileExprWithDefsAndGlobals(expr, defs, typeDefs)
	if err != nil {
		return err
	}

	vars := map[string] interface{} {
		"DefA": adef,
		"DefB": bdef,
		"Expr": expr,
		"Errors": compileErrs,
		"Op": elts[1],
	}

	return body.Execute(w, &vars)
}

