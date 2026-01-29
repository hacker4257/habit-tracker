# Habit Tracker

一个简洁的个人习惯追踪应用，支持日历视图和活动频率热力图。

## 技术栈

### 后端
- **Go 1.21+** - 高性能后端语言
- **SQLite/MySQL** - 可切换的数据库支持
- **标准库 net/http** - 轻量级HTTP服务

### 前端
- **React 18** - 现代化UI框架
- **纯CSS** - 无额外依赖

### 架构设计

```
backend/
├── cmd/server/          # 应用入口
├── internal/
│   ├── config/          # 配置管理
│   ├── handler/         # HTTP处理器
│   ├── middleware/      # 中间件（CORS、日志、恢复）
│   ├── model/           # 数据模型
│   ├── repository/      # 数据访问层
│   └── service/         # 业务逻辑层
└── pkg/logger/          # 日志工具
```

## 快速开始

### 本地开发

```bash
# 后端
cd backend
go mod download
go run ./cmd/server

# 前端
cd frontend
npm install
npm start
```

### Docker 部署

```bash
docker-compose up -d
```

访问 http://localhost:3000

## 配置

通过环境变量配置：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| SERVER_PORT | 8080 | 服务端口 |
| CORS_ORIGINS | * | CORS允许的源 |
| DB_DRIVER | sqlite | 数据库类型 (sqlite/mysql) |
| DB_DSN | data.db | 数据库连接字符串 |

### MySQL 配置示例

```bash
export DB_DRIVER=mysql
export DB_DSN="user:password@tcp(localhost:3306)/habit_tracker?parseTime=true"
```

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/records | 获取所有记录 |
| POST | /api/records | 创建记录 |
| GET | /api/records/:id | 获取单条记录 |
| PUT | /api/records/:id | 更新记录 |
| DELETE | /api/records/:id | 删除记录 |
| GET | /api/stats | 获取统计数据 |
| GET | /health | 健康检查 |

## 功能特性

- 日历视图：按月浏览，标记有记录的日期
- 频率热力图：类似GitHub贡献图，展示活动频率
- 统计面板：总记录数、总时长、本周/本月统计
- 响应式设计：支持移动端访问
- 数据持久化：SQLite（默认）或MySQL

## 为什么不用 Redis/Kafka？

这是一个个人追踪应用，设计原则是**够用就好**：

- **无高并发**：单用户场景，无需缓存层
- **无消息队列需求**：同步操作足够，无异步处理需求
- **简单部署**：SQLite零配置，单文件数据库

如需扩展为多用户SaaS，可考虑添加：
- Redis：会话管理、API限流
- MySQL：多用户数据隔离
- 用户认证系统

## License

MIT
