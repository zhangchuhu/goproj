
# go环境配置

    export GOPATH=/data/vip/diamond/goproj
    export GOBIN=$GOPATH/bin
    export PATH=$GOPATH/bin:$PATH

# 依赖第三方库

    github.com/go-sql-driver/mysql
    github.com/astaxie/beego               # mysql的orm
    github.com/cihub/seelog                # 日志
    github.com/garyburd/redigo             # redis
    github.com/valyala/gorpc               # 进程间通信

# To-do

* redis hash分片
* redis 客户端封装测试
* http长连接接口封装
* 好友模块
* 聊天模块
* 客户端demo
