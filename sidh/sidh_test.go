package sidh

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"testing"

	. "github.com/cloudflare/p751sidh/internal/isogeny"
)

const (
	// PrA - Alice's Private Key: 2*randint(0,2^371)
	PrA = "C09957CC83045FB4C3726384D784476ACB6FFD92E5B15B3C2D451BA063F1BD4CED8FBCF682A98DD0954D3" +
		"7BCAF730E00"
	// PrB - Bob's Private Key: 3*randint(0,3^238)
	PrB = "393E8510E78A16D2DC1AACA9C9D17E7E78DB630881D8599C7040D05BB5557ECAE8165C45D5366ECB37B00" +
		"969740AF201"
	PkA = "74D8EF08CB74EC99BF08B6FBE4FB3D048873B67F018E44988B9D70C564D058401D20E093C7DF0C66F022C" +
		"823E5139D2EA0EE137804B4820E950B046A90B0597759A0B6A197C56270128EA089FA1A2007DDE3430B37" +
		"A3E6350BD47B7F513863741C125FA63DEDAFC475C13DB59E533055B7CBE4B2F32672DF2DF97E03E29617B" +
		"0E9B6A35B58ABB26527A721142701EB147C7050E1D9125DA577B08CD51C8BB50627B8B47FACFC9C7C07DD" +
		"00DD75115DD83719FD5F96115DED23ECAA50B1044C6BF3F27442DA284BA4A272D850F414FB185801BF2EF" +
		"7E628EDB5643E35694B992CF30A2C5120CAF9434F09ACFCA3645B3FFC3A308901FAC7B8955FD5C98576AE" +
		"FD03F5806CB7430F75B3431B75BEC080596ABCA26E637E6E8D4C25175A8C052C9CBE77900A863F83FAB00" +
		"95B32D9C3858EF8A35B9F163D429E71DBA47539EB4791D117FE39DDE94EA7801A42DB12D84DE4740ACF51" +
		"CD7C32BB854569D7D94E11E69D9663CC7ED02E78CF48F4069DF3D3E86198B307095C6B11D46C0DC849F9D" +
		"94C7693209E5B3848AFAA6DA6A8D73362D779CBC43515902ED2BCE3A748C537DE2FCF092FD3E91B790AF5" +
		"4E1092C5E5B89BE5BE23B955A52F769D97277EF69F820109042F28C316AC90AE69EB374C9280300B816E6" +
		"2494B2E01072D1CA96E4B284D2BE1368D6969744B614FACBC8C165864E26E33481D4FDC47B6E523954A25" +
		"C1A096A37CD23FB81AE64FB11BD0A439609F1CE40673B06DD96F698A910E935219D840F3D411EDFB00D98" +
		"065AB9868C32D3DA05FF415"
	PkB = "F6C260C4141E418457CB442E11F0F5558375437576E55D211D19EF83E2839E51D07A82765D8E7B6366FA7" +
		"0B56CDE3AD3B629ACF542A433369496EDA51EDFBE16EFA1B8DEE1CE46B37820ECBD0CD674AACD4F21FABC" +
		"2436651E3AF604356FF3EB2CA87976890E34A56FAEC9A2ACD9559B1BB67B69AC1A521342E1E787DA5D709" +
		"32B0F5842ECA1C99B269DB6C2ED8397F0FC49F114CF8B5AF327A698C0251575CDD1D67732668109A91A3B" +
		"FA5B47D413C7FAB8817FCBEBFE9BDD9C0B1F3B1934A7028A65233E8B58A92E7E9F66B68B2057ECBF7E44A" +
		"0EF6EFCC3C8AA5414E100FA0C24F7545324AD17062FC11377A2A4749DEE27E192460E099DBDA8E840EA11" +
		"AD9D5C83DF065AF77030E7FE18CE24CFC71D356B9B9601811B93676C12CB6B41747133D5259E7A20CC065" +
		"FAB99DF944FDB34ABB9A374F9E9CC8F9C186BD2181DC2771F69C02629C3E4801A7E7C21F6F3CFF7D257E2" +
		"257C88C015F0CC8DC0E7FB3373CF4ED6A786AB329E7F16895CA147AD91F6EAE1DFE38116580DF52381599" +
		"E4246278CB1848FE4A56ABF98652E9E7C2E681551A3D78FA033D932087D8B6567D779A56B726B153033D7" +
		"2231A1B5C16ED7DC4458308D6B64AF6723CC0F52C94E04C58FCA9739E890AA40CC05E22321F10129D2B59" +
		"1F317102034C109A56D711591E5B44C717CFC9C9B9461894767CAFA42D2B394194B03999C2A9EF48868F3" +
		"FB03D1A40F596613AF97F4ED7643A1C2D12692E959C6DEB8E72403ADC0E42204DBCE5056EEF0CC60B0C6E" +
		"83B8B55AC01F6C85644EE49"
)

