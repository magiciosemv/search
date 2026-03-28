# Solana Wallet Monitor

实时监控 Solana 钱包余额变化，支持多渠道告警（Telegram、Email）。

## 功能特性

### ✅ 已完成功能

**后端 (Go + Gin)**
- SQLite 数据库存储（地址、规则、通知、告警）
- Solana RPC 集成（获取钱包余额）
- REST API 完整 CRUD 接口
- 后台监控服务（30秒轮询）
- Telegram 通知发送
- Email (SMTP) 通知发送

**前端 (Vue 3 + Tailwind)**
- Dashboard 统计面板（监控状态、总地址数、告警数）
- 地址管理（添加、删除、刷新余额）
- 规则管理（创建、启用/禁用、删除）
- 通知渠道管理（添加、测试、删除）
- 告警历史记录

**API 接口**
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/addresses | 获取地址列表 |
| POST | /api/addresses | 添加钱包地址 |
| GET | /api/addresses/:id | 获取单个地址 |
| PUT | /api/addresses/:id | 更新地址 |
| DELETE | /api/addresses/:id | 删除地址 |
| POST | /api/addresses/:id/refresh | 手动刷新余额 |
| GET | /api/rules | 获取规则列表 |
| POST | /api/rules | 创建规则 |
| PUT | /api/rules/:id | 更新规则 |
| DELETE | /api/rules/:id | 删除规则 |
| GET | /api/notifications | 获取通知列表 |
| POST | /api/notifications | 添加通知渠道 |
| PUT | /api/notifications/:id | 更新通知 |
| DELETE | /api/notifications/:id | 删除通知 |
| POST | /api/notifications/:id/test | 测试通知发送 |
| GET | /api/alerts | 获取告警历史 |
| GET | /api/alerts/stats | 获取告警统计 |
| GET | /api/stats | 获取统计信息 |
| GET | /api/health | 健康检查 |

**告警规则类型**
- `balance_change` - 余额变化 > X SOL
- `balance_change_percent` - 余额变化 > X%
- `threshold_above` - 余额 > X SOL
- `threshold_below` - 余额 < X SOL
- `large_incoming` - 转入 > X SOL
- `large_outgoing` - 转出 > X SOL

---

### ❌ 缺失功能

| 功能 | 说明 |
|------|------|
| Docker 构建配置 | 需要完善 Dockerfile 和构建脚本 |
| 配置文件支持 | 未实现 config.yaml 配置文件 |
| 输入验证 | 地址格式、阈值范围未验证 |
| Token 余额查询 | 仅支持 SOL 余额，未支持 SPL Token |
| 多链支持 | 仅支持 Solana |
| WebSocket 实时推送 | 需轮询获取告警 |
| 用户系统 | 无登录注册功能 |
| Discord/Webhook | 代码提到但未实现 |

---

### ⚠️ 已知问题

1. **CGO 依赖** - `github.com/mattn/go-sqlite3` 需要 CGO 编译，Windows Docker 构建需要 gcc
2. **WSL2 网络** - Telegram API 在 WSL2 下可能超时，需在 Windows 端运行
3. **前端代理** - 开发模式下需配置 Vite 代理（已配置）

---

## 快速开始

### 方式一：Docker 部署（推荐）

```bash
# 1. 复制环境配置
cp .env.example .env

# 2. 编辑 .env 文件，配置以下必填项：
# - RPC_URL: Solana RPC 节点（如 https://mainnet.helius-rpc.com/?api-key=YOUR_KEY）
# - TELEGRAM_BOT_TOKEN: Telegram Bot Token（可选）
# - SMTP 相关配置（可选）

# 3. 启动服务
docker-compose up -d

# 4. 访问
# 前端: http://localhost:3000
# 后端: http://localhost:8080
```

### 方式二：本地开发运行

**后端：**
```bash
cd backend

# 安装 Go 依赖
go mod download

# 运行（需要配置环境变量）
export PORT=8080
export RPC_URL="https://mainnet.helius-rpc.com/?api-key=YOUR_KEY"
export DATABASE_PATH="./data/monitor.db"
export TELEGRAM_BOT_TOKEN="YOUR_BOT_TOKEN"

go run cmd/server/main.go
```

**前端：**
```bash
cd frontend

# 安装依赖
npm install

# 开发模式运行（已配置 API 代理）
npm run dev

# 访问 http://localhost:5173
```

---

## 环境变量配置

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| PORT | 是 | 服务端口 | 8080 |
| DATABASE_PATH | 是 | 数据库文件路径 | ./data/monitor.db |
| RPC_URL | 是 | Solana RPC 节点 | https://api.mainnet-beta.solana.com |
| RPC_API_KEY | 否 | RPC API Key | - |
| TELEGRAM_BOT_TOKEN | 否 | Telegram Bot Token | - |
| SMTP_HOST | 否 | SMTP 服务器 | smtp.gmail.com |
| SMTP_PORT | 否 | SMTP 端口 | 587 |
| SMTP_USER | 否 | SMTP 用户名 | - |
| SMTP_PASS | 否 | SMTP 密码 | - |
| SMTP_FROM | 否 | 发件人地址 | noreply@example.com |

---

## 项目结构

```
AIvibe/
├── backend/                  # Go 后端
│   ├── cmd/server/          # 入口文件
│   ├── internal/
│   │   ├── config/         # 配置加载
│   │   ├── handlers/       # HTTP 处理
│   │   ├── models/         # 数据模型 & 数据库
│   │   └── services/       # 业务逻辑
│   ├── go.mod
│   └── Dockerfile
│
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── views/          # 页面组件
│   │   ├── components/     # 公共组件
│   │   ├── assets/         # 静态资源
│   │   ├── App.vue         # 根组件
│   │   └── main.js         # 入口文件
│   ├── vite.config.js      # Vite 配置（含代理）
│   ├── package.json
│   └── Dockerfile
│
├── data/                    # 数据目录（SQLite 数据库）
├── docker-compose.yml       # Docker 编排
├── .env                     # 环境变量
└── README.md
```

---

## 后续开发计划

### 短期
- [ ] 完善 Docker 构建（配置多阶段构建）
- [ ] 添加 config.yaml 配置文件支持
- [ ] 输入验证增强（Solana 地址格式校验）
- [ ] 错误处理完善

### 中期
- [ ] SPL Token 余额查询
- [ ] WebSocket 实时推送告警
- [ ] Discord/Webhook 通知渠道
- [ ] 余额历史图表

### 长期
- [ ] 用户系统（登录注册）
- [ ] 多链支持（Ethereum、Tron）
- [ ] 移动端适配
- [ ] 多语言支持

---

## 常见问题

**Q: Docker 构建失败？**
A: 检查网络是否能访问 Docker Hub，可配置国内镜像加速。

**Q: Telegram 通知发不出去？**
A: WSL2 网络限制导致外网请求超时，建议在 Windows 端运行 Docker。

**Q: 数据库报错 sqlite3 unknown driver？**
A: 确保 `go.mod` 中有 `_ "github.com/mattn/go-sqlite3"` 导入。

**Q: 前端无法请求后端 API？**
A: 开发模式下 Vite 已配置代理，生产环境需配置 Nginx 反向代理。

---

## 技术栈

- **后端**: Go 1.21, Gin, SQLite (go-sqlite3)
- **前端**: Vue 3, Tailwind CSS, Vite
- **数据库**: SQLite
- **RPC**: Helius / Solana 官方 RPC
- **通知**: Telegram Bot API, SMTP