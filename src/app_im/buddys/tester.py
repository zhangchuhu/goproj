#!/usr/bin/env python
#-*- coding:utf8 -*-
import requests
import json

def Post(uid, url, req):
    print "// 请求" 
    print json.dumps(req, indent=4)

    r = requests.get('http://58.215.180.211:8091/token/get?uid=' + str(uid))
    headers = {'token' : r.text}
    r = requests.post('http://172.26.38.31:10067' + url, data=json.dumps(req), headers=headers)
    #r = requests.post('http://127.0.0.1:10067' + url, data=json.dumps(req), headers=headers)
    #r = requests.post('http://106.38.255.199:9098' + url, data=json.dumps(req), headers=headers)
    print "// 响应"
    print json.dumps(json.loads(r.text), indent=4)

# 请求加好友
reqAddBuddy = {
    'uid' : 1234,
    'bid' : 1235,
    'msg' : "hello",
}

# 同步好友列表
fullSyncBuddy = {
    'uid' : 3,
}

# 同意 / 拒绝加好友
addBuddyAnswer = {
    'uid' : 1234,
    'bid' : 1235,
    'agree' : 1,
    'seqId' : 0,
}

# 删除好友
delBuddy = {
    'uid' : 1234,
    'bid' : 1235,
}

# 创建群
createGrp = {
    'grp_name' : u"群聊5",
    'logo' : 'http://dl.vip.yy.com/xxx',
    'privilege' : 1,
    'password' : 'xxx',
    'invite_type' : 0,
}

# 加入群
joinGrp = {
    'gid' : 33,
}

# 退出群
exitGrp = {
    'gid' : 33,
}

# 拉取群成员
syncGrpMember = {
    'gid' : 33, 
}

# 同步群列表
syncGrpList = {
}

# 群邀请
inviteJoinGrp = {
    'uid' : 1234,
    'to_uid' : 12345,
    'gid' : 1,
}

# 群管理员验证结果
joinCheckGrp = {
    'uid' : 5,
    'gid' : 33,
    'agree' : 1,
}

#Post('/buddy/addbuddy/', reqAddBuddy)

#addBuddyAnswer['bid'] = 1236
#Post('/buddy/addbuddy_answer', addBuddyAnswer)

#Post('/buddy/del_buddy/', delBuddy)
#Post('/buddy/full_sync/', fullSyncBuddy)

#Post(5, '/group/create_grp/', createGrp)
#Post(5, '/group/join_grp/', joinGrp)
#Post(2, '/group/check_join_grp/', joinCheckGrp)
#Post('/group/invite/', inviteJoinGrp)
#Post(5, '/group/exit_grp/', exitGrp)
#Post(6, '/group/sync_member/', syncGrpMember)
Post(5, '/group/full_sync/', syncGrpList)
