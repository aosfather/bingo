package service

import (
	"github.com/go-redis/redis"
	"time"
	"fmt"
)

/**
分布式锁
1、基于redis的set nx功能 并设置一个过期时间防止客户端死掉
2、处理完成通过del来释放锁
3、通过设置一个钥匙，来保证拥有钥匙的删除锁

可能存在的问题，拥有钥匙的一方卡死，等到其能删除锁的时候，错误的删除了另外获取资源的锁拥有者新建的锁
*/
const (
	_delscript=`
if redis.call("get", KEYS[1]) == ARGV[1] then
	return redis.call("del", KEYS[1])
else
	return 0
end`
)
type RedisLockFactory struct {
	client *redis.Client
	expire int64
}

func (this *RedisLockFactory) Init(c *redis.Client,e int64){
	this.client=c
	this.expire=e
}

func (this *RedisLockFactory) Create(r string) *RedisLock {
	return &RedisLock{this.client,r,this.getKey(),this.getExpire()}
}

func (this *RedisLockFactory) getExpire() time.Duration {
	return time.Duration(this.expire)*time.Microsecond
}

func (this *RedisLockFactory) getKey() string {
	return fmt.Sprintf("lock%d",time.Now().UnixNano())
}

type RedisLock struct {
	client *redis.Client
	resource string //资源
	key string//密钥
	expire time.Duration
}

func (this *RedisLock) Lock() {
    TheLock:
   	 exist,err:=this.client.SetNX(this.resource,this.key,this.expire).Result()
   	 if err!=nil {
   	 	println(err.Error())
	 }

     if exist {
   	  time.Sleep(time.Microsecond*10)
      goto TheLock
     }

}
func (this *RedisLock) Unlock(){
	this.client.Eval(_delscript,[]string{this.resource},this.key)
}