package p751toolbox

import (
	"bytes"
	"fmt"
)

func printHex(vector []byte) (out string) {
	var buffer bytes.Buffer
	buffer.WriteString("0x")
	len := len(vector)
	for i := len - 1; i >= 0; i-- {
		buffer.WriteString(fmt.Sprintf("%02x", vector[i]))
	}
	buffer.WriteString("\n")
	return buffer.String()
}

func (element Fp751Element) String() string {
	var out [94]byte
	element.toBytesFromMontgomeryForm(out[:])
	return fmt.Sprintf("%s", printHex(out[:]))
}

func (primeElement PrimeFieldElement) String() string {
	return fmt.Sprintf("%s", primeElement.A.String())
}

func (extElement ExtensionFieldElement) String() string {
	var out [188]byte
	extElement.ToBytes(out[:])
	return fmt.Sprintf("A: %sB: %s", printHex(out[:94]), printHex(out[94:]))
}

func (point ProjectivePoint) String() string {
	return fmt.Sprintf("X:\n%sZ:\n%s", point.X.String(), point.Z.String())
}

func (point ProjectivePrimeFieldPoint) String() string {
	return fmt.Sprintf("X:\n%sZ:\n%s", point.X.String(), point.Z.String())
}
