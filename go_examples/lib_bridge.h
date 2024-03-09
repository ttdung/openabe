#pragma once
#ifdef __cplusplus
extern "C" {
#endif

void* LIB_NewABE(char* name);
void LIB_generateParams(void* abe);
char* LIB_encrypt(void* abe, char* att, char* pt);
char* LIB_decrypt(void* abe, char* key, char* ct);
int LIB_keygen(void* abe, char* att, char* key);
char* LIB_importUserKey(void* abe,char* key);
char* LIB_exportUserKey(void* abe,char* key);
char* LIB_exportMPK(void* abe);
char* LIB_exportMSK(void* abe);
int LIB_importMPK(void* abe, char* mpk);
int LIB_importMSK(void* abe, char* msk);

#ifdef __cplusplus
}  // extern "C"
#endif
