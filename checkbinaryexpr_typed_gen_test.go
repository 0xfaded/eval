package eval

import (
	"testing"
	"reflect"
)

// Test Int8 + Int
func TestCheckBinaryTypedExprInt8AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + 4`, env,
		`constant 131 overflows int8`,
	)

}

// Test Int8 + Rune
func TestCheckBinaryTypedExprInt8AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + '@'`, env,
		`constant 191 overflows int8`,
	)

}

// Test Int8 + Float
func TestCheckBinaryTypedExprInt8AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + 2.0`, env,
		`constant 129 overflows int8`,
	)

}

// Test Int8 + Complex
func TestCheckBinaryTypedExprInt8AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int8 + Bool
func TestCheckBinaryTypedExprInt8AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 + true (mismatched types int8 and bool)`,
	)

}

// Test Int8 + String
func TestCheckBinaryTypedExprInt8AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 + "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 + Nil
func TestCheckBinaryTypedExprInt8AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) + nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int8 - Int
func TestCheckBinaryTypedExprInt8SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) - 4`, env, int8(0x7f) - 4, reflect.TypeOf(int8(0x7f) - 4))
}

// Test Int8 - Rune
func TestCheckBinaryTypedExprInt8SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) - '@'`, env, int8(0x7f) - '@', reflect.TypeOf(int8(0x7f) - '@'))
}

// Test Int8 - Float
func TestCheckBinaryTypedExprInt8SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) - 2.0`, env, int8(0x7f) - 2.0, reflect.TypeOf(int8(0x7f) - 2.0))
}

// Test Int8 - Complex
func TestCheckBinaryTypedExprInt8SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int8 - Bool
func TestCheckBinaryTypedExprInt8SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) - true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 - true (mismatched types int8 and bool)`,
	)

}

// Test Int8 - String
func TestCheckBinaryTypedExprInt8SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) - "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 - "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 - Nil
func TestCheckBinaryTypedExprInt8SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) - nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int8 & Int
func TestCheckBinaryTypedExprInt8AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) & 4`, env, int8(0x7f) & 4, reflect.TypeOf(int8(0x7f) & 4))
}

// Test Int8 & Rune
func TestCheckBinaryTypedExprInt8AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) & '@'`, env, int8(0x7f) & '@', reflect.TypeOf(int8(0x7f) & '@'))
}

// Test Int8 & Float
func TestCheckBinaryTypedExprInt8AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) & 2.0`, env, int8(0x7f) & 2.0, reflect.TypeOf(int8(0x7f) & 2.0))
}

// Test Int8 & Complex
func TestCheckBinaryTypedExprInt8AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int8 & Bool
func TestCheckBinaryTypedExprInt8AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) & true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 & true (mismatched types int8 and bool)`,
	)

}

// Test Int8 & String
func TestCheckBinaryTypedExprInt8AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) & "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 & "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 & Nil
func TestCheckBinaryTypedExprInt8AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) & nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int8 % Int
func TestCheckBinaryTypedExprInt8RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) % 4`, env, int8(0x7f) % 4, reflect.TypeOf(int8(0x7f) % 4))
}

// Test Int8 % Rune
func TestCheckBinaryTypedExprInt8RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) % '@'`, env, int8(0x7f) % '@', reflect.TypeOf(int8(0x7f) % '@'))
}

// Test Int8 % Float
func TestCheckBinaryTypedExprInt8RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) % 2.0`, env, int8(0x7f) % 2.0, reflect.TypeOf(int8(0x7f) % 2.0))
}

// Test Int8 % Complex
func TestCheckBinaryTypedExprInt8RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Int8 % Bool
func TestCheckBinaryTypedExprInt8RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) % true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 % true (mismatched types int8 and bool)`,
	)

}

// Test Int8 % String
func TestCheckBinaryTypedExprInt8RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) % "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 % "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 % Nil
func TestCheckBinaryTypedExprInt8RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) % nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int8 == Int
func TestCheckBinaryTypedExprInt8EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) == 4`, env, int8(0x7f) == 4, reflect.TypeOf(int8(0x7f) == 4))
}

// Test Int8 == Rune
func TestCheckBinaryTypedExprInt8EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) == '@'`, env, int8(0x7f) == '@', reflect.TypeOf(int8(0x7f) == '@'))
}

// Test Int8 == Float
func TestCheckBinaryTypedExprInt8EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) == 2.0`, env, int8(0x7f) == 2.0, reflect.TypeOf(int8(0x7f) == 2.0))
}

// Test Int8 == Complex
func TestCheckBinaryTypedExprInt8EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int8 == Bool
func TestCheckBinaryTypedExprInt8EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) == true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 == true (mismatched types int8 and bool)`,
	)

}

// Test Int8 == String
func TestCheckBinaryTypedExprInt8EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) == "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 == "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 == Nil
func TestCheckBinaryTypedExprInt8EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) == nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int8 > Int
func TestCheckBinaryTypedExprInt8GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) > 4`, env, int8(0x7f) > 4, reflect.TypeOf(int8(0x7f) > 4))
}

// Test Int8 > Rune
func TestCheckBinaryTypedExprInt8GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) > '@'`, env, int8(0x7f) > '@', reflect.TypeOf(int8(0x7f) > '@'))
}

// Test Int8 > Float
func TestCheckBinaryTypedExprInt8GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int8(0x7f) > 2.0`, env, int8(0x7f) > 2.0, reflect.TypeOf(int8(0x7f) > 2.0))
}

// Test Int8 > Complex
func TestCheckBinaryTypedExprInt8GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int8 > Bool
func TestCheckBinaryTypedExprInt8GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) > true`, env,
		`cannot convert true to type int8`,
		`invalid operation: 127 > true (mismatched types int8 and bool)`,
	)

}

// Test Int8 > String
func TestCheckBinaryTypedExprInt8GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) > "abc"`, env,
		`cannot convert "abc" to type int8`,
		`invalid operation: 127 > "abc" (mismatched types int8 and string)`,
	)

}

// Test Int8 > Nil
func TestCheckBinaryTypedExprInt8GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int8(0x7f) > nil`, env,
		`cannot convert nil to type int8`,
	)

}

// Test Int16 + Int
func TestCheckBinaryTypedExprInt16AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + 4`, env,
		`constant 32771 overflows int16`,
	)

}

// Test Int16 + Rune
func TestCheckBinaryTypedExprInt16AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + '@'`, env,
		`constant 32831 overflows int16`,
	)

}

// Test Int16 + Float
func TestCheckBinaryTypedExprInt16AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + 2.0`, env,
		`constant 32769 overflows int16`,
	)

}

// Test Int16 + Complex
func TestCheckBinaryTypedExprInt16AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int16 + Bool
func TestCheckBinaryTypedExprInt16AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 + true (mismatched types int16 and bool)`,
	)

}

// Test Int16 + String
func TestCheckBinaryTypedExprInt16AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 + "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 + Nil
func TestCheckBinaryTypedExprInt16AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) + nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int16 - Int
func TestCheckBinaryTypedExprInt16SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) - 4`, env, int16(0x7fff) - 4, reflect.TypeOf(int16(0x7fff) - 4))
}

// Test Int16 - Rune
func TestCheckBinaryTypedExprInt16SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) - '@'`, env, int16(0x7fff) - '@', reflect.TypeOf(int16(0x7fff) - '@'))
}

// Test Int16 - Float
func TestCheckBinaryTypedExprInt16SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) - 2.0`, env, int16(0x7fff) - 2.0, reflect.TypeOf(int16(0x7fff) - 2.0))
}

// Test Int16 - Complex
func TestCheckBinaryTypedExprInt16SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int16 - Bool
func TestCheckBinaryTypedExprInt16SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) - true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 - true (mismatched types int16 and bool)`,
	)

}

// Test Int16 - String
func TestCheckBinaryTypedExprInt16SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) - "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 - "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 - Nil
func TestCheckBinaryTypedExprInt16SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) - nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int16 & Int
func TestCheckBinaryTypedExprInt16AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) & 4`, env, int16(0x7fff) & 4, reflect.TypeOf(int16(0x7fff) & 4))
}

// Test Int16 & Rune
func TestCheckBinaryTypedExprInt16AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) & '@'`, env, int16(0x7fff) & '@', reflect.TypeOf(int16(0x7fff) & '@'))
}

// Test Int16 & Float
func TestCheckBinaryTypedExprInt16AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) & 2.0`, env, int16(0x7fff) & 2.0, reflect.TypeOf(int16(0x7fff) & 2.0))
}

// Test Int16 & Complex
func TestCheckBinaryTypedExprInt16AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int16 & Bool
func TestCheckBinaryTypedExprInt16AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) & true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 & true (mismatched types int16 and bool)`,
	)

}

// Test Int16 & String
func TestCheckBinaryTypedExprInt16AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) & "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 & "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 & Nil
func TestCheckBinaryTypedExprInt16AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) & nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int16 % Int
func TestCheckBinaryTypedExprInt16RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) % 4`, env, int16(0x7fff) % 4, reflect.TypeOf(int16(0x7fff) % 4))
}

// Test Int16 % Rune
func TestCheckBinaryTypedExprInt16RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) % '@'`, env, int16(0x7fff) % '@', reflect.TypeOf(int16(0x7fff) % '@'))
}

// Test Int16 % Float
func TestCheckBinaryTypedExprInt16RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) % 2.0`, env, int16(0x7fff) % 2.0, reflect.TypeOf(int16(0x7fff) % 2.0))
}

// Test Int16 % Complex
func TestCheckBinaryTypedExprInt16RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Int16 % Bool
func TestCheckBinaryTypedExprInt16RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) % true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 % true (mismatched types int16 and bool)`,
	)

}

// Test Int16 % String
func TestCheckBinaryTypedExprInt16RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) % "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 % "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 % Nil
func TestCheckBinaryTypedExprInt16RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) % nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int16 == Int
func TestCheckBinaryTypedExprInt16EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) == 4`, env, int16(0x7fff) == 4, reflect.TypeOf(int16(0x7fff) == 4))
}

// Test Int16 == Rune
func TestCheckBinaryTypedExprInt16EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) == '@'`, env, int16(0x7fff) == '@', reflect.TypeOf(int16(0x7fff) == '@'))
}

// Test Int16 == Float
func TestCheckBinaryTypedExprInt16EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) == 2.0`, env, int16(0x7fff) == 2.0, reflect.TypeOf(int16(0x7fff) == 2.0))
}

// Test Int16 == Complex
func TestCheckBinaryTypedExprInt16EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int16 == Bool
func TestCheckBinaryTypedExprInt16EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) == true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 == true (mismatched types int16 and bool)`,
	)

}

// Test Int16 == String
func TestCheckBinaryTypedExprInt16EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) == "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 == "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 == Nil
func TestCheckBinaryTypedExprInt16EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) == nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int16 > Int
func TestCheckBinaryTypedExprInt16GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) > 4`, env, int16(0x7fff) > 4, reflect.TypeOf(int16(0x7fff) > 4))
}

// Test Int16 > Rune
func TestCheckBinaryTypedExprInt16GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) > '@'`, env, int16(0x7fff) > '@', reflect.TypeOf(int16(0x7fff) > '@'))
}

// Test Int16 > Float
func TestCheckBinaryTypedExprInt16GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int16(0x7fff) > 2.0`, env, int16(0x7fff) > 2.0, reflect.TypeOf(int16(0x7fff) > 2.0))
}

// Test Int16 > Complex
func TestCheckBinaryTypedExprInt16GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int16 > Bool
func TestCheckBinaryTypedExprInt16GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) > true`, env,
		`cannot convert true to type int16`,
		`invalid operation: 32767 > true (mismatched types int16 and bool)`,
	)

}

