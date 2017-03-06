// define
var vipim = {
    
    /* init
     * listeners {
     *      onSyncBuddy : function(Arrays[]),           好友列表回调
     *      log:          function(level, msg),         日志
     * }
     */
    init: function(uid, pushId, listeners) {},

    // 发送消息
    sendMsg: function(to_uid, msg) {},

    // ajax
    apiCall: function(url, data, callBack) {},
};

// imp
(function(vipim){

    // push-server 地址
    var serverAddr = "http://106.38.255.199:20001";
    function genUrl(url) {
        return url;
    }

    // 日志函数
    var log = null;
    var Debug = 0;
    var Info = 1;

    // 消息类型
    var CHAT_TYPE_BUDDY = 0;
    var CHAT_TYPE_GROUP = 1;

    // 上下文
    var ctx = {
        selfUid: null,
        pushId: null,
    };

    // 好友
    var BuddyManager = new function() {
        console.log("new BuddyManager");
        this.buddys = new Array();
        this.version = 0;           // 版本号

        // 同步好友列表回调
        var onSyncBuddy = null;
        this.onFullSyncBuddyCb = function(res, ok) {
            if (!ok) {
                log(Info, "sync buddy faild.");
                return;
            }     
            this.version = res.version;
            this.buddys = res.bids; 
            if (onSyncBuddy) {
                onSyncBuddy(this.buddys); 
            }
            log(Info, "sync buddy ok. version: " + res.version + " buddy: " + res.bids);
        }

        this.clear = function() {
            this.version = 0;
            this.buddys = new Array();
        };

        this.init = function(listeners) {
            onSyncBuddy = listeners.onSyncBuddy;
        }

        // 同步好友列表
        this.sync = function() {
            log(Info, "sync buddy start. ver: " + this.version);
            var req = {
                'uid' : ctx.selfUid,
                'version' : this.version,
            };
            apiCall(genUrl('/buddy/full_sync'), req, function(res, ok) {
                // 这里处于全局作用域
                vipim.BuddyManager.onFullSyncBuddyCb(res, ok);
            });
        }
    };

    // clear 
    function clearSdk() {
        // 清空正在请求的ajax 
        // to-do
        
        ctx = {
            selfUid: null,
            pushId: null,
        };
        BuddyManager.clear();
    }

    // 发送请求
    function apiCall(url, data, callBack) {
        log(Debug, "==== ajax call. url: " + url);
        $.ajax({
            url: url,
            type: "POST",
            data: data,
            success: function(res, textStatus, jqXHR) {
                callBack(res, true);
            },
            error: function(jqXHR, textStatus, err) {
                console.log(err, textStatus, jqXHR);
                log(Debug, "==== ajax faild. url: " + url);
                callBack(data, false);
            },
            dataType: "json",
        });
    }
    
    // init
    vipim.init = function(uid, pushId, listeners) {
        // 清空上下文
        if (ctx.selfUid && ctx.selfUid != uid) {
            clearSdk();
        }

        BuddyManager.init(listeners);

        ctx.selfUid = uid;
        ctx.pushId = pushId;
        log = listeners.log;
        var pushClient = PushClient(serverAddr, {
                transports: ['websocket', 'polling'],
                pushId: pushId
                });

        // 收到服务器发回的数据
        pushClient.on('push', function(data) {
            log(Debug, "====" + JSON.stringify(data)); 
        });

        // 和长连接重连回调
        pushClient.on('connect',function(data){
            log(Info, "connect push server ok. pushId: " + pushId);
            
            // 绑定pushid
            var pk = {};
            pk.uid = ctx.selfUid;
            pk.pushid = ctx.pushId;
            apiCall(genUrl('/buddy/bind_push'),
                    pk,
                    function(res, ok) {
                        log(Info, "bind push id " + (ok ? "ok" : "faild"));
                    });         

            // 同步好友列表
            BuddyManager.sync();
        });

        pushClient.on('disconnect',function(){
            log(Info, "==== disconnected. pushId:  " + pushId);
        });
    }

    // 发送消息
    vipim.sendMsg = function(id, msgType, msg) {
        if (msgType == CHAT_TYPE_BUDDY) {
            var req = {
                "From":  ctx.selfUid,
                "To":    id,
                "Msg":   msg,
            }
            apiCall(genUrl("/send_msg"), req, function(res, ok) {
                log(Info, "====sendMsg " + (ok ? "ok" : "faild"));
            });
        } else if (msgType == CHAT_TYPE_GROUP) {
        }
        log(Info, "send msg to to " + id);
    }

    vipim.apiCall = apiCall;
    vipim.BuddyManager = BuddyManager;

})(vipim);
