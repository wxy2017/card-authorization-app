# 功能卡片授权系统

一个专为情侣和朋友设计的互动卡片应用，支持创建、发送和使用功能卡片。

## 功能特性

- **创建卡片**：用户可以创建自定义功能卡片，设置名称、描述和有效期
- **发送卡片**：将卡片发送给其他用户
- **使用卡片**：收到卡片后可以使用，触发相应任务
- **卡片管理**：查看和管理所有创建和收到的卡片
- **移动端优化**：响应式设计，完美适配手机使用

## 技术栈

- **后端**：Go + Gin框架 + GORM + SQLite
- **前端**：HTML5 + CSS3 + 原生JavaScript
- **认证**：JWT Token
- **部署**：单文件可执行程序

## 项目结构

```
card-authorization-app/
├── backend/
│   ├── main.go              # 主程序入口
│   ├── go.mod              # Go模块配置
│   ├── database/
│   │   └── database.go     # 数据库初始化
│   ├── models/
│   │   ├── user.go         # 用户模型
│   │   └── card.go         # 卡片模型
│   ├── handlers/
│   │   ├── auth.go         # 认证处理
│   │   ├── card.go         # 卡片处理
│   │   └── page.go         # 页面处理
│   └── middleware/
│       └── auth.go         # 认证中间件
├── frontend/
│   ├── static/
│   │   ├── css/
│   │   │   └── style.css   # 样式文件
│   │   └── js/
│   │       ├── auth.js     # 认证相关JS
│   │       ├── dashboard.js # 控制台JS
│   │       ├── create_card.js # 创建卡片JS
│   │       └── cards.js    # 卡片管理JS
│   └── templates/          # HTML模板
│       ├── index.html      # 首页
│       ├── login.html      # 登录页
│       ├── register.html   # 注册页
│       ├── dashboard.html  # 控制台
│       ├── create_card.html # 创建卡片
│       └── cards.html      # 卡片管理
└── README.md               # 项目说明
```

## 快速开始

### 环境要求
- Go 1.21 或更高版本

### 安装依赖
```bash
cd backend
go mod tidy
```

### 运行项目
```bash
cd backend
go run main.go
```

访问 http://localhost:8080 开始使用

## API接口

### 认证相关
- `POST /api/register` - 用户注册
- `POST /api/login` - 用户登录
- `GET /api/profile` - 获取用户信息

### 卡片相关
- `POST /api/cards` - 创建卡片
- `GET /api/cards` - 获取我创建的卡片
- `GET /api/cards/received` - 获取收到的卡片
- `POST /api/cards/:id/use` - 使用卡片
- `POST /api/cards/:id/send` - 发送卡片
- `GET /api/users/search` - 搜索用户

## 使用示例

1. **注册账号**：访问 `/register` 创建新账号
2. **创建卡片**：登录后访问 `/cards/create` 创建新卡片
3. **发送卡片**：在卡片列表中选择卡片并发送给好友
4. **使用卡片**：收到卡片后点击"使用"按钮

## 卡片状态

- **active**：可用状态，可以发送或使用
- **used**：已使用状态，卡片已注销
- **expired**：已过期状态，超过有效期

## 安全特性

- 密码使用bcrypt加密存储
- JWT Token认证
- 输入验证和错误处理
- 跨站请求伪造保护

## 移动端优化

- 响应式设计，适配各种屏幕尺寸
- 触摸友好的交互设计
- 底部导航栏，方便单手操作
- 加载动画和状态反馈

## 部署说明

项目编译为单文件可执行程序，包含所有静态资源：

```bash
# 编译
go build -o card-authorization main.go

# 运行
./card-authorization
```

数据库文件 `card_authorization.db` 会自动创建在项目目录下。
