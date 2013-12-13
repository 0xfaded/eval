package interactive

import (
	"math/big"
)

// Big complex behaves like a *big.Re, but has an imaginary component
// and separate implementation for + - * /
type BigComplex struct {
	Re big.Rat
	Im big.Rat
}

func (z *BigComplex) Add(x, y *BigComplex) *BigComplex {
	z.Re.Add(&x.Re, &y.Re)
	z.Im.Add(&x.Im, &y.Im)
	return z
}

func (z *BigComplex) Sub(x, y *BigComplex) *BigComplex {
	z.Re.Sub(&x.Re, &y.Re)
	z.Im.Sub(&x.Im, &y.Im)
	return z
}

func (z *BigComplex) Mul(x, y *BigComplex) *BigComplex {
	re := new(big.Rat).Mul(&x.Re, &y.Re)
	re.Sub(re, new(big.Rat).Mul(&x.Im, &y.Im))

	im := new(big.Rat).Mul(&x.Re, &y.Im)
	im.Add(im, new(big.Rat).Mul(&x.Im, &y.Re))

	z.Re = *re
	z.Im = *im
	return z
}

func (z *BigComplex) Quo(x, y *BigComplex) *BigComplex {
	// a+bi   ac+bd   bc-ad
	// ---- = ----- + ----- i
	// c+di   cc+dd   cc+dd

	cc := new(big.Rat).Mul(&y.Re, &y.Re)
	dd := new(big.Rat).Mul(&y.Im, &y.Im)
	ccdd := new(big.Rat).Add(cc, dd)

	ac := new(big.Rat).Mul(&x.Re, &y.Re)
	ad := new(big.Rat).Mul(&x.Re, &y.Im)
	bc := new(big.Rat).Mul(&x.Im, &y.Re)
	bd := new(big.Rat).Mul(&x.Im, &y.Im)

	z.Re.Add(ac, bd)
	z.Re.Quo(&z.Re, ccdd)

	z.Im.Sub(bc, ad)
	z.Im.Quo(&z.Im, ccdd)

	return z
}

func (z *BigComplex) IsZero() bool {
	return z.Re.Num().BitLen() == 0 && z.Im.Num().BitLen() == 0
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
	num := integer.Re.Num()

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
	num := integer.Re.Num()

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
	f, exact = z.Re.Float64()
	return f, exact, !z.IsReal()
}

// Return a complex128 representation of z. 
// exact will be true if the conversion was completed without loss of precision
func (z *BigComplex) Complex128() (_ complex128, exact bool) {
	r, re := z.Re.Float64()
	i, ie := z.Im.Float64()
	return complex(r, i), re && ie
}

// Return a representation of z truncated to be a integer value.
// Second return is true if a truncation occured.
func (z *BigComplex) Integer() (_ *BigComplex, truncation bool) {
	if z.IsInteger() {
		return z, false
	} else {
		trunc := new(BigComplex)
		trunc.Re.SetInt(z.Re.Num())
		trunc.Re.Num().Div(trunc.Re.Num(), z.Re.Denom())
		return trunc, true
	}
}

// Return a representation of z truncated to be a real value.
// Second return is true if a truncation occured.
func (z *BigComplex) Real() (_ *BigComplex, truncation bool) {
	if z.IsReal() {
		return z, false
	} else {
		return &BigComplex{Re: z.Re}, true
	}
}

func (z *BigComplex) IsInteger() bool {
	return z.Re.IsInt() && z.Im.Num().BitLen() == 0
}

func (z *BigComplex) IsReal() bool {
	return z.Im.Num().BitLen() == 0
}

func (z *BigComplex) Equals(other *BigComplex) bool {
	return new(BigComplex).Sub(z, other).IsZero()
}

func (z *BigComplex) String() string {
	var s string
	if z.Re.IsInt() {
		s += z.Re.Num().String()
	} else {
		s += z.Re.FloatString(5)
	}
	if !z.IsReal() {
		s += "+"
		if z.Im.IsInt() {
			s += z.Im.Num().String()
		} else {
			s += z.Im.FloatString(5)
		}
		s += "i"
	}
	return s
}

