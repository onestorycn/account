package rediscli

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"account/library/config"
	"fmt"
	"errors"
	"sync"
)

//connection pool
var (
	// 定义常量
	RedisClients = make(map[string]*redis.Pool, 1)
)

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Password string `json:"password"`
}
var lock *sync.Mutex = &sync.Mutex{}

func LoadRedisConn(poolName string) (conn redis.Conn){
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println(err)
			conn = nil
		}
	}()
	if RedisClients[poolName] == nil{
		lock.Lock()
		defer lock.Unlock()
		if RedisClients[poolName] == nil{
			loadRedisPool(poolName)
		}
	}
	return RedisClients[poolName].Get()
}

func init(){
	envDefault := config.GetServiceEnv("env")
	if envDefault == nil {
		envDefault = "online"
	}
	loadRedisPool(envDefault.(string))
}

func loadRedisPool(poolName string)  {
	// 从配置文件获取redis的ip以及db

	redisConfig := new(RedisConfig)
	errConf := config.GetConfigMapObj("db", redisConfig, "redis", poolName)
	if errConf != nil {
		panic(errConf.Error())
	}

	host := redisConfig.Host
	port := redisConfig.Port
	if len(host) == 0 || len(port) == 0 {
		panic(errors.New("get redis conf fail"))
	}

	var password = redisConfig.Password

	REDIS_HOST := host + ":" + port

	RedisClients[poolName] = &redis.Pool{
		// 从配置文件获取maxidle以及maxactive，取不到则用后面的默认值
		MaxIdle:     10,
		MaxActive:   30,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial(
				"tcp",
				REDIS_HOST,
				redis.DialConnectTimeout(1*time.Second),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second),
			)
			if err != nil {
				return nil, err
			}

			if password != "" && len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			// 选择db
			return c, nil
		},
	}
}

