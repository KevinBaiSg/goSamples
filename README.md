# etcd client v3
Fabric 1.4 开始使用 etcd 代替 kafka 来实现排序服务，案例中作为学习 etcd 使用

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