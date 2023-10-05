# OpenTreeHole Backend Project

开源树洞后端综合体

- 单体架构，避免重复代码
- 基于 [nunu](https://github.com/go-nunu/nunu) 的框架设计，分层架构更清晰
- 基于 wire 的依赖注入设计

## 构建工具

[![Go][go.dev]][go-url]
[![Swagger][swagger.io]][swagger-url]

## 构建本项目

### 本地构建并运行

```shell
# 克隆本项目
git clone https://github.com/OpenTreeHole/backend.git
cd chatdan_backend

# 安装 swaggo 并且生成 API 文档
go install github.com/swaggo/swag/cmd/swag@latest
swag init -d cmd,internal/handler,internal/schema -p snakecase -o internal\docs

# 将 config/config_default.json 复制为 config/config.json, 并且按照需求修改配置，否则会使用默认配置
cp config/config_default.json config/config.json

# 运行
go run ./cmd

# 使用 nunu 运行本项目，支持文件更新热重启
go install github.com/go-nunu/nunu@latest
nunu run
```

API 文档详见启动项目之后的 http://localhost:8000/docs

### 生产部署

#### 使用 docker 部署

```shell
docker run -d \
  --name opentreehole_backend \
  -p 8000:8000 \
  -e MODULES_CURRICULUM_BOARD=true \
  -v opentreehole_data:/app/data \
  -v opentreehole_config:/app/config \
  opentreehole/backend:latest
```

#### 使用 docker-compose 部署

```yaml
version: '3'

services:
  backend:
    image: opentreehole/backend:latest
    restart: unless-stopped
    environment:
      - DB_TYPE=mysql
      - DB_URL=opentreehole:${MYSQL_PASSWORD}@tcp(mysql:3306)/opentreehole?parseTime=true&loc=Asia%2FShanghai
      - CACHE_TYPE=redis
      - CACHE_URL=redis:6379
      - MODULES_CURRICULUM_BOARD=true
    volumes:
      - data:/app/data
      - config:/app/config

  mysql:
    image: mysql:8.0.34
    restart: always
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - MYSQL_USER=opentreehole
      - MYSQL_PASSWORD=${MYSQL_PASSWORD}
      - MYSQL_DATABASE=opentreehole
      - TZ=Asia/Shanghai
    volumes:
      - mysql_data:/var/lib/mysql
      - mysql_config:/etc/mysql/conf.d

  redis:
    container_name: fduhole_redis
    image: redis:7.0.11-alpine
    restart: always

volumes:
  data:
  config:
```

环境变量：

1. `TZ`: 生产环境默认 `Asia/Shanghai`
2. `MODE`: 开发环境默认 `dev`, 生产环境默认 `production`, 可选 `test`, `bench`
3. `LOG_LEVEL`: 开发环境默认 `debug`, 生产环境默认 `info`, 可选 `warn`, `error`, `panic`, `fatal`
4. `PORT`: 默认 8000
5. `DB_TYPE`: 默认 `sqlite`, 可选 `mysql`, `postgres`
6. `DB_DSN`: 默认 `data/sqlite.db`
7. `MODULES_{AUTH/NOTIFICATION/TREEHOLE/CURRICULUM_BOARD}`: 开启模块，默认为 `false`

数据卷：

1. `/app/data`: 数据库文件存放位置
2. `/app/config`: 配置文件存放位置

注：环境变量设置仅在程序启动时生效，后续可以修改 `config/config.json` 动态修改配置

### 开发指南

1. 使用 wire 作为依赖注入框架。如果创建了新的依赖项，需要在 `cmd/wire/wire.go` 中注册依赖项的构造函数，之后运行 `nunu wire`
   生成新的 `cmd/wire/wire_gen.go`
2. 基于数据流的分层架构。本项目分层为 `handler` -> `service` -> `repository`。
    1. `hander` 结构必须组合 `*Handler` 类型。`handler` 只负责接口接收、鉴权和响应，并且把控制权交给 `service`
    2. `service` 结构必须组合 `Service` 接口。`service` 负责主要业务逻辑，其中数据库操作调用 `repository`
       的接口，`Service.Transaction` 方法可以把多个 `repository` 接口调用合并到一个事务中。
    3. `repository` 结构必须组合 `Repository` 接口。与数据库、缓存的CURD操作相关的必须放在 `repository` 的绑定函数中实现。
3. 分离接口**请求和响应模型** `schema` 和**数据库模型** `model`。
    1. `handler` 接受和发送都只能使用 `schema`，`service` 负责 `schema` 和 `model` 的转换，`repository` 只使用 `model`
    2. 如果 `model` 和 `schema` 之间的转换逻辑冗余，可以在 `schema`
       中定义类似 `func (s *Schema) FromModel(m *Model) *Schema` 和 `func (s *Schema) ToModel() *Model`
       绑定函数，保证模块的引用顺序是 `schema` -> `model`，避免循环引用。
    3. 避免在 `model` 和 `schema` 中作数据库、缓存的CURD，这些都应该在 `service` 层中调用 `repository`
       层的接口来完成。如果模型转换时需要用到数据库接口，需要作为函数参数传入转换函数。

## 计划路径

- [ ] 迁移旧项目到此处

## 贡献列表

<a href="https://github.com/OpenTreeHole/backend/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=OpenTreeHole/backend"  alt="contributors"/>
</a>

## 联系方式

JingYiJun - jingyijun@fduhole.com

Danxi-Dev - dev@fduhole.com

## 项目链接

[https://github.com/OpenTreeHole/backend](https://github.com/OpenTreeHole/backend)

[//]: # (https://www.markdownguide.org/basic-syntax/#reference-style-links)

[contributors-shield]: https://img.shields.io/github/contributors/OpenTreeHole/backend.svg?style=for-the-badge

[contributors-url]: https://github.com/OpenTreeHole/backend/graphs/contributors

[forks-shield]: https://img.shields.io/github/forks/OpenTreeHole/backend.svg?style=for-the-badge

[forks-url]: https://github.com/OpenTreeHole/backend/network/members

[stars-shield]: https://img.shields.io/github/stars/OpenTreeHole/backend.svg?style=for-the-badge

[stars-url]: https://github.com/OpenTreeHole/backend/stargazers

[issues-shield]: https://img.shields.io/github/issues/OpenTreeHole/backend.svg?style=for-the-badge

[issues-url]: https://github.com/OpenTreeHole/backend/issues

[license-shield]: https://img.shields.io/github/license/OpenTreeHole/backend.svg?style=for-the-badge

[license-url]: https://github.com/OpenTreeHole/backend/blob/main/LICENSE

[go.dev]: https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white

[go-url]: https://go.dev

[swagger.io]: https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white

[swagger-url]: https://swagger.io