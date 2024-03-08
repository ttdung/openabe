#pragma once
#ifdef __cplusplus
extern "C" {
#endif

void* LIB_NewABE(char* name);
void LIB_generateParams(void* abe);
char* LIB_encrypt(void* abe, char* att, char* pt);
char* LIB_decrypt(void* abe, char* key, char* ct);
int LIB_keygen(void* abe, char* att, char* key);

#ifdef __cplusplus
}  // extern "C"
#endif
