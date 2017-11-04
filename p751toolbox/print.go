package p751toolbox

import (
	"reflect"
	"fmt"
	"bytes"
)

func PrintHex(obj interface{}) (out string) {
	var buffer bytes.Buffer
	vector := reflect.ValueOf(obj)
	switch vector.Kind() {
		case reflect.Array, reflect.Slice:
			buffer.WriteString("0x")
			format := ""
			switch vector.Index(0).Kind() {
				case reflect.Uint8:
					format = "%02x"
				case reflect.Uint64:
					format = "%016x"
			}
			len := vector.Len()
			for i:=len-1; i>=0; i-- {
				word := vector.Index(i)
				buffer.WriteString(fmt.Sprintf(format,word))
			}
			buffer.WriteString("\n")
	}
	return buffer.String()
}

func (element Fp751Element) String() string {
	var out [94]byte
	element.toBytesFromMontgomeryForm(out[:])
	return fmt.Sprintf("%s",PrintHex(out[:]))
}

func (primeElement PrimeFieldElement) String()  string {
	return fmt.Sprintf("%s",primeElement.A.String())
}

func (extElement ExtensionFieldElement) String() string {
	var out [188]byte
	extElement.ToBytes(out[:])
	return fmt.Sprintf("A: %sB: %s", PrintHex(out[:94]), PrintHex(out[94:]))
}

func (point ProjectivePoint) String() string {
	return fmt.Sprintf("X:\n%sZ:\n%s", point.X.String(), point.Z.String())
}

func (point ProjectivePrimeFieldPoint) String() string {
	return fmt.Sprintf("X:\n%sZ:\n%s", point.X.String(), point.Z.String())
}

