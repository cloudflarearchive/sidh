package sike

import (
	"testing"

	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"
	"strings"

	rand "crypto/rand"
	. "github.com/cloudflare/p751sidh/sidh"
)

const (
	PkB    = "7C55E268665504B9A11A1B30B4363A4957960AD015A7B74DF39FB0141A95CC51A4BEBBB48452EF0C881220D68CB5FF904C0549F05F06BF49A520E684DD610A7E121B420C751B789BDCDB8B6EC136BA0CE74EB6904906057EA7343839EA35FAF2C3D7BE76C81DCA4DF0850CE5F111FF9FF97242EC5520310D7F90A004BACFD75408CBFE8948232A9CCF035136DE3691D9BEF110C3081AADF0D2328CE2CC94998D8AE94D6575083FAFA045F50201FCE841D01C214CC8BBEFCC701484215EA70518204C76A0DA89BEAF0B066F6FD9E78A2C908CF0AFF74E0B55477190F918397F0CF3A537B7911DA846196AD914114A15C2F3C1062D78B19D23348C3D3D4A9C2B2018B382CC44544DA2FA263EB6212D2D13F254216DE002D4AEA55C75C5349A681D7A809BCC29C4CAE1168AC790321FF7429FAAC2FC09465F93E10B9DD970901A1B1D045DDAC9D7B901E00F29AA9F2C87C8EF848E80B7B290ECF85D6BB4C7E975A939A7AFB63069F900A75C9B7B71C2E7472C21A87AB604B6372D4EBEC5974A711281A819636D8FA3E6608F2B81F35599BBB4A1EB5CBD8F743587550F8CE3A809F5C9C399DD52B2D15F217A36F3218C772FD4E67F67D526DEBE1D31FEC4634927A873A1A6CFE55FF1E35AB72EBBD22E3CDD9D2640813345015BB6BD25A6977D0391D4D78998DD178155FEBF247BED3A9F83EAF3346BA90098B908B2359B60491C94330626709D235D1CFB7C87DCA779CFBA23DA280DC06FAEA0FDB3773B0C6391F889D803B7C04AC6AB27375B440336789823176C57"
	PrB    = "00010203040506070809000102030405060708090001020304050607080901028626ED79D451140800E03B59B956F8210E556067407D13DC90FA9E8B872BFB8FAB0A7289852106E40538D3575C500201"
	KAT644 = "../etc/PQCkemKAT_644.rsp"
)

var params = Params(FP_751)

// Fail if err !=nil. Display msg as an error message
func checkErr(t testing.TB, err error, msg string) {
	if err != nil {
		t.Error(msg)
	}
}

// Encrypt, Decrypt, check if input/output plaintext is the same
func TestPKERoundTrip(t *testing.T) {
	// Message to be encrypted
	var msg [32]byte
	for i, _ := range msg {
		msg[i] = byte(i)
	}

	// Import keys
	pkB := NewPublicKey(params.Id, KeyVariant_SIKE)
	skB := NewPrivateKey(params.Id, KeyVariant_SIKE)
	pk_hex, err := hex.DecodeString(PkB)
	if err != nil {
		t.Fatal(err)
	}
	sk_hex, err := hex.DecodeString(PrB)
	if err != nil {
		t.Fatal(err)
	}
	if pkB.Import(pk_hex) != nil || skB.Import(sk_hex) != nil {
		t.Error("Import")
	}

	ct, err := Encrypt(rand.Reader, pkB, msg[:])
	if err != nil {
		t.Fatal(err)
	}
	pt, err := Decrypt(skB, ct)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(pt[:], msg[:]) {
		t.Errorf("Decryption failed \n got : %X\n exp : %X", pt, msg)
	}
}

