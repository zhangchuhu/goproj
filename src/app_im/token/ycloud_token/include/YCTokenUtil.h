/*
 * YCTokenUtil.h
 *
 *  Created on: Nov 15, 2014
 *      Author: wanggb
 */

#ifndef YCTOKENUTIL_H_
#define YCTOKENUTIL_H_

#include <stdint.h>

namespace yctoken {
	const uint16_t us_flag = 1;
	//is little end byte order
	const bool little_end_flag = *((uint8_t*)&us_flag) == 1;

	template<typename T>
	void little_end_to_host(T& to,T* from)
	{
		uint8_t byteLen = sizeof(T);

		if(little_end_flag){
			to = *from;
			return;
		}
		else{
			char* to_char =  (char*)&to;
		    //char* from_char = &from;
			for(int i=0;i<byteLen;i++){
				to_char[i] = from[byteLen-i-1];
			}
			return;
		}
	}


	template<typename T>
	void host_to_little_end(T& to,T* from)
	{
		little_end_to_host(to,from);
	}

}
#endif /* YCTOKENUTIL_H_ */
