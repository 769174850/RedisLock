package lock

import (
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Lock struct {
	Client     *redis.Client
	lockKey    string
	expireTime time.Duration
} //锁的结构

func NewLock(Client *redis.Client, lockKey string, expireTime time.Duration) *Lock {
	return &Lock{
		Client:     Client,
		lockKey:    lockKey,
		expireTime: expireTime,
	}
} //创建一个锁

func (l *Lock) Lock(id string) (bool, error) {
	return l.Client.SetNX(l.lockKey, id, l.expireTime).Result()
} //使用setNX添加一个锁，并将调用锁的用户存储下来

const unlockScript = `
if (redis.call("get",KEYS[1]) == KEYS[2]) then
	redis.call("del",KEYS[1])
	return true 
end
return false
` //Lua脚本检验锁是否被调用，当key存在且value等于id时，删除key，返回true，否则返回false

func (l *Lock) Unlock(id string) error {
	_, err := l.Client.Eval(unlockScript, []string{l.lockKey, id}).Result()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}//当调用锁的用户id等于key时，解除锁并删除锁
