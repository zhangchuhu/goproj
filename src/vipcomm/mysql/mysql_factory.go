package mysql

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"github.com/bitly/go-simplejson"
	log "github.com/cihub/seelog"
	_ "github.com/go-sql-driver/mysql"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	ZK_DB_PRE = "/db"
	DB_PASSWD = map[string]string{
		"app_im":    "im-db-admin:dsqazxsw21",
		"vip_logic": "im-db-admin:dsqazxsw21",
	}
)

/*
* names 数据库名
 */
func InitDbOrm(conn *zk.Conn, names []string, defdb string, groupId string) (err error) {
	for _, name := range names {
		passwd, ok := DB_PASSWD[name]
		if ok == false {
			return errors.New("not found mysql passwd. name: " + name)
		}

		path := ZK_DB_PRE + "/" + name
		value, _, err := conn.Get(path)
		if err != nil {
			return err
		}

		js, err := simplejson.NewJson([]byte(value))
		if err != nil {
			log.Errorf("unmarsh json faild. name:%v val:%v err:%v", name, string(value), err)
			return err
		}

		hosts, err := js.Get("host").Get(groupId).StringArray()
		if err != nil {
			log.Errorf("get hosts faild. name:%v group:%v val:%v err:%v", groupId, name, string(value), err)
			return err
		}

		addr := passwd + "@tcp(" + hosts[0] + ")/" + name + "?charset=utf8"
		connNum := 20
		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase(name, "mysql", addr)
		orm.SetMaxOpenConns(name, connNum)
		log.Infof("new db ok. name:%v addr:%v maxconn:%v", name, addr, connNum)

		// default db
		if defdb == name {
			orm.RegisterDriver("mysql", orm.DRMySQL)
			orm.RegisterDataBase("default", "mysql", addr)
			orm.SetMaxOpenConns("default", connNum)
			log.Infof("new db ok. name:default addr:%v maxconn:%v", addr, connNum)
		}
	}
	return nil
}