// Test Int16 > String
func TestCheckBinaryTypedExprInt16GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) > "abc"`, env,
		`cannot convert "abc" to type int16`,
		`invalid operation: 32767 > "abc" (mismatched types int16 and string)`,
	)

}

// Test Int16 > Nil
func TestCheckBinaryTypedExprInt16GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int16(0x7fff) > nil`, env,
		`cannot convert nil to type int16`,
	)

}

// Test Int32 + Int
func TestCheckBinaryTypedExprInt32AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + 4`, env,
		`constant 2147483651 overflows int32`,
	)

}

// Test Int32 + Rune
func TestCheckBinaryTypedExprInt32AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + '@'`, env,
		`constant 2147483711 overflows int32`,
	)

}

// Test Int32 + Float
func TestCheckBinaryTypedExprInt32AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + 2.0`, env,
		`constant 2147483649 overflows int32`,
	)

}

// Test Int32 + Complex
func TestCheckBinaryTypedExprInt32AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int32 + Bool
func TestCheckBinaryTypedExprInt32AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 + true (mismatched types int32 and bool)`,
	)

}

// Test Int32 + String
func TestCheckBinaryTypedExprInt32AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 + "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 + Nil
func TestCheckBinaryTypedExprInt32AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) + nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int32 - Int
func TestCheckBinaryTypedExprInt32SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) - 4`, env, int32(0x7fffffff) - 4, reflect.TypeOf(int32(0x7fffffff) - 4))
}

// Test Int32 - Rune
func TestCheckBinaryTypedExprInt32SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) - '@'`, env, int32(0x7fffffff) - '@', reflect.TypeOf(int32(0x7fffffff) - '@'))
}

// Test Int32 - Float
func TestCheckBinaryTypedExprInt32SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) - 2.0`, env, int32(0x7fffffff) - 2.0, reflect.TypeOf(int32(0x7fffffff) - 2.0))
}

// Test Int32 - Complex
func TestCheckBinaryTypedExprInt32SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int32 - Bool
func TestCheckBinaryTypedExprInt32SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) - true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 - true (mismatched types int32 and bool)`,
	)

}

// Test Int32 - String
func TestCheckBinaryTypedExprInt32SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) - "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 - "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 - Nil
func TestCheckBinaryTypedExprInt32SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) - nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int32 & Int
func TestCheckBinaryTypedExprInt32AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) & 4`, env, int32(0x7fffffff) & 4, reflect.TypeOf(int32(0x7fffffff) & 4))
}

// Test Int32 & Rune
func TestCheckBinaryTypedExprInt32AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) & '@'`, env, int32(0x7fffffff) & '@', reflect.TypeOf(int32(0x7fffffff) & '@'))
}

// Test Int32 & Float
func TestCheckBinaryTypedExprInt32AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) & 2.0`, env, int32(0x7fffffff) & 2.0, reflect.TypeOf(int32(0x7fffffff) & 2.0))
}

// Test Int32 & Complex
func TestCheckBinaryTypedExprInt32AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int32 & Bool
func TestCheckBinaryTypedExprInt32AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) & true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 & true (mismatched types int32 and bool)`,
	)

}

// Test Int32 & String
func TestCheckBinaryTypedExprInt32AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) & "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 & "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 & Nil
func TestCheckBinaryTypedExprInt32AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) & nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int32 % Int
func TestCheckBinaryTypedExprInt32RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) % 4`, env, int32(0x7fffffff) % 4, reflect.TypeOf(int32(0x7fffffff) % 4))
}

// Test Int32 % Rune
func TestCheckBinaryTypedExprInt32RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) % '@'`, env, int32(0x7fffffff) % '@', reflect.TypeOf(int32(0x7fffffff) % '@'))
}

// Test Int32 % Float
func TestCheckBinaryTypedExprInt32RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) % 2.0`, env, int32(0x7fffffff) % 2.0, reflect.TypeOf(int32(0x7fffffff) % 2.0))
}

// Test Int32 % Complex
func TestCheckBinaryTypedExprInt32RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Int32 % Bool
func TestCheckBinaryTypedExprInt32RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) % true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 % true (mismatched types int32 and bool)`,
	)

}

// Test Int32 % String
func TestCheckBinaryTypedExprInt32RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) % "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 % "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 % Nil
func TestCheckBinaryTypedExprInt32RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) % nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int32 == Int
func TestCheckBinaryTypedExprInt32EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) == 4`, env, int32(0x7fffffff) == 4, reflect.TypeOf(int32(0x7fffffff) == 4))
}

// Test Int32 == Rune
func TestCheckBinaryTypedExprInt32EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) == '@'`, env, int32(0x7fffffff) == '@', reflect.TypeOf(int32(0x7fffffff) == '@'))
}

// Test Int32 == Float
func TestCheckBinaryTypedExprInt32EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) == 2.0`, env, int32(0x7fffffff) == 2.0, reflect.TypeOf(int32(0x7fffffff) == 2.0))
}

// Test Int32 == Complex
func TestCheckBinaryTypedExprInt32EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int32 == Bool
func TestCheckBinaryTypedExprInt32EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) == true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 == true (mismatched types int32 and bool)`,
	)

}

// Test Int32 == String
func TestCheckBinaryTypedExprInt32EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) == "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 == "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 == Nil
func TestCheckBinaryTypedExprInt32EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) == nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int32 > Int
func TestCheckBinaryTypedExprInt32GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) > 4`, env, int32(0x7fffffff) > 4, reflect.TypeOf(int32(0x7fffffff) > 4))
}

// Test Int32 > Rune
func TestCheckBinaryTypedExprInt32GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) > '@'`, env, int32(0x7fffffff) > '@', reflect.TypeOf(int32(0x7fffffff) > '@'))
}

// Test Int32 > Float
func TestCheckBinaryTypedExprInt32GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int32(0x7fffffff) > 2.0`, env, int32(0x7fffffff) > 2.0, reflect.TypeOf(int32(0x7fffffff) > 2.0))
}

// Test Int32 > Complex
func TestCheckBinaryTypedExprInt32GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int32 > Bool
func TestCheckBinaryTypedExprInt32GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) > true`, env,
		`cannot convert true to type int32`,
		`invalid operation: 2147483647 > true (mismatched types int32 and bool)`,
	)

}

// Test Int32 > String
func TestCheckBinaryTypedExprInt32GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) > "abc"`, env,
		`cannot convert "abc" to type int32`,
		`invalid operation: 2147483647 > "abc" (mismatched types int32 and string)`,
	)

}

// Test Int32 > Nil
func TestCheckBinaryTypedExprInt32GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int32(0x7fffffff) > nil`, env,
		`cannot convert nil to type int32`,
	)

}

// Test Int64 + Int
func TestCheckBinaryTypedExprInt64AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + 4`, env,
		`constant 9223372036854775811 overflows int64`,
	)

}

// Test Int64 + Rune
func TestCheckBinaryTypedExprInt64AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + '@'`, env,
		`constant 9223372036854775871 overflows int64`,
	)

}

// Test Int64 + Float
func TestCheckBinaryTypedExprInt64AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + 2.0`, env,
		`constant 9223372036854775809 overflows int64`,
	)

}

// Test Int64 + Complex
func TestCheckBinaryTypedExprInt64AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int64 + Bool
func TestCheckBinaryTypedExprInt64AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 + true (mismatched types int64 and bool)`,
	)

}

// Test Int64 + String
func TestCheckBinaryTypedExprInt64AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 + "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 + Nil
func TestCheckBinaryTypedExprInt64AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) + nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Int64 - Int
func TestCheckBinaryTypedExprInt64SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) - 4`, env, int64(0x7fffffffffffffff) - 4, reflect.TypeOf(int64(0x7fffffffffffffff) - 4))
}

// Test Int64 - Rune
func TestCheckBinaryTypedExprInt64SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) - '@'`, env, int64(0x7fffffffffffffff) - '@', reflect.TypeOf(int64(0x7fffffffffffffff) - '@'))
}

// Test Int64 - Float
func TestCheckBinaryTypedExprInt64SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) - 2.0`, env, int64(0x7fffffffffffffff) - 2.0, reflect.TypeOf(int64(0x7fffffffffffffff) - 2.0))
}

// Test Int64 - Complex
func TestCheckBinaryTypedExprInt64SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int64 - Bool
func TestCheckBinaryTypedExprInt64SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) - true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 - true (mismatched types int64 and bool)`,
	)

}

// Test Int64 - String
func TestCheckBinaryTypedExprInt64SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) - "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 - "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 - Nil
func TestCheckBinaryTypedExprInt64SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) - nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Int64 & Int
func TestCheckBinaryTypedExprInt64AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) & 4`, env, int64(0x7fffffffffffffff) & 4, reflect.TypeOf(int64(0x7fffffffffffffff) & 4))
}

// Test Int64 & Rune
func TestCheckBinaryTypedExprInt64AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) & '@'`, env, int64(0x7fffffffffffffff) & '@', reflect.TypeOf(int64(0x7fffffffffffffff) & '@'))
}

// Test Int64 & Float
func TestCheckBinaryTypedExprInt64AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) & 2.0`, env, int64(0x7fffffffffffffff) & 2.0, reflect.TypeOf(int64(0x7fffffffffffffff) & 2.0))
}

// Test Int64 & Complex
func TestCheckBinaryTypedExprInt64AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int64 & Bool
func TestCheckBinaryTypedExprInt64AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) & true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 & true (mismatched types int64 and bool)`,
	)

}

// Test Int64 & String
func TestCheckBinaryTypedExprInt64AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) & "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 & "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 & Nil
func TestCheckBinaryTypedExprInt64AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) & nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Int64 % Int
func TestCheckBinaryTypedExprInt64RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) % 4`, env, int64(0x7fffffffffffffff) % 4, reflect.TypeOf(int64(0x7fffffffffffffff) % 4))
}

// Test Int64 % Rune
func TestCheckBinaryTypedExprInt64RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) % '@'`, env, int64(0x7fffffffffffffff) % '@', reflect.TypeOf(int64(0x7fffffffffffffff) % '@'))
}

// Test Int64 % Float
func TestCheckBinaryTypedExprInt64RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) % 2.0`, env, int64(0x7fffffffffffffff) % 2.0, reflect.TypeOf(int64(0x7fffffffffffffff) % 2.0))
}

// Test Int64 % Complex
func TestCheckBinaryTypedExprInt64RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Int64 % Bool
func TestCheckBinaryTypedExprInt64RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) % true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 % true (mismatched types int64 and bool)`,
	)

}

// Test Int64 % String
func TestCheckBinaryTypedExprInt64RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) % "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 % "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 % Nil
func TestCheckBinaryTypedExprInt64RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) % nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Int64 == Int
func TestCheckBinaryTypedExprInt64EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) == 4`, env, int64(0x7fffffffffffffff) == 4, reflect.TypeOf(int64(0x7fffffffffffffff) == 4))
}

// Test Int64 == Rune
func TestCheckBinaryTypedExprInt64EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) == '@'`, env, int64(0x7fffffffffffffff) == '@', reflect.TypeOf(int64(0x7fffffffffffffff) == '@'))
}

// Test Int64 == Float
func TestCheckBinaryTypedExprInt64EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) == 2.0`, env, int64(0x7fffffffffffffff) == 2.0, reflect.TypeOf(int64(0x7fffffffffffffff) == 2.0))
}

// Test Int64 == Complex
func TestCheckBinaryTypedExprInt64EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int64 == Bool
func TestCheckBinaryTypedExprInt64EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) == true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 == true (mismatched types int64 and bool)`,
	)

}

