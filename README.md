SEQSVR
---------

#### High performance unique number generator powered by Go

[中文 README](README-zh_CN.md)

### Features

* Distributed: Can be scaled horizontally
* High performance: Allocation ID only accesses memory (up to the upper limit will request the database once)
* Ease of use: Provide as HTTP service
* Unique: MySQL auto increment ID, never repeat

### Requirement

We are using these awesome projects as necessary libraries.

* gopkg.in/yaml.v2
* github.com/go-sql-driver/mysql
* github.com/satori/go.uuid

### Installation

You can install this service in the following three ways.

**Note: You need to create the database and modify the configuration of the database in the configuration file before starting.**

go get:

```shell
go get github.com/qichengzx/seqsvr
seqsvr
```
Compile By Yourself:

```shell
git clone git@github.com:qichengzx/seqsvr.git
cd seqsvr
go build .
./seqsvr
```
Docker:

The Docker "multi-stage" build feature is used, be ensure the Docker version is 17.05 or above. See:[Use multi-stage builds](https://docs.docker.com/develop/develop-images/multistage-build/)

```shell
git clone git@github.com:qichengzx/seqsvr.git
cd seqsvr
docker build -t seqsvr:latest .
docker run -p 8000:8000 seqsvr:latest
```

### Create Database

The database name can be customized, modify config.yml.

Then import the following SQL to generate the data table.

```mysql
CREATE TABLE `generator_table` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `uuid` char(36) NOT NULL COMMENT 'Machine identification',
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `stub_UNIQUE` (`uuid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
```

### Configuration

The configuration file is using [YAML](http://yaml.org/).

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

### Usage

```shell
curl http://localhost:8000/new

{"code":0,"msg":"ok","data":{"id":101}}
```