#include <iostream>
#include "TokenValidateSvr.hpp"
#include <string>
#include "YCTokenAppSecretProvider.h"
#include "YCTokenPropertyProvider.h"
#include "YCToken.h"
#include "YCTokenBuilder.h"
#include "Base64.h"

using namespace std;
//c++ sdk的名称空间为yctoken
using namespace yctoken;
//实现YCloudAppSecretProvider,只有一个getAppsecret方法
class MyAppSecretProvider:public YCloudAppSecretProvider
{
public:
    MyAppSecretProvider()
    {
    }
    /*根据appKey返回有效的密钥map，map的key为密钥版本号，值为密钥。这个 方法在每次生成token时均要调用，要保证返回的token是没有过期的*/
    std::map<uint16_t,std::string> getAppsecret(uint32_t& appKey)
    {
        /*此处只是示例，实际应用中应根据密钥的存储和加载&刷新方案来获取有效密钥*/
        std::map<uint16_t,std::string> secretMap;
        secretMap.insert(pair<int,string>(1,"83e199e9_e2c"));
        return secretMap;
    }
};


int64_t TokenValidateSvr::validate(string tokenStr) {
    MyAppSecretProvider myAppsecretProvider;
    // 初始化YCTokenBuilder实例，该实例用于构造token
    YCTokenBuilder builder(myAppsecretProvider);
    //扩展属性名(uid)和属性值
    try
    {
        //cout << tokenStr <<endl;
        // 客户的appKey
        uint32_t myAppKey = 1436775582;
        uint32_t expire = 1484735992;

        uint32_t uid = 0;

        if(tokenStr.length() % 4 == 2)
        {
            tokenStr += "==";
        }
        else if (tokenStr.length() % 4 == 3)
        {
            tokenStr += "=";
        }

        string res;
        ZBase64::UrlDecode(tokenStr.c_str(), tokenStr.length(), res);
        YCToken *pToken = builder.validateTokenBytes(res);

        //cout<< "f" << endl;

        if(pToken->isExpired() )
        {
            return -2;
        }

        pToken->fetchExtendPropertyValue("UID", uid);
        //cout<< "uid: " << uid <<endl;

        return uid;
    }
    catch(YCTokenException &e)
    {
        //cout << e.what() <<endl;
        //cout<< "errorCode : " << e.errorCode() <<endl;
        return -1;
    }
}
