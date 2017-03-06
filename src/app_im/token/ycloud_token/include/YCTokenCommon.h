/*
 * YCTokenCommon.h
 *
 *  Created on: Feb 12, 2015
 *      Author: wanggb
 */
#ifndef YCTOKENCOMMON_H_
#define YCTOKENCOMMON_H_

#ifdef WIN32
	#ifdef YCLOUD_TOKEN_EXPORTS
		#define YCLOUD_TOKEN_API __declspec(dllexport)
	#else
		#define YCLOUD_TOKEN_API __declspec(dllimport)
	#endif
#else
    #define YCLOUD_TOKEN_API
#endif

#endif /* YCTOKENCOMMON_H_ */