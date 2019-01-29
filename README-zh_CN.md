SEQSVR
---------

#### Go实现的全局唯一序列号生成服务

[English](README.md)

### 特性

* 分布式：可任意横向扩展
* 高性能：分配 ID 只访问内存(到达上限会请求数据库一次)
* 易用性：对外提供 HTTP 服务
* 唯一性：MySQL 自增 ID，永不重复
* 高可靠：MySQL 持久化

### 依赖项

本项目使用下列优秀的项目作为必要组件。

* gopkg.in/yaml.v2
* github.com/go-sql-driver/mysql
* github.com/satori/go.uuid

### 安装

**注意:需要在启动之前创建数据库并修改配置文件中数据库的配置。**

go get 方式：

需保证 $GOPATH/bin 在系统 PATH 中。

```shell
go get github.com/qichengzx/seqsvr
seqsvr
```
单独编译：

```shell
git clone git@github.com:qichengzx/seqsvr.git
cd seqsvr
go build .
./seqsvr
```
Docker 方式：

Dockerfile 使用了 Docker 多阶段构建功能，需保证 Docker 版本在 17.05 及以上。详见：[Use multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/)

```shell
git clone git@github.com:qichengzx/seqsvr.git
cd seqsvr
docker build -t seqsvr:latest .
docker run -p 8000:8000 seqsvr:latest
```

### 初始化数据库

数据库名称可以自定义，修改 config.yml 即可。

然后导入以下 SQL 生成数据表。

```mysql
CREATE TABLE `generator_table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL COMMENT '机器识别码',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `stub_UNIQUE` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
```

### 修改配置

配置文件使用 [YAML](http://yaml.org/) 格式。

```yaml
#app
port: ':8000'

#service
step: 100

#db
mysql:
  user: 'root'
  password: ''
  host: 'tcp(localhost:3306)'
  database: 'sequence'
```

可修改端口号及 MySQL 的配置。

### 使用

```shell
curl http://localhost:8000/new

{"code":0,"msg":"ok","data":{"id":101}}
```

### 原理

本项目设计原理来自 携程技术中心 的[干货 | 分布式架构系统生成全局唯一序列号的一个思路](https://mp.weixin.qq.com/s/F7WTNeC3OUr76sZARtqRjw)。

服务初始化后第一次请求会在 MySQL 数据库中插入一条数据，以生成初始 ID。

后续的请求，都会在内存中进行自增返回，并且保证返回的 ID 不会超过设置的上限，到达上限后会再次从 MySQL 中更新数据，返回新的初始 ID 。

##### 核心SQL

```mysql
REPLACE INTO `generator_table` (uuid) VALUES ("54f5a3e2-e04c-4664-81db-d7f6a1259d01");
```
