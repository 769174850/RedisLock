# RedisLock
简单实现了使用redis的setNX与Lua脚本的分布式锁
并且进行了验证
## redis的SetNX命令
SETEX key seconds value
将值 value 关联到 key ，并将 key 的过期时间设为 seconds (以秒为单位)
引用:

[Golang的分布式应用](https://juejin.cn/post/7124281261371162661?from=search-suggest#heading-0)
