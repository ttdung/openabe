#include <iostream>

#include "lib_bridge.h"
#include "openabe/openabe.h"
#include <openabe/zsymcrypto.h>

using namespace std;
using namespace oabe;
using namespace oabe::crypto;

// make chean

void* LIB_InitializeOpenABE() {

    InitializeOpenABE();
    return 0;
}

void* LIB_ShutdownOpenABE() {

    ShutdownOpenABE();
    return 0;
}

void* LIB_NewABE(char* name) {

    auto cpabe = new OpenABECryptoContext(name);
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

char* LIB_exportMSK(void* abe) {

    std::string sk;
    AsAbe(abe)->exportSecretParams(sk);

    char *skc = new char[sk.size() + 1];
    std::strcpy(skc, sk.c_str());

    return skc;
}

int LIB_importMSK(void* abe, char* msk) {

    AsAbe(abe)->importSecretParams(std::string((char*)msk));
    return 0;
}


char* LIB_exportMPK(void* abe) {

    std::string pk;
    AsAbe(abe)->exportPublicParams(pk);

    char *pkc = new char[pk.size() + 1];
    std::strcpy(pkc, pk.c_str());

    return pkc;
}

int LIB_importMPK(void* abe, char* mpk) {

    AsAbe(abe)->importPublicParams(std::string((char*)mpk));
    return 0;
}

char* LIB_exportUserKey(void* abe, char* key) {

    std::string k;
    AsAbe(abe)->exportUserKey(std::string((char*)key), k);

    char *kc = new char[k.size() + 1];
    std::strcpy(kc, k.c_str());

    return kc;
}

char* LIB_importUserKey(void* abe,char* index, char* key) {

    std::string k(index);
    AsAbe(abe)->importUserKey(std::string((char*)index), std::string((char*)key));
    
    char *skc = new char[k.size() + 1];
    std::strcpy(skc, k.c_str());

    return skc;
}

char* LIB_ImportAndDecrypt(void* abe,char* key, char* ct) {

    std::string k;
    AsAbe(abe)->importUserKey(k, std::string((char*)key));
    
    std::string pt1;
    AsAbe(abe)->decrypt(k, std::string((char*)ct), pt1);

    char *pt = new char[pt1.size() + 1];
    std::strcpy(pt, pt1.c_str());

    return pt;

}

char* LIB_KeygenAndDecrypt(void* abe, char* att, char* ct) {

    std::string key;
    AsAbe(abe)->keygen(std::string((char*)att), key);

    std::string pt1;
    AsAbe(abe)->decrypt(key, std::string((char*)ct), pt1);

    char *pt = new char[pt1.size() + 1];
    std::strcpy(pt, pt1.c_str());

    return pt;
}