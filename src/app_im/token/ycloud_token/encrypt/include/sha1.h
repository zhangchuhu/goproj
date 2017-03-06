/*
 * sha1.h
 *
 *  Created on: Oct 27, 2014
 *      Author: wanggb
 */

#ifndef SHA1_H_
#define SHA1_H_

#include <stdint.h>

namespace yctoken {
	#define SHA1_CTX		librad_SHA1_CTX
	#define SHA1Transform	librad_SHA1Transform
	#define SHA1Init		librad_SHA1Init
	#define SHA1Update		librad_SHA1Update
	#define SHA1Final       librad_SHA1Final

	typedef struct {
		uint32_t state[5];
		uint32_t count[2];
		uint8_t buffer[64];
	} SHA1_CTX;

	void SHA1Transform(uint32_t state[5], const uint8_t buffer[64]);
	void SHA1Init(SHA1_CTX* context);
	void SHA1Update(SHA1_CTX* context, const uint8_t* data, unsigned int len);
	void SHA1Final(uint8_t digest[20], SHA1_CTX* context);

	/*
	 * this version implements a raw SHA1 transform, no length is appended,
	 * nor any 128s out to the block size.
	 */
	void SHA1FinalNoLen(uint8_t digest[20], SHA1_CTX* context);

	/*
	 * FIPS 186-2 PRF based upon SHA1.
	 */
	extern void fips186_2prf(uint8_t mk[20], uint8_t finalkey[160]);

}

#endif /* SHA1_H_ */
