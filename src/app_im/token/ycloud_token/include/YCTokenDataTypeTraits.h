/*
 * YCTokenDataTypeTraits.h
 *
 *  Created on: Nov 15, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENDATATYPETRAITS_H_
#define YCTOKENDATATYPETRAITS_H_
#include <stdint.h>
#include "YCTokenUtil.h"

namespace yctoken {
template <typename T> struct DataTypeTraits;
    template<>
    struct DataTypeTraits<uint8_t>
    {
    	enum {TypeId = 1};
    	static bool convert(uint8_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
		{
    		if(valueLen != sizeof(uint8_t) || TypeId != type){
    			return false;
    		}
    		toValue = *(uint8_t*)fromRaw;
    		return true;
		}
    };

    template<>
    struct DataTypeTraits<uint16_t>
    {
      	enum {TypeId = 2};
      	static bool convert(uint16_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
      	{
      		if(valueLen != sizeof(uint16_t) || TypeId != type){
      			return false;
      		}
      		little_end_to_host(toValue,(uint16_t*)fromRaw);
      		return true;
      	}
    };

    template<>
    struct DataTypeTraits<uint32_t>
    {
    	enum {TypeId = 3};
    	static bool convert(uint32_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(uint32_t) || TypeId != type ){
    			return false;
    		}
    		little_end_to_host(toValue,(uint32_t*)fromRaw);
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<uint64_t>
    {
    	enum {TypeId = 4};
    	static bool convert(uint64_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(uint64_t) || TypeId != type){
    			return false;
    		}
    		little_end_to_host(toValue,(uint64_t*)fromRaw);
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<int8_t>
    {
    	enum {TypeId = 5};
    	static bool convert(int8_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(int8_t) || TypeId != type ){
    			return false;
    		}
    		toValue = *(int8_t*)fromRaw;
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<int16_t>
    {
    	enum {TypeId = 6};
    	static bool convert(int16_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(int16_t) || TypeId != type ){
    			return false;
    		}
    		little_end_to_host(toValue,(int16_t*)fromRaw);
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<int32_t>
    {
    	enum {TypeId = 7};
    	static bool convert(int32_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(int32_t) || TypeId != type ){
    			return false;
    		}
    		little_end_to_host(toValue,(int32_t*)fromRaw);
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<int64_t>
    {
    	enum {TypeId = 8};
    	static bool convert(int64_t& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(valueLen != sizeof(int64_t) || TypeId != type ){
    			return false;
    		}
    		little_end_to_host(toValue,(int64_t*)fromRaw);
    		return true;
    	}
    };

    template<>
    struct DataTypeTraits<std::string>
    {
    	enum {TypeId = 9};
    	static bool convert(std::string& toValue,void*& fromRaw,uint16_t& valueLen,uint8_t& type)
    	{
    		if(TypeId != type ){
    			return false;
    		}
    		if(valueLen == 0){
    			toValue = "";
    			return true;
    		}
    		toValue.assign((char *)fromRaw,valueLen);
    		return true;
    	}
    };

}

#endif /* YCTOKENDATATYPETRAITS_H_ */
