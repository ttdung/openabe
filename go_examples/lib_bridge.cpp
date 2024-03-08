#include <iostream>

#include "lib_bridge.h"
#include "openabe/openabe.h"
#include <openabe/zsymcrypto.h>

using namespace std;
using namespace oabe;
using namespace oabe::crypto;

// make chean


void* LIB_NewABE(char* name) {

    // std::cout << "[c++ bridge] LIB_NewABE" << name << ")" << std::endl;
    
    InitializeOpenABE();

   auto cpabe = new OpenABECryptoContext(name);

//   std::cout << "[c++ bridge] LIB_NewABE(" << name << ") will return pointer "
            // << cpabe << std::endl;

  return cpabe;
}

// Utility function local to the bridge's implementation
OpenABECryptoContext* AsAbe(void* abe) { return reinterpret_cast<OpenABECryptoContext*>(abe); }

void LIB_generateParams(void* abe) {
    AsAbe(abe)->generateParams();
}

int LIB_keygen(void* abe, char* att, char* key) {

    AsAbe(abe)->keygen(std::string((char*)att), std::string((char*)key));
    return 0;
}

char* LIB_encrypt(void* abe, char* att, char* pt) {

    std::string ct1;
    AsAbe(abe)->encrypt(std::string((char*)att), std::string((char*)pt), ct1);

    char *ct = new char[ct1.size() + 1];
    std::strcpy(ct, ct1.c_str());

    return ct;
}

char* LIB_decrypt(void* abe, char* key, char* ct) {

    std::string pt1;
    AsAbe(abe)->decrypt(std::string((char*)key), std::string((char*)ct), pt1);

    char *pt = new char[pt1.size() + 1];
    std::strcpy(pt, pt1.c_str());

    return pt;
}