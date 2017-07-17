package cln16sidh

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

var quickCheckConfig = &quick.Config{MaxCount: (1 << 16)}

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

func (x *Fp751X2) toBigInt() *big.Int {
	return radix64ToBigInt(x[:])
}

func generateFp751(rand *rand.Rand) Fp751Element {
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

	return Fp751Element{
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
	}
}

func (x Fp751Element) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(generateFp751(rand))
}

func (x FieldElement) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(FieldElement{A: generateFp751(rand), B: generateFp751(rand)})
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

func TestFieldElementMulIsAssociative(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	is_associative := func(x, y, z FieldElement) bool {
		// Compute t1 = (x*y)*z
		t1 := new(FieldElement)
		t1.Mul(&x, &y)
		t1.Mul(t1, &z)

		// Compute t2 = (y*z)*x
		t2 := new(FieldElement)
		t2.Mul(&y, &z)
		t2.Mul(t2, &x)

		a1 := t1.A.toBigInt()
		a1.Mod(a1, p)
		a2 := t2.A.toBigInt()
		a2.Mod(a2, p)
		b1 := t1.B.toBigInt()
		b1.Mod(b1, p)
		b2 := t2.B.toBigInt()
		b2.Mod(b2, p)

		return (a1.Cmp(a2) == 0) && (b1.Cmp(b2) == 0)
	}

	if err := quick.Check(is_associative, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestFp751StrongReduceVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	reductionIsCorrect := func(x Fp751Element) bool {
		xOrig := x.toBigInt()
		xOrig.Mod(xOrig, p)

		Fp751StrongReduce(&x)

		xRed := x.toBigInt()

		return xRed.Cmp(xOrig) == 0
	}

	if err := quick.Check(reductionIsCorrect, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestFp751AddReducedVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	// Returns true if computing x + y in this implementation matches
	// computing x + y using big.Int
	assertion := func(x, y Fp751Element) bool {
		z := new(Fp751Element)
		// Compute z = x + y using Fp751AddReduced
		Fp751AddReduced(z, &x, &y)

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

	if err := quick.Check(assertion, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestFp751SubReducedVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))

	// Returns true if computing x - y in this implementation matches
	// computing x - y using big.Int
	assertion := func(x, y Fp751Element) bool {
		z := new(Fp751Element)
		// Compute z = x - y using Fp751SubReduced
		Fp751SubReduced(z, &x, &y)

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

	if err := quick.Check(assertion, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestFp751MulReduceVersusBigInt(t *testing.T) {
	// The CLN16-SIDH prime
	p := new(big.Int)
	p.UnmarshalText(([]byte)("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831"))
	// The inverse of the Montgomery constant 1/(2^768) (mod p)
	Rprime := new(big.Int)
	Rprime.UnmarshalText(([]byte)("1518725603824737389819053798918035007730761831988339415201705393383230549778130460683426736321139507281132743779724182711261647601379839529950740166978024981554536824910527087746254630334120179536606671225338947864384536159295"))

	// Returns true if computing x * y in this implementation matches
	// computing x * y using big.Int
	assertion := func(x, y Fp751Element) bool {
		// Compute z = x * y using Fp751Mul
		z := new(Fp751X2)
		Fp751Mul(z, &x, &y)
		zReduced := new(Fp751Element)
		Fp751MontgomeryReduce(zReduced, z)
		zBig := zReduced.toBigInt()

		// Compute z = x * y * Rprime using big.Int
		tmp := new(big.Int)
		xBig := x.toBigInt()
		yBig := y.toBigInt()
		tmp.Mul(xBig, yBig)
		tmp.Mul(tmp, Rprime)

		// Reduce both mod p and check that they are equal.
		zBig.Mod(zBig, p)
		tmp.Mod(tmp, p)
		return zBig.Cmp(tmp) == 0
	}

	if err := quick.Check(assertion, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Package-level storage for this field element is intended to deter
// compiler optimizations.
var benchmarkFp751Element Fp751Element
var benchmarkFp751X2 Fp751X2

func BenchmarkFp751Multiply(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	for n := 0; n < b.N; n++ {
		Fp751Mul(&benchmarkFp751X2, &x, &y)
	}
}

func BenchmarkFp751MontgomeryReduce(b *testing.B) {
	z := Fp751X2{1595347748594595712, 10854920567160033970, 16877102267020034574, 12435724995376660096, 3757940912203224231, 8251999420280413600, 3648859773438820227, 17622716832674727914, 11029567000887241528, 11216190007549447055, 17606662790980286987, 4720707159513626555, 12887743598335030915, 14954645239176589309, 14178817688915225254, 1191346797768989683, 12629157932334713723, 6348851952904485603, 16444232588597434895, 7809979927681678066, 14642637672942531613, 3092657597757640067, 10160361564485285723, 240071237}

	// This benchmark actually computes garbage, because Fp751 mangles its
	// input, but since it's constant-time that shouldn't matter for the
	// benchmarks.
	for n := 0; n < b.N; n++ {
		Fp751MontgomeryReduce(&benchmarkFp751Element, &z)
	}
}

func BenchmarkFp751AddReduced(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	for n := 0; n < b.N; n++ {
		Fp751AddReduced(&benchmarkFp751Element, &x, &y)
	}
}

func BenchmarkFp751SubReduced(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	for n := 0; n < b.N; n++ {
		Fp751SubReduced(&benchmarkFp751Element, &x, &y)
	}
}

func BenchmarkFieldElementMultiply(b *testing.B) {
	x := Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}

	y := Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}

	z := &FieldElement{A: x, B: y}
	w := new(FieldElement)

	for n := 0; n < b.N; n++ {
		w.Mul(z, z)
	}
}
