// +build amd64,!noasm

package p503

import (
	. "github.com/cloudflare/sidh/internal/isogeny"
	cpu "github.com/cloudflare/sidh/internal/utils"
	"reflect"
	"testing"
	"testing/quick"
)

type OptimFlag uint

const (
	kUse_MUL     OptimFlag = 1 << 0
	kUse_MULX              = 1 << 1
	kUse_MULXADX           = 1 << 2
)

// Utility function used for testing Mul implementations. Tests caller provided
// mulFunc against mul()
func testMul(t *testing.T, f1, f2 OptimFlag) {
	doMulTest := func(multiplier, multiplicant FpElement) bool {
		defer recognizecpu()
		var resMulRef, resMulOptim FpElementX2

		// Compute multiplier*multiplicant with first implementation
		useMULX = (kUse_MULX & f1) == kUse_MULX
		useADXMULX = (kUse_MULXADX & f1) == kUse_MULXADX
		fp503Mul(&resMulOptim, &multiplier, &multiplicant)

		// Compute multiplier*multiplicant with second implementation
		useMULX = (kUse_MULX & f2) == kUse_MULX
		useADXMULX = (kUse_MULXADX & f2) == kUse_MULXADX
		fp503Mul(&resMulRef, &multiplier, &multiplicant)

		// Compare results
		return reflect.DeepEqual(resMulRef, resMulOptim)
	}

	if err := quick.Check(doMulTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Utility function used for testing REDC implementations. Tests caller provided
// redcFunc against redc()
func testRedc(t *testing.T, f1, f2 OptimFlag) {
	doRedcTest := func(aRR FpElementX2) bool {
		defer recognizecpu()
		var resRedcF1, resRedcF2 FpElement
		var aRRcpy = aRR

		// Compute redc with first implementation
		useMULX = (kUse_MULX & f1) == kUse_MULX
		useADXMULX = (kUse_MULXADX & f1) == kUse_MULXADX
		fp503MontgomeryReduce(&resRedcF1, &aRR)

		// Compute redc with second implementation
		useMULX = (kUse_MULX & f2) == kUse_MULX
		useADXMULX = (kUse_MULXADX & f2) == kUse_MULXADX
		fp503MontgomeryReduce(&resRedcF2, &aRRcpy)

		// Compare results
		return reflect.DeepEqual(resRedcF2, resRedcF1)
	}

	if err := quick.Check(doRedcTest, quickCheckConfig); err != nil {
		t.Error(err)
	}
}

// Ensures corretness of implementation of mul operation which uses MULX
func TestMulWithMULX(t *testing.T) {
	defer recognizecpu()
	if !cpu.HasBMI2 {
		t.Skip("MULX not supported by the platform")
	}
	testMul(t, kUse_MULX, kUse_MUL)
}

// Ensures corretness of implementation of mul operation which uses MULX and ADOX/ADCX
func TestMulWithMULXADX(t *testing.T) {
	defer recognizecpu()
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testMul(t, kUse_MULXADX, kUse_MUL)
}

// Ensures corretness of implementation of mul operation which uses MULX and ADOX/ADCX
func TestMulWithMULXADXAgainstMULX(t *testing.T) {
	defer recognizecpu()
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testMul(t, kUse_MULX, kUse_MULXADX)
}

// Ensures corretness of Montgomery reduction implementation which uses MULX
func TestRedcWithMULX(t *testing.T) {
	defer recognizecpu()
	if !cpu.HasBMI2 {
		t.Skip("MULX not supported by the platform")
	}
	testRedc(t, kUse_MULX, kUse_MUL)
}

// Ensures corretness of Montgomery reduction implementation which uses MULX
// and ADX
func TestRedcWithMULXADX(t *testing.T) {
	defer recognizecpu()
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testRedc(t, kUse_MULXADX, kUse_MUL)
}

// Ensures corretness of Montgomery reduction implementation which uses MULX
// and ADX.
func TestRedcWithMULXADXAgainstMULX(t *testing.T) {
	defer recognizecpu()
	if !(cpu.HasADX && cpu.HasBMI2) {
		t.Skip("MULX, ADCX and ADOX not supported by the platform")
	}
	testRedc(t, kUse_MULXADX, kUse_MULX)
}
