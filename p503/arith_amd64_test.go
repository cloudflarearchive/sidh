// +build amd64,!noasm

package p503

import (
	. "github.com/cloudflare/sidh/internal/isogeny"
	cpu "github.com/cloudflare/sidh/internal/utils"
	"reflect"
	"testing"
	"testing/quick"
)

// Utility function used for testing Mul implementations. Tests caller provided
// mulFunc against mul()
func testMul(t *testing.T, mulFunc func(z *FpElementX2, x, y *FpElement)) {
	doMulTest := func(multiplier, multiplicant FpElement) bool {
		var resMulRef, resMulOptim FpElementX2
		mul(&resMulRef, &multiplier, &multiplicant)
		mulFunc(&resMulOptim, &multiplier, &multiplicant)
		return reflect.DeepEqual(resMulRef, resMulOptim)
	}

	if err := quick.Check(doMulTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Utility function used for testing REDC implementations. Tests caller provided
// redcFunc against redc()
func testRedc(t *testing.T, redcFunc func(z *FpElement, x *FpElementX2)) {
	doRedcTest := func(aRR FpElementX2) bool {
		var resRedcRef, resRedcOptim FpElement
		var aRRcpy = aRR
		redc(&resRedcRef, &aRR)
		redcFunc(&resRedcOptim, &aRRcpy)
		return reflect.DeepEqual(resRedcRef, resRedcOptim)
	}

	if err := quick.Check(doRedcTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Ensures corretness of implementation of mul operation which uses MULX
func TestMulWithMULX(t *testing.T) {
	if !cpu.HasBMI2 {
		t.Skip("MULX not supported by the platform")
	}
	testMul(t, mulWithMULX)
}

// Ensures corretness of implementation of mul operation which uses MULX and ADOX/ADCX
func TestMulWithMULXADX(t *testing.T) {
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testMul(t, mulWithMULXADX)
}

// Ensures corretness of Montgomery reduction implementation which uses MULX
func TestRedcWithMULX(t *testing.T) {
	if !cpu.HasBMI2 {
		t.Skip("MULX not supported by the platform")
	}
	testRedc(t, redcWithMULX)
}

// Ensures corretness of Montgomery reduction implementation which uses MULX
func TestRedcWithMULXADX(t *testing.T) {
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testRedc(t, redcWithMULXADX)
}
