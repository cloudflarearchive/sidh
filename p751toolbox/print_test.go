package p751toolbox

// Tools used for testing and debugging

import (
	"fmt"
)

func (primeElement PrimeFieldElement) String() string {
   return fmt.Sprintf("%X", primeElement.A.toBigInt().String())
}

func (extElement ExtensionFieldElement) String() string {
	return fmt.Sprintf("\nA: %X\nB: %X", extElement.A.toBigInt().String(), extElement.B.toBigInt().String())
}

func (point ProjectivePoint) String() string {
	return fmt.Sprintf("X:\n%sZ:\n%s", point.X.String(), point.Z.String())
}
