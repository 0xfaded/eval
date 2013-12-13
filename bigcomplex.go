package interactive

import (
	"math/big"
)

// Big complex behaves like a *big.Rat, but has an imaginary component
// and separate implementation for + - * /
type BigComplex struct {
	big.Rat
	Imag big.Rat
}

// Use with token.INT ast.BasicLit
func NewBigInteger(i string) (*BigComplex, bool) {
	z := new(BigComplex)
	z.Rat.Denom().SetInt64(1)
	_, ok := z.Num().SetString(i, 0)
	return z, ok
}

// Use with token.FLOAT ast.BasicLit
func NewBigReal(r string) (*BigComplex, bool) {
	z := new(BigComplex)
	_, ok := z.Rat.SetString(r)
	return z, ok
}

// Use with token.IMAG ast.BasicLit
func NewBigImag(i string) (*BigComplex, bool) {
	z := new(BigComplex)
	ok := i[len(i)-1] == 'i'
	if ok {
		_, ok = z.Imag.SetString(i[:len(i)-1])
	}
	return z, ok
}

// Use with token.CHAR ast.BasicLit
func NewBigRune(n rune) *BigComplex {
	z := new(BigComplex)
	z.Rat.Denom().SetInt64(1)
	z.Num().SetInt64(int64(n))
	return z
}

func NewBigInt64(i int64) *BigComplex {
	z := new(BigComplex)
	z.Rat.Denom().SetInt64(1)
	z.Num().SetInt64(i)
	return z
}

func NewBigUint64(u uint64) *BigComplex {
	z := new(BigComplex)
	z.Rat.Denom().SetInt64(1)
	z.Num().SetUint64(u)
	return z
}

func NewBigFloat64(f float64) *BigComplex {
	z := new(BigComplex)
	z.Rat.SetFloat64(f)
	return z
}

func NewBigComplex128(c complex128) *BigComplex {
	z := new(BigComplex)
	z.Rat.SetFloat64(real(c))
	z.Imag.SetFloat64(imag(c))
	return z
}

func (z *BigComplex) Add(x, y *BigComplex) *BigComplex {
	z.Rat.Add(&x.Rat, &y.Rat)
	z.Imag.Add(&x.Imag, &y.Imag)
	return z
}

func (z *BigComplex) Sub(x, y *BigComplex) *BigComplex {
	z.Rat.Sub(&x.Rat, &y.Rat)
	z.Imag.Sub(&x.Imag, &y.Imag)
	return z
}

func (z *BigComplex) Mul(x, y *BigComplex) *BigComplex {
	r := new(big.Rat).Mul(&x.Rat, &y.Rat)
	r.Sub(r, new(big.Rat).Mul(&x.Imag, &y.Imag))

	i := new(big.Rat).Mul(&x.Rat, &y.Imag)
	i.Add(i, new(big.Rat).Mul(&x.Imag, &y.Rat))

	z.Rat = *r
	z.Imag = *i
	return z
}

func (z *BigComplex) Quo(x, y *BigComplex) *BigComplex {
	// a+bi   ac+bd   bc-ad
	// ---- = ----- + ----- i
	// c+di   cc+dd   cc+dd

	cc := new(big.Rat).Mul(&y.Rat, &y.Rat)
	dd := new(big.Rat).Mul(&y.Imag, &y.Imag)
	ccdd := new(big.Rat).Add(cc, dd)

	ac := new(big.Rat).Mul(&x.Rat, &y.Rat)
	ad := new(big.Rat).Mul(&x.Rat, &y.Imag)
	bc := new(big.Rat).Mul(&x.Imag, &y.Rat)
	bd := new(big.Rat).Mul(&x.Imag, &y.Imag)

	z.Rat.Add(ac, bd)
	z.Rat.Quo(&z.Rat, ccdd)

	z.Imag.Sub(bc, ad)
	z.Imag.Quo(&z.Imag, ccdd)

	return z
}

