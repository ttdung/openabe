package main

// #cgo LDFLAGS: -L. -llibrary
// #include "lib_bridge.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

type ABE struct {
	ptr unsafe.Pointer
}

func InitializeOpenABE() {
	C.LIB_InitializeOpenABE()
}

func ShutdownABE() {
	C.LIB_ShutdownOpenABE()
}

func NewCString(input string) *C.char {
	for {
		output := C.CString(input)
		str := fmt.Sprintf("%p", output)
		if len(str) == 14 {
			fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!! output addr ", output)
			C.free(unsafe.Pointer(output))
			continue
		}

		fmt.Println("output addr ", output)
		return output
	}
}

func NewABE(abename string) ABE {
	for {
		name := NewCString(abename)
		defer C.free(unsafe.Pointer(name))

		var abe ABE
		abe.ptr = C.LIB_NewABE(C.CString(abename))

		str := fmt.Sprintf("%p", abe.ptr)
		if len(str) == 14 {
			C.free(unsafe.Pointer(abe.ptr))
			continue
		}

		fmt.Println("abe.ptr ", abe.ptr)
		return abe
	}
}

func (abe *ABE) GenerateParams() {
	C.LIB_generateParams(abe.ptr)
}

func (abe *ABE) Genkey(att string, key string) {
	latt := NewCString(att)
	lkey := NewCString(key)

	defer C.free(unsafe.Pointer(latt))
	defer C.free(unsafe.Pointer(lkey))

	C.LIB_keygen(abe.ptr, latt, lkey)
}

func (abe *ABE) Encrypt(att string, pt string) string {
	latt := NewCString(att)
	lpt := NewCString(pt)

	defer C.free(unsafe.Pointer(latt))
	defer C.free(unsafe.Pointer(lpt))

	return C.GoString(C.LIB_encrypt(abe.ptr, latt, lpt))
}

func (abe *ABE) Decrypt(key string, ct string) string {
	lkey := NewCString(key)
	lct := NewCString(ct)

	defer C.free(unsafe.Pointer(lkey))
	defer C.free(unsafe.Pointer(lct))

	return C.GoString(C.LIB_decrypt(abe.ptr, lkey, lct))

}

func (abe *ABE) ExportUserKey(key string) string {
	lkey := NewCString(key)

	defer C.free(unsafe.Pointer(lkey))

	return C.GoString(C.LIB_exportUserKey(abe.ptr, lkey))
}

func (abe *ABE) ImportUserKey(index string, key string) string {
	lkey := NewCString(key)
	lidx := NewCString(index)

	defer C.free(unsafe.Pointer(lkey))
	defer C.free(unsafe.Pointer(lidx))

	//LIB_exportUserKey
	return C.GoString(C.LIB_importUserKey(abe.ptr, lidx, lkey))

}

func (abe *ABE) ExportMPK() string {
	return C.GoString(C.LIB_exportMPK(abe.ptr))
}

func (abe *ABE) ExportMSK() string {
	return C.GoString(C.LIB_exportMSK(abe.ptr))
}

func (abe *ABE) ImportMSK(key string) {
	lkey := NewCString(key)

	defer C.free(unsafe.Pointer(lkey))

	C.LIB_importMSK(abe.ptr, lkey)
}

func (abe *ABE) ImportMPK(key string) {
	lkey := NewCString(key)

	defer C.free(unsafe.Pointer(lkey))

	C.LIB_importMPK(abe.ptr, lkey)
}

func (abe *ABE) ImportAndDecrypt(key string, ct string) string {
	lkey := NewCString(key)
	lct := NewCString(ct)

	defer C.free(unsafe.Pointer(lkey))
	defer C.free(unsafe.Pointer(lct))

	return C.GoString(C.LIB_ImportAndDecrypt(abe.ptr, lkey, lct))
}

func (abe *ABE) KeygenAndDecrypt(att string, ct string) string {
	latt := NewCString(att)
	lct := NewCString(ct)

	defer C.free(unsafe.Pointer(latt))
	defer C.free(unsafe.Pointer(lct))

	return C.GoString(C.LIB_KeygenAndDecrypt(abe.ptr, latt, lct))
}

