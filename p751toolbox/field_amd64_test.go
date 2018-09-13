// +build amd64,!noasm

package p751toolbox

import (
	"golang.org/x/sys/cpu"
	"testing"
	"testing/quick"
)

func TestFp751MontgomeryReduce(t *testing.T) {
	// First make sure that at least one value with a known result reduces
	// correctly as defined in TestPrimeFieldElementToBigInt.
	fp751MontgomeryReduce = fp751MontgomeryReduceFallback
	t.Run("PrimeFieldElementToBigInt", TestPrimeFieldElementToBigInt)

	if !cpu.X86.HasBMI2 {
		return
	}

	fp751MontgomeryReduce = fp751MontgomeryReduceBMI2
	t.Run("PrimeFieldElementToBigInt", TestPrimeFieldElementToBigInt)

	// Also check that the BMI2 implementation produces the same results
	// as the fallback implementation.
	compareMontgomeryReduce := func(x, y PrimeFieldElement) bool {
		var z, zbackup fp751X2
		var zred1, zred2 Fp751Element

		fp751Mul(&z, &x.A, &y.A)
		zbackup = z

		fp751MontgomeryReduceFallback(&zred1, &z)
		// z may be destroyed.
		z = zbackup
		fp751MontgomeryReduceBMI2(&zred2, &z)

		return zred1 == zred2
	}

	if err := quick.Check(compareMontgomeryReduce, quickCheckConfig); err != nil {
		t.Error(err)
	}

	if !cpu.X86.HasADX {
		return
	}

	fp751MontgomeryReduce = fp751MontgomeryReduceBMI2ADX
	t.Run("PrimeFieldElementToBigInt", TestPrimeFieldElementToBigInt)

	// Check that the BMI2ADX implementation produces the same results as
	// the BMI2 implementation. By transitivity, it should also produce the
	// same results as the fallback implementation.
	compareMontgomeryReduce = func(x, y PrimeFieldElement) bool {
		var z, zbackup fp751X2
		var zred1, zred2 Fp751Element

		fp751Mul(&z, &x.A, &y.A)
		zbackup = z

		fp751MontgomeryReduceBMI2(&zred1, &z)
		// z may be destroyed.
		z = zbackup
		fp751MontgomeryReduceBMI2ADX(&zred2, &z)

		return zred1 == zred2
	}

	if err := quick.Check(compareMontgomeryReduce, quickCheckConfig); err != nil {
		t.Error(err)
	}
}
