# Open Tree Hole 后端

使用 Django 框架

## 开发须知

0. 克隆本仓库并安装依赖

   ```shell
   git clone git@github.com:OpenTreeHole/backend.git OpenTreeHole
   cd OpenTreeHole
   # 安装 python 依赖
   pipenv install
   # 安装系统依赖（以Debian为例）
   sudo apt install redis-server
   sudo apt install libmagic1
   # 执行数据库迁移并预加载数据
   pipenv shell
   python manage.py migrate
   python manage.py loaddata init_data  

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
