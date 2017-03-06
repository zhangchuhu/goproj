/*
 * YCTokenException.h
 *
 *  Created on: Nov 16, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENEXCEPTION_H_
#define YCTOKENEXCEPTION_H_

#include <exception>
//#include <inttypes.h>
#include <stdint.h>

#include "YCTokenCommon.h"

namespace yctoken {

extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_PACKET_LEN_OVERFLOW_ECODE;
extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_PACKET_BADFORMAT_ECODE;
extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_PACKET_BADAPPKEY_ECODE;

extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_INVALIDATE_ECODE;
extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_EXPIRE_ECODE;

extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_EXTENDPROP_BADTYPE_ECODE;
extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_EXTENDPROP_NAMELEN_OVERFLOW_ECODE;
extern YCLOUD_TOKEN_API  const uint16_t  TOKEN_EXTENDPROP_NAME_DUPLICATE_ECODE;

class YCLOUD_TOKEN_API YCTokenException: public std::exception {
public:
	YCTokenException(const uint16_t& errorCode);
	virtual ~YCTokenException()  throw();
	const uint16_t& errorCode();
	const char* what() const throw();
private:
	uint16_t _errorCode;

public:

};

} /* namespace yctoken */

#endif /* YCTOKENEXCEPTION_H_ */
