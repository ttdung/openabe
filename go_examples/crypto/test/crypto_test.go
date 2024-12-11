package test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"

	"kscm.kasikornbank.com/10000/abe/sds-client/conf"
	"kscm.kasikornbank.com/10000/abe/sds-client/internal/app/crypto"
)

type CryptoTestSuite struct {
	suite.Suite
	config *conf.Configuration
}

// =================================================================
// Test suite: Test_ABE_Encrypt_Decrypt, Test_AES_Encrypt_Decrypt, Test_RSA_Encrypt_Decrypt
// =================================================================
func TestCryptoTestSuite(t *testing.T) {
	suite.Run(t, new(CryptoTestSuite))
}

// this function executes before the test suite begins execution
func (suite *CryptoTestSuite) SetupSuite() {
	fmt.Println("[Crypto] Start crypto test suite")
	suite.config = conf.GetConfiguration()
}

// this function executes after all tests executed
func (suite *CryptoTestSuite) TearDownSuite() {
	fmt.Println("Finish crypto test suite")
}

// =============================================================
// Test cases: Test_ABE_Encrypt_Decrypt
// Description:
//  1. Encrypt data using ABE
//  2. Decrypt data using ABE
//  3. Verify the decrypted data
//
// =============================================================
func (suite *CryptoTestSuite) Test_ABE_Encrypt_Decrypt() {
	fmt.Println("\n[Crypto] Start ABE test case")

	message := "Hello World"
	accesstree := "admin or employee"

	enc_message, err := crypto.AbeEncrypt(suite.config.Mpk, accesstree, message)
	suite.Nil(err)

	// Generate key
	abekey, err := crypto.GenerateABEKey("admin", "testabekey")
	suite.Nil(err)

	abeFailedKey, err := crypto.GenerateABEKey("leader", "testfailedabekey")
	suite.Nil(err)

	decrypt_message, err := crypto.AbeDecrypt(suite.config.Mpk, "testabekey", abekey, enc_message)
	suite.Nil(err)

	suite.Equal(message, decrypt_message)

	decrypt_message, err = crypto.AbeDecrypt(suite.config.Mpk, "testfailedabekey", abeFailedKey, enc_message)
	suite.Nil(err)

	suite.NotEqual(message, decrypt_message)
}

// =============================================================
// Test cases: Test_AES_Encrypt_Decrypt
// Description:
//  1. Encrypt data using AES
//  2. Decrypt data using AES
//  3. Verify the decrypted data
//
// =============================================================
func (suite *CryptoTestSuite) Test_AES_Encrypt_Decrypt() {
	fmt.Println("\n[Crypto] Start AES test case")

	plaintext := "Hello World"

	cipher_text, err := crypto.AESEncrypt(plaintext, []byte(suite.config.SecretKey))
	suite.Nil(err)
	suite.NotEmpty(cipher_text)

	decrypt_message, err := crypto.AESDecrypt(cipher_text, []byte(suite.config.SecretKey))

	suite.Nil(err)
	suite.Equal(plaintext, decrypt_message)
}

// =============================================================
// Test cases: Test_AES_Encrypt_EmptyPlaintext
// Description:
//  1. Encrypt empty data using AES
//  2. Verify the encrypted data
//  3. Expect error
//
// =============================================================
func (suite *CryptoTestSuite) TestAESEncrypt_EmptyPlaintext() {
	fmt.Println("\n[Crypto] Start AES empty plaintext test case")

	plaintext := ""
	secretKey := []byte("0123456789abcdef0123456789abcdef") // 32-byte key

	ciphertext, err := crypto.AESEncrypt(plaintext, secretKey)
	suite.NotNil(err)
	suite.Empty(ciphertext)
}

