package main

// #cgo LDFLAGS: -L. -llibrary
// #include "lib_bridge.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type ABE struct {
	ptr unsafe.Pointer
}

func NewABE(abename string) ABE {
	var abe ABE
	abe.ptr = C.LIB_NewABE(C.CString(abename))
	return abe
}

func (abe ABE) generateParams() {
	C.LIB_generateParams(abe.ptr)
}

func (abe ABE) genkey(att string, key string) {
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

// d, err := os.ReadFile("./sample.json")
// if err != nil {
// 	panic(err)
// }
// data := string(d)

func main() {
	abe := NewABE("CP-ABE")

	abe.generateParams() // (MPK, MSK)

	abe.genkey("student|math", "key_alice")
	abe.genkey("student|CS", "key_bob")

	data := "hello world"

	ct := abe.encrypt("(student) and (math or EE)", data)

	pt := abe.decrypt("key_alice", ct)

	if pt == data {
		fmt.Printf("Decrypt Successful pt = %v \n", pt)
	} else {
		fmt.Println("Fail to decrypt")
	}
}
