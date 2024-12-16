#include "lib_bridge.h"
#include "openabe/openabe.h"
#include <openabe/zsymcrypto.h>

using namespace std;
using namespace oabe;
using namespace oabe::crypto;


int main() {
    
	//LIB_InitializeOpenABE();

	InitializeOpenABE();

	auto abe = LIB_NewABE("CP-ABE");

	// abe.GenerateParams() // (MPK, MSK)

	// msk := abe.ExportMSK()
	// mpk := abe.ExportMPK()

	char* msk = "AAAAFqpvyQIZL3nR3yOZR9Px5By5mHBtc2sAAAB2oQVhbHBoYaEisQAfjGiaHVVxMOk/id6sYAoIsmCq7ySFq3pVU5/STat3+KEDZzJhoUSzoUECGJpkzPFeDuoErI4CJowASS86VZhG2X6vW8j7IZ2Za2gPNimuppqdqIJfdBVyugJ1cyceOMi4b8jmW6BN3gT2FQ==";
	char* mpk = "AAAAFqpvyd/crzM1+7CsJYDbP4E7anltcGsAAAHToQFBsgEEtLIBABcLKfKn1/PVsOT3HjO7T2eZz5fPwbNBKmHhAFChQtIFCWNIOg20QHuj12yQit/+ezMALP509rIaIhiRs6Ijf7wXwScYoabFyv2Ezpu6y0P0Sk9dXh48cAqLCL8DVQPOYgPTq1IQdGa+kizgP6IbPyqMhX3dkimqKQXx6a8QPhN8EPQDbhE7uDZxFKN8S/k5pmuGAMTPVb3qOWmnLUgFivMHRul5iD7176c6BOagdKXj9qiQaK+G7r04dMAsedR9Oh+O0YKEcEDkhxcfX8nTKeWpq3q5FAC1LYzuKy3MZHqlI6kvc5OKQaE2iVIG6k2LgVwWH/eRAbGWO6KIBwtvcuGhAmcxoSSyoSECDfXYPLhSTnAePm+FkAz7D/YOtnTMccNO9zkTDPYlrxmhA2cxYaEksqEhAiEraf5oInnjDMrpNgZml1JD/++krtTtSLbRpoW+W8GAoQJnMqFEs6FBAhEJuCjd23TEUK+/rupQNCgYudXQivRLKajQqhH2Qq1OG8e3+fKKRw6ETk8/8QmnYQvv4yCPRNWbD/81VxGlkSWhAWuhJR0AAAAgIuvHIm519bNW3v5yQ8Kb9VOl1bOG0h2k9TWfdY0I3EE=";
	// fmt.Println("msk:", msk)
	// fmt.Println("mpk:", mpk)
	LIB_importMSK(abe, msk);
	LIB_importMPK(abe, mpk);

	LIB_keygen(abe, "student|math", "key");
	char* alice_key = LIB_exportUserKey(abe, "key");
	LIB_keygen(abe, "student|CS", "key");
	char* bob_key = LIB_exportUserKey(abe, "key");
	LIB_keygen(abe, "student|EE", "key");
	char* carol_key = LIB_exportUserKey(abe,"key");

	// bob_key := abe.ExportUserKey("key_bob")
	// carol_key := abe.ExportUserKey("key_carol")

	auto abe1 = LIB_NewABE("CP-ABE");
	LIB_importMPK(abe1, mpk);
	LIB_importUserKey(abe1, alice_key);

	char* data = "hello world";

	char* ct = LIB_encrypt(abe1, "(student) and (math or EE)", data);

	//pt := abe.decrypt("key_alice", ct)

	auto abe2 = LIB_NewABE("CP-ABE");
	LIB_importMPK(abe2, mpk);
	// akey := abe2.ImportUserKey(alice_key)
	// pt := abe2.Decrypt(akey, ct)
	char* pt = LIB_ImportAndDecrypt(abe2, alice_key, ct);

	if (strcmp(pt,data) == 0) {
		printf("Alice Decrypt Successful pt = %v \n", pt);
	} else {
		printf("Alice Fail to decrypt");
	}
	return 0;

}