// =============================================================
// Test cases: Test_RSA_Encrypt_Decrypt
// Description:
//  1. Encrypt data using RSA
//  2. Decrypt data using RSA
//  3. Verify the decrypted data
//
// =============================================================
func (suite *CryptoTestSuite) Test_RSA_Encrypt_Decrypt() {

	fmt.Println("\n[Crypto] Start RSA test case")

	plaintext := "Hello World"

	sender_pubkey := "p9esDSIhBAG/xZk0cGTm368HoLRh6ccrYzPeCRY7f1k="
	sender_private := "XZxvaveCopgzefMq7IOXMxrITtQ9wYSVt3bDmNGYdJk="

	receiver_pubkey := suite.config.AdminPubkey
	receiver_privkey := suite.config.AdminPrivateKey

	//=============================================================
	// Encrypt data
	//=============================================================
	var receiver_pubkey32 [32]byte
	rpk, err := base64.StdEncoding.DecodeString(receiver_pubkey)
	suite.Nil(err)
	copy(receiver_pubkey32[:], rpk[:])

	var sender_private32 [32]byte
	ssk, err := base64.StdEncoding.DecodeString(sender_private)
	suite.Nil(err)
	copy(sender_private32[:], ssk[:])

	cipher_text, err := crypto.RSAEncrypt(plaintext, &receiver_pubkey32, &sender_private32)
	suite.Nil(err)

	//=============================================================
	// Decrypt data
	//=============================================================
	var sender_pubkey32 [32]byte
	spk, err := base64.StdEncoding.DecodeString(sender_pubkey)
	suite.Nil(err)
	copy(sender_pubkey32[:], spk[:])

	var receiver_private32 [32]byte
	rsk, err := base64.StdEncoding.DecodeString(receiver_privkey)
	suite.Nil(err)
	copy(receiver_private32[:], rsk[:])

	decrypt_message, err := crypto.RSADecrypt(cipher_text, &sender_pubkey32, &receiver_private32)

	suite.Nil(err)
	suite.Equal(plaintext, string(decrypt_message))
}

// =============================================================
// Test cases: Test_CenerateAES256Key
// Description:
//  1. Generate 2 AES 256 keys
//  2. Generate 2 ABE keys
//  2. Verify the randomness of generated keys
//
// =============================================================
func (suite *CryptoTestSuite) Test_Genkey() {

	fmt.Println("\n[Crypto] Start test genrate keys")

	sk1, err := crypto.GenerateAES256Key()
	suite.Nil(err)

	suite.NotNil(sk1)

	fmt.Println("sk1:", base64.StdEncoding.EncodeToString(sk1[:]))

	sk2, err := crypto.GenerateAES256Key()
	suite.Nil(err)

	fmt.Println("sk2:", base64.StdEncoding.EncodeToString(sk2[:]))
	suite.NotEqual(fmt.Sprintf("%v", sk1), fmt.Sprintf("%v", sk2))

	abe1, err := crypto.GenerateABEKey("hello", "")
	suite.Nil(err)

	fmt.Println("abe1:", abe1)
	abe2, err := crypto.GenerateABEKey("hello", "")
	suite.Nil(err)

	fmt.Println("abe2:", abe2)

	suite.NotEqual(fmt.Sprintf("%v", abe1), fmt.Sprintf("%v", abe2))
}

func (suite *CryptoTestSuite) Test_new_ABE_Encrypt_Decrypt() {

	fmt.Println("\n[Crypto] Start New ABE test case")

	crypto.InitializeOpenABE()

	message := "Hello World"

	accesstree := "admin or employee"

	// GEN KEY PHASE

	param := crypto.AbeParams{

		Config: suite.config,
	}

	abeKGC := crypto.NewKGCAbe(param)

	// Generate key

	alice_key, err := abeKGC.GenerateABEKey("admin", "alice")
	suite.Nil(err)

	bob_key, err := abeKGC.GenerateABEKey("manager", "bob")
	suite.Nil(err)

	// TEST ENCRYPT DECRYPT PHASE

	abeClient1 := crypto.NewClientAbe(param)

	enc_message, err := abeClient1.AbeEncrypt(context.Background(), accesstree, message)

	// _, err = abeObj1.AbeEncrypt(context.Background(), accesstree, message)

	suite.Nil(err)

	abeClient2 := crypto.NewClientAbe(param)

	decrypt_message, err := abeClient2.AbeDecrypt(context.Background(), "alice", alice_key, enc_message)

	suite.Nil(err)

	suite.Equal(message, decrypt_message)

	decrypt_message, err = abeClient2.AbeDecrypt(context.Background(), "bob", bob_key, enc_message)

	suite.Nil(err)

	suite.NotEqual(message, decrypt_message)

	// abeObj.ShuttingDown()

	// abeObj = crypto.InitClientAbe(param, "testfailedabekey", abeFailedKey)

	// decrypt_message, err = abeObj.AbeDecrypt(context.Background(), "testfailedabekey", abeFailedKey, enc_message)

	// suite.Nil(err)

	// suite.NotEqual(message, decrypt_message)

	crypto.ShutdownABE()

}
