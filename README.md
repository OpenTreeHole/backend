# OpenTreeHole Backend Project

DanXi & OpenTreeHole 后端

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
cd danke && swag init --pd --parseDepth 1

# 运行
go run danke/.

# 在 docker 中运行
docker compose up
```

API 文档详见启动项目之后的 http://localhost:8000/docs

### 生产部署

#### 使用 docker 部署

```shell
docker run -d \
  --name opentreehole_backend \
  -p 8000:8000 \
  -v opentreehole_data:/app/data \
  opentreehole/backend:latest
```

#### 使用 docker-compose 部署

修改 `docker-compose.yml`

环境变量：

1. `TZ`: 生产环境默认 `Asia/Shanghai`
2. `MODE`: 开发环境默认 `dev`, 生产环境默认 `production`, 可选 `test`, `bench`
3. `LOG_LEVEL`: 开发环境默认 `debug`, 生产环境默认 `info`, 可选 `warn`, `error`, `panic`, `fatal`
4. `PORT`: 默认 8000
5. `DB_TYPE`: 默认 `sqlite`, 可选 `mysql`, `postgres`
6. `DB_URL`: 默认 `sqlite.db`
7. `CACHE_TYPE`: 默认 `memory`, 可选 `redis`
8. `CACHE_URL`: 默认为空

数据卷：

1. `/app/data`: 数据库文件存放位置

### 开发指南

1. 基于数据流的分层架构。本项目分层为 `api` -> `model`。
   1. `api` 中包含所有的路由处理函数，负责接收请求和返回响应，其中的请求和响应模型定义在 `schema` 中。
   2. `model` 中包含所有的数据库模型定义，负责数据库的 CURD 操作。
   3. 简单的 CURD 操作可以直接在 `api` 中完成，复杂的、共享的业务逻辑应该在 `model` 中完成。
2. 分离接口**请求和响应模型** `schema` 和**数据库模型** `model`。
   1. `api` 接受和发送都只能使用 `schema`，数据库、缓存操作只能使用 `model`。
   2. 如果 `model` 和 `schema` 之间的转换逻辑冗余，可以在 `schema`
      中定义类似 `func (s *Schema) FromModel(m *Model) *Schema` 和 `func (s *Schema) ToModel() *Model`
      绑定函数，保证模块的引用顺序是 `schema` -> `model`，避免循环引用。
   3. 避免在 `model` 和 `schema` 中作数据库、缓存的CURD，这些都应该在加载 `model` 时完成。
   4. 如果模型转换时需要用到其他辅助数据，需要作为函数参数传入转换函数。

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