# Open Tree Hole 后端
使用 Django 框架
## 使用须知
0. 克隆本仓库并使用 `pipenv` 安装依赖
   ```shell
   git clone git@github.com:OpenTreeHole/backend.git
   cd backend
   pipenv install
   pipenv shell
   
1. 设置环境变量 `ENV`
   
   开发环境为 `development`，生产环境为 `production`
   ```shell
   # *nix
   export ENV=development
   # windows
   此电脑 -> 属性 -> 高级系统设置 -> 高级 -> 环境变量 -> 用户变量 -> 新建

2. 在 backend/OpenTreeHole 中新建配置文件 `secret.py`,并配置好**安全的权限**
    ```python
    SECRET_KEY = ""             # 足够长的密码，供 Django 安全机制
    
   DATABASE_HOST = ""          # 数据库主机
   DATABASE_PORT = 3306        # 数据库端口
   DATABASE_USER = ""          # 数据库用户
   DATABASE_PASSWORD = ""      # 数据库密码
   
   GITHUB_TOKEN = ""           # 采用 Github 图床
   API_KEY = [""]
   EMAIL_PSSWORD = ""
3. 采用**测试导向**开发模式，首先编写单元测试（在 /tests 目录），在**新分支**上开发，确保通过测试后提交至 `dev` 分支
   
    

