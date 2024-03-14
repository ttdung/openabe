package main

// #cgo LDFLAGS: -L. -llibrary
// #include "lib_bridge.h"
import "C"
import (
	"fmt"
	"unsafe"

	p "github.com/nxhieu3102/openabe/policy"
)

type ABE struct {
	ptr unsafe.Pointer
}

func NewABE(abename string) ABE {
	return ABE{
		ptr: C.LIB_NewABE(C.CString(abename)),
	}
}

// re format token with format
// char "(" : 4, char ")": 5
// string "and": 2, string "or": 3
// value exist in attribute map: 1
// value not exist in attribute map: 0

func (abe ABE) generateParams() {
	C.LIB_generateParams(abe.ptr)
}

func (abe *ABE) genkey(att string, key string) {
	latt := C.CString(att)
	lkey := C.CString(key)

	C.LIB_keygen(abe.ptr, latt, lkey)
}

func (abe ABE) encrypt(att string, pt string) string {
	latt := C.CString(att)
	lpt := C.CString(pt)

	return C.GoString(C.LIB_encrypt(abe.ptr, latt, lpt))
}

func (abe ABE) decrypt(key string, ct string) string {

	lkey := C.CString(key)
	lct := C.CString(ct)

	return C.GoString(C.LIB_decrypt(abe.ptr, lkey, lct))

}

func (abe ABE) exportUserKey(key string) string {

	lkey := C.CString(key)
	return C.GoString(C.LIB_exportUserKey(abe.ptr, lkey))

}

func (abe ABE) importUserKey(key string) string {

	lkey := C.CString(key)
	//LIB_exportUserKey
	return C.GoString(C.LIB_importUserKey(abe.ptr, lkey))

}

// d, err := os.ReadFile("./sample.json")
// if err != nil {
// 	panic(err)
// }
// data := string(d)

func main() {

	policyValidator := p.NewValidator("(student) and (math or EE)")
	if !policyValidator.Validate("student|math") {
		fmt.Println("Attribute not valid with policy")
		return
	}

	abe := NewABE("CP-ABE")
	abe.generateParams() // (MPK, MSK)
	abe.genkey("student|math", "key_alice")
	abe.genkey("student|CS", "key_bob")

	ekey := abe.exportUserKey("key_alice")

	data := "hello world"

	ct := abe.encrypt("(student) and (math or EE)", data)

	//pt := abe.decrypt("key_alice", ct)

	ikey := abe.importUserKey(ekey)
	pt := abe.decrypt(ikey, ct)

	if pt == data {
		fmt.Printf("Decrypt Successful pt = %v \n", pt)
	} else {
		fmt.Println("Fail to decrypt")
	}
}