func (z *BigComplex) IsZero() bool {
	return z.Rat.Num().BitLen() == 0 && z.Imag.Num().BitLen() == 0
}

// Return a representation of z truncated to be an int of length bits.
// Valid values for bits are 8, 16, 32, 64. Result is otherwise undefined
// If a truncation occurs, the decimal part is dropped and the conversion
// continues as usual. truncation will be true
// If an overflow occurs, the result is equivelant to a cast of the form
// int32(x). overflow will be true.
func (z *BigComplex) Int(bits int) (_ int64, truncation, overflow bool) {
	var integer *BigComplex
	integer, truncation = z.Integer()
	num := integer.Num()

	// Numerator must fit in bits - 1, with 1 bit left for sign
	if overflow = num.BitLen() > bits - 1; overflow {
		var mask uint64 = ^uint64(0) >> uint(64 - bits)
		num.And(num, new(big.Int).SetUint64(mask))
	}
	return num.Int64(), truncation, overflow
}

// Return a representation of z truncated to be a uint of length bits.
// Valid values for bits are 0, 8, 16, 32, 64. Result is otherwise undefined
// If a truncation occurs, the decimal part is dropped and the conversion
// continues as usual. truncation will be true
// If an overflow occurs, the result is equivelant to a cast of the form
// uint32(x). overflow will be true.
func (z *BigComplex) Uint(bits int) (_ uint64, truncation, overflow bool) {
	var integer *BigComplex
	integer, truncation = z.Integer()
	num := integer.Num()

	var mask uint64 = ^uint64(0) >> uint(64 - bits)
	if overflow = num.BitLen() > bits; overflow {
		num.And(num, new(big.Int).SetUint64(mask))
	}

	r := num.Uint64()
	if num.Sign() < 0 {
		overflow = true
		r = (^r + 1) & mask
	}
	return r, truncation, overflow
}

// Return a representation of z truncated to a float64
// If a truncation from a complex occurs, the imaginary part is dropped
// and the conversion continues as usual. truncation will be true
// exact will be true if the conversion was completed without loss of precision
func (z *BigComplex) Float64() (f float64, truncation, exact bool) {
	f, exact = z.Rat.Float64()
	return f, exact, !z.IsReal()
}

// Return a complex128 representation of z. 
// exact will be true if the conversion was completed without loss of precision
func (z *BigComplex) Complex128() (_ complex128, exact bool) {
	r, re := z.Rat.Float64()
	i, ie := z.Imag.Float64()
	return complex(r, i), re && ie
}

// Return a representation of z truncated to be a integer value.
// Second return is true if a truncation occured.
func (z *BigComplex) Integer() (_ *BigComplex, truncation bool) {
	if z.IsInteger() {
		return z, false
	} else {
		trunc := new(BigComplex)
		trunc.SetInt(z.Num())
		trunc.Num().Div(trunc.Num(), z.Denom())
		return trunc, true
	}
}

// Return a representation of z truncated to be a real value.
// Second return is true if a truncation occured.
func (z *BigComplex) Real() (_ *BigComplex, truncation bool) {
	if z.IsReal() {
		return z, false
	} else {
		return &BigComplex{Rat: z.Rat}, true
	}
}

func (z *BigComplex) IsInteger() bool {
	return z.Rat.IsInt() && z.Imag.Num().BitLen() == 0
}

func (z *BigComplex) IsReal() bool {
	return z.Imag.Num().BitLen() == 0
}

func (z *BigComplex) Equals(other *BigComplex) bool {
	return new(BigComplex).Sub(z, other).IsZero()
}

func (z *BigComplex) String() string {
	var s string
	if z.Rat.IsInt() {
		s += z.Rat.Num().String()
	} else {
		s += z.Rat.FloatString(5)
	}
	if !z.IsReal() {
		s += "+"
		if z.Imag.IsInt() {
			s += z.Imag.Num().String()
		} else {
			s += z.Imag.FloatString(5)
		}
		s += "i"
	}
	return s
}

