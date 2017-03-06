package redisf

import (
	"errors"
	"strings"
	"time"
	"vipcomm/zkclient"

	"github.com/chasex/redis-go-cluster"
	log "github.com/cihub/seelog"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/samuel/go-zookeeper/zk"
)

var (
	name2Pool           = make(map[string]*RedisPool, 0)
	cluster2Pool        = make(map[string]*redis.Cluster, 0)
	ZK_REDIS_PRE        = "/redis"
	ZK_REDISCLUSTER_PRE = "/rediscluster"
)

// 实例
type RedisPool struct {
	pools    []*redigo.Pool
	HashFunc func(hashKey interface{}) int64
}

func initRedisByName(conn *zk.Conn, name string, groupId string) (err error) {
	path := ZK_REDIS_PRE + "/" + name + "/" + groupId
	childrens, _, err := conn.Children(path)
	if err != nil {
		return err
	}

	length := len(childrens)
	if length == 0 {
		return errors.New("redis config is empty. path:" + path)
	}

	// 通过额外节点，实现redis多实例轮询
	idx, err := zkclient.ZkGetIndex(conn, path)
	if err != nil {
		idx = int32(time.Now().Unix())
	}
	strAddr := childrens[idx%int32(length)]
	log.Debugf("get redis hosts:%v idx:%v all:%v", strAddr, idx, childrens)

	// 创建实例
	rpool := &RedisPool{
		pools:    []*redigo.Pool{},
		HashFunc: defaultHashFunc, // 自定义hash函数 to-do
	}

	// 实例有多个切片组成
	addrs := strings.Split(strAddr, ",")
	for _, addr := range addrs {
		pool := NewPool(addr, "4xT&2b}zblP")
		c := pool.Get()
		defer c.Close()

		rpool.pools = append(rpool.pools, pool)
		log.Infof("new redis. name:%v addr:%v err:%v ",
			name, addr, c.Err())
	}
	name2Pool[name] = rpool
	return
}

func initRedisClusterByName(conn *zk.Conn, name string) (err error) {
	path := ZK_REDISCLUSTER_PRE + "/" + name
	value, _, err := conn.Get(path)
	if err != nil {
		return err
	}
	addrs := strings.Split(string(value), ",")
	cluster, err := redis.NewCluster(
		&redis.Options{
			StartNodes:   addrs,
			ConnTimeout:  50 * time.Millisecond,
			ReadTimeout:  50 * time.Millisecond,
			WriteTimeout: 50 * time.Millisecond,
			KeepAlive:    16,
			AliveTime:    60 * time.Second,
		})

	if err != nil {
		log.Infof("redis.New error: %s", err.Error())
	}
	cluster2Pool[name] = cluster
	return
}

func InitRedis(conn *zk.Conn, names []string, groupId string) (err error) {
	for _, name := range names {
		if err = initRedisByName(conn, name, groupId); err != nil {
			log.Errorf("init redis faild. name:%v group:%v", name, groupId)
			return err
		}
	}
	return nil
}

func InitRedisCluster(conn *zk.Conn, names []string) (err error) {
	for _, name := range names {
		if err = initRedisClusterByName(conn, name); err != nil {
			log.Errorf("init redis faild. name:%v ", name)
			return err
		}
	}
	return nil
}

/*
 * 执行redis命令
 * @param  name   redis实例配置名
 */
func Do(name, cmd string, args ...interface{}) (reply interface{}, err error) {
	v, ok := name2Pool[name]
	if !ok || len(v.pools) == 0 {
		log.Error("redis not init. name: ", name)
		return nil, errors.New("redis: " + name + "not init")
	}

	conn := v.pools[0].Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

/*
 * 执行redis命令
 * @param  name   rediscluster实例配置名
 */
func ClusterDo(name, cmd string, args ...interface{}) (reply interface{}, err error) {
	cluster, ok := cluster2Pool[name]
	if !ok || (cluster == nil) {
		log.Error("redis not init. name: ", name)
		return nil, errors.New("redis: " + name + "not init")
	}
	return cluster.Do(cmd, args...)
}

/*
 * 执行redis命令
 * @param  name     redis实例配置名
 * @param  hashKey  用来计算redis-server的slot，默认hash函数只支持int64
 */
func HashDo(name string, hashKey interface{}, cmd string, args ...interface{}) (reply interface{}, err error) {

	v, ok := name2Pool[name]
	if !ok || len(v.pools) == 0 {
		log.Error("redis not init. name: ", name)
		return nil, errors.New("redis not init")
	}

	k := v.HashFunc(hashKey)
	idx := k % int64(len(v.pools))
	conn := v.pools[idx].Get()
	defer conn.Close()
	return conn.Do(cmd, args...)
}

// 创建redis poll
func NewPool(server, password string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle: 20,
		Dial: func() (redigo.Conn, error) {
			//c, err := redigo.Dial("tcp", server)
			c, err := redigo.DialTimeout("tcp",
				server,
				10*time.Second, // conn timeout
				30*time.Second, // read
				30*time.Second) // write
			if err != nil {
				return nil, err
			}
			if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// 默认hash函数
func defaultHashFunc(args interface{}) int64 {
	v, ok := args.(int64)
	if !ok {
		panic(errors.New("default hash func mush use int64"))
	}
	return v
}
