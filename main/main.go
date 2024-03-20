package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"redis_lock/lock"
	"sync"
	"time"
)

var rdb *redis.Client

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := rdb.Ping().Result() // 检验是否链接成功
	if err != nil {
		log.Println(err) // 输出错误日志
		return
	}

	defer rdb.Close()//关闭redis

	newLock := lock.NewLock(rdb, "locker" , 30*time.Second)//创建一个锁

	counter := 0

	worker := func(i int) {
		for {
			id := fmt.Sprintf("worker%d", i)
			ok, err := newLock.Lock(id)//调用生成锁的
			if err != nil {
				log.Println(err)//返回错误日志
				return
			}
			if ok {
				log.Printf("worker %d obtain newLock", i)
			}
			if !ok {
				log.Printf("worker %d can not obtain newLock", i)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			newLock.Unlock(id)//解除锁
			counter++
			log.Printf("worker %d, add counter %d", i, counter)
			break
		}
	}

	wg := sync.WaitGroup{}//启用协程
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		id := i
		go func() {
			defer wg.Done()
			worker(id)//多个协程调用锁，检验是否成功
		}()
	}

	wg.Wait()
}
