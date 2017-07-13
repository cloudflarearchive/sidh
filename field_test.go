package cln16sidh

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// Convert an Fp751Element to a big.Int for testing.  Because this is only
// for testing, no big.Int to Fp751Element conversion is provided.

func radix64ToBigInt(x []uint64) *big.Int {
	radix := new(big.Int)
	// 2^64
	radix.UnmarshalText(([]byte)("18446744073709551616"))

	base := new(big.Int).SetUint64(1)
	val := new(big.Int).SetUint64(0)
	tmp := new(big.Int)

	for _, xi := range x {
		tmp.SetUint64(xi)
		tmp.Mul(tmp, base)
		val.Add(val, tmp)
		base.Mul(base, radix)
	}

	return val
}

func (x *Fp751Element) toBigInt() *big.Int {
	return radix64ToBigInt(x[:])
}

func (x Fp751Element) Generate(rand *rand.Rand, size int) reflect.Value {
	// Generation strategy: low limbs taken from [0,2^64); high limb
	// taken from smaller range
	//
	// Size hint is ignored since all elements are fixed size.
	//
	// Field elements taken in range [0,2p).  Emulate this by capping
	// the high limb by the top digit of 2*p-1:
	//
	// sage: (2*p-1).digits(2^64)[-1]
	// 246065832128056
	//
	// This still allows generating values >= 2p, but hopefully that
	// excess is OK (and if it's not, we'll find out, because it's for
	// testing...)
	//
	highLimb := rand.Uint64() % 246065832128056

	return reflect.ValueOf(
		Fp751Element{
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			rand.Uint64(),
			highLimb,
		})
}

func TestFp751ElementToBigInt(t *testing.T) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}
	// Generated using:
	// from_radix_64 = lambda xs: sum((xi * (2**64)**i for i,xi in enumerate(xs)))
	xValueFromPython := new(big.Int)
	xValueFromPython.UnmarshalText(([]byte)("2306321702539636385951781662938907950351672982559361883294272447056580097879230968214879481380598127719197355217101680007303754413021721583686358737065759814651127760172383109266658276160240175619638712486220677169799330645607"))

	xValue := x.toBigInt()

	if xValue.Cmp(xValueFromPython) != 0 {
		t.Error("Expected", xValueFromPython, "found", xValue)
	}
}

func TestFp751AddVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	// Returns true if computing x + y in this implementation matches
	// computing x + y using big.Int
	assertion := func(x, y Fp751Element) bool {
		z := new(Fp751Element)
		// Compute z = x + y using Fp751Add
		Fp751Add(z, &x, &y)

		xBig := x.toBigInt()
		yBig := y.toBigInt()
		zBig := z.toBigInt()

		// Compute z = x + y using big.Int
		tmp := new(big.Int)
		tmp.Add(xBig, yBig)

		// Reduce both mod p and check that they are equal.
		zBig.Mod(zBig, p)
		tmp.Mod(tmp, p)
		return zBig.Cmp(tmp) == 0
	}

	// Run 1M tests
	config := &quick.Config{MaxCount: (1 << 20)}
	if err := quick.Check(assertion, config); err != nil {
		t.Error(err)
	}
}

func TestFp751SubVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	// Returns true if computing x - y in this implementation matches
	// computing x - y using big.Int
	assertion := func(x, y Fp751Element) bool {
		z := new(Fp751Element)
		// Compute z = x - y using Fp751Sub
		Fp751Sub(z, &x, &y)

		xBig := x.toBigInt()
		yBig := y.toBigInt()
		zBig := z.toBigInt()

		// Compute z = x - y using big.Int
		tmp := new(big.Int)
		tmp.Sub(xBig, yBig)

		// Reduce both mod p and check that they are equal.
		zBig.Mod(zBig, p)
		tmp.Mod(tmp, p)
		return zBig.Cmp(tmp) == 0
	}

	// Run 1M tests
	config := &quick.Config{MaxCount: (1 << 20)}
	if err := quick.Check(assertion, config); err != nil {
		t.Error(err)
	}
}

// Package-level storage for this field element is intended to deter
// compiler optimizations.
var benchmarkFp751Element Fp751Element

func BenchmarkFp751Add(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	for n := 0; n < b.N; n++ {
		Fp751Add(&benchmarkFp751Element, &x, &y)
	}
}

func BenchmarkFp751Sub(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	for n := 0; n < b.N; n++ {
		Fp751Sub(&benchmarkFp751Element, &x, &y)
	}
}
