/*
 * YCToken.cpp
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */
#include <time.h>
#include "YCToken.h"

namespace yctoken {

YCToken::YCToken(uint32_t& appKey,uint64_t& expireTime,uint64_t& timeStamp):
		_appKey(appKey),_expireTime(expireTime),_timeStamp(timeStamp)
{

}

YCToken::~YCToken()
{
	//todo release pointer in map
	std::map<std::string,YCTokenExtendProperty<void*>* >::iterator iter;
	for(iter = _extendPropsMap.begin(); iter != _extendPropsMap.end();) {
		std::map<std::string,YCTokenExtendProperty<void*>* >::iterator temp = iter;
		iter++;
		delete temp->second;
	}
	_extendPropsMap.clear();
}

void YCToken::setDigest(std::string& digest)
{
	_digest = digest;
}

uint32_t YCToken::getAppKey()
{
	return _appKey;
}

uint64_t YCToken::getTimeStamp()
{
	return _timeStamp;
}

uint64_t YCToken::getExpireTime()
{
	return _expireTime;
}
std::string YCToken::getDigest()
{
	return _digest;
}

bool YCToken::isExpired()
{
	return _expireTime < (uint64_t)time(NULL);
}


} /* namespace yctoken */