// Test Int64 == String
func TestCheckBinaryTypedExprInt64EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) == "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 == "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 == Nil
func TestCheckBinaryTypedExprInt64EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) == nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Int64 > Int
func TestCheckBinaryTypedExprInt64GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) > 4`, env, int64(0x7fffffffffffffff) > 4, reflect.TypeOf(int64(0x7fffffffffffffff) > 4))
}

// Test Int64 > Rune
func TestCheckBinaryTypedExprInt64GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) > '@'`, env, int64(0x7fffffffffffffff) > '@', reflect.TypeOf(int64(0x7fffffffffffffff) > '@'))
}

// Test Int64 > Float
func TestCheckBinaryTypedExprInt64GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `int64(0x7fffffffffffffff) > 2.0`, env, int64(0x7fffffffffffffff) > 2.0, reflect.TypeOf(int64(0x7fffffffffffffff) > 2.0))
}

// Test Int64 > Complex
func TestCheckBinaryTypedExprInt64GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Int64 > Bool
func TestCheckBinaryTypedExprInt64GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) > true`, env,
		`cannot convert true to type int64`,
		`invalid operation: 9223372036854775807 > true (mismatched types int64 and bool)`,
	)

}

// Test Int64 > String
func TestCheckBinaryTypedExprInt64GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) > "abc"`, env,
		`cannot convert "abc" to type int64`,
		`invalid operation: 9223372036854775807 > "abc" (mismatched types int64 and string)`,
	)

}

// Test Int64 > Nil
func TestCheckBinaryTypedExprInt64GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `int64(0x7fffffffffffffff) > nil`, env,
		`cannot convert nil to type int64`,
	)

}

// Test Uint8 + Int
func TestCheckBinaryTypedExprUint8AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + 4`, env,
		`constant 259 overflows uint8`,
	)

}

// Test Uint8 + Rune
func TestCheckBinaryTypedExprUint8AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + '@'`, env,
		`constant 319 overflows uint8`,
	)

}

// Test Uint8 + Float
func TestCheckBinaryTypedExprUint8AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + 2.0`, env,
		`constant 257 overflows uint8`,
	)

}

// Test Uint8 + Complex
func TestCheckBinaryTypedExprUint8AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint8 + Bool
func TestCheckBinaryTypedExprUint8AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 + true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 + String
func TestCheckBinaryTypedExprUint8AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 + "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 + Nil
func TestCheckBinaryTypedExprUint8AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) + nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint8 - Int
func TestCheckBinaryTypedExprUint8SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) - 4`, env, uint8(0xff) - 4, reflect.TypeOf(uint8(0xff) - 4))
}

// Test Uint8 - Rune
func TestCheckBinaryTypedExprUint8SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) - '@'`, env, uint8(0xff) - '@', reflect.TypeOf(uint8(0xff) - '@'))
}

// Test Uint8 - Float
func TestCheckBinaryTypedExprUint8SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) - 2.0`, env, uint8(0xff) - 2.0, reflect.TypeOf(uint8(0xff) - 2.0))
}

// Test Uint8 - Complex
func TestCheckBinaryTypedExprUint8SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint8 - Bool
func TestCheckBinaryTypedExprUint8SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) - true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 - true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 - String
func TestCheckBinaryTypedExprUint8SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) - "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 - "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 - Nil
func TestCheckBinaryTypedExprUint8SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) - nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint8 & Int
func TestCheckBinaryTypedExprUint8AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) & 4`, env, uint8(0xff) & 4, reflect.TypeOf(uint8(0xff) & 4))
}

// Test Uint8 & Rune
func TestCheckBinaryTypedExprUint8AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) & '@'`, env, uint8(0xff) & '@', reflect.TypeOf(uint8(0xff) & '@'))
}

// Test Uint8 & Float
func TestCheckBinaryTypedExprUint8AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) & 2.0`, env, uint8(0xff) & 2.0, reflect.TypeOf(uint8(0xff) & 2.0))
}

// Test Uint8 & Complex
func TestCheckBinaryTypedExprUint8AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint8 & Bool
func TestCheckBinaryTypedExprUint8AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) & true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 & true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 & String
func TestCheckBinaryTypedExprUint8AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) & "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 & "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 & Nil
func TestCheckBinaryTypedExprUint8AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) & nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint8 % Int
func TestCheckBinaryTypedExprUint8RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) % 4`, env, uint8(0xff) % 4, reflect.TypeOf(uint8(0xff) % 4))
}

// Test Uint8 % Rune
func TestCheckBinaryTypedExprUint8RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) % '@'`, env, uint8(0xff) % '@', reflect.TypeOf(uint8(0xff) % '@'))
}

// Test Uint8 % Float
func TestCheckBinaryTypedExprUint8RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) % 2.0`, env, uint8(0xff) % 2.0, reflect.TypeOf(uint8(0xff) % 2.0))
}

// Test Uint8 % Complex
func TestCheckBinaryTypedExprUint8RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Uint8 % Bool
func TestCheckBinaryTypedExprUint8RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) % true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 % true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 % String
func TestCheckBinaryTypedExprUint8RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) % "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 % "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 % Nil
func TestCheckBinaryTypedExprUint8RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) % nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint8 == Int
func TestCheckBinaryTypedExprUint8EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) == 4`, env, uint8(0xff) == 4, reflect.TypeOf(uint8(0xff) == 4))
}

// Test Uint8 == Rune
func TestCheckBinaryTypedExprUint8EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) == '@'`, env, uint8(0xff) == '@', reflect.TypeOf(uint8(0xff) == '@'))
}

// Test Uint8 == Float
func TestCheckBinaryTypedExprUint8EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) == 2.0`, env, uint8(0xff) == 2.0, reflect.TypeOf(uint8(0xff) == 2.0))
}

// Test Uint8 == Complex
func TestCheckBinaryTypedExprUint8EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint8 == Bool
func TestCheckBinaryTypedExprUint8EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) == true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 == true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 == String
func TestCheckBinaryTypedExprUint8EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) == "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 == "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 == Nil
func TestCheckBinaryTypedExprUint8EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) == nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint8 > Int
func TestCheckBinaryTypedExprUint8GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) > 4`, env, uint8(0xff) > 4, reflect.TypeOf(uint8(0xff) > 4))
}

// Test Uint8 > Rune
func TestCheckBinaryTypedExprUint8GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) > '@'`, env, uint8(0xff) > '@', reflect.TypeOf(uint8(0xff) > '@'))
}

// Test Uint8 > Float
func TestCheckBinaryTypedExprUint8GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint8(0xff) > 2.0`, env, uint8(0xff) > 2.0, reflect.TypeOf(uint8(0xff) > 2.0))
}

// Test Uint8 > Complex
func TestCheckBinaryTypedExprUint8GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint8 > Bool
func TestCheckBinaryTypedExprUint8GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) > true`, env,
		`cannot convert true to type uint8`,
		`invalid operation: 255 > true (mismatched types uint8 and bool)`,
	)

}

// Test Uint8 > String
func TestCheckBinaryTypedExprUint8GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) > "abc"`, env,
		`cannot convert "abc" to type uint8`,
		`invalid operation: 255 > "abc" (mismatched types uint8 and string)`,
	)

}

// Test Uint8 > Nil
func TestCheckBinaryTypedExprUint8GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint8(0xff) > nil`, env,
		`cannot convert nil to type uint8`,
	)

}

// Test Uint16 + Int
func TestCheckBinaryTypedExprUint16AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + 4`, env,
		`constant 65539 overflows uint16`,
	)

}

// Test Uint16 + Rune
func TestCheckBinaryTypedExprUint16AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + '@'`, env,
		`constant 65599 overflows uint16`,
	)

}

// Test Uint16 + Float
func TestCheckBinaryTypedExprUint16AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + 2.0`, env,
		`constant 65537 overflows uint16`,
	)

}

// Test Uint16 + Complex
func TestCheckBinaryTypedExprUint16AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint16 + Bool
func TestCheckBinaryTypedExprUint16AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 + true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 + String
func TestCheckBinaryTypedExprUint16AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 + "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 + Nil
func TestCheckBinaryTypedExprUint16AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) + nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint16 - Int
func TestCheckBinaryTypedExprUint16SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) - 4`, env, uint16(0xffff) - 4, reflect.TypeOf(uint16(0xffff) - 4))
}

// Test Uint16 - Rune
func TestCheckBinaryTypedExprUint16SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) - '@'`, env, uint16(0xffff) - '@', reflect.TypeOf(uint16(0xffff) - '@'))
}

// Test Uint16 - Float
func TestCheckBinaryTypedExprUint16SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) - 2.0`, env, uint16(0xffff) - 2.0, reflect.TypeOf(uint16(0xffff) - 2.0))
}

// Test Uint16 - Complex
func TestCheckBinaryTypedExprUint16SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint16 - Bool
func TestCheckBinaryTypedExprUint16SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) - true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 - true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 - String
func TestCheckBinaryTypedExprUint16SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) - "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 - "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 - Nil
func TestCheckBinaryTypedExprUint16SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) - nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint16 & Int
func TestCheckBinaryTypedExprUint16AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) & 4`, env, uint16(0xffff) & 4, reflect.TypeOf(uint16(0xffff) & 4))
}

// Test Uint16 & Rune
func TestCheckBinaryTypedExprUint16AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) & '@'`, env, uint16(0xffff) & '@', reflect.TypeOf(uint16(0xffff) & '@'))
}

// Test Uint16 & Float
func TestCheckBinaryTypedExprUint16AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) & 2.0`, env, uint16(0xffff) & 2.0, reflect.TypeOf(uint16(0xffff) & 2.0))
}

// Test Uint16 & Complex
func TestCheckBinaryTypedExprUint16AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint16 & Bool
func TestCheckBinaryTypedExprUint16AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) & true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 & true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 & String
func TestCheckBinaryTypedExprUint16AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) & "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 & "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 & Nil
func TestCheckBinaryTypedExprUint16AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) & nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint16 % Int
func TestCheckBinaryTypedExprUint16RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) % 4`, env, uint16(0xffff) % 4, reflect.TypeOf(uint16(0xffff) % 4))
}

// Test Uint16 % Rune
func TestCheckBinaryTypedExprUint16RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) % '@'`, env, uint16(0xffff) % '@', reflect.TypeOf(uint16(0xffff) % '@'))
}

// Test Uint16 % Float
func TestCheckBinaryTypedExprUint16RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) % 2.0`, env, uint16(0xffff) % 2.0, reflect.TypeOf(uint16(0xffff) % 2.0))
}

// Test Uint16 % Complex
func TestCheckBinaryTypedExprUint16RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Uint16 % Bool
func TestCheckBinaryTypedExprUint16RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) % true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 % true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 % String
func TestCheckBinaryTypedExprUint16RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) % "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 % "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 % Nil
func TestCheckBinaryTypedExprUint16RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) % nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint16 == Int
func TestCheckBinaryTypedExprUint16EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) == 4`, env, uint16(0xffff) == 4, reflect.TypeOf(uint16(0xffff) == 4))
}

// Test Uint16 == Rune
func TestCheckBinaryTypedExprUint16EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) == '@'`, env, uint16(0xffff) == '@', reflect.TypeOf(uint16(0xffff) == '@'))
}

// Test Uint16 == Float
func TestCheckBinaryTypedExprUint16EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) == 2.0`, env, uint16(0xffff) == 2.0, reflect.TypeOf(uint16(0xffff) == 2.0))
}