var params *SidhParams

// Use init() function to initialize params in order to avoid
// static initialization order fiasco.
func init() {
	params = Params(FP_751)
}

// Fail if err !=nil. Display msg as an error message
func checkErr(t testing.TB, err error, msg string) {
	if err != nil {
		t.Error(msg)
	}
}

// Converts string to private key
func convToPrv(s string, v KeyVariant) *PrivateKey {
	key := NewPrivateKey(params.Id, v)
	hex, e := hex.DecodeString(s)
	if e != nil {
		panic("non-hex number provided")
	}
	e = key.Import(hex)
	if e != nil {
		panic("Can't import private key")
	}
	return key
}

// Converts string to public key
func convToPub(s string, v KeyVariant) *PublicKey {
	key := NewPublicKey(params.Id, v)
	hex, e := hex.DecodeString(s)
	if e != nil {
		panic("non-hex number provided")
	}
	e = key.Import(hex)
	if e != nil {
		panic("Can't import public key")
	}
	return key
}

func testKeygen(s *SidhParams, t *testing.T) {
	alicePrivate := convToPrv(PrA, KeyVariant_SIDH_A)
	bobPrivate := convToPrv(PrB, KeyVariant_SIDH_B)
	expPubA := convToPub(PkA, KeyVariant_SIDH_A)
	expPubB := convToPub(PkB, KeyVariant_SIDH_B)

	pubA := alicePrivate.GeneratePublicKey()
	pubB := bobPrivate.GeneratePublicKey()

	if !bytes.Equal(pubA.Export(), expPubA.Export()) {
		t.Fatalf("unexpected value of public key A")
	}
	if !bytes.Equal(pubB.Export(), expPubB.Export()) {
		t.Fatalf("unexpected value of public key B")
	}
}

func testRoundtrip(s *SidhParams, t *testing.T) {
	var err error

	prvA := NewPrivateKey(params.Id, KeyVariant_SIDH_A)
	prvB := NewPrivateKey(params.Id, KeyVariant_SIDH_B)

	// Generate private keys
	err = prvA.Generate(rand.Reader)
	checkErr(t, err, "key generation failed")
	err = prvB.Generate(rand.Reader)
	checkErr(t, err, "key generation failed")

	// Generate public keys
	pubA := prvA.GeneratePublicKey()
	pubB := prvB.GeneratePublicKey()

	// Derive shared secret
	s1, err := DeriveSecret(prvB, pubA)
	checkErr(t, err, "")

	s2, err := DeriveSecret(prvA, pubB)
	checkErr(t, err, "")

	if !bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("Tthe two shared keys: \n%X, \n%X do not match", s1, s2)
	}
}

func testKeyAgreement(s *SidhParams, t testing.TB, pkA, pkB, prA, prB string) {
	var e error

	// KeyPairs
	alicePublic := convToPub(pkA, KeyVariant_SIDH_A)
	bobPublic := convToPub(pkB, KeyVariant_SIDH_B)
	alicePrivate := convToPrv(prA, KeyVariant_SIDH_A)
	bobPrivate := convToPrv(prB, KeyVariant_SIDH_B)

	// Do actual test
	s1, e := DeriveSecret(bobPrivate, alicePublic)
	checkErr(t, e, "derivation s1")
	s2, e := DeriveSecret(alicePrivate, bobPublic)
	checkErr(t, e, "derivation s1")

	if !bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("two shared keys: %d, %d do not match", s1, s2)
	}

	// Negative case
	dec, e := hex.DecodeString(PkA)
	if e != nil {
		t.FailNow()
	}
	dec[0] = ^dec[0]
	e = alicePublic.Import(dec)
	if e != nil {
		t.FailNow()
	}

	s1, e = DeriveSecret(bobPrivate, alicePublic)
	checkErr(t, e, "derivation of s1 failed")
	s2, e = DeriveSecret(alicePrivate, bobPublic)
	checkErr(t, e, "derivation of s2 failed")

	if bytes.Equal(s1[:], s2[:]) {
		t.Fatalf("The two shared keys: %d, %d match", s1, s2)
	}
}

func TestKeygenP751(t *testing.T) {
	testKeygen(Params(FP_751), t)
}

func TestKeyAgreementP751(t *testing.T) {
	testKeyAgreement(Params(FP_751), t, PkA, PkB, PrA, PrB)
}

func TestRoundtripP751(t *testing.T) {
	testRoundtrip(Params(FP_751), t)
}

