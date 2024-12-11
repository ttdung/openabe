package crypto

// #cgo LDFLAGS: -L../../../lib -llibrary
// #include "../../../include/lib_bridge.h"
// #include <stdlib.h>
import "C"

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	rand1 "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/crypto/nacl/box"
	"kscm.kasikornbank.com/10000/abe/sds-client/conf"
)

type ABE struct {
	ptr unsafe.Pointer
}

func InitializeOpenABE() {
	C.LIB_InitializeOpenABE()
}

func ShutdownABE() {
	C.LIB_ShutdownABE()
}

func NewABE(abename string) ABE {
	var abe ABE
	abe.ptr = C.LIB_NewABE(C.CString(abename))
	return abe
}

func (abe *ABE) GenerateParams() {
	C.LIB_generateParams(abe.ptr)
}

func (abe *ABE) Genkey(att string, key string) {
	latt := C.CString(att)
	lkey := C.CString(key)

	defer C.free(unsafe.Pointer(latt))
	defer C.free(unsafe.Pointer(lkey))

	C.LIB_keygen(abe.ptr, latt, lkey)
}

func (abe *ABE) Encrypt(att string, pt string) string {
	latt := C.CString(att)
	lpt := C.CString(pt)

	defer C.free(unsafe.Pointer(latt))
	defer C.free(unsafe.Pointer(lpt))

	return C.GoString(C.LIB_encrypt(abe.ptr, latt, lpt))
}

func (abe *ABE) Decrypt(key string, ct string) string {
	lkey := C.CString(key)
	lct := C.CString(ct)

	defer C.free(unsafe.Pointer(lkey))
	defer C.free(unsafe.Pointer(lct))

	return C.GoString(C.LIB_decrypt(abe.ptr, lkey, lct))

}

func (abe *ABE) ExportUserKey(key string) string {
	lkey := C.CString(key)

	defer C.free(unsafe.Pointer(lkey))

	return C.GoString(C.LIB_exportUserKey(abe.ptr, lkey))
}

func (abe *ABE) ImportUserKey(key string) string {
	lkey := C.CString(key)

	defer C.free(unsafe.Pointer(lkey))

	//LIB_exportUserKey
	return C.GoString(C.LIB_importUserKey(abe.ptr, lkey))

}

func (abe *ABE) ExportMPK() string {
	return C.GoString(C.LIB_exportMPK(abe.ptr))
}

func (abe *ABE) ExportMSK() string {
	return C.GoString(C.LIB_exportMSK(abe.ptr))
}

func (abe *ABE) ImportMSK(key string) {
	lkey := C.CString(key)

	defer C.free(unsafe.Pointer(lkey))

	C.LIB_importMSK(abe.ptr, lkey)
}

func (abe *ABE) ImportMPK(key string) {
	lkey := C.CString(key)

	defer C.free(unsafe.Pointer(lkey))

	C.LIB_importMPK(abe.ptr, lkey)
}

func AbeEncrypt(mpk string, accesstree string, data string) (string, error) {

	if len(accesstree) == 0 {
		return "", errors.New("accesstree is empty")
	}
	if len(data) == 0 {
		return "", errors.New("data is empty")
	}

	InitializeOpenABE()

	abe := NewABE("CP-ABE")

	abe.ImportMPK(mpk)

	ct := abe.Encrypt(accesstree, data)

	ShutdownABE()

	return ct, nil
}

func AbeDecrypt(mpk string, idx string, ekey string, ct string) (string, error) {

	if len(mpk) == 0 {
		return "", errors.New("MPK is empty")
	}
	if len(ekey) == 0 {
		return "", errors.New("ABE secret key is empty")
	}
	if len(ct) == 0 {
		return "", nil
	}

	InitializeOpenABE()

	abe := NewABE("CP-ABE")

	abe.ImportMPK(mpk)

	key := abe.ImportUserKey(ekey)

	pt := abe.Decrypt(key, ct)

	ShutdownABE()

	return pt, nil
}

