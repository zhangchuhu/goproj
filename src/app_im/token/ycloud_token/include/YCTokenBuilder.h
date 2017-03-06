/*
 * YCTokenBuilder.h
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENBUILDER_H_
#define YCTOKENBUILDER_H_

#include <string>
#include <stdint.h>
#include "YCTokenCommon.h"
#include "YCTokenPropertyProvider.h"
#include "YCTokenAppSecretProvider.h"
#include "YCToken.h"
#include "YCTokenException.h"
#include "YCTokenInvalidateException.h"
#include "YCTokenExpireException.h"

namespace yctoken {

const uint16_t DIGEST_LEN = 20;

#define MIN_TOKEN_LEN (50)
#define MAX_TOKEN_LEN (32000)
/*
 * thread safe,can use singleton instance
 * */

class YCLOUD_TOKEN_API YCTokenBuilder {
public:
	YCTokenBuilder(YCloudAppSecretProvider& secretProvider);
	~YCTokenBuilder();
	/*for third party application to build token*/
	std::string buildBinaryToken(YCTokenPropertyProvider& propProvider) throw(YCTokenException);
	/*for duowan service to validate token*/
	YCToken* validateTokenBytes(const std::string& binaryToken)  throw(YCTokenInvalidateException,YCTokenExpireException,YCTokenException);
private:
	YCloudAppSecretProvider& _secretProvider;
private:
	void buildBinaryTokenPacketHead(std::ostringstream& oss,YCTokenPropertyProvider& propProvider);
	void buildBinaryTokenPacketFixedProperties(std::ostringstream& oss,YCTokenPropertyProvider& propProvider,uint16_t& secretVersion);
	void buildBinaryTokenPacketExtendProperties(std::ostringstream& oss,YCTokenPropertyProvider& propProvider);
	void buildBinaryTokenPacketDigest(std::ostringstream& oss,std::string& secret);
	void getLatestSecret(uint32_t& appKey,uint16_t& secretVersion,std::string& secret);
	bool validateExtendPropValueLen(uint8_t& dataType,uint16_t& valueLen);
};

} /* namespace yctoken */

#endif /* YCTOKENBUILDER_H_ */
