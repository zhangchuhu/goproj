#ifndef _MY_PACKAGE_FOO_H_
#define _MY_PACKAGE_FOO_H_

#ifdef __cplusplus
extern "C" {
#endif
#include <stdint.h>

	typedef void* TokenValidateSvrPtr;
	TokenValidateSvrPtr Init(void);
	void Free(TokenValidateSvrPtr);
    int64_t Validate(TokenValidateSvrPtr, char *pToken, int len);

#ifdef __cplusplus
}
#endif

#endif
