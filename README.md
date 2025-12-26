# Blog Service

Go + Gin 实现的博客后台 API 项目骨架，预留 MySQL/GORM 支持，便于后续扩展文章、分类、评论等模块。

## 快速开始
- 准备 Go 1.24+ 环境，本地 MySQL 可选。
- 复制配置：`cp .env.example .env`，按需修改（`MYSQL_DSN` 留空可无数据库运行）。
- 安装依赖：`go mod download`（首次需要）。
- 启动服务：`go run ./cmd/server`，默认监听 `:8080`，会自动创建上传目录（`UPLOAD_DIR`）。

## 配置
- `APP_ADDR`：服务监听地址，默认 `:8080`。
- `MYSQL_DSN`：MySQL 连接串，留空则跳过数据库连接。
- `JWT_SECRET`：JWT 签名密钥，默认开发值，部署前务必修改。
- `UPLOAD_DIR`：上传目录路径，默认 `./uploads`。
- `COMMENT_MODERATION`：是否开启评论审核布尔值，默认 `false`。

## 可用接口（当前）
- `GET /healthz`：健康检查；配置了 `MYSQL_DSN` 时会同时 ping 数据库。
- `GET /api/v1/ping`：基础连通性探活，返回 `{"message":"pong"}`。

## 目录结构
```text
.
|-- cmd/server/main.go       # 程序入口，加载配置、初始化依赖
|-- internal/config          # 环境变量配置加载
|-- internal/router          # 路由注册
|-- internal/handlers        # 业务处理器（健康检查等）
|-- internal/middleware      # 通用中间件（统一错误等）
|-- internal/db              # MySQL/GORM 初始化
|-- uploads/                 # 默认上传目录（运行时自动创建）
|-- .env.example             # 配置示例
```

## 后续计划
- 接入用户鉴权、JWT 登录。
- 博客文章/分类/标签 CRUD 与分页检索。
- 上传文件与评论审核链路。
- 补充单元测试与 CI 检查。
