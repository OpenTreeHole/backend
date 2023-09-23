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

### 开发指南

1. 使用 wire 作为依赖注入框架。如果创建了新的依赖项，需要在 `cmd/wire/wire.go` 中注册依赖项的构造函数，之后运行 `nunu wire`
   生成新的 `cmd/wire/wire_gen.go`
2. 利用分层架构。本项目分层为 handler -> service -> repository。
    1. `hander` 结构必须组合 `*Handler` 类型。`handler` 只负责接口接收、鉴权和响应，并且把控制权交给 `service`
    2. `service` 结构必须组合 `Service` 接口。`service` 负责主要业务逻辑，其中数据库操作调用 `repository`
       的接口，`Service.Transaction` 方法可以把多个 `repository` 接口调用合并到一个事务中。
    3. `repository` 结构必须组合 `Repository` 接口。与数据库、缓存的CURD操作相关的必须放在 `repository` 的绑定函数中实现。

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