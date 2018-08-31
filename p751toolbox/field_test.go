package p751toolbox

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

var quickCheckScaleFactor = uint8(3)
var quickCheckConfig = &quick.Config{MaxCount: (1 << (12 + quickCheckScaleFactor))}

var cln16prime, _ = new(big.Int).SetString("10354717741769305252977768237866805321427389645549071170116189679054678940682478846502882896561066713624553211618840202385203911976522554393044160468771151816976706840078913334358399730952774926980235086850991501872665651576831", 10)

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

func VartimeEq(x,y *PrimeFieldElement) bool {
	return x.A.vartimeEq(y.A)
}

func (x *PrimeFieldElement) toBigInt() *big.Int {
	// Convert from Montgomery form
	return x.A.toBigIntFromMontgomeryForm()
}

func (x *Fp751Element) toBigIntFromMontgomeryForm() *big.Int {
	// Convert from Montgomery form
	a := Fp751Element{}
	aR := fp751X2{}
	copy(aR[:], x[:])              // = a*R
	fp751MontgomeryReduce(&a, &aR) // = a mod p  in [0,2p)
	fp751StrongReduce(&a)          // = a mod p  in [0,p)
	return radix64ToBigInt(a[:])
}

func TestPrimeFieldElementToBigInt(t *testing.T) {
	// Chosen so that p < xR < 2p
	x := PrimeFieldElement{A: Fp751Element{
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 140737488355328,
	}}
	// Computed using Sage:
	// sage: p = 2^372 * 3^239 - 1
	// sage: R = 2^768
	// sage: from_radix_64 = lambda xs: sum((xi * (2**64)**i for i,xi in enumerate(xs)))
	// sage: xR = from_radix_64([1]*11 + [2^47])
	// sage: assert(p < xR)
	// sage: assert(xR < 2*p)
	// sage: (xR / R) % p
	xBig, _ := new(big.Int).SetString("4469946751055876387821312289373600189787971305258234719850789711074696941114031433609871105823930699680637820852699269802003300352597419024286385747737509380032982821081644521634652750355306547718505685107272222083450567982240", 10)
	if xBig.Cmp(x.toBigInt()) != 0 {
		t.Error("Expected", xBig, "found", x.toBigInt())
	}
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

func (x PrimeFieldElement) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(PrimeFieldElement{A: generateFp751(rand)})
}

func (x ExtensionFieldElement) Generate(rand *rand.Rand, size int) reflect.Value {
	return reflect.ValueOf(ExtensionFieldElement{A: generateFp751(rand), B: generateFp751(rand)})
}

//------------------------------------------------------------------------------
// Extension Field
//------------------------------------------------------------------------------

func TestOneExtensionFieldToBytes(t *testing.T) {
	var x ExtensionFieldElement
	var xBytes [188]byte

	x.One()
	x.ToBytes(xBytes[:])

	if xBytes[0] != 1 {
		t.Error("Expected 1, got", xBytes[0])
	}
	for i := 1; i < 188; i++ {
		if xBytes[i] != 0 {
			t.Error("Expected 0, got", xBytes[0])
		}
	}
}

