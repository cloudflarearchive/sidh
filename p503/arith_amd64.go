// +build amd64,!noasm

package p503

import cpu "github.com/cloudflare/p751sidh/internal/utils"

// There couple of reasons for having those variables here:
// * to have an access to them from assembly
//   TODO(kk): Is there a way to access variable from different package?
//             If it is then probably this file could be moved to internal
//             and we don't need to have many copies of that
// * make it easy to vendor the library
// * make it possible to test all functionalities
var useMULX bool
var useADXMULX bool

func recognizecpu() {
	useMULX = cpu.HasBMI2
	useADXMULX = cpu.HasADX && cpu.HasBMI2
}

func init() {
	recognizecpu()
}
