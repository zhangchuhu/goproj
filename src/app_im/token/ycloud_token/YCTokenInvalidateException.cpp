/*
 * YCTokenInvalidateException.cpp
 *
 *  Created on: Nov 22, 2014
 *      Author: wanggb
 */

#include "YCTokenInvalidateException.h"

namespace yctoken {
	YCTokenInvalidateException::YCTokenInvalidateException(uint32_t& appKey,uint16_t& secretVersion):YCTokenException(TOKEN_INVALIDATE_ECODE)
			,_appKey(appKey),_secretVersion(secretVersion)
	{
		// TODO Auto-generated constructor stub

	}

	YCTokenInvalidateException::~YCTokenInvalidateException()   throw()
	{
		// TODO Auto-generated destructor stub
	}

	uint32_t YCTokenInvalidateException::getAppKey()
	{
		return _appKey;
	}

	uint16_t YCTokenInvalidateException::getSecretVersion()
	{
		return 	_secretVersion;
	}
}
