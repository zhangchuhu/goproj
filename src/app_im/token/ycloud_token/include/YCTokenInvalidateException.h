/*
 * YCTokenInValidateException.h
 *
 *  Created on: Nov 22, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENINVALIDATEEXCEPTION_H_
#define YCTOKENINVALIDATEEXCEPTION_H_
#include <stdint.h>
#include "YCTokenCommon.h"
#include "YCTokenException.h"
namespace yctoken {
	class YCLOUD_TOKEN_API YCTokenInvalidateException: public yctoken::YCTokenException {
	public:
		YCTokenInvalidateException(uint32_t& appKey,uint16_t& secretVersion);
		virtual ~YCTokenInvalidateException()   throw();
		uint32_t getAppKey();
		uint16_t getSecretVersion();
	private:
		uint32_t _appKey;
		uint16_t _secretVersion;
	};
}
#endif /* YCTOKENINVALIDATEEXCEPTION_H_ */
