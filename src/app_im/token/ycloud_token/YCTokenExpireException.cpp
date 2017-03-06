/*
 * YCTokenExpireException.cpp
 *
 *  Created on: Nov 22, 2014
 *      Author: wanggb
 */

#include "YCTokenExpireException.h"
namespace yctoken {
	YCTokenExpireException::YCTokenExpireException(uint32_t& appKey,uint16_t& secretVersion,uint64_t& expire):YCTokenException(TOKEN_EXPIRE_ECODE)
			,_appKey(appKey),_secretVersion(secretVersion),_expire(expire)
	{
		// TODO Auto-generated constructor stub

	}

	YCTokenExpireException::~YCTokenExpireException()   throw()
	{
		// TODO Auto-generated destructor stub
	}

	uint32_t YCTokenExpireException::getAppKey()
	{
		return _appKey;
	}

	uint16_t YCTokenExpireException::getSecretVersion()
	{
		return 	_secretVersion;
	}

	uint64_t YCTokenExpireException::getExpire()
	{
		return 	_expire;
	}

}
