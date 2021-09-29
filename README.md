# Open Tree Hole 后端

使用 Django 框架

## 使用须知

本项目使用 Docker 持续集成部署

### 使用 docker-compose 安装：

1. 下载 [docker-compose.yaml](https://github.com/OpenTreeHole/backend/blob/master/docker-compose.yaml)

2. 填写 `environment` 块下的若干环境变量，完整的列表及说明请参见 [配置文件](https://github.com/OpenTreeHole/backend/blob/master/OpenTreeHole/config.py)

   注意：环境变量**不应**以引号包裹，否则会无法解析

3. 运行 `docker-compose up -d`

若成功，项目可以在 80 端口访问

### 注意：

- 域名和 CORS 等配置应在 nginx 等反向代理服务器中完成，请自行配置相关项

- 项目初始化时会自动创建管理员账户，邮箱为 admin@opentreehole.org，密码为 admin，须尽快登录至管理后台修改管理员信息

## 开发须知

0. 克隆本仓库并安装依赖

   ```shell
   git clone git@github.com:OpenTreeHole/backend.git OpenTreeHole
   cd OpenTreeHole
   # 安装系统依赖（以 Debian 为例）
   sudo apt install python3 python3-pip redis-server libmagic1
   pip3 install pipenv
   # 安装 python 依赖
   pipenv install --dev
   # 执行数据库迁移并预加载数据
   pipenv shell
   python manage.py migrate
   python manage.py loaddata init_data
   python start.py
   # 运行开发服务器
   python manage.py runserver  

1. 设置环境变量 `HOLE_ENV`

   开发环境为 `development`，生产环境为 `production`
   ```shell
   # *nix
   export HOLE_ENV=development
   # windows
   此电脑 -> 属性 -> 高级系统设置 -> 高级 -> 环境变量 -> 用户变量 -> 新建

2. 在 OpenTreeHole/OpenTreeHole/config.py 中填写需要的配置, 或设置同名的环境变量以覆盖


3. 启动celery

   ```shell
   celery -A OpenTreeHole worker -l info -P eventlet
   ```

4. 采用**测试导向**开发模式，首先编写单元测试（在 /tests 目录），在**新分支**上开发，确保通过测试后提交至 `dev` 分支