// Test Uint16 == Complex
func TestCheckBinaryTypedExprUint16EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint16 == Bool
func TestCheckBinaryTypedExprUint16EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) == true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 == true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 == String
func TestCheckBinaryTypedExprUint16EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) == "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 == "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 == Nil
func TestCheckBinaryTypedExprUint16EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) == nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint16 > Int
func TestCheckBinaryTypedExprUint16GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) > 4`, env, uint16(0xffff) > 4, reflect.TypeOf(uint16(0xffff) > 4))
}

// Test Uint16 > Rune
func TestCheckBinaryTypedExprUint16GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) > '@'`, env, uint16(0xffff) > '@', reflect.TypeOf(uint16(0xffff) > '@'))
}

// Test Uint16 > Float
func TestCheckBinaryTypedExprUint16GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint16(0xffff) > 2.0`, env, uint16(0xffff) > 2.0, reflect.TypeOf(uint16(0xffff) > 2.0))
}

// Test Uint16 > Complex
func TestCheckBinaryTypedExprUint16GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint16 > Bool
func TestCheckBinaryTypedExprUint16GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) > true`, env,
		`cannot convert true to type uint16`,
		`invalid operation: 65535 > true (mismatched types uint16 and bool)`,
	)

}

// Test Uint16 > String
func TestCheckBinaryTypedExprUint16GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) > "abc"`, env,
		`cannot convert "abc" to type uint16`,
		`invalid operation: 65535 > "abc" (mismatched types uint16 and string)`,
	)

}

// Test Uint16 > Nil
func TestCheckBinaryTypedExprUint16GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint16(0xffff) > nil`, env,
		`cannot convert nil to type uint16`,
	)

}

// Test Uint32 + Int
func TestCheckBinaryTypedExprUint32AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + 4`, env,
		`constant 4294967299 overflows uint32`,
	)

}

// Test Uint32 + Rune
func TestCheckBinaryTypedExprUint32AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + '@'`, env,
		`constant 4294967359 overflows uint32`,
	)

}

// Test Uint32 + Float
func TestCheckBinaryTypedExprUint32AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + 2.0`, env,
		`constant 4294967297 overflows uint32`,
	)

}

// Test Uint32 + Complex
func TestCheckBinaryTypedExprUint32AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint32 + Bool
func TestCheckBinaryTypedExprUint32AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 + true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 + String
func TestCheckBinaryTypedExprUint32AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 + "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 + Nil
func TestCheckBinaryTypedExprUint32AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) + nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint32 - Int
func TestCheckBinaryTypedExprUint32SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) - 4`, env, uint32(0xffffffff) - 4, reflect.TypeOf(uint32(0xffffffff) - 4))
}

// Test Uint32 - Rune
func TestCheckBinaryTypedExprUint32SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) - '@'`, env, uint32(0xffffffff) - '@', reflect.TypeOf(uint32(0xffffffff) - '@'))
}

// Test Uint32 - Float
func TestCheckBinaryTypedExprUint32SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) - 2.0`, env, uint32(0xffffffff) - 2.0, reflect.TypeOf(uint32(0xffffffff) - 2.0))
}

// Test Uint32 - Complex
func TestCheckBinaryTypedExprUint32SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint32 - Bool
func TestCheckBinaryTypedExprUint32SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) - true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 - true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 - String
func TestCheckBinaryTypedExprUint32SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) - "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 - "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 - Nil
func TestCheckBinaryTypedExprUint32SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) - nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint32 & Int
func TestCheckBinaryTypedExprUint32AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) & 4`, env, uint32(0xffffffff) & 4, reflect.TypeOf(uint32(0xffffffff) & 4))
}

// Test Uint32 & Rune
func TestCheckBinaryTypedExprUint32AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) & '@'`, env, uint32(0xffffffff) & '@', reflect.TypeOf(uint32(0xffffffff) & '@'))
}

// Test Uint32 & Float
func TestCheckBinaryTypedExprUint32AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) & 2.0`, env, uint32(0xffffffff) & 2.0, reflect.TypeOf(uint32(0xffffffff) & 2.0))
}

// Test Uint32 & Complex
func TestCheckBinaryTypedExprUint32AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint32 & Bool
func TestCheckBinaryTypedExprUint32AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) & true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 & true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 & String
func TestCheckBinaryTypedExprUint32AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) & "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 & "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 & Nil
func TestCheckBinaryTypedExprUint32AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) & nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint32 % Int
func TestCheckBinaryTypedExprUint32RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) % 4`, env, uint32(0xffffffff) % 4, reflect.TypeOf(uint32(0xffffffff) % 4))
}

// Test Uint32 % Rune
func TestCheckBinaryTypedExprUint32RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) % '@'`, env, uint32(0xffffffff) % '@', reflect.TypeOf(uint32(0xffffffff) % '@'))
}

// Test Uint32 % Float
func TestCheckBinaryTypedExprUint32RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) % 2.0`, env, uint32(0xffffffff) % 2.0, reflect.TypeOf(uint32(0xffffffff) % 2.0))
}

// Test Uint32 % Complex
func TestCheckBinaryTypedExprUint32RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Uint32 % Bool
func TestCheckBinaryTypedExprUint32RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) % true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 % true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 % String
func TestCheckBinaryTypedExprUint32RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) % "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 % "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 % Nil
func TestCheckBinaryTypedExprUint32RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) % nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint32 == Int
func TestCheckBinaryTypedExprUint32EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) == 4`, env, uint32(0xffffffff) == 4, reflect.TypeOf(uint32(0xffffffff) == 4))
}

// Test Uint32 == Rune
func TestCheckBinaryTypedExprUint32EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) == '@'`, env, uint32(0xffffffff) == '@', reflect.TypeOf(uint32(0xffffffff) == '@'))
}

// Test Uint32 == Float
func TestCheckBinaryTypedExprUint32EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) == 2.0`, env, uint32(0xffffffff) == 2.0, reflect.TypeOf(uint32(0xffffffff) == 2.0))
}

// Test Uint32 == Complex
func TestCheckBinaryTypedExprUint32EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint32 == Bool
func TestCheckBinaryTypedExprUint32EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) == true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 == true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 == String
func TestCheckBinaryTypedExprUint32EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) == "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 == "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 == Nil
func TestCheckBinaryTypedExprUint32EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) == nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint32 > Int
func TestCheckBinaryTypedExprUint32GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) > 4`, env, uint32(0xffffffff) > 4, reflect.TypeOf(uint32(0xffffffff) > 4))
}

// Test Uint32 > Rune
func TestCheckBinaryTypedExprUint32GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) > '@'`, env, uint32(0xffffffff) > '@', reflect.TypeOf(uint32(0xffffffff) > '@'))
}

// Test Uint32 > Float
func TestCheckBinaryTypedExprUint32GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint32(0xffffffff) > 2.0`, env, uint32(0xffffffff) > 2.0, reflect.TypeOf(uint32(0xffffffff) > 2.0))
}

// Test Uint32 > Complex
func TestCheckBinaryTypedExprUint32GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint32 > Bool
func TestCheckBinaryTypedExprUint32GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) > true`, env,
		`cannot convert true to type uint32`,
		`invalid operation: 4294967295 > true (mismatched types uint32 and bool)`,
	)

}

// Test Uint32 > String
func TestCheckBinaryTypedExprUint32GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) > "abc"`, env,
		`cannot convert "abc" to type uint32`,
		`invalid operation: 4294967295 > "abc" (mismatched types uint32 and string)`,
	)

}

// Test Uint32 > Nil
func TestCheckBinaryTypedExprUint32GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint32(0xffffffff) > nil`, env,
		`cannot convert nil to type uint32`,
	)

}

// Test Uint64 + Int
func TestCheckBinaryTypedExprUint64AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + 4`, env,
		`constant 18446744073709551619 overflows uint64`,
	)

}

// Test Uint64 + Rune
func TestCheckBinaryTypedExprUint64AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + '@'`, env,
		`constant 18446744073709551679 overflows uint64`,
	)

}

// Test Uint64 + Float
func TestCheckBinaryTypedExprUint64AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + 2.0`, env,
		`constant 18446744073709551617 overflows uint64`,
	)

}

// Test Uint64 + Complex
func TestCheckBinaryTypedExprUint64AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint64 + Bool
func TestCheckBinaryTypedExprUint64AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 + true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 + String
func TestCheckBinaryTypedExprUint64AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 + "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 + Nil
func TestCheckBinaryTypedExprUint64AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) + nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Uint64 - Int
func TestCheckBinaryTypedExprUint64SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) - 4`, env, uint64(0xffffffffffffffff) - 4, reflect.TypeOf(uint64(0xffffffffffffffff) - 4))
}

// Test Uint64 - Rune
func TestCheckBinaryTypedExprUint64SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) - '@'`, env, uint64(0xffffffffffffffff) - '@', reflect.TypeOf(uint64(0xffffffffffffffff) - '@'))
}

// Test Uint64 - Float
func TestCheckBinaryTypedExprUint64SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) - 2.0`, env, uint64(0xffffffffffffffff) - 2.0, reflect.TypeOf(uint64(0xffffffffffffffff) - 2.0))
}

// Test Uint64 - Complex
func TestCheckBinaryTypedExprUint64SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint64 - Bool
func TestCheckBinaryTypedExprUint64SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) - true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 - true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 - String
func TestCheckBinaryTypedExprUint64SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) - "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 - "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 - Nil
func TestCheckBinaryTypedExprUint64SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) - nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Uint64 & Int
func TestCheckBinaryTypedExprUint64AndInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) & 4`, env, uint64(0xffffffffffffffff) & 4, reflect.TypeOf(uint64(0xffffffffffffffff) & 4))
}

// Test Uint64 & Rune
func TestCheckBinaryTypedExprUint64AndRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) & '@'`, env, uint64(0xffffffffffffffff) & '@', reflect.TypeOf(uint64(0xffffffffffffffff) & '@'))
}

// Test Uint64 & Float
func TestCheckBinaryTypedExprUint64AndFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) & 2.0`, env, uint64(0xffffffffffffffff) & 2.0, reflect.TypeOf(uint64(0xffffffffffffffff) & 2.0))
}

// Test Uint64 & Complex
func TestCheckBinaryTypedExprUint64AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint64 & Bool
func TestCheckBinaryTypedExprUint64AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) & true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 & true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 & String
func TestCheckBinaryTypedExprUint64AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) & "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 & "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 & Nil
func TestCheckBinaryTypedExprUint64AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) & nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Uint64 % Int
func TestCheckBinaryTypedExprUint64RemInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) % 4`, env, uint64(0xffffffffffffffff) % 4, reflect.TypeOf(uint64(0xffffffffffffffff) % 4))
}

// Test Uint64 % Rune
func TestCheckBinaryTypedExprUint64RemRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) % '@'`, env, uint64(0xffffffffffffffff) % '@', reflect.TypeOf(uint64(0xffffffffffffffff) % '@'))
}

// Test Uint64 % Float
func TestCheckBinaryTypedExprUint64RemFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) % 2.0`, env, uint64(0xffffffffffffffff) % 2.0, reflect.TypeOf(uint64(0xffffffffffffffff) % 2.0))
}

// Test Uint64 % Complex
func TestCheckBinaryTypedExprUint64RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Uint64 % Bool
func TestCheckBinaryTypedExprUint64RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) % true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 % true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 % String
func TestCheckBinaryTypedExprUint64RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) % "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 % "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 % Nil
func TestCheckBinaryTypedExprUint64RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) % nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Uint64 == Int
func TestCheckBinaryTypedExprUint64EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) == 4`, env, uint64(0xffffffffffffffff) == 4, reflect.TypeOf(uint64(0xffffffffffffffff) == 4))
}