func main() {
	InitializeOpenABE()

	abe := NewABE("CP-ABE")

	// abe.GenerateParams() // (MPK, MSK)

	// msk := abe.ExportMSK()
	// mpk := abe.ExportMPK()

	msk := "AAAAFqpvyQIZL3nR3yOZR9Px5By5mHBtc2sAAAB2oQVhbHBoYaEisQAfjGiaHVVxMOk/id6sYAoIsmCq7ySFq3pVU5/STat3+KEDZzJhoUSzoUECGJpkzPFeDuoErI4CJowASS86VZhG2X6vW8j7IZ2Za2gPNimuppqdqIJfdBVyugJ1cyceOMi4b8jmW6BN3gT2FQ=="
	mpk := "AAAAFqpvyd/crzM1+7CsJYDbP4E7anltcGsAAAHToQFBsgEEtLIBABcLKfKn1/PVsOT3HjO7T2eZz5fPwbNBKmHhAFChQtIFCWNIOg20QHuj12yQit/+ezMALP509rIaIhiRs6Ijf7wXwScYoabFyv2Ezpu6y0P0Sk9dXh48cAqLCL8DVQPOYgPTq1IQdGa+kizgP6IbPyqMhX3dkimqKQXx6a8QPhN8EPQDbhE7uDZxFKN8S/k5pmuGAMTPVb3qOWmnLUgFivMHRul5iD7176c6BOagdKXj9qiQaK+G7r04dMAsedR9Oh+O0YKEcEDkhxcfX8nTKeWpq3q5FAC1LYzuKy3MZHqlI6kvc5OKQaE2iVIG6k2LgVwWH/eRAbGWO6KIBwtvcuGhAmcxoSSyoSECDfXYPLhSTnAePm+FkAz7D/YOtnTMccNO9zkTDPYlrxmhA2cxYaEksqEhAiEraf5oInnjDMrpNgZml1JD/++krtTtSLbRpoW+W8GAoQJnMqFEs6FBAhEJuCjd23TEUK+/rupQNCgYudXQivRLKajQqhH2Qq1OG8e3+fKKRw6ETk8/8QmnYQvv4yCPRNWbD/81VxGlkSWhAWuhJR0AAAAgIuvHIm519bNW3v5yQ8Kb9VOl1bOG0h2k9TWfdY0I3EE="
	// fmt.Println("msk:", msk)
	// fmt.Println("mpk:", mpk)
	abe.ImportMSK(msk)
	abe.ImportMPK(mpk)

	abe.Genkey("student|math", "key")
	alice_key := abe.ExportUserKey("key")
	abe.Genkey("student|CS", "key")
	bob_key := abe.ExportUserKey("key")
	abe.Genkey("student|EE", "key")
	carol_key := abe.ExportUserKey("key")

	abe1 := NewABE("CP-ABE")
	abe1.ImportMPK(mpk)
	// abe1.ImportUserKey("alice", alice_key)

	data := "hello world"

	ct := abe1.Encrypt("(student) and (math or EE)", data)

	//pt := abe.decrypt("key_alice", ct)

	abe2 := NewABE("CP-ABE")
	abe2.ImportMPK(mpk)
	akey := abe2.ImportUserKey("alice", alice_key)
	pt := abe2.Decrypt(akey, ct)
	// pt := abe2.ImportAndDecrypt(alice_key, ct)
	// pt := abe2.KeygenAndDecrypt("student|math", ct)

	if pt == data {
		fmt.Printf("Alice Decrypt Successful pt = %v \n", pt)
	} else {
		fmt.Println("Alice Fail to decrypt")
	}

	bkey := abe2.ImportUserKey("bob", bob_key)
	pt = abe2.Decrypt(bkey, ct)
	// pt = abe2.KeygenAndDecrypt("student|CS", ct)

	if pt == data {
		fmt.Printf("Bob Decrypt Successful pt = %v \n", pt)
	} else {
		fmt.Println("Bob Fail to decrypt")
	}

	ckey := abe2.ImportUserKey("carol", carol_key)
	pt = abe2.Decrypt(ckey, ct)

	// pt = abe2.KeygenAndDecrypt("student|EE", ct)
	if pt == data {
		fmt.Printf("Carol Decrypt Successful pt = %v \n", pt)
	} else {
		fmt.Println("Carol Fail to decrypt")
	}

	ShutdownABE()
}
