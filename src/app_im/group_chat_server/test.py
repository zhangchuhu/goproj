#!/usr/bin/env python
#-*- coding:utf8 -*-
import requests
import json
import sys

def Post(url, req, token):
    print "// 请求" 
    print json.dumps(req, indent=2)
    headers = {'token' : token}
    #r = requests.post('http://172.26.38.31:10067' + url, data=json.dumps(req), headers=headers)
    r = requests.post('http://106.38.255.199:9098' + url, data=json.dumps(req), headers=headers)
    print "// 响应"
    print r.text


reqSendGroupMsg = {
        'groupId' : 33,
        'msgId'  : '12b11fa1fd7606698b63b7d42',
        'msgType' : 1,
        'msgBody' : "fuck"
}

reqPullGroupMsg = {
        'groupId' : 1234,
        'fromSeqId' : 1,
        'toSeqId' : 3
}

if __name__ == '__main__':
    token = ""
    if len(sys.argv) != 1:
        print "please input token"
    else:
        token = sys.argv[0]
        #Post('/group_chat/pullGroupMsg/', reqSendGroupMsg)
        Post('/group_chat/sendGroupMsg/', reqPullGroupMsg, token)
