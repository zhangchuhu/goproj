/*
 * YCTokenAppSecretProvider.h
 *
 *  Created on: Nov 16, 2014
 *      Author: wanggb
 */

#ifndef YCTOKEN_APPSECRETPROVIDER_H_
#define YCTOKEN_APPSECRETPROVIDER_H_

#include <string>
#include <stdint.h>
#include <map>
#include <stdint.h>
#include "YCTokenCommon.h"

namespace yctoken {

class YCLOUD_TOKEN_API YCloudAppSecretProvider {
public:
	YCloudAppSecretProvider();
	virtual ~YCloudAppSecretProvider();
	virtual std::map<uint16_t,std::string> getAppsecret(uint32_t& appKey) = 0;
};

} /* namespace yctoken */

#endif /* YCTOKEN_APPSECRETPROVIDER_H_ */
