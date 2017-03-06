#include "TokenValidateSvr.hpp"
#include "TokenValidateSvr.h"
#include <string>
#include <iostream>
#include <stdint.h>


TokenValidateSvrPtr Init() {
	TokenValidateSvr * ret = new TokenValidateSvr();
	return (void*)ret;
}

void Free(TokenValidateSvrPtr f) {
	TokenValidateSvr * foo = (TokenValidateSvr*)f;
	delete foo;
}

int64_t Validate(TokenValidateSvrPtr f, char *pToken, int len) {
    std::string token(pToken, len);
	TokenValidateSvr * foo = (TokenValidateSvr*)f;
    return foo->validate(token);
}
