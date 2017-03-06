/*
 * YCToken.h
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */

#ifndef YCTOKEN_H_
#define YCTOKEN_H_

#include <map>
#include <stdint.h>
#include "YCTokenCommon.h"
#include "YCTokenExtendProperty.h"
#include "YCTokenException.h"

namespace yctoken {

class YCLOUD_TOKEN_API YCToken {
public:
	virtual ~YCToken();

	template<typename value_type>
	bool fetchExtendPropertyValue(std::string propName,value_type& value)  throw(YCTokenException)
	{
		YCTokenExtendProperty<void*> *prop =  _extendPropsMap[propName];
		if(NULL == prop){
			return false;
		}

		prop->getValue(value);
		return true;
	}

	uint32_t getAppKey();
	uint64_t getTimeStamp();
	uint64_t getExpireTime();
	std::string getDigest();
	bool isExpired();
private:
	YCToken(uint32_t& appKey,uint64_t& expireTime,uint64_t& timeStamp);
	void setDigest(std::string& digest);
	uint32_t _appKey;
    uint64_t _expireTime;
    uint64_t _timeStamp;
    std::string _digest;

	std::map<std::string,YCTokenExtendProperty<void*>* > _extendPropsMap;

	friend class YCTokenBuilder;
};

} /* namespace yctoken */

#endif /* YCTOKEN_H_ */
