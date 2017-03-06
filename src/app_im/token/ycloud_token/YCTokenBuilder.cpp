/*
 * YCTokenBuilder.cpp
 *
 *  Created on: Nov 14, 2014
 *      Author: wanggb
 */
#include <memory>
#include <time.h>
#include <string.h>
#include "YCTokenBuilder.h"
#include "hmac_sha1.h"

namespace yctoken {

YCTokenBuilder::YCTokenBuilder(YCloudAppSecretProvider& secretProvider):_secretProvider(secretProvider) {
	// TODO Auto-generated constructor stub

}

YCTokenBuilder::~YCTokenBuilder() {
	// TODO Auto-generated destructor stub
}

std::string YCTokenBuilder::buildBinaryToken(YCTokenPropertyProvider& propProvider)  throw(YCTokenException)
{
	std::ostringstream oss;
	uint16_t latestSecretVersion=0;
	std::string  secret;

	getLatestSecret(propProvider._appKey,latestSecretVersion,secret);
	if(secret.empty()){
		//throw
		return "";
	}

	buildBinaryTokenPacketHead(oss,propProvider);
	buildBinaryTokenPacketFixedProperties(oss,propProvider,latestSecretVersion);
	buildBinaryTokenPacketExtendProperties(oss,propProvider);
	buildBinaryTokenPacketDigest(oss,secret);
	return oss.str();
}

YCToken* YCTokenBuilder::validateTokenBytes(const std::string& binaryToken)  throw(YCTokenInvalidateException,YCTokenExpireException,YCTokenException)
{
	if(binaryToken.size()< MIN_TOKEN_LEN || binaryToken.size() > MAX_TOKEN_LEN){
		//todo throw exception
		throw YCTokenException(TOKEN_PACKET_LEN_OVERFLOW_ECODE);
	}
	//parse packet head
	const char* p = binaryToken.data();
	//wangguobo(2014-11-24):byte align problem(core dump) ,should use memcpy safely.
	//uint16_t little_end_totalLen = *(uint16_t*)p;
	uint16_t little_end_totalLen = 0;
	uint16_t totalLen = 0;
	memcpy(&little_end_totalLen,p,sizeof(little_end_totalLen));

	little_end_to_host(totalLen,&little_end_totalLen);
	if(totalLen != binaryToken.size()){
		//todo throw exception
		throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
	}

	p += 2;

	//uint16_t little_end_extendPropsLen = *(uint16_t*)p;
	uint16_t little_end_extendPropsLen = 0;
	memcpy(&little_end_extendPropsLen,p,sizeof(little_end_extendPropsLen));
	uint16_t extendPropsLen = 0;
	little_end_to_host(extendPropsLen,&little_end_extendPropsLen);

	//packet head 4 byte
	if(totalLen != 4+ FIXED_PROPS_LEN + extendPropsLen + DIGEST_LEN){
		//todo throw exception
		throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
	}

	p += 2;

	//parse fixed property and validate exprie time
	//uint32_t little_end_appKey = *(uint32_t*)p;
	uint32_t little_end_appKey = 0;
	memcpy(&little_end_appKey,p,sizeof(little_end_appKey));
	uint32_t appKey = 0;
	little_end_to_host(appKey,&little_end_appKey);

	p += sizeof(appKey);

	//uint16_t little_end_version = *(uint16_t*)p;
	uint16_t little_end_version = 0;
	memcpy(&little_end_version,p,sizeof(little_end_version));
	uint16_t secretVersion = 0;
	little_end_to_host(secretVersion,&little_end_version);

	p += sizeof(secretVersion);

	//uint64_t little_end_expireTime = *(uint64_t*)p;
	uint64_t little_end_expireTime = 0;
	memcpy(&little_end_expireTime,p,sizeof(little_end_expireTime));
	uint64_t expireTime = 0;
	little_end_to_host(expireTime,&little_end_expireTime);

	p += 8;

	//uint64_t little_end_timeStamp = *(uint64_t*)p;
	uint64_t little_end_timeStamp = 0;
	memcpy(&little_end_timeStamp,p,sizeof(little_end_timeStamp));
	uint64_t timeStamp = 0;
	little_end_to_host(timeStamp,&little_end_timeStamp);
	p += 8;


	//skip 4 byte random num
	p += 4;

	//parse extend property
	//checkpoint during parse extend property
	uint16_t checkPoint =0;
	uint8_t extendPropNameLen;
	std::string extendPropName;
	uint8_t extendPropDataType =0;
	uint16_t little_end_extendPropLen = 0;
	uint16_t extendPropLen = 0;

	//exception safe ,use smart pointer;
	YCToken* token = new YCToken(appKey,expireTime,timeStamp);
	std::auto_ptr<YCToken> auto_Token(token);

	while(checkPoint != extendPropsLen){
		//extendPropNameLen = *(uint8_t*)p;
		memcpy(&extendPropNameLen,p,sizeof(extendPropNameLen));
		if(0 == extendPropNameLen){
			//todo throw exception
			throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
		}
		checkPoint += (1+extendPropNameLen);
		if(checkPoint >= extendPropsLen){
			//todo throw exception
			throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
		}
		p +=1 ;
		extendPropName.assign(p,extendPropNameLen);
		p += extendPropNameLen;

		checkPoint += 3;
		if(checkPoint > extendPropsLen){
			//todo throw exception
			throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
		}

		//extendPropDataType = *(uint8_t*)p;
		memcpy(&extendPropDataType,p,sizeof(extendPropDataType));
		p += 1;

		//little_end_extendPropLen = *(uint16_t*)p;
		memcpy(&little_end_extendPropLen,p,sizeof(little_end_extendPropLen));
		little_end_to_host(extendPropLen,&little_end_extendPropLen);

		checkPoint += extendPropLen;
		if(checkPoint > extendPropsLen || !validateExtendPropValueLen(extendPropDataType,extendPropLen)){
			//todo throw exception
			throw YCTokenException(TOKEN_PACKET_BADFORMAT_ECODE);
		}

		p +=2;

		if(token->_extendPropsMap.find(extendPropName) != token->_extendPropsMap.end()){
			//todo throw exception
			throw YCTokenException(TOKEN_EXTENDPROP_NAME_DUPLICATE_ECODE);
		}
        token->_extendPropsMap[extendPropName] =  new YCTokenExtendProperty<void*>(extendPropName,extendPropDataType,p,extendPropLen);

		p += extendPropLen;
	}

	//validate digest
	std::map<uint16_t,std::string> secretMap =  _secretProvider.getAppsecret(appKey);
	std::map<uint16_t,std::string>::iterator it;
	it=secretMap.find(secretVersion);
	if(it == secretMap.end()){
		throw YCTokenException(TOKEN_PACKET_BADAPPKEY_ECODE);
	}

	std::string appSecret = it->second;

	uint8_t myDigest[DIGEST_LEN];
	ycloud_hmac_sha1((const unsigned char*)binaryToken.data(),totalLen-DIGEST_LEN,(const unsigned char*)appSecret.data(),appSecret.size(),myDigest);

	std::string peerDigest(binaryToken,totalLen-DIGEST_LEN,DIGEST_LEN);
	//compare two digest
	if(0 != peerDigest.compare(std::string((char*)myDigest,DIGEST_LEN))){
		throw YCTokenInvalidateException(appKey,secretVersion);
	}

	//check whether expire
	if(expireTime < (uint64_t)time(NULL)){
		throw YCTokenExpireException(appKey,secretVersion,expireTime);
	}

	token->setDigest(peerDigest);
	return auto_Token.release();
}

/*
 *     private member
 *
 * */
void YCTokenBuilder::buildBinaryTokenPacketHead(std::ostringstream& oss,YCTokenPropertyProvider& propProvider)
{
	uint16_t fixedLen = propProvider.fixedPropertiesLen();
	uint16_t extendLen = propProvider.extendPropertiesLen();

	//total len: 2 byte , extend property len: 2 byte
	uint16_t totalLen = 2+ 2 + fixedLen + extendLen + DIGEST_LEN;
	uint16_t little_end_totalLen = 0;
	uint16_t little_end_extendLen = 0;

	host_to_little_end(little_end_totalLen,&totalLen);
	host_to_little_end(little_end_extendLen,&extendLen);

	oss.write((char*)&little_end_totalLen,sizeof(little_end_totalLen));
	oss.write((char*)&little_end_extendLen,sizeof(little_end_extendLen));
	return;
}

void YCTokenBuilder::buildBinaryTokenPacketFixedProperties(std::ostringstream& oss,YCTokenPropertyProvider& propProvider,uint16_t& secretVersion)
{
	propProvider.buildFixedPropertiesPart(oss,secretVersion);
}

void YCTokenBuilder::buildBinaryTokenPacketExtendProperties(std::ostringstream& oss,YCTokenPropertyProvider& propProvider)
{
	propProvider.buildExtendPropertiesPart(oss);
}

void YCTokenBuilder::buildBinaryTokenPacketDigest(std::ostringstream& oss,std::string& secret)
{
	std::string rawData = oss.str();
	uint8_t digest[DIGEST_LEN];
	ycloud_hmac_sha1((const unsigned char*)rawData.data(),rawData.size(),(const unsigned char*)secret.data(),secret.size(), digest);
	oss.write((char*)digest,DIGEST_LEN);
}

void YCTokenBuilder::getLatestSecret(uint32_t& appKey,uint16_t& secretVersion,std::string& secret)
{
	std::map<uint16_t,std::string> secretMap =  _secretProvider.getAppsecret(appKey);
	if(secretMap.empty()){
		secret = "";
		secretVersion = 0;
		return ;
	}

	std::map<uint16_t,std::string>::iterator iter;
	uint16_t minVersion = 0;
	uint16_t maxVersion = 0;

	for(iter = secretMap.begin( );iter != secretMap.end( ); iter++ ){
		if(0 == minVersion ){
			minVersion = iter->first;
		}
		else if(minVersion > iter->first){
			minVersion = iter->first;
		}

		if(maxVersion < iter->first){
			maxVersion = iter->first;
		}
	}

	if(maxVersion-minVersion > 60000){
		secretVersion = minVersion;
	}
	else{
		secretVersion = maxVersion;
	}

	secret = secretMap[secretVersion];
	return;
}

bool YCTokenBuilder::validateExtendPropValueLen(uint8_t& dataType,uint16_t& valueLen)
{
	switch(dataType)
	{
	//uint8_t
	case 1:
		return 1 == valueLen;
	//uint16_t
	case 2:
		return 2 == valueLen;
	//uint32_t
	case 3:
		return 4 == valueLen;
	//uint64_t
	case 4:
		return 8 == valueLen;
	//int8_t
	case 5:
		return 1 == valueLen;
	//int16_t
	case 6:
		return 2 == valueLen;
	//int32_t
	case 7:
		return 4 == valueLen;
	//int64_t
	case 8:
		return 8 == valueLen;
	default:
		return true;
	}
}
} /* namespace yctoken */