func TestExtensionFieldElementToBytesRoundTrip(t *testing.T) {
	roundTrips := func(x ExtensionFieldElement) bool {
		var xBytes [188]byte
		var xPrime ExtensionFieldElement
		x.ToBytes(xBytes[:])
		xPrime.FromBytes(xBytes[:])

		return x.VartimeEq(&xPrime)
	}

	if err := quick.Check(roundTrips, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestExtensionFieldElementMulDistributesOverAdd(t *testing.T) {
	mulDistributesOverAdd := func(x, y, z ExtensionFieldElement) bool {
		// Compute t1 = (x+y)*z
		t1 := new(ExtensionFieldElement)
		t1.Add(&x, &y)
		t1.Mul(t1, &z)

		// Compute t2 = x*z + y*z
		t2 := new(ExtensionFieldElement)
		t3 := new(ExtensionFieldElement)
		t2.Mul(&x, &z)
		t3.Mul(&y, &z)
		t2.Add(t2, t3)

		return t1.VartimeEq(t2)
	}

	if err := quick.Check(mulDistributesOverAdd, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestExtensionFieldElementMulIsAssociative(t *testing.T) {
	isAssociative := func(x, y, z ExtensionFieldElement) bool {
		// Compute t1 = (x*y)*z
		t1 := new(ExtensionFieldElement)
		t1.Mul(&x, &y)
		t1.Mul(t1, &z)

		// Compute t2 = (y*z)*x
		t2 := new(ExtensionFieldElement)
		t2.Mul(&y, &z)
		t2.Mul(t2, &x)

		return t1.VartimeEq(t2)
	}

	if err := quick.Check(isAssociative, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestExtensionFieldElementSquareMatchesMul(t *testing.T) {
	sqrMatchesMul := func(x ExtensionFieldElement) bool {
		// Compute t1 = (x*x)
		t1 := new(ExtensionFieldElement)
		t1.Mul(&x, &x)

		// Compute t2 = x^2
		t2 := new(ExtensionFieldElement)
		t2.Square(&x)

		return t1.VartimeEq(t2)
	}

	if err := quick.Check(sqrMatchesMul, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestExtensionFieldElementInv(t *testing.T) {
	inverseIsCorrect := func(x ExtensionFieldElement) bool {
		z := new(ExtensionFieldElement)
		z.Inv(&x)

		// Now z = (1/x), so (z * x) * x == x
		z.Mul(z, &x)
		z.Mul(z, &x)

		return z.VartimeEq(&x)
	}

	// This is more expensive; run fewer tests
	var quickCheckConfig = &quick.Config{MaxCount: (1 << (8 + quickCheckScaleFactor))}
	if err := quick.Check(inverseIsCorrect, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestExtensionFieldElementBatch3Inv(t *testing.T) {
	batchInverseIsCorrect := func(x1, x2, x3 ExtensionFieldElement) bool {
		var x1Inv, x2Inv, x3Inv ExtensionFieldElement
		x1Inv.Inv(&x1)
		x2Inv.Inv(&x2)
		x3Inv.Inv(&x3)

		var y1, y2, y3 ExtensionFieldElement
		ExtensionFieldBatch3Inv(&x1, &x2, &x3, &y1, &y2, &y3)

		return (y1.VartimeEq(&x1Inv) && y2.VartimeEq(&x2Inv) && y3.VartimeEq(&x3Inv))
	}

	// This is more expensive; run fewer tests
	var quickCheckConfig = &quick.Config{MaxCount: (1 << (5 + quickCheckScaleFactor))}
	if err := quick.Check(batchInverseIsCorrect, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

//------------------------------------------------------------------------------
// Prime Field
//------------------------------------------------------------------------------
func TestPrimeFieldElementMulVersusBigInt(t *testing.T) {
	mulMatchesBigInt := func(x, y PrimeFieldElement) bool {
		z := new(PrimeFieldElement)
		z.Mul(&x, &y)

		check := new(big.Int)
		check.Mul(x.toBigInt(), y.toBigInt())
		check.Mod(check, cln16prime)

		return check.Cmp(z.toBigInt()) == 0
	}

	if err := quick.Check(mulMatchesBigInt, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

func TestPrimeFieldElementP34VersusBigInt(t *testing.T) {
	var p34, _ = new(big.Int).SetString("2588679435442326313244442059466701330356847411387267792529047419763669735170619711625720724140266678406138302904710050596300977994130638598261040117192787954244176710019728333589599932738193731745058771712747875468166412894207", 10)
	p34MatchesBigInt := func(x PrimeFieldElement) bool {
		z := new(PrimeFieldElement)
		z.P34(&x)

		check := x.toBigInt()
		check.Exp(check, p34, cln16prime)

		return check.Cmp(z.toBigInt()) == 0
	}

	// This is more expensive; run fewer tests
	var quickCheckConfig = &quick.Config{MaxCount: (1 << (8 + quickCheckScaleFactor))}
	if err := quick.Check(p34MatchesBigInt, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Package-level storage for this field element is intended to deter
// compiler optimizations.
var benchmarkFp751Element Fp751Element
var benchmarkFp751X2 fp751X2
var bench_x = Fp751Element{17026702066521327207, 5108203422050077993, 10225396685796065916, 11153620995215874678, 6531160855165088358, 15302925148404145445, 1248821577836769963, 9789766903037985294, 7493111552032041328, 10838999828319306046, 18103257655515297935, 27403304611634}
var bench_y = Fp751Element{4227467157325093378, 10699492810770426363, 13500940151395637365, 12966403950118934952, 16517692605450415877, 13647111148905630666, 14223628886152717087, 7167843152346903316, 15855377759596736571, 4300673881383687338, 6635288001920617779, 30486099554235}
var bench_z = fp751X2{1595347748594595712, 10854920567160033970, 16877102267020034574, 12435724995376660096, 3757940912203224231, 8251999420280413600, 3648859773438820227, 17622716832674727914, 11029567000887241528, 11216190007549447055, 17606662790980286987, 4720707159513626555, 12887743598335030915, 14954645239176589309, 14178817688915225254, 1191346797768989683, 12629157932334713723, 6348851952904485603, 16444232588597434895, 7809979927681678066, 14642637672942531613, 3092657597757640067, 10160361564485285723, 240071237}

func BenchmarkExtensionFieldElementMul(b *testing.B) {
	z := &ExtensionFieldElement{A: bench_x, B: bench_y}
	w := new(ExtensionFieldElement)

	for n := 0; n < b.N; n++ {
		w.Mul(z, z)
	}
}

func BenchmarkExtensionFieldElementInv(b *testing.B) {
	z := &ExtensionFieldElement{A: bench_x, B: bench_y}
	w := new(ExtensionFieldElement)

	for n := 0; n < b.N; n++ {
		w.Inv(z)
	}
}

func BenchmarkExtensionFieldElementSquare(b *testing.B) {
	z := &ExtensionFieldElement{A: bench_x, B: bench_y}
	w := new(ExtensionFieldElement)

	for n := 0; n < b.N; n++ {
		w.Square(z)
	}
}

func BenchmarkExtensionFieldElementAdd(b *testing.B) {
	z := &ExtensionFieldElement{A: bench_x, B: bench_y}
	w := new(ExtensionFieldElement)

	for n := 0; n < b.N; n++ {
		w.Add(z, z)
	}
}

func BenchmarkExtensionFieldElementSub(b *testing.B) {
	z := &ExtensionFieldElement{A: bench_x, B: bench_y}
	w := new(ExtensionFieldElement)

	for n := 0; n < b.N; n++ {
		w.Sub(z, z)
	}
}

func BenchmarkPrimeFieldElementMul(b *testing.B) {
	z := &PrimeFieldElement{A: bench_x}
	w := new(PrimeFieldElement)

	for n := 0; n < b.N; n++ {
		w.Mul(z, z)
	}
}

// --- field operation functions

func BenchmarkFp751Multiply(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fp751Mul(&benchmarkFp751X2, &bench_x, &bench_y)
	}
}

func BenchmarkFp751MontgomeryReduce(b *testing.B) {
	z := bench_z

	// This benchmark actually computes garbage, because
	// fp751MontgomeryReduce mangles its input, but since it's
	// constant-time that shouldn't matter for the benchmarks.
	for n := 0; n < b.N; n++ {
		fp751MontgomeryReduce(&benchmarkFp751Element, &z)
	}
}

func BenchmarkFp751AddReduced(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fp751AddReduced(&benchmarkFp751Element, &bench_x, &bench_y)
	}
}

func BenchmarkFp751SubReduced(b *testing.B) {
	for n := 0; n < b.N; n++ {
		fp751SubReduced(&benchmarkFp751Element, &bench_x, &bench_y)
	}
}

func BenchmarkFp751ConditionalSwap(b *testing.B) {
	x, y := bench_x, bench_y
	for n := 0; n < b.N; n++ {
		fp751ConditionalSwap(&x, &y, 1)
		fp751ConditionalSwap(&x, &y, 0)
	}
}

func BenchmarkFp751StrongReduce(b *testing.B) {
	x := bench_x
	for n := 0; n < b.N; n++ {
		fp751StrongReduce(&x)
	}
}

func BenchmarkFp751AddLazy(b *testing.B) {
	var z Fp751Element
	x, y := bench_x, bench_y
	for n := 0; n < b.N; n++ {
		fp751AddLazy(&z, &x, &y)
	}
}

func BenchmarkFp751X2AddLazy(b *testing.B) {
	x, y, z := bench_z, bench_z, bench_z
	for n := 0; n < b.N; n++ {
		fp751X2AddLazy(&x, &y, &z)
	}
}

func BenchmarkFp751X2SubLazy(b *testing.B) {
	x, y, z := bench_z, bench_z, bench_z
	for n := 0; n < b.N; n++ {
		fp751X2SubLazy(&x, &y, &z)
	}
}
