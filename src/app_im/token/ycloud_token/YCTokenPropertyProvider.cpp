/*
 * YCTokenPropertyProvider.cpp
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */
//#include <stdlib.h>
#include <cstdlib>
#include <time.h>
#include "YCTokenPropertyProvider.h"

namespace yctoken {

YCTokenPropertyProvider::YCTokenPropertyProvider(uint32_t& appKey,uint32_t& tokenExpireTime):
		_appKey(appKey),_tokenExpireTime(tokenExpireTime)
{

}

YCTokenPropertyProvider::~YCTokenPropertyProvider()
{
	//delete pointer in extendPropsSerializerList
	while(!extendPropsSerializerList.empty())
	{
	    delete extendPropsSerializerList.front();
	    extendPropsSerializerList.pop_front();
	}

}

uint16_t YCTokenPropertyProvider::fixedPropertiesLen()
{
	// delete :appkey len 2 byte;expireDate 8 byte,timestamp 8 byte;random(nonce) 4 byte
	// add: appkey 4 byte; expireDate 8 byte;timestamp 8 byte;random(nonce) 4 byte
	//return 2 + _appKey.size() + 8 + 8 + 4;
	return FIXED_PROPS_LEN;
}

uint16_t YCTokenPropertyProvider::extendPropertiesLen()
{
	uint16_t totalLen = 0;
	for (std::list<YcTokenPropertySerializable*>::iterator it=extendPropsSerializerList.begin(); it != extendPropsSerializerList.end(); ++it){
		totalLen += (*it)->binaryLength();
	}
	return totalLen;
}

void YCTokenPropertyProvider::buildFixedPropertiesPart(std::ostringstream& oss,uint16_t& secretVersion)
{
	time_t now = time(NULL);
	uint64_t expireTime = now + _tokenExpireTime;
	uint64_t timeStamp = now;
	//should init call srandom(time(NULL));
    #ifndef _WINDOWS
		uint32_t randNum = random();
	#else
		uint32_t randNum = rand();
	#endif

	uint32_t little_end_appKey = 0;
	uint16_t little_end_version = 0;
	uint64_t little_end_expire = 0;
	uint64_t little_end_timeStamp = 0;
	uint32_t littele_end_randNum = 0;

	host_to_little_end(little_end_appKey,&_appKey);
	host_to_little_end(little_end_version,&secretVersion);
	host_to_little_end(little_end_expire,&expireTime);
	host_to_little_end(little_end_timeStamp,&timeStamp);
	host_to_little_end(littele_end_randNum,&randNum);

	oss.write((char*)&little_end_appKey,sizeof(little_end_appKey));
	oss.write((char*)&little_end_version,sizeof(little_end_version));
	oss.write((char*)&little_end_expire,sizeof(little_end_expire));
	oss.write((char*)&little_end_timeStamp,sizeof(little_end_timeStamp));
	oss.write((char*)&littele_end_randNum,sizeof(littele_end_randNum));
	return;
}

void YCTokenPropertyProvider::buildExtendPropertiesPart(std::ostringstream& oss)
{
	for(std::list<YcTokenPropertySerializable*>::iterator it=extendPropsSerializerList.begin(); it != extendPropsSerializerList.end(); ++it){
		(*it)->serialize(oss);
	}
	return;
}

} /* namespace yctoken */