// Generate key and check if can encrypt
func TestPKEKeyGeneration(t *testing.T) {
	// Message to be encrypted
	var msg [32]byte
	var err error
	for i, _ := range msg {
		msg[i] = byte(i)
	}

	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	err = sk.Generate(rand.Reader)
	checkErr(t, err, "PEK key generation")
	pk := sk.GeneratePublicKey()

	// Try to encrypt
	ct, err := Encrypt(rand.Reader, pk, msg[:])
	checkErr(t, err, "PEK encryption")
	pt, err := Decrypt(sk, ct)
	checkErr(t, err, "PEK key decryption")

	if !bytes.Equal(pt[:], msg[:]) {
		t.Fatalf("Decryption failed \n got : %X\n exp : %X", pt, msg)
	}
}

func TestNegativePKE(t *testing.T) {
	var msg [40]byte
	var err error

	// Generate key
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	err = sk.Generate(rand.Reader)
	checkErr(t, err, "key generation")

	pk := sk.GeneratePublicKey()

	ct, err := Encrypt(rand.Reader, pk, msg[:39])
	if err == nil {
		t.Fatal("Error hasn't been returned")
	}
	if ct != nil {
		t.Fatal("Ciphertext must be nil")
	}

	pt, err := Decrypt(sk, msg[:23])
	if err == nil {
		t.Fatal("Error hasn't been returned")
	}
	if pt != nil {
		t.Fatal("Ciphertext must be nil")
	}
}

func testKEMRoundTrip(pkB, skB []byte) bool {
	// Import keys
	pk := NewPublicKey(params.Id, KeyVariant_SIKE)
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	if pk.Import(pkB) != nil || sk.Import(skB) != nil {
		return false
	}

	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	if err != nil {
		return false
	}

	ss_d, err := Decapsulate(sk, pk, ct)
	if err != nil {
		return false
	}
	return bytes.Equal(ss_e, ss_d)
}

func TestKEMRoundTrip(t *testing.T) {
	pk, err := hex.DecodeString(PkB)
	checkErr(t, err, "public key B not a number")
	sk, err := hex.DecodeString(PrB)
	checkErr(t, err, "private key B not a number")
	if !testKEMRoundTrip(pk, sk) {
		t.Error("kem roundtrip failed")
	}
}

func TestKEMKeyGeneration(t *testing.T) {
	// Generate key
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	// calculated shared secret
	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "encapsulation failed")
	ss_d, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "decapsulation failed")

	if !bytes.Equal(ss_e, ss_d) {
		t.Fatalf("KEM failed \n encapsulated: %X\n decapsulated: %X", ss_d, ss_e)
	}
}

func TestNegativeKEM(t *testing.T) {
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	ct, ss_e, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "pre-requisite for a test failed")

	ct[0] = ct[0] - 1
	ss_d, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "decapsulation returns error when invalid ciphertext provided")

	if bytes.Equal(ss_e, ss_d) {
		// no idea how this could ever happen, but it would be very bad
		t.Error("critical error")
	}

	// Try encapsulating with SIDH key
	pkSidh := NewPublicKey(params.Id, KeyVariant_SIDH_B)
	prSidh := NewPrivateKey(params.Id, KeyVariant_SIDH_B)
	_, _, err = Encapsulate(rand.Reader, pkSidh)
	if err == nil {
		t.Error("encapsulation accepts SIDH public key")
	}
	// Try decapsulating with SIDH key
	_, err = Decapsulate(prSidh, pk, ct)
	if err == nil {
		t.Error("decapsulation accepts SIDH private key key")
	}
}

