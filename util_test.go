package eval

import (
	"testing"
)

func errorPosEqual(a, b []string) bool {
	if len(a) != len(b) { return false }
	for i, aVal := range a {
		if aVal != b[i] { return false }
	}
	return true
}

func TestFormatErrorPos(t *testing.T) {
	source  := `split(os.Args ", )")`
	errmsg  := `1:15: expected ')', found 'STRING' ", "`
	results := FormatErrorPos(source, errmsg)
	expect  := []string { source,  "--------------^" }
	if !errorPosEqual(expect, results) {
		t.Fatalf("Expected %v, got %v", expect, results)
	}

	source  = "`"
	errmsg  = `1:1: string not terminated`
	results = FormatErrorPos(source, errmsg)
	expect  = []string { source,  "^" }

	source  = "y("
	errmsg  = `1:3: expected ')', found 'EOF'`
	results = FormatErrorPos(source, errmsg)
	expect  = []string { source,  "--^" }

}
