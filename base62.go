// Package base62 implements base62 encoding
package base62

import (
	"math"
	"math/big"
	"strings"
)

const base = 62

type Encoding struct {
	encode string
}

const encodeStd = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// NewEncoding returns a new Encoding defined by the given alphabet
func NewEncoding(encoder string) *Encoding {
	return &Encoding{
		encode: encoder,
	}
}

// StdEncoding is the standard base62 encoding
var StdEncoding = NewEncoding(encodeStd)

// EncodeInt64 returns the base62 encoding of n using the StdEncoding
func EncodeInt64(n int64) string {
	return StdEncoding.EncodeInt64(n)
}

// EncodeInt64 returns the base62 encoding of n
func (e *Encoding) EncodeInt64(n int64) string {
	var (
		s   = make([]byte, 0)
		rem int64
	)

	// Progressively divide by base, store remainder each time
	// Prepend as an additional character is the higher power
	for n > 0 {
		rem = n % base
		n = n / base
		s = append([]byte{e.encode[rem]}, s...)
	}

	return string(s)
}

// EncodeBigInt returns the base62 encoding of an arbitrary precision integer using the StdEncoding
func EncodeBigInt(n *big.Int) string {
	return StdEncoding.EncodeBigInt(n)
}

// EncodeBigInt returns the base62 encoding of an arbitrary precision integer
func (e *Encoding) EncodeBigInt(n *big.Int) string {
	var (
		s    = make([]byte, 0)
		rem  = new(big.Int)
		bse  = new(big.Int)
		zero = new(big.Int)
	)
	bse.SetInt64(base)
	zero.SetInt64(0)

	// Progressively divide by base, until we hit zero
	// store remainder each time
	// Prepend as an additional character is the higher power
	for n.Cmp(zero) == 1 {
		n, rem = n.DivMod(n, bse, rem)
		s = append([]byte{e.encode[rem.Int64()]}, s...)
	}

	return string(s)
}

// DecodeToInt64 decodes a base62 encoded string using the StdEncoding
func DecodeToInt64(s string) int64 {
	return StdEncoding.DecodeToInt64(s)
}

// DecodeToInt64 decodes a base62 encoded string
func (e *Encoding) DecodeToInt64(s string) int64 {
	var (
		n     int64
		c     int64
		idx   int
		power int
	)

	for i, v := range s {
		idx = strings.IndexRune(e.encode, v)

		// Work downwards through powers of our base
		power = len(s) - (i + 1)

		// Calculate value at this position and add
		c = int64(idx) * int64(math.Pow(float64(base), float64(power)))
		n = n + c
	}

	return int64(n)
}

// DecodeToBigInt returns an arbitrary precision integer from the base62
// encoded string using the StdEncoding
func DecodeToBigInt(s string) *big.Int {
	return StdEncoding.DecodeToBigInt(s)
}

// DecodeToBigInt returns an arbitrary precision integer from the base62 encoded string
func (e *Encoding) DecodeToBigInt(s string) *big.Int {
	var (
		n = new(big.Int)

		c     = new(big.Int)
		idx   = new(big.Int)
		power = new(big.Int)
		exp   = new(big.Int)
		bse   = new(big.Int)
	)
	bse.SetInt64(base)

	// Run through each character to decode
	for i, v := range s {
		// Get index/position of the rune as a big int
		idx.SetInt64(int64(strings.IndexRune(e.encode, v)))

		// Work downwards through exponents
		exp.SetInt64(int64(len(s) - (i + 1)))

		// Calculate power for this exponent
		power.Exp(bse, exp, nil)

		// Multiplied by our index, gives us the value for this character
		c = c.Mul(idx, power)

		// Finally add to running total
		n.Add(n, c)
	}

	return n
}