// Test Uint64 == Rune
func TestCheckBinaryTypedExprUint64EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) == '@'`, env, uint64(0xffffffffffffffff) == '@', reflect.TypeOf(uint64(0xffffffffffffffff) == '@'))
}

// Test Uint64 == Float
func TestCheckBinaryTypedExprUint64EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) == 2.0`, env, uint64(0xffffffffffffffff) == 2.0, reflect.TypeOf(uint64(0xffffffffffffffff) == 2.0))
}

// Test Uint64 == Complex
func TestCheckBinaryTypedExprUint64EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint64 == Bool
func TestCheckBinaryTypedExprUint64EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) == true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 == true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 == String
func TestCheckBinaryTypedExprUint64EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) == "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 == "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 == Nil
func TestCheckBinaryTypedExprUint64EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) == nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Uint64 > Int
func TestCheckBinaryTypedExprUint64GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) > 4`, env, uint64(0xffffffffffffffff) > 4, reflect.TypeOf(uint64(0xffffffffffffffff) > 4))
}

// Test Uint64 > Rune
func TestCheckBinaryTypedExprUint64GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) > '@'`, env, uint64(0xffffffffffffffff) > '@', reflect.TypeOf(uint64(0xffffffffffffffff) > '@'))
}

// Test Uint64 > Float
func TestCheckBinaryTypedExprUint64GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `uint64(0xffffffffffffffff) > 2.0`, env, uint64(0xffffffffffffffff) > 2.0, reflect.TypeOf(uint64(0xffffffffffffffff) > 2.0))
}

// Test Uint64 > Complex
func TestCheckBinaryTypedExprUint64GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Uint64 > Bool
func TestCheckBinaryTypedExprUint64GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) > true`, env,
		`cannot convert true to type uint64`,
		`invalid operation: 18446744073709551615 > true (mismatched types uint64 and bool)`,
	)

}

// Test Uint64 > String
func TestCheckBinaryTypedExprUint64GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) > "abc"`, env,
		`cannot convert "abc" to type uint64`,
		`invalid operation: 18446744073709551615 > "abc" (mismatched types uint64 and string)`,
	)

}

// Test Uint64 > Nil
func TestCheckBinaryTypedExprUint64GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `uint64(0xffffffffffffffff) > nil`, env,
		`cannot convert nil to type uint64`,
	)

}

// Test Float32 + Int
func TestCheckBinaryTypedExprFloat32AddInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) + 4`, env, float32(0xffffffff) + 4, reflect.TypeOf(float32(0xffffffff) + 4))
}

// Test Float32 + Rune
func TestCheckBinaryTypedExprFloat32AddRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) + '@'`, env, float32(0xffffffff) + '@', reflect.TypeOf(float32(0xffffffff) + '@'))
}

// Test Float32 + Float
func TestCheckBinaryTypedExprFloat32AddFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) + 2.0`, env, float32(0xffffffff) + 2.0, reflect.TypeOf(float32(0xffffffff) + 2.0))
}

// Test Float32 + Complex
func TestCheckBinaryTypedExprFloat32AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float32 + Bool
func TestCheckBinaryTypedExprFloat32AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) + true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 + true (mismatched types float32 and bool)`,
	)

}

// Test Float32 + String
func TestCheckBinaryTypedExprFloat32AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) + "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 + "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 + Nil
func TestCheckBinaryTypedExprFloat32AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) + nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float32 - Int
func TestCheckBinaryTypedExprFloat32SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) - 4`, env, float32(0xffffffff) - 4, reflect.TypeOf(float32(0xffffffff) - 4))
}

// Test Float32 - Rune
func TestCheckBinaryTypedExprFloat32SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) - '@'`, env, float32(0xffffffff) - '@', reflect.TypeOf(float32(0xffffffff) - '@'))
}

// Test Float32 - Float
func TestCheckBinaryTypedExprFloat32SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) - 2.0`, env, float32(0xffffffff) - 2.0, reflect.TypeOf(float32(0xffffffff) - 2.0))
}

// Test Float32 - Complex
func TestCheckBinaryTypedExprFloat32SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float32 - Bool
func TestCheckBinaryTypedExprFloat32SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) - true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 - true (mismatched types float32 and bool)`,
	)

}

// Test Float32 - String
func TestCheckBinaryTypedExprFloat32SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) - "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 - "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 - Nil
func TestCheckBinaryTypedExprFloat32SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) - nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float32 & Int
func TestCheckBinaryTypedExprFloat32AndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & 4`, env,
		`invalid operation: 4.29497e+09 & 4 (operator & not defined on float32)`,
	)

}

// Test Float32 & Rune
func TestCheckBinaryTypedExprFloat32AndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & '@'`, env,
		`invalid operation: 4.29497e+09 & 64 (operator & not defined on float32)`,
	)

}

// Test Float32 & Float
func TestCheckBinaryTypedExprFloat32AndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & 2.0`, env,
		`invalid operation: 4.29497e+09 & 2 (operator & not defined on float32)`,
	)

}

// Test Float32 & Complex
func TestCheckBinaryTypedExprFloat32AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
		`invalid operation: 4.29497e+09 & 0 (operator & not defined on float32)`,
	)

}

// Test Float32 & Bool
func TestCheckBinaryTypedExprFloat32AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 & true (mismatched types float32 and bool)`,
	)

}

// Test Float32 & String
func TestCheckBinaryTypedExprFloat32AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 & "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 & Nil
func TestCheckBinaryTypedExprFloat32AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) & nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float32 % Int
func TestCheckBinaryTypedExprFloat32RemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % 4`, env,
		`invalid operation: 4.29497e+09 % 4 (operator % not defined on float32)`,
	)

}

// Test Float32 % Rune
func TestCheckBinaryTypedExprFloat32RemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % '@'`, env,
		`invalid operation: 4.29497e+09 % 64 (operator % not defined on float32)`,
	)

}

// Test Float32 % Float
func TestCheckBinaryTypedExprFloat32RemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % 2.0`, env,
		`invalid operation: 4.29497e+09 % 2 (operator % not defined on float32)`,
	)

}

// Test Float32 % Complex
func TestCheckBinaryTypedExprFloat32RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`invalid operation: 4.29497e+09 % 0 (operator % not defined on float32)`,
	)

}

// Test Float32 % Bool
func TestCheckBinaryTypedExprFloat32RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 % true (mismatched types float32 and bool)`,
	)

}

// Test Float32 % String
func TestCheckBinaryTypedExprFloat32RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 % "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 % Nil
func TestCheckBinaryTypedExprFloat32RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) % nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float32 == Int
func TestCheckBinaryTypedExprFloat32EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) == 4`, env, float32(0xffffffff) == 4, reflect.TypeOf(float32(0xffffffff) == 4))
}

// Test Float32 == Rune
func TestCheckBinaryTypedExprFloat32EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) == '@'`, env, float32(0xffffffff) == '@', reflect.TypeOf(float32(0xffffffff) == '@'))
}

// Test Float32 == Float
func TestCheckBinaryTypedExprFloat32EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) == 2.0`, env, float32(0xffffffff) == 2.0, reflect.TypeOf(float32(0xffffffff) == 2.0))
}

// Test Float32 == Complex
func TestCheckBinaryTypedExprFloat32EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float32 == Bool
func TestCheckBinaryTypedExprFloat32EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) == true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 == true (mismatched types float32 and bool)`,
	)

}

// Test Float32 == String
func TestCheckBinaryTypedExprFloat32EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) == "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 == "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 == Nil
func TestCheckBinaryTypedExprFloat32EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) == nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float32 > Int
func TestCheckBinaryTypedExprFloat32GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) > 4`, env, float32(0xffffffff) > 4, reflect.TypeOf(float32(0xffffffff) > 4))
}

// Test Float32 > Rune
func TestCheckBinaryTypedExprFloat32GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) > '@'`, env, float32(0xffffffff) > '@', reflect.TypeOf(float32(0xffffffff) > '@'))
}

// Test Float32 > Float
func TestCheckBinaryTypedExprFloat32GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float32(0xffffffff) > 2.0`, env, float32(0xffffffff) > 2.0, reflect.TypeOf(float32(0xffffffff) > 2.0))
}

// Test Float32 > Complex
func TestCheckBinaryTypedExprFloat32GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float32 > Bool
func TestCheckBinaryTypedExprFloat32GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) > true`, env,
		`cannot convert true to type float32`,
		`invalid operation: 4.29497e+09 > true (mismatched types float32 and bool)`,
	)

}

// Test Float32 > String
func TestCheckBinaryTypedExprFloat32GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) > "abc"`, env,
		`cannot convert "abc" to type float32`,
		`invalid operation: 4.29497e+09 > "abc" (mismatched types float32 and string)`,
	)

}

// Test Float32 > Nil
func TestCheckBinaryTypedExprFloat32GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float32(0xffffffff) > nil`, env,
		`cannot convert nil to type float32`,
	)

}

// Test Float64 + Int
func TestCheckBinaryTypedExprFloat64AddInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) + 4`, env, float64(0xffffffff) + 4, reflect.TypeOf(float64(0xffffffff) + 4))
}

// Test Float64 + Rune
func TestCheckBinaryTypedExprFloat64AddRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) + '@'`, env, float64(0xffffffff) + '@', reflect.TypeOf(float64(0xffffffff) + '@'))
}

// Test Float64 + Float
func TestCheckBinaryTypedExprFloat64AddFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) + 2.0`, env, float64(0xffffffff) + 2.0, reflect.TypeOf(float64(0xffffffff) + 2.0))
}

// Test Float64 + Complex
func TestCheckBinaryTypedExprFloat64AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) + 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float64 + Bool
func TestCheckBinaryTypedExprFloat64AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) + true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 + true (mismatched types float64 and bool)`,
	)

}

// Test Float64 + String
func TestCheckBinaryTypedExprFloat64AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) + "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 + "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 + Nil
func TestCheckBinaryTypedExprFloat64AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) + nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Float64 - Int
func TestCheckBinaryTypedExprFloat64SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) - 4`, env, float64(0xffffffff) - 4, reflect.TypeOf(float64(0xffffffff) - 4))
}

// Test Float64 - Rune
func TestCheckBinaryTypedExprFloat64SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) - '@'`, env, float64(0xffffffff) - '@', reflect.TypeOf(float64(0xffffffff) - '@'))
}

// Test Float64 - Float
func TestCheckBinaryTypedExprFloat64SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) - 2.0`, env, float64(0xffffffff) - 2.0, reflect.TypeOf(float64(0xffffffff) - 2.0))
}

// Test Float64 - Complex
func TestCheckBinaryTypedExprFloat64SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) - 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float64 - Bool
func TestCheckBinaryTypedExprFloat64SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) - true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 - true (mismatched types float64 and bool)`,
	)

}

// Test Float64 - String
func TestCheckBinaryTypedExprFloat64SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) - "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 - "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 - Nil
func TestCheckBinaryTypedExprFloat64SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) - nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Float64 & Int
func TestCheckBinaryTypedExprFloat64AndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & 4`, env,
		`invalid operation: 4.29497e+09 & 4 (operator & not defined on float64)`,
	)

}

// Test Float64 & Rune
func TestCheckBinaryTypedExprFloat64AndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & '@'`, env,
		`invalid operation: 4.29497e+09 & 64 (operator & not defined on float64)`,
	)

}

// Test Float64 & Float
func TestCheckBinaryTypedExprFloat64AndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & 2.0`, env,
		`invalid operation: 4.29497e+09 & 2 (operator & not defined on float64)`,
	)

}

// Test Float64 & Complex
func TestCheckBinaryTypedExprFloat64AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & 8.0i`, env,
		`constant 0+8i truncated to real`,
		`invalid operation: 4.29497e+09 & 0 (operator & not defined on float64)`,
	)

}

// Test Float64 & Bool
func TestCheckBinaryTypedExprFloat64AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 & true (mismatched types float64 and bool)`,
	)

}

// Test Float64 & String
func TestCheckBinaryTypedExprFloat64AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 & "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 & Nil
func TestCheckBinaryTypedExprFloat64AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) & nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Float64 % Int
func TestCheckBinaryTypedExprFloat64RemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % 4`, env,
		`invalid operation: 4.29497e+09 % 4 (operator % not defined on float64)`,
	)

}

// Test Float64 % Rune
func TestCheckBinaryTypedExprFloat64RemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % '@'`, env,
		`invalid operation: 4.29497e+09 % 64 (operator % not defined on float64)`,
	)

}

// Test Float64 % Float
func TestCheckBinaryTypedExprFloat64RemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % 2.0`, env,
		`invalid operation: 4.29497e+09 % 2 (operator % not defined on float64)`,
	)

}

// Test Float64 % Complex
func TestCheckBinaryTypedExprFloat64RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % 8.0i`, env,
		`constant 0+8i truncated to real`,
		`invalid operation: 4.29497e+09 % 0 (operator % not defined on float64)`,
	)

}

// Test Float64 % Bool
func TestCheckBinaryTypedExprFloat64RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 % true (mismatched types float64 and bool)`,
	)

}

// Test Float64 % String
func TestCheckBinaryTypedExprFloat64RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 % "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 % Nil
func TestCheckBinaryTypedExprFloat64RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) % nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Float64 == Int
func TestCheckBinaryTypedExprFloat64EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) == 4`, env, float64(0xffffffff) == 4, reflect.TypeOf(float64(0xffffffff) == 4))
}

// Test Float64 == Rune
func TestCheckBinaryTypedExprFloat64EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) == '@'`, env, float64(0xffffffff) == '@', reflect.TypeOf(float64(0xffffffff) == '@'))
}

// Test Float64 == Float
func TestCheckBinaryTypedExprFloat64EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) == 2.0`, env, float64(0xffffffff) == 2.0, reflect.TypeOf(float64(0xffffffff) == 2.0))
}

// Test Float64 == Complex
func TestCheckBinaryTypedExprFloat64EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) == 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float64 == Bool
func TestCheckBinaryTypedExprFloat64EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) == true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 == true (mismatched types float64 and bool)`,
	)

}

// Test Float64 == String
func TestCheckBinaryTypedExprFloat64EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) == "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 == "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 == Nil
func TestCheckBinaryTypedExprFloat64EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) == nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Float64 > Int
func TestCheckBinaryTypedExprFloat64GtrInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) > 4`, env, float64(0xffffffff) > 4, reflect.TypeOf(float64(0xffffffff) > 4))
}

// Test Float64 > Rune
func TestCheckBinaryTypedExprFloat64GtrRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) > '@'`, env, float64(0xffffffff) > '@', reflect.TypeOf(float64(0xffffffff) > '@'))
}

// Test Float64 > Float
func TestCheckBinaryTypedExprFloat64GtrFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `float64(0xffffffff) > 2.0`, env, float64(0xffffffff) > 2.0, reflect.TypeOf(float64(0xffffffff) > 2.0))
}

// Test Float64 > Complex
func TestCheckBinaryTypedExprFloat64GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) > 8.0i`, env,
		`constant 0+8i truncated to real`,
	)

}

// Test Float64 > Bool
func TestCheckBinaryTypedExprFloat64GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) > true`, env,
		`cannot convert true to type float64`,
		`invalid operation: 4.29497e+09 > true (mismatched types float64 and bool)`,
	)

}

// Test Float64 > String
func TestCheckBinaryTypedExprFloat64GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) > "abc"`, env,
		`cannot convert "abc" to type float64`,
		`invalid operation: 4.29497e+09 > "abc" (mismatched types float64 and string)`,
	)

}

// Test Float64 > Nil
func TestCheckBinaryTypedExprFloat64GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `float64(0xffffffff) > nil`, env,
		`cannot convert nil to type float64`,
	)

}

// Test Complex64 + Int
func TestCheckBinaryTypedExprComplex64AddInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) + 4`, env, complex64(0xffffffff + 0xffffffff * 1i) + 4, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) + 4))
}

// Test Complex64 + Rune
func TestCheckBinaryTypedExprComplex64AddRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) + '@'`, env, complex64(0xffffffff + 0xffffffff * 1i) + '@', reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) + '@'))
}

// Test Complex64 + Float
func TestCheckBinaryTypedExprComplex64AddFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) + 2.0`, env, complex64(0xffffffff + 0xffffffff * 1i) + 2.0, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) + 2.0))
}

// Test Complex64 + Complex
func TestCheckBinaryTypedExprComplex64AddComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) + 8.0i`, env, complex64(0xffffffff + 0xffffffff * 1i) + 8.0i, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) + 8.0i))
}

// Test Complex64 + Bool
func TestCheckBinaryTypedExprComplex64AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) + true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) + true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 + String
func TestCheckBinaryTypedExprComplex64AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) + "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) + "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 + Nil
func TestCheckBinaryTypedExprComplex64AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) + nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex64 - Int
func TestCheckBinaryTypedExprComplex64SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) - 4`, env, complex64(0xffffffff + 0xffffffff * 1i) - 4, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) - 4))
}

// Test Complex64 - Rune
func TestCheckBinaryTypedExprComplex64SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) - '@'`, env, complex64(0xffffffff + 0xffffffff * 1i) - '@', reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) - '@'))
}

// Test Complex64 - Float
func TestCheckBinaryTypedExprComplex64SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) - 2.0`, env, complex64(0xffffffff + 0xffffffff * 1i) - 2.0, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) - 2.0))
}

// Test Complex64 - Complex
func TestCheckBinaryTypedExprComplex64SubComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) - 8.0i`, env, complex64(0xffffffff + 0xffffffff * 1i) - 8.0i, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) - 8.0i))
}

// Test Complex64 - Bool
func TestCheckBinaryTypedExprComplex64SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) - true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) - true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 - String
func TestCheckBinaryTypedExprComplex64SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) - "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) - "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 - Nil
func TestCheckBinaryTypedExprComplex64SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) - nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex64 & Int
func TestCheckBinaryTypedExprComplex64AndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 4 (operator & not defined on complex64)`,
	)

}

// Test Complex64 & Rune
func TestCheckBinaryTypedExprComplex64AndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 64 (operator & not defined on complex64)`,
	)

}

// Test Complex64 & Float
func TestCheckBinaryTypedExprComplex64AndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 2 (operator & not defined on complex64)`,
	)

}

// Test Complex64 & Complex
func TestCheckBinaryTypedExprComplex64AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 8i (operator & not defined on complex64)`,
	)

}

// Test Complex64 & Bool
func TestCheckBinaryTypedExprComplex64AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) & true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 & String
func TestCheckBinaryTypedExprComplex64AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) & "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 & Nil
func TestCheckBinaryTypedExprComplex64AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) & nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex64 % Int
func TestCheckBinaryTypedExprComplex64RemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 4 (operator % not defined on complex64)`,
	)

}

// Test Complex64 % Rune
func TestCheckBinaryTypedExprComplex64RemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 64 (operator % not defined on complex64)`,
	)

}

// Test Complex64 % Float
func TestCheckBinaryTypedExprComplex64RemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 2 (operator % not defined on complex64)`,
	)

}

// Test Complex64 % Complex
func TestCheckBinaryTypedExprComplex64RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 8i (operator % not defined on complex64)`,
	)

}

// Test Complex64 % Bool
func TestCheckBinaryTypedExprComplex64RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) % true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 % String
func TestCheckBinaryTypedExprComplex64RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) % "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 % Nil
func TestCheckBinaryTypedExprComplex64RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) % nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex64 == Int
func TestCheckBinaryTypedExprComplex64EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) == 4`, env, complex64(0xffffffff + 0xffffffff * 1i) == 4, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) == 4))
}

// Test Complex64 == Rune
func TestCheckBinaryTypedExprComplex64EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) == '@'`, env, complex64(0xffffffff + 0xffffffff * 1i) == '@', reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) == '@'))
}

// Test Complex64 == Float
func TestCheckBinaryTypedExprComplex64EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) == 2.0`, env, complex64(0xffffffff + 0xffffffff * 1i) == 2.0, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) == 2.0))
}

// Test Complex64 == Complex
func TestCheckBinaryTypedExprComplex64EqlComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex64(0xffffffff + 0xffffffff * 1i) == 8.0i`, env, complex64(0xffffffff + 0xffffffff * 1i) == 8.0i, reflect.TypeOf(complex64(0xffffffff + 0xffffffff * 1i) == 8.0i))
}

// Test Complex64 == Bool
func TestCheckBinaryTypedExprComplex64EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) == true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) == true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 == String
func TestCheckBinaryTypedExprComplex64EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) == "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) == "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 == Nil
func TestCheckBinaryTypedExprComplex64EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) == nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex64 > Int
func TestCheckBinaryTypedExprComplex64GtrInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 4 (operator > not defined on complex64)`,
	)

}

// Test Complex64 > Rune
func TestCheckBinaryTypedExprComplex64GtrRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 64 (operator > not defined on complex64)`,
	)

}

// Test Complex64 > Float
func TestCheckBinaryTypedExprComplex64GtrFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 2 (operator > not defined on complex64)`,
	)

}

// Test Complex64 > Complex
func TestCheckBinaryTypedExprComplex64GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 8i (operator > not defined on complex64)`,
	)

}

// Test Complex64 > Bool
func TestCheckBinaryTypedExprComplex64GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > true`, env,
		`cannot convert true to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) > true (mismatched types complex64 and bool)`,
	)

}

// Test Complex64 > String
func TestCheckBinaryTypedExprComplex64GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > "abc"`, env,
		`cannot convert "abc" to type complex64`,
		`invalid operation: (4.29497e+09+4.29497e+09i) > "abc" (mismatched types complex64 and string)`,
	)

}

// Test Complex64 > Nil
func TestCheckBinaryTypedExprComplex64GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex64(0xffffffff + 0xffffffff * 1i) > nil`, env,
		`cannot convert nil to type complex64`,
	)

}

// Test Complex128 + Int
func TestCheckBinaryTypedExprComplex128AddInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) + 4`, env, complex128(0xffffffff + 0xffffffff * 1i) + 4, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) + 4))
}

// Test Complex128 + Rune
func TestCheckBinaryTypedExprComplex128AddRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) + '@'`, env, complex128(0xffffffff + 0xffffffff * 1i) + '@', reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) + '@'))
}

// Test Complex128 + Float
func TestCheckBinaryTypedExprComplex128AddFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) + 2.0`, env, complex128(0xffffffff + 0xffffffff * 1i) + 2.0, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) + 2.0))
}

// Test Complex128 + Complex
func TestCheckBinaryTypedExprComplex128AddComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) + 8.0i`, env, complex128(0xffffffff + 0xffffffff * 1i) + 8.0i, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) + 8.0i))
}

