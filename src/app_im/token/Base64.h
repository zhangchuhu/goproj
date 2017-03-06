#pragma once
#include <string>

class ZBase64
{
public:

	static void UrlEncode(const char* Data, int DataByte, std::string& strEncode);

	static void UrlDecode(const char* Data, int DataByte, std::string& out);
};
