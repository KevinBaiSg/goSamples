
## etcd v3 服务注册与发现

### 流程 
* service 启动时向 etcd 注册自己的信息，即注册到services/  这个目录 
* service 可能异常推出，需要维护一个TTL(V3 使用 lease实现)，类似于心跳，挂掉了，master可以监听到  
* master监听 services/ 目录下的所有服务，根据不同 action（V3有put/delete），进行处理