// In case invalid ciphertext is provided, SIKE's decapsulation must
// return same (but unpredictable) result for a given key.
func TestNegativeKEMSameWrongResult(t *testing.T) {
	sk := NewPrivateKey(params.Id, KeyVariant_SIKE)
	checkErr(t, sk.Generate(rand.Reader), "error: key generation")
	pk := sk.GeneratePublicKey()

	ct, encSs, err := Encapsulate(rand.Reader, pk)
	checkErr(t, err, "pre-requisite for a test failed")

	// make ciphertext wrong
	ct[0] = ct[0] - 1
	decSs1, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "pre-requisite for a test failed")

	// second decapsulation must be done with same, but imported private key
	expSk := sk.Export()

	// creat new private key
	sk = NewPrivateKey(params.Id, KeyVariant_SIKE)
	err = sk.Import(expSk)
	checkErr(t, err, "import failed")

	// try decapsulating again. ss2 must be same as ss1 and different than
	// original plaintext
	decSs2, err := Decapsulate(sk, pk, ct)
	checkErr(t, err, "pre-requisite for a test failed")

	if !bytes.Equal(decSs1, decSs2) {
		t.Error("decapsulation is insecure")
	}

	if bytes.Equal(encSs, decSs1) || bytes.Equal(encSs, decSs2) {
		// this test requires that decapsulation returns wrong result
		t.Errorf("test implementation error")
	}
}

func readAndCheckLine(r *bufio.Reader) []byte {
	// Read next line from buffer
	line, isPrefix, err := r.ReadLine()
	if err != nil || isPrefix {
		panic("Wrong format of input file")
	}

	// Function expects that line is in format "KEY = HEX_VALUE". Get
	// value, which should be a hex string
	hexst := strings.Split(string(line), "=")[1]
	hexst = strings.TrimSpace(hexst)
	// Convert value to byte string
	ret, err := hex.DecodeString(hexst)
	if err != nil {
		panic("Wrong format of input file")
	}
	return ret
}

func testKeygen(pk, sk []byte) bool {
	// Import provided private key
	var prvKey = NewPrivateKey(params.Id, KeyVariant_SIKE)
	if prvKey.Import(sk) != nil {
		panic("sike test: can't load KAT")
	}

	// Generate public key
	pubKey := prvKey.GeneratePublicKey()
	return bytes.Equal(pubKey.Export(), pk)
}

func testDecapsulation(pk, sk, ct, ssExpected []byte) bool {
	var pubKey = NewPublicKey(params.Id, KeyVariant_SIKE)
	var prvKey = NewPrivateKey(params.Id, KeyVariant_SIKE)
	if pubKey.Import(pk) != nil || prvKey.Import(sk) != nil {
		panic("sike test: can't load KAT")
	}

	ssGot, err := Decapsulate(prvKey, pubKey, ct)
	if err != nil {
		panic("sike test: can't perform decapsulation KAT")
	}

	if err != nil {
		return false
	}
	return bytes.Equal(ssGot, ssExpected)
}

func TestSIKE_KAT(t *testing.T) {
	f, err := os.Open(KAT644)
	if err != nil {
		t.Fatal(err)
	}

	r := bufio.NewReader(f)
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil || isPrefix {
			if err == io.EOF {
				break
			} else {
				t.Fatal(err)
			}
		}
		if len(strings.TrimSpace(string(line))) == 0 || line[0] == '#' {
			continue
		}

		// count
		count := strings.Split(string(line), "=")[1]
		// seed
		_ = readAndCheckLine(r)
		// pk
		pk := readAndCheckLine(r)
		// sk (secret key in test vector is concatenation of
		// MSG + SECRET_BOB_KEY + PUBLIC_BOB_KEY. We use only MSG+SECRET_BOB_KEY
		sk := readAndCheckLine(r)
		sk = sk[:params.MsgLen+uint(params.B.SecretByteLen)]
		// ct
		ct := readAndCheckLine(r)
		// ss
		ss := readAndCheckLine(r)

		if !testKeygen(pk, sk) {
			t.Fatalf("KAT keygen form private failed at %s\n", count)
		}

		if !testDecapsulation(pk, sk, ct, ss) {
			t.Fatalf("KAT decapsulation failed at %s\n", count)
		}

		// aditionally test roundtrip with a keypair
		if !testKEMRoundTrip(pk, sk) {
			t.Fatalf("KAT roundtrip failed at %s\n", count)
		}
	}
}
