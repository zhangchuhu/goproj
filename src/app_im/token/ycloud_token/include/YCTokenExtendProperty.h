/*
 * YCTokenExtendProperty.h
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENEXTENDPROPERTY_H_
#define YCTOKENEXTENDPROPERTY_H_

//#include <inttypes.h>
#include <stdint.h>
#include <string.h>
#include <sstream>
//#include <iostream>
#include "YCTokenCommon.h"
#include "YCTokenUtil.h"
#include "YCTokenDataTypeTraits.h"
#include "YCTokenException.h"

namespace yctoken {

	class YcTokenPropertySerializable
	{
		public:
			virtual void serialize(std::ostringstream& oss) = 0;
			virtual uint16_t binaryLength() =0 ;
			virtual ~YcTokenPropertySerializable()
			{

			}
			//virtual void deSerialize() =0;
	};

	template<typename value_type>
	class YCTokenExtendProperty:public YcTokenPropertySerializable {
	public:
		YCTokenExtendProperty(const std::string& name,const value_type& value):_name(name),_value(value)
		{
			this->_valueLen = sizeof(value_type);
		}

		~YCTokenExtendProperty()
		{

		}

		void serialize(std::ostringstream& oss)
		{
			uint8_t   nameLen = (uint8_t)_name.size();
			oss.write((const char*)&nameLen,sizeof(uint8_t));

			oss.write(_name.data(),_name.size());

			uint8_t typeId = (uint8_t)DataTypeTraits<value_type>::TypeId;
			oss.write((const char*)&typeId,sizeof(typeId));

			uint16_t little_end_valueLen = 0;
			host_to_little_end(little_end_valueLen,&_valueLen);
			oss.write((const char*)&little_end_valueLen,sizeof(little_end_valueLen));

			value_type little_end_value;
			host_to_little_end(little_end_value,&_value);
			oss.write((const char*)&little_end_value,sizeof(little_end_value));
		}

		uint16_t binaryLength()
		{
			// name length 1 byte ,data type 1 byte
			return 1 +_name.size() + 1 + 2 + sizeof(value_type);
		}
		/*
		YCTokenPara(uint8_t id,void* value)
		{
			this->id= id;
			this->value= value;
			this->valueLen = sizeof(valueLen);
		}
        */

		std::string const& getName()
		{
			return _name;
		}

		void getValue()
		{
			//return _value;
		}


		/*
		uint16_t getValueLen()
		{
			return _valueLen;
		}
        */

	private:
		std::string _name;
		value_type _value;
		uint16_t _valueLen;
	};

	template<>
	class YCTokenExtendProperty<void*> {
	public:
		YCTokenExtendProperty(std::string name,uint8_t& valueDataType,const void* value,uint16_t& valueLen):_name(name),_valueDataType(valueDataType),_valueLen(valueLen)
		{
			_value = NULL;
			if(valueLen>0){
				this->_value = new char[valueLen];
				memcpy(this->_value,value,valueLen);
			}
		}

		~YCTokenExtendProperty()
		{
			if(NULL != _value){
				delete [](char*)_value;
			}
		}

		std::string getName()
		{
			return _name;
		}

		template<typename T>
		void getValue(T& value)  throw(YCTokenException)
		{
			if(DataTypeTraits<T>::TypeId != _valueDataType){
				// throw exception;
				throw YCTokenException(TOKEN_EXTENDPROP_BADTYPE_ECODE);
			}
			bool ret = DataTypeTraits<T>::convert(value,_value,_valueLen,_valueDataType);
			if(!ret){
				throw YCTokenException(TOKEN_EXTENDPROP_BADTYPE_ECODE);
			}
		}

		private:
			std::string _name;
			uint8_t _valueDataType;
			void* _value;
			uint16_t _valueLen;
	};


	template<>
		class YCTokenExtendProperty<std::string>:public YcTokenPropertySerializable {
		public:
		YCTokenExtendProperty(const std::string& name,const std::string& value):_name(name),_value(value)
		{
			this->_valueLen = value.size();
		}

		void serialize(std::ostringstream& oss)
		{
			uint8_t   nameLen = (uint8_t)_name.size();
			oss.write((const char*)&nameLen,sizeof(uint8_t));

			oss.write(_name.data(),_name.size());

			uint8_t typeId = (uint8_t)DataTypeTraits<std::string>::TypeId;
			oss.write((const char*)&typeId,sizeof(typeId));

			uint16_t little_end_valueLen = 0;
			host_to_little_end(little_end_valueLen,&_valueLen);

			oss.write((const char*)&little_end_valueLen,sizeof(little_end_valueLen));
			oss.write(_value.data(),_valueLen);
		}

		uint16_t binaryLength()
		{
			// name length 1 byte ,data type 1 byte
			return 1 +_name.size() + 1 + 2 + _value.size();
		}


		private:
			std::string _name;
		    std::string _value;
			uint16_t _valueLen;
		};

} /* namespace yctoken */


#endif /* YCTOKENEXTENDPROPERTY_H_ */
