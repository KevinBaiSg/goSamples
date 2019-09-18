# 环境
Debian 9 + 10.1.41-MariaDB

# Q&A
## 设置 root password

```shell script
$sudo mysql -u root

mysql> use mysql;
​mysql> update user set authentication_string=password('your_username') where user='root';
​mysql> flush privileges;
​mysql> quit
```

## Create develop account

### 创建用户
```shell script
mysql> CREATE USER 'your_username'@'%' IDENTIFIED BY 'your_password'; //新建一个用户
mysql> GRANT ALL PRIVILEGES ON *.* TO 'your_username'@'%'; // 分配所有权限
mysql> FLUSH PRIVILEGES; // 刷新权限
```

### 修改 bind address 

查看文件 `/etc/mysql/mariadb.conf.d/50-server.cnf` 并修改

```shell script
# By default we only accept connections from localhost
- bind-address    = 127.0.0.1
+ bind-address    = 0.0.0.0
```

## java.io.EOFException: unexpected end of stream, read 0 bytes from 4 (socket was closed by server).


# 参考 
[How to set, change, and recover a MySQL root password](https://www.techrepublic.com/article/how-to-set-change-and-recover-a-mysql-root-password/)
