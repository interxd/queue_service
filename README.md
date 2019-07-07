## queue_service
  - 第一次接触golang，基于golang实现了简单的排队系统
  - 排队系统本质上是一个限流系统，虽然有多种限流的实现机制，对于游戏业务的排队系统来说，较合适的实现是限流漏斗（Leaky bucket）算法
  - 为了便于扩展，抽象了限流队列的interface（RateLimitQueue），并且实现了简单的限流漏斗类（LeakyBucketQueue）
#### 服务器
  - 配置文件为server/conf.ini，各字段含义配置文件中有说明
#### 客户端
  - 配置文件为client/conf.ini,各字段含义配置文件中有说明


## 测试
### 单测
  - server/queue中包含LeakyBucketQueue的简单单测，在该目录中 go test -v 即可启动单测
  
### 压测
#### 环境
- 系统：MacOs-10.14 
- CPU: 4核 Intel(R) Core(TM) i7-4870HQ CPU @ 2.50GHz
- 内存: 16G
- go版本： 1.12.6

#### 测试启动
- 服务器启动：
    - `go run server_main.go`
- 客户端启动
    - `go run client_main.go`

#### 测试结果
- 为了让用户能尽快知道目前在队列中的位置，目前策略为出队一批后全局广播当前位置，此操作非常消耗CPU，且出队后瞬间CPU会急速飙升至100%多
- 如果去掉广播操作，服务器CPU负载保持在30%左右


## TODO LIST
- 不同游戏可实现不同策略的排队机制，譬如可以使用多个队列，在基本漏斗限流的基础上加上优先级权重，让vip用户优先拿到token
- 目前的位置广播测量非常低效，对于一个高效来说，服务器端无需每次队列有变化时都进行全量广播，我们可以使用不同的策略来处理队列中的不同用户：
  - 对于排在队列前端的用户，其位置信息比较敏感，位置发生变化时实时广播，比如前200名的用户
  - 对于排在靠后的位置的用户，其位置信息的准确性要求不用那么高，服务端可以不用实时广播其位置，客户端可以根据以下三个数据自己进行模拟计算：
    - 登录是入队的位置
    - 服务器发放token时间间隔
    - 服务器每次发放token数量
  
