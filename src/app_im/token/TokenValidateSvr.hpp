#ifndef _MY_PACKAGE_FOO_HPP_
#define _MY_PACKAGE_FOO_HPP_

#include <string>
#include <stdint.h>

class TokenValidateSvr {
public:
	TokenValidateSvr() {};
	~TokenValidateSvr(){};
    int64_t validate(std::string t);
};

#endif
