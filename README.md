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