func TestKeyAgreementP751_AliceEvenNumber(t *testing.T) {
	// even alice
	prE := "C09957CC83045FB4C3726384D784476ACB6FFD92E5B15B3C2D451BA063F1BD4CED8FBCF682A98DD0954D37BCAF730F00"
	pkE := "8A2DE6FD963C475F7829B689C8B8306FC0917A39EBBC35CA171546269A85698FEC0379E2E1A3C567BE1B8EF5639F81F304889737E6CC444DBED4579DB204DC8C7928F5CBB1ECDD682A1B5C48C0DAF34208C06BF201BE4E6063B1BFDC42413B0537F8E76BEE645C1A24118301BAB17EB8D6E0F283BCB16EFB833E4BB3463953C93165A0DDAC55B385059F27FF7228486D0A733812C81C792BE9EC3A16A5DB0EB099EEA76AC0E59612251A3AD19F7CC567DA2AEBD7733171F48E471D17648692355164E27B515D2A47D7BA34B3B48A047BE7C09C4ABEE2FCC9ACA7396C8A8C9E73E29533FC7369094DFA7988778E55E53F309922C6E233F8F9C7936C3D29CEA640406FCA06450AA1978FF39F227BF06B1E072F1763447C6F513B23CDF3B0EC0379070AEE5A02D9AD8E0EB023461D631F4A9643A4C79921334945F6B33DDFC11D9703BD06B047B4DA404AB12EFD2C3A49E5C42D10DA063352748B21DE41C32A5693FE1C0DCAB111F4990CD58BECADB1892EE7A7E99C9DB4DA4E69C96E57138B99038BC9B877ECE75914EFB98DD08B9E4A2DCCB948A8F7D2F26678A9952BA0EFAB1E9CF6E51B557480DEC2BA30DE0FE4AFE30A6B30765EE75EF64F678316D81C72755AD2CFA0B8C7706B07BFA52FBC3DB84EF9E79796C0089305B1E13C78660779E0FF2A13820CE141104F976B1678990F85B2D3D2B89CD5BC4DD52603A5D24D3EFEDA44BAA0F38CDB75A220AF45EAB70F2799875D435CE50FC6315EDD4BB7AA7260AFD7CD0561B69B4FA3A817904322661C3108DA24"
	testKeyAgreement(Params(FP_751), t, pkE, PkB, prE, PrB)
}

func TestImportExport(t *testing.T) {
	var err error
	a := NewPublicKey(params.Id, KeyVariant_SIDH_A)
	b := NewPublicKey(params.Id, KeyVariant_SIDH_B)

	// Import keys
	a_hex, err := hex.DecodeString(PkA)
	checkErr(t, err, "invalid hex-number provided")

	err = a.Import(a_hex)
	checkErr(t, err, "import failed")

	b_hex, err := hex.DecodeString(PkB)
	checkErr(t, err, "invalid hex-number provided")

	err = b.Import(b_hex)
	checkErr(t, err, "import failed")

	// Export and check if same
	if !bytes.Equal(b.Export(), b_hex) || !bytes.Equal(a.Export(), a_hex) {
		t.Fatalf("export/import failed")
	}

	if (len(b.Export()) != b.Size()) || (len(a.Export()) != a.Size()) {
		t.Fatalf("wrong size of exported keys")
	}
}

func TestMultiplyByThree(t *testing.T) {
	// sage: repr((3^238 -1).digits(256))
	var three238minus1 = [48]byte{
		248, 132, 131, 130, 138, 113, 205, 237, 20, 122, 66, 212, 191, 53, 59, 115, 56, 207,
		215, 148, 207, 41, 130, 248, 214, 42, 124, 12, 153, 108, 197, 99, 199, 34, 66, 143,
		126, 168, 88, 184, 245, 234, 37, 181, 198, 201, 84, 2}
	// sage: repr((3*(3^238 -1)).digits(256))
	var threeTimesThree238minus1 = [48]byte{
		232, 142, 138, 135, 159, 84, 104, 201, 62, 110, 199, 124, 63, 161, 177, 89, 169, 109,
		135, 190, 110, 125, 134, 233, 132, 128, 116, 37, 203, 69, 80, 43, 86, 104, 198, 173,
		123, 249, 9, 41, 225, 192, 113, 31, 84, 93, 254, 6}

	multiplyByThree(three238minus1[:])

	for i := 0; i < 48; i++ {
		if three238minus1[i] != threeTimesThree238minus1[i] {
			t.Error("Digit", i, "error: found", three238minus1[i],
				"expected", threeTimesThree238minus1[i])
		}
	}
}

