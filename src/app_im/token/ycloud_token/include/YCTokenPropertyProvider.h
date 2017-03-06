/*
 * YCTokenPropertyProvider.h
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENPROPERTYPROVIDER_H_
#define YCTOKENPROPERTYPROVIDER_H_

#include <list>
//#include <inttypes.h>
#include <stdint.h>
#include "YCTokenCommon.h"
#include "YCTokenExtendProperty.h"
#include "YCTokenException.h"

namespace yctoken {

#define MAX_EXTEND_PROP_NAME_LEN 255
#define MAX_EXTEND_PROP_VALUE_LEN 32000
#define MAX_APP_KEY_LEN 1280

#define FIXED_PROPS_LEN 26

class YCLOUD_TOKEN_API YCTokenPropertyProvider
{
public:
	/* parameter:
	 * tokenExpireTime : token expire time since apply time, unit is second.
	 */
	YCTokenPropertyProvider(uint32_t& appKey,uint32_t& tokenExpireTime);

	virtual ~YCTokenPropertyProvider();


	template<typename T>
	YCTokenPropertyProvider& addTokenExtendProperty(std::string& propName,T& propValue)
	{
		if(propName.size()> MAX_EXTEND_PROP_NAME_LEN){
			throw YCTokenException(TOKEN_EXTENDPROP_NAMELEN_OVERFLOW_ECODE);
		}
		YCTokenExtendProperty<T> *pExtendProp = new YCTokenExtendProperty<T>(propName,propValue);
		extendPropsSerializerList.push_back(pExtendProp);
		return *this;
	}

	/*
	void populateExtendProperty()
	{

	}

	std::list<YcTokenPropertySerializable*> getExtendProps()
	{

	 	 return extendPropsSerializerList;
	}
	*/

private:
    uint32_t _appKey;
    //unit:second
    uint32_t _tokenExpireTime;
    std::list<YcTokenPropertySerializable*> extendPropsSerializerList;

private:
    uint16_t fixedPropertiesLen();
    uint16_t extendPropertiesLen();
    void buildFixedPropertiesPart(std::ostringstream& oss,uint16_t& secretVersion);
    void buildExtendPropertiesPart(std::ostringstream& oss);

    friend class YCTokenBuilder;
};

} /* namespace yctoken */

#endif /* YCTOKENPROPERTYPROVIDER_H_ */
