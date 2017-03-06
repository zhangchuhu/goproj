/*
 * hmac_sha1.h
 *
 *  Created on: Oct 28, 2014
 *      Author: wanggb
 */

#ifndef HMAC_SHA1_H_
#define HMAC_SHA1_H_
namespace yctoken {

	void ycloud_hmac_sha1(const unsigned char *text, int text_len,
			   const unsigned char *key, int key_len,
			   unsigned char *digest);

}
#endif /* HMAC_SHA1_H_ */
