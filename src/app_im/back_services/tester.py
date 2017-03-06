#!/usr/bin/env python
#-*- coding:utf8 -*-
import requests
import json

def Get(url, req):
    print "// 请求" 
    print json.dumps(req, indent=4)
    print ""

    r = requests.get("http://127.0.0.1:10070" + url, params= req)
    print "// 响应"
    print json.dumps(json.loads(r.text), indent=4)

#Get("/back_services/get_grp_role/", { "uid" : 2, "gid" : 33, })
Get("/back_services/get_grp_role/?uid=2&gid=33", {})

