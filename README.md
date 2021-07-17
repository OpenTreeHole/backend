# Open Tree Hole 后端

使用 Django 框架

## 使用须知

0. 克隆本仓库并使用 `pipenv` 安装依赖
   ```shell
   git clone git@github.com:OpenTreeHole/backend.git
   cd backend
   pipenv install
   pipenv shell
   python manage.py migrate auth       # 迁移数据库，这两条必须，否则会报错
   python manage.py migrate authtoken
   python manage.py migrate 

1. 设置环境变量 `ENV`

   开发环境为 `development`，生产环境为 `production`
   ```shell
   # *nix
   export ENV=development
   # windows
   此电脑 -> 属性 -> 高级系统设置 -> 高级 -> 环境变量 -> 用户变量 -> 新建

2. 在 backend/OpenTreeHole 中新建配置文件 `config.py`,并配置好**安全的权限**
    ```python
   # 站点配置
   SITE_NAME = 'Open Tree Hole'
   EMAIL_WHITELIST = ["test.com"]  # 允许注册树洞的邮箱域名
   MIN_PASSWORD_LENGTH = 8  # 允许的最短用户密码长度
   VALIDATION_CODE_EXPIRE_TIME = 5  # 验证码失效时间（分钟）
   MAX_TAGS = 5
   MAX_TAG_LENGTH = 8
   NAME_LIST = []  # 随机昵称列表
   SECRET_KEY = ""  # 足够长的密码，供 Django 安全机制
   # 数据库配置
   DATABASE_HOST = ""          # 数据库主机
   DATABASE_PORT = 3306        # 数据库端口
   DATABASE_NAME = ""          # 数据库名称
   DATABASE_USER = ""          # 数据库用户
   DATABASE_PASSWORD = ""      # 数据库密码
   # 邮件配置
   EMAIL_HOST = ''
   EMAIL_PORT = 587
   EMAIL_HOST_USER = ''
   EMAIL_HOST_PASSWORD = ''
   EMAIL_USE_TLS = 
   EMAIL_USE_SSL = 
   DEFAULT_FROM_EMAIL = ''  # 默认发件人地址
   
   GITHUB_TOKEN = ""           # 采用 Github 图床
   EMAIL_PSSWORD = ""
3. 采用**测试导向**开发模式，首先编写单元测试（在 /tests 目录），在**新分支**上开发，确保通过测试后提交至 `dev` 分支
   
    