func AESEncrypt(plaintext string, secretKey []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", errors.New("plaintext is empty")
	}

	aes, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	// We need a 12-byte nonce for GCM(modifiable if you use cipher[GCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand1.Read(nonce)
	if err != nil {
		return "", err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AESDecrypt(ciphertext string, secretKey []byte) (string, error) {
	if len(ciphertext) == 0 {
		return "", nil
	}

	ct, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	aes, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, cipherText := ct[:nonceSize], ct[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(cipherText), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func GenerateABEKey(att string, idx string) (string, error) {

	if len(att) == 0 {
		return "", errors.New("attritbutes are empty")
	}

	config := conf.GetConfiguration()

	// Init ABE Lib
	InitializeOpenABE()
	abe := NewABE("CP-ABE")

	// Import keys
	abe.ImportMSK(config.Msk)
	abe.ImportMPK(config.Mpk)

	// Generate key
	abe.Genkey(att, idx)

	// Export key
	abekey := abe.ExportUserKey(idx)

	// Shutdown ABE lib
	ShutdownABE()

	return abekey, nil
}

// Purpose: decrypt encrypted data which is encrypted by a RSA pubkey
// encdata: encrypted data in base64 encoding
// senderPubkey: sender public key
// privatekey: private key of receiver
func RSADecrypt(encryptdata string, senderPubkey *[32]byte, receiverPrivatekey *[32]byte) ([]byte, error) {

	encdata, err := base64.StdEncoding.DecodeString(encryptdata)
	if err != nil {
		return nil, err
	}

	// The recipient can decrypt the message using their private key and the
	// sender's public key. When you decrypt, you must use the same nonce you
	// used to encrypt the message. One way to achieve this is to store the
	// nonce alongside the encrypted message. Above, we stored the nonce in the
	// first 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], encdata[:24])

	data, ok := box.Open(nil, encdata[24:], &decryptNonce, senderPubkey, receiverPrivatekey)

	if !ok {
		return nil, errors.New("RSADecrypt: failed to decrypt data")
	}

	return data, nil
}

func RSAEncrypt(data string, receiverRSAPubkeyKey *[32]byte, senderRSAPrivatekey *[32]byte) (string, error) {

	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	var nonce [24]byte
	if _, err := io.ReadFull(rand1.Reader, nonce[:]); err != nil {
		return "", err
	}

	encryptedData := box.Seal(nonce[:], []byte(data), &nonce, receiverRSAPubkeyKey, senderRSAPrivatekey)

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func GenerateRSAKeyPair() (publicKey *[32]byte, privateKey *[32]byte, err error) {
	return box.GenerateKey(rand1.Reader)
}

// generateEncryptionKey generates a 256 bit (32 byte) AES encryption key and
// prints the base64 representation.
func GenerateAES256Key() (*[32]byte, error) {
	var key [32]byte
	if _, err := rand1.Read(key[:]); err != nil {
		return nil, fmt.Errorf("rand.Read: %w", err)
	}

	return &key, nil
}

// Create function AES256Encryption to encrypts the plaintext using the provided key and returns the ciphertext.using best secure library nacl/box. Note do not re-use function AESEncrypt
func AES256Encryption(plaintext string, key *[32]byte) (string, error) {
	return AESEncrypt(plaintext, key[:])
}

type AbeParams struct {

	// fx.In

	Config *conf.Configuration
}

type abeObj struct {
	// ptr unsafe.Pointer,
	abe    ABE
	Config *conf.Configuration
}

type AbeObj interface {
	AbeEncrypt(ctx context.Context, accesstree string, data string) (string, error)
	AbeDecrypt(ctx context.Context, idx string, ekey string, ct string) (string, error)
	GenerateABEKey(att string, idx string) (string, error)
}

func NewKGCAbe(p AbeParams) AbeObj {

	// pre_ini_state
	var abe ABE
	abe.ptr = C.LIB_NewABE(C.CString("CP-ABE"))

	// Import MSK keys
	abe.ImportMSK(p.Config.Msk)

	abe.ImportMPK(p.Config.Mpk)

	return &abeObj{
		abe:    abe,
		Config: p.Config,
	}
}

func NewClientAbe(p AbeParams) AbeObj {

	// pre_ini_state
	var abe ABE
	abe.ptr = C.LIB_NewABE(C.CString("CP-ABE"))

	// Import MSK keys
	abe.ImportMPK(p.Config.Mpk)

	return &abeObj{
		abe:    abe,
		Config: p.Config,
	}
}

func (obj *abeObj) GenerateABEKey(att string, idx string) (string, error) {

	if len(att) == 0 {
		return "", errors.New("attritbutes are empty")
	}

	if len(idx) == 0 {
		idx = "key"
	}

	// Generate key
	obj.abe.Genkey(att, idx)

	// Export key
	abekey := obj.abe.ExportUserKey(idx)

	return abekey, nil
}

func (obj *abeObj) AbeEncrypt(ctx context.Context, accesstree string, data string) (string, error) {

	// panic recover
	if len(accesstree) == 0 {
		return "", errors.New("accesstree is empty")
	}

	if len(data) == 0 {
		return "", errors.New("data is empty")
	}

	ct := obj.abe.Encrypt(accesstree, data)

	return ct, nil
}

func (obj *abeObj) AbeDecrypt(
	ctx context.Context,
	idx string,
	ekey string,
	ct string,
) (string, error) {

	if len(ekey) == 0 {
		return "", errors.New("ABE secret key is empty")
	}

	if len(ct) == 0 {
		return "", nil
	}

	// concern about pointer abe in pointer have many call from another one
	key := obj.abe.ImportUserKey(ekey)

	pt := obj.abe.Decrypt(key, ct)

	return pt, nil
}
