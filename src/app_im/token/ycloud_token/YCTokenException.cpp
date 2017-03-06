/*
 * YCTokenException.cpp
 *
 *  Created on: Nov 16, 2014
 *      Author: wanggb
 */

#include "YCTokenException.h"

namespace yctoken {

const uint16_t  TOKEN_PACKET_LEN_OVERFLOW_ECODE               = 1;
const uint16_t  TOKEN_PACKET_BADFORMAT_ECODE                  = 2;
const uint16_t  TOKEN_PACKET_BADAPPKEY_ECODE                  = 3;

const uint16_t  TOKEN_INVALIDATE_ECODE                        = 8;
const uint16_t  TOKEN_EXPIRE_ECODE                            = 9;

const uint16_t  TOKEN_EXTENDPROP_BADTYPE_ECODE                = 20;
const uint16_t  TOKEN_EXTENDPROP_NAMELEN_OVERFLOW_ECODE       = 21;
const uint16_t  TOKEN_EXTENDPROP_NAME_DUPLICATE_ECODE         = 22;

YCTokenException::YCTokenException(const uint16_t& errorCode):_errorCode(errorCode)
{
	// TODO Auto-generated constructor stub

}

YCTokenException::~YCTokenException()  throw()
{
	// TODO Auto-generated destructor stub
}

const uint16_t& YCTokenException::errorCode()
{
	return _errorCode;
}

const char* YCTokenException::what() const throw()
{
	return "ycloud token exception!";
}

} /* namespace yctoken */