func TestCheckLessThanThree238(t *testing.T) {
	var three238minus1 = [48]byte{
		248, 132, 131, 130, 138, 113, 205, 237, 20, 122, 66, 212, 191, 53, 59, 115,
		56, 207, 215, 148, 207, 41, 130, 248, 214, 42, 124, 12, 153, 108, 197, 99,
		199, 34, 66, 143, 126, 168, 88, 184, 245, 234, 37, 181, 198, 201, 84, 2}
	var three238 = [48]byte{
		249, 132, 131, 130, 138, 113, 205, 237, 20, 122, 66, 212, 191, 53, 59, 115,
		56, 207, 215, 148, 207, 41, 130, 248, 214, 42, 124, 12, 153, 108, 197, 99, 199,
		34, 66, 143, 126, 168, 88, 184, 245, 234, 37, 181, 198, 201, 84, 2}
	var three238plus1 = [48]byte{250, 132, 131, 130, 138, 113, 205, 237, 20, 122, 66,
		212, 191, 53, 59, 115, 56, 207, 215, 148, 207, 41, 130, 248, 214, 42, 124, 12,
		153, 108, 197, 99, 199, 34, 66, 143, 126, 168, 88, 184, 245, 234, 37, 181, 198,
		201, 84, 2}
	// makes second 64-bit digits bigger than in three238. checks if carries are correctly propagated
	var three238plus2power65 = [48]byte{249, 132, 131, 130, 138, 113, 205, 237, 22, 122,
		66, 212, 191, 53, 59, 115, 56, 207, 215, 148, 207, 41, 130, 248, 214, 42, 124, 12,
		153, 108, 197, 99, 199, 34, 66, 143, 126, 168, 88, 184, 245, 234, 37, 181, 198,
		201, 84, 2}

	var result uint8

	result = checkLessThanThree238(three238minus1[:])
	if result != 0 {
		t.Error("expected 0, got", result)
	}
	result = checkLessThanThree238(three238[:])
	if result != 1 {
		t.Error("expected nonzero, got", result)
	}
	result = checkLessThanThree238(three238plus1[:])
	if result != 1 {
		t.Error("expected nonzero, got", result)
	}
	result = checkLessThanThree238(three238plus2power65[:])
	if result != 1 {
		t.Error("expected nonzero, got", result)
	}
}

func BenchmarkSidhKeyAgreement(b *testing.B) {
	// KeyPairs
	alicePublic := convToPub(PkA, KeyVariant_SIDH_A)
	bobPublic := convToPub(PkB, KeyVariant_SIDH_B)
	alicePrivate := convToPrv(PrA, KeyVariant_SIDH_A)
	bobPrivate := convToPrv(PrB, KeyVariant_SIDH_B)

	for i := 0; i < b.N; i++ {
		// Derive shared secret
		DeriveSecret(bobPrivate, alicePublic)
		DeriveSecret(alicePrivate, bobPublic)
	}
}

func BenchmarkAliceKeyGenPrv(b *testing.B) {
	prv := NewPrivateKey(params.Id, KeyVariant_SIDH_A)
	for n := 0; n < b.N; n++ {
		prv.Generate(rand.Reader)
	}
}

func BenchmarkBobKeyGenPrv(b *testing.B) {
	prv := NewPrivateKey(params.Id, KeyVariant_SIDH_B)
	for n := 0; n < b.N; n++ {
		prv.Generate(rand.Reader)
	}
}

func BenchmarkAliceKeyGenPub(b *testing.B) {
	prv := NewPrivateKey(params.Id, KeyVariant_SIDH_A)
	prv.Generate(rand.Reader)
	for n := 0; n < b.N; n++ {
		prv.GeneratePublicKey()
	}
}

func BenchmarkBobKeyGenPub(b *testing.B) {
	prv := NewPrivateKey(params.Id, KeyVariant_SIDH_B)
	prv.Generate(rand.Reader)
	for n := 0; n < b.N; n++ {
		prv.GeneratePublicKey()
	}
}

func BenchmarkSharedSecretAlice(b *testing.B) {
	aPr := convToPrv(PrA, KeyVariant_SIDH_A)
	bPk := convToPub(PkB, KeyVariant_SIDH_B)
	for n := 0; n < b.N; n++ {
		DeriveSecret(aPr, bPk)
	}
}

func BenchmarkSharedSecretBob(b *testing.B) {
	// m_B = 3*randint(0,3^238)
	aPk := convToPub(PkA, KeyVariant_SIDH_A)
	bPr := convToPrv(PrB, KeyVariant_SIDH_B)
	for n := 0; n < b.N; n++ {
		DeriveSecret(bPr, aPk)
	}
}
