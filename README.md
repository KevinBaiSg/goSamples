# etcd 
Fabric 1.4 开始使用 etcd 代替 kafka 来实现排序服务，案例中作为学习使用

### TODO
* 增加 查询 删除 修改
* watch
* 事务
* 服务发现

### Running a local etcd cluster
First install goreman, which manages Procfile-based applications.

```bash  
goreman -f etcd_Procfile start
```