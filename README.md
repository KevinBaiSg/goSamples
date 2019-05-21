# etcd client v3
Fabric 1.4 增加了 raft 协议来实现区块排序服务，并且使用了 etcd v3 版本的实现。

### TODO
* PUT DELETE
* WATCH
* 事务
* 服务发现
* 分布式锁

### Running a local etcd cluster
First install goreman, which manages Procfile-based applications.

```bash  
goreman -f etcd_Procfile start
```