// Test Complex128 + Bool
func TestCheckBinaryTypedExprComplex128AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) + true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) + true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 + String
func TestCheckBinaryTypedExprComplex128AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) + "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) + "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 + Nil
func TestCheckBinaryTypedExprComplex128AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) + nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Complex128 - Int
func TestCheckBinaryTypedExprComplex128SubInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) - 4`, env, complex128(0xffffffff + 0xffffffff * 1i) - 4, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) - 4))
}

// Test Complex128 - Rune
func TestCheckBinaryTypedExprComplex128SubRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) - '@'`, env, complex128(0xffffffff + 0xffffffff * 1i) - '@', reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) - '@'))
}

// Test Complex128 - Float
func TestCheckBinaryTypedExprComplex128SubFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) - 2.0`, env, complex128(0xffffffff + 0xffffffff * 1i) - 2.0, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) - 2.0))
}

// Test Complex128 - Complex
func TestCheckBinaryTypedExprComplex128SubComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) - 8.0i`, env, complex128(0xffffffff + 0xffffffff * 1i) - 8.0i, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) - 8.0i))
}

// Test Complex128 - Bool
func TestCheckBinaryTypedExprComplex128SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) - true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) - true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 - String
func TestCheckBinaryTypedExprComplex128SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) - "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) - "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 - Nil
func TestCheckBinaryTypedExprComplex128SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) - nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Complex128 & Int
func TestCheckBinaryTypedExprComplex128AndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 4 (operator & not defined on complex128)`,
	)

}

// Test Complex128 & Rune
func TestCheckBinaryTypedExprComplex128AndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 64 (operator & not defined on complex128)`,
	)

}

// Test Complex128 & Float
func TestCheckBinaryTypedExprComplex128AndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 2 (operator & not defined on complex128)`,
	)

}

// Test Complex128 & Complex
func TestCheckBinaryTypedExprComplex128AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) & 8i (operator & not defined on complex128)`,
	)

}

// Test Complex128 & Bool
func TestCheckBinaryTypedExprComplex128AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) & true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 & String
func TestCheckBinaryTypedExprComplex128AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) & "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 & Nil
func TestCheckBinaryTypedExprComplex128AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) & nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Complex128 % Int
func TestCheckBinaryTypedExprComplex128RemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 4 (operator % not defined on complex128)`,
	)

}

// Test Complex128 % Rune
func TestCheckBinaryTypedExprComplex128RemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 64 (operator % not defined on complex128)`,
	)

}

// Test Complex128 % Float
func TestCheckBinaryTypedExprComplex128RemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 2 (operator % not defined on complex128)`,
	)

}

// Test Complex128 % Complex
func TestCheckBinaryTypedExprComplex128RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) % 8i (operator % not defined on complex128)`,
	)

}

// Test Complex128 % Bool
func TestCheckBinaryTypedExprComplex128RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) % true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 % String
func TestCheckBinaryTypedExprComplex128RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) % "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 % Nil
func TestCheckBinaryTypedExprComplex128RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) % nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Complex128 == Int
func TestCheckBinaryTypedExprComplex128EqlInt(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) == 4`, env, complex128(0xffffffff + 0xffffffff * 1i) == 4, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) == 4))
}

// Test Complex128 == Rune
func TestCheckBinaryTypedExprComplex128EqlRune(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) == '@'`, env, complex128(0xffffffff + 0xffffffff * 1i) == '@', reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) == '@'))
}

// Test Complex128 == Float
func TestCheckBinaryTypedExprComplex128EqlFloat(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) == 2.0`, env, complex128(0xffffffff + 0xffffffff * 1i) == 2.0, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) == 2.0))
}

// Test Complex128 == Complex
func TestCheckBinaryTypedExprComplex128EqlComplex(t *testing.T) {
	env := makeEnv()

	expectConst(t, `complex128(0xffffffff + 0xffffffff * 1i) == 8.0i`, env, complex128(0xffffffff + 0xffffffff * 1i) == 8.0i, reflect.TypeOf(complex128(0xffffffff + 0xffffffff * 1i) == 8.0i))
}

// Test Complex128 == Bool
func TestCheckBinaryTypedExprComplex128EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) == true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) == true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 == String
func TestCheckBinaryTypedExprComplex128EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) == "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) == "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 == Nil
func TestCheckBinaryTypedExprComplex128EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) == nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Complex128 > Int
func TestCheckBinaryTypedExprComplex128GtrInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > 4`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 4 (operator > not defined on complex128)`,
	)

}

// Test Complex128 > Rune
func TestCheckBinaryTypedExprComplex128GtrRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > '@'`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 64 (operator > not defined on complex128)`,
	)

}

// Test Complex128 > Float
func TestCheckBinaryTypedExprComplex128GtrFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > 2.0`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 2 (operator > not defined on complex128)`,
	)

}

// Test Complex128 > Complex
func TestCheckBinaryTypedExprComplex128GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > 8.0i`, env,
		`invalid operation: (4.29497e+09+4.29497e+09i) > 8i (operator > not defined on complex128)`,
	)

}

// Test Complex128 > Bool
func TestCheckBinaryTypedExprComplex128GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > true`, env,
		`cannot convert true to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) > true (mismatched types complex128 and bool)`,
	)

}

// Test Complex128 > String
func TestCheckBinaryTypedExprComplex128GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > "abc"`, env,
		`cannot convert "abc" to type complex128`,
		`invalid operation: (4.29497e+09+4.29497e+09i) > "abc" (mismatched types complex128 and string)`,
	)

}

// Test Complex128 > Nil
func TestCheckBinaryTypedExprComplex128GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `complex128(0xffffffff + 0xffffffff * 1i) > nil`, env,
		`cannot convert nil to type complex128`,
	)

}

// Test Rune32 + Int
func TestCheckBinaryTypedExprRune32AddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + 4`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967299 overflows rune`,
	)

}

// Test Rune32 + Rune
func TestCheckBinaryTypedExprRune32AddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + '@'`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967359 overflows rune`,
	)

}

// Test Rune32 + Float
func TestCheckBinaryTypedExprRune32AddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + 2.0`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967297 overflows rune`,
	)

}

// Test Rune32 + Complex
func TestCheckBinaryTypedExprRune32AddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 + Bool
func TestCheckBinaryTypedExprRune32AddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) + true (mismatched types rune and bool)`,
	)

}

// Test Rune32 + String
func TestCheckBinaryTypedExprRune32AddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) + "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 + Nil
func TestCheckBinaryTypedExprRune32AddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) + nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test Rune32 - Int
func TestCheckBinaryTypedExprRune32SubInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - 4`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967291 overflows rune`,
	)

}

// Test Rune32 - Rune
func TestCheckBinaryTypedExprRune32SubRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - '@'`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967231 overflows rune`,
	)

}

// Test Rune32 - Float
func TestCheckBinaryTypedExprRune32SubFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - 2.0`, env,
		`constant 4294967295 overflows rune`,
		`constant 4294967293 overflows rune`,
	)

}

// Test Rune32 - Complex
func TestCheckBinaryTypedExprRune32SubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 - Bool
func TestCheckBinaryTypedExprRune32SubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) - true (mismatched types rune and bool)`,
	)

}

// Test Rune32 - String
func TestCheckBinaryTypedExprRune32SubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) - "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 - Nil
func TestCheckBinaryTypedExprRune32SubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) - nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test Rune32 & Int
func TestCheckBinaryTypedExprRune32AndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & 4`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 & Rune
func TestCheckBinaryTypedExprRune32AndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & '@'`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 & Float
func TestCheckBinaryTypedExprRune32AndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & 2.0`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 & Complex
func TestCheckBinaryTypedExprRune32AndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
	)

}

// Test Rune32 & Bool
func TestCheckBinaryTypedExprRune32AndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) & true (mismatched types rune and bool)`,
	)

}

// Test Rune32 & String
func TestCheckBinaryTypedExprRune32AndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) & "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 & Nil
func TestCheckBinaryTypedExprRune32AndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) & nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test Rune32 % Int
func TestCheckBinaryTypedExprRune32RemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % 4`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 % Rune
func TestCheckBinaryTypedExprRune32RemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % '@'`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 % Float
func TestCheckBinaryTypedExprRune32RemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % 2.0`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 % Complex
func TestCheckBinaryTypedExprRune32RemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
		`division by zero`,
	)

}

// Test Rune32 % Bool
func TestCheckBinaryTypedExprRune32RemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) % true (mismatched types rune and bool)`,
	)

}

// Test Rune32 % String
func TestCheckBinaryTypedExprRune32RemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) % "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 % Nil
func TestCheckBinaryTypedExprRune32RemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) % nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test Rune32 == Int
func TestCheckBinaryTypedExprRune32EqlInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == 4`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 == Rune
func TestCheckBinaryTypedExprRune32EqlRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == '@'`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 == Float
func TestCheckBinaryTypedExprRune32EqlFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == 2.0`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 == Complex
func TestCheckBinaryTypedExprRune32EqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
	)

}

// Test Rune32 == Bool
func TestCheckBinaryTypedExprRune32EqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) == true (mismatched types rune and bool)`,
	)

}

// Test Rune32 == String
func TestCheckBinaryTypedExprRune32EqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) == "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 == Nil
func TestCheckBinaryTypedExprRune32EqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) == nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test Rune32 > Int
func TestCheckBinaryTypedExprRune32GtrInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > 4`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 > Rune
func TestCheckBinaryTypedExprRune32GtrRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > '@'`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 > Float
func TestCheckBinaryTypedExprRune32GtrFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > 2.0`, env,
		`constant 4294967295 overflows rune`,
	)

}

// Test Rune32 > Complex
func TestCheckBinaryTypedExprRune32GtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > 8.0i`, env,
		`constant 4294967295 overflows rune`,
		`constant 0+8i truncated to real`,
	)

}

// Test Rune32 > Bool
func TestCheckBinaryTypedExprRune32GtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > true`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert true to type rune`,
		`invalid operation: rune(4294967295) > true (mismatched types rune and bool)`,
	)

}

// Test Rune32 > String
func TestCheckBinaryTypedExprRune32GtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > "abc"`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert "abc" to type rune`,
		`invalid operation: rune(4294967295) > "abc" (mismatched types rune and string)`,
	)

}

// Test Rune32 > Nil
func TestCheckBinaryTypedExprRune32GtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `rune(0xffffffff) > nil`, env,
		`constant 4294967295 overflows rune`,
		`cannot convert nil to type rune`,
	)

}

// Test StringT + Int
func TestCheckBinaryTypedExprStringTAddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" + 4 (mismatched types string and int)`,
	)

}

// Test StringT + Rune
func TestCheckBinaryTypedExprStringTAddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" + rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT + Float
func TestCheckBinaryTypedExprStringTAddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" + 2 (mismatched types string and float64)`,
	)

}

// Test StringT + Complex
func TestCheckBinaryTypedExprStringTAddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" + 8i (mismatched types string and complex128)`,
	)

}

// Test StringT + Bool
func TestCheckBinaryTypedExprStringTAddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" + true (mismatched types string and bool)`,
	)

}

// Test StringT + String
func TestCheckBinaryTypedExprStringTAddString(t *testing.T) {
	env := makeEnv()

	expectConst(t, `string("abc") + "abc"`, env, string("abc") + "abc", reflect.TypeOf(string("abc") + "abc"))
}

// Test StringT + Nil
func TestCheckBinaryTypedExprStringTAddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") + nil`, env,
		`invalid operation: "abc" + nil (mismatched types string and nil)`,
	)

}

// Test StringT - Int
func TestCheckBinaryTypedExprStringTSubInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" - 4 (mismatched types string and int)`,
	)

}

// Test StringT - Rune
func TestCheckBinaryTypedExprStringTSubRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" - rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT - Float
func TestCheckBinaryTypedExprStringTSubFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" - 2 (mismatched types string and float64)`,
	)

}

// Test StringT - Complex
func TestCheckBinaryTypedExprStringTSubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" - 8i (mismatched types string and complex128)`,
	)

}

// Test StringT - Bool
func TestCheckBinaryTypedExprStringTSubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" - true (mismatched types string and bool)`,
	)

}

// Test StringT - String
func TestCheckBinaryTypedExprStringTSubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - "abc"`, env,
		`invalid operation: "abc" - "abc" (operator - not defined on string)`,
	)

}

// Test StringT - Nil
func TestCheckBinaryTypedExprStringTSubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") - nil`, env,
		`invalid operation: "abc" - nil (mismatched types string and nil)`,
	)

}

// Test StringT & Int
func TestCheckBinaryTypedExprStringTAndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" & 4 (mismatched types string and int)`,
	)

}

// Test StringT & Rune
func TestCheckBinaryTypedExprStringTAndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" & rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT & Float
func TestCheckBinaryTypedExprStringTAndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" & 2 (mismatched types string and float64)`,
	)

}

// Test StringT & Complex
func TestCheckBinaryTypedExprStringTAndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" & 8i (mismatched types string and complex128)`,
	)

}

// Test StringT & Bool
func TestCheckBinaryTypedExprStringTAndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" & true (mismatched types string and bool)`,
	)

}

// Test StringT & String
func TestCheckBinaryTypedExprStringTAndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & "abc"`, env,
		`invalid operation: "abc" & "abc" (operator & not defined on string)`,
	)

}

// Test StringT & Nil
func TestCheckBinaryTypedExprStringTAndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") & nil`, env,
		`invalid operation: "abc" & nil (mismatched types string and nil)`,
	)

}

// Test StringT % Int
func TestCheckBinaryTypedExprStringTRemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" % 4 (mismatched types string and int)`,
	)

}

// Test StringT % Rune
func TestCheckBinaryTypedExprStringTRemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" % rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT % Float
func TestCheckBinaryTypedExprStringTRemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" % 2 (mismatched types string and float64)`,
	)

}

// Test StringT % Complex
func TestCheckBinaryTypedExprStringTRemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" % 8i (mismatched types string and complex128)`,
	)

}

// Test StringT % Bool
func TestCheckBinaryTypedExprStringTRemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" % true (mismatched types string and bool)`,
	)

}

// Test StringT % String
func TestCheckBinaryTypedExprStringTRemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % "abc"`, env,
		`invalid operation: "abc" % "abc" (operator % not defined on string)`,
	)

}

// Test StringT % Nil
func TestCheckBinaryTypedExprStringTRemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") % nil`, env,
		`invalid operation: "abc" % nil (mismatched types string and nil)`,
	)

}

// Test StringT == Int
func TestCheckBinaryTypedExprStringTEqlInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" == 4 (mismatched types string and int)`,
	)

}

// Test StringT == Rune
func TestCheckBinaryTypedExprStringTEqlRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" == rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT == Float
func TestCheckBinaryTypedExprStringTEqlFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" == 2 (mismatched types string and float64)`,
	)

}

// Test StringT == Complex
func TestCheckBinaryTypedExprStringTEqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" == 8i (mismatched types string and complex128)`,
	)

}

// Test StringT == Bool
func TestCheckBinaryTypedExprStringTEqlBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" == true (mismatched types string and bool)`,
	)

}

// Test StringT == String
func TestCheckBinaryTypedExprStringTEqlString(t *testing.T) {
	env := makeEnv()

	expectConst(t, `string("abc") == "abc"`, env, string("abc") == "abc", reflect.TypeOf(string("abc") == "abc"))
}

// Test StringT == Nil
func TestCheckBinaryTypedExprStringTEqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") == nil`, env,
		`invalid operation: "abc" == nil (mismatched types string and nil)`,
	)

}

// Test StringT > Int
func TestCheckBinaryTypedExprStringTGtrInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > 4`, env,
		`cannot convert 4 to type string`,
		`invalid operation: "abc" > 4 (mismatched types string and int)`,
	)

}

// Test StringT > Rune
func TestCheckBinaryTypedExprStringTGtrRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > '@'`, env,
		`cannot convert '@' to type string`,
		`invalid operation: "abc" > rune(64) (mismatched types string and rune)`,
	)

}

// Test StringT > Float
func TestCheckBinaryTypedExprStringTGtrFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > 2.0`, env,
		`cannot convert 2 to type string`,
		`invalid operation: "abc" > 2 (mismatched types string and float64)`,
	)

}

// Test StringT > Complex
func TestCheckBinaryTypedExprStringTGtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > 8.0i`, env,
		`cannot convert 8i to type string`,
		`invalid operation: "abc" > 8i (mismatched types string and complex128)`,
	)

}

// Test StringT > Bool
func TestCheckBinaryTypedExprStringTGtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > true`, env,
		`cannot convert true to type string`,
		`invalid operation: "abc" > true (mismatched types string and bool)`,
	)

}

// Test StringT > String
func TestCheckBinaryTypedExprStringTGtrString(t *testing.T) {
	env := makeEnv()

	expectConst(t, `string("abc") > "abc"`, env, string("abc") > "abc", reflect.TypeOf(string("abc") > "abc"))
}

// Test StringT > Nil
func TestCheckBinaryTypedExprStringTGtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `string("abc") > nil`, env,
		`invalid operation: "abc" > nil (mismatched types string and nil)`,
	)

}

// Test BoolT + Int
func TestCheckBinaryTypedExprBoolTAddInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true + 4 (mismatched types bool and int)`,
	)

}

// Test BoolT + Rune
func TestCheckBinaryTypedExprBoolTAddRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true + rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT + Float
func TestCheckBinaryTypedExprBoolTAddFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true + 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT + Complex
func TestCheckBinaryTypedExprBoolTAddComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true + 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT + Bool
func TestCheckBinaryTypedExprBoolTAddBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + true`, env,
		`invalid operation: true + true (operator + not defined on bool)`,
	)

}

// Test BoolT + String
func TestCheckBinaryTypedExprBoolTAddString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true + "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT + Nil
func TestCheckBinaryTypedExprBoolTAddNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) + nil`, env,
		`cannot convert nil to type bool`,
	)

}

// Test BoolT - Int
func TestCheckBinaryTypedExprBoolTSubInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true - 4 (mismatched types bool and int)`,
	)

}

// Test BoolT - Rune
func TestCheckBinaryTypedExprBoolTSubRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true - rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT - Float
func TestCheckBinaryTypedExprBoolTSubFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true - 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT - Complex
func TestCheckBinaryTypedExprBoolTSubComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true - 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT - Bool
func TestCheckBinaryTypedExprBoolTSubBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - true`, env,
		`invalid operation: true - true (operator - not defined on bool)`,
	)

}

// Test BoolT - String
func TestCheckBinaryTypedExprBoolTSubString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true - "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT - Nil
func TestCheckBinaryTypedExprBoolTSubNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) - nil`, env,
		`cannot convert nil to type bool`,
	)

}

// Test BoolT & Int
func TestCheckBinaryTypedExprBoolTAndInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true & 4 (mismatched types bool and int)`,
	)

}

// Test BoolT & Rune
func TestCheckBinaryTypedExprBoolTAndRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true & rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT & Float
func TestCheckBinaryTypedExprBoolTAndFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true & 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT & Complex
func TestCheckBinaryTypedExprBoolTAndComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true & 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT & Bool
func TestCheckBinaryTypedExprBoolTAndBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & true`, env,
		`invalid operation: true & true (operator & not defined on bool)`,
	)

}

// Test BoolT & String
func TestCheckBinaryTypedExprBoolTAndString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true & "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT & Nil
func TestCheckBinaryTypedExprBoolTAndNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) & nil`, env,
		`cannot convert nil to type bool`,
	)

}

// Test BoolT % Int
func TestCheckBinaryTypedExprBoolTRemInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true % 4 (mismatched types bool and int)`,
	)

}

// Test BoolT % Rune
func TestCheckBinaryTypedExprBoolTRemRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true % rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT % Float
func TestCheckBinaryTypedExprBoolTRemFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true % 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT % Complex
func TestCheckBinaryTypedExprBoolTRemComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true % 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT % Bool
func TestCheckBinaryTypedExprBoolTRemBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % true`, env,
		`invalid operation: true % true (operator % not defined on bool)`,
	)

}

// Test BoolT % String
func TestCheckBinaryTypedExprBoolTRemString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true % "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT % Nil
func TestCheckBinaryTypedExprBoolTRemNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) % nil`, env,
		`cannot convert nil to type bool`,
	)

}

// Test BoolT == Int
func TestCheckBinaryTypedExprBoolTEqlInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true == 4 (mismatched types bool and int)`,
	)

}

// Test BoolT == Rune
func TestCheckBinaryTypedExprBoolTEqlRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true == rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT == Float
func TestCheckBinaryTypedExprBoolTEqlFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true == 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT == Complex
func TestCheckBinaryTypedExprBoolTEqlComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true == 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT == Bool
func TestCheckBinaryTypedExprBoolTEqlBool(t *testing.T) {
	env := makeEnv()

	expectConst(t, `bool(true) == true`, env, bool(true) == true, reflect.TypeOf(bool(true) == true))
}

// Test BoolT == String
func TestCheckBinaryTypedExprBoolTEqlString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true == "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT == Nil
func TestCheckBinaryTypedExprBoolTEqlNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) == nil`, env,
		`cannot convert nil to type bool`,
	)

}

// Test BoolT > Int
func TestCheckBinaryTypedExprBoolTGtrInt(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > 4`, env,
		`cannot convert 4 to type bool`,
		`invalid operation: true > 4 (mismatched types bool and int)`,
	)

}

// Test BoolT > Rune
func TestCheckBinaryTypedExprBoolTGtrRune(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > '@'`, env,
		`cannot convert '@' to type bool`,
		`invalid operation: true > rune(64) (mismatched types bool and rune)`,
	)

}

// Test BoolT > Float
func TestCheckBinaryTypedExprBoolTGtrFloat(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > 2.0`, env,
		`cannot convert 2 to type bool`,
		`invalid operation: true > 2 (mismatched types bool and float64)`,
	)

}

// Test BoolT > Complex
func TestCheckBinaryTypedExprBoolTGtrComplex(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > 8.0i`, env,
		`cannot convert 8i to type bool`,
		`invalid operation: true > 8i (mismatched types bool and complex128)`,
	)

}

// Test BoolT > Bool
func TestCheckBinaryTypedExprBoolTGtrBool(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > true`, env,
		`invalid operation: true > true (operator > not defined on bool)`,
	)

}

// Test BoolT > String
func TestCheckBinaryTypedExprBoolTGtrString(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > "abc"`, env,
		`cannot convert "abc" to type bool`,
		`invalid operation: true > "abc" (mismatched types bool and string)`,
	)

}

// Test BoolT > Nil
func TestCheckBinaryTypedExprBoolTGtrNil(t *testing.T) {
	env := makeEnv()

	expectCheckError(t, `bool(true) > nil`, env,
		`cannot convert nil to type bool`,
	)

}
