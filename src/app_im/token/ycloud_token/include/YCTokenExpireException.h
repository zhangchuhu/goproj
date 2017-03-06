/*
 * YCTokenExpireException.h
 *
 *  Created on: Nov 22, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENEXPIREEXCEPTION_H_
#define YCTOKENEXPIREEXCEPTION_H_
#include <stdint.h>
#include "YCTokenCommon.h"
#include "YCTokenException.h"

namespace yctoken {
	class YCLOUD_TOKEN_API YCTokenExpireException: public yctoken::YCTokenException {
	public:
		YCTokenExpireException(uint32_t& appKey,uint16_t& secretVersion,uint64_t& expire);
		virtual ~YCTokenExpireException()   throw();
		uint32_t getAppKey();
		uint16_t getSecretVersion();
		uint64_t getExpire();
	private:
		uint32_t _appKey;
		uint16_t _secretVersion;
		uint64_t _expire;
	};
}
#endif /* YCTOKENEXPIREEXCEPTION_H_ */
