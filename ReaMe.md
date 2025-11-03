## 一、推荐 Go 后端目录结构

```
backend/
├─ main.go                    # 主程序入口
├─ config/                    # 配置文件
│   └─ config.go              # 加载 Redis / WebSocket / API 配置
├─ internal/
│   ├─ api/                   # HTTP REST API
│   │   ├─ router.go          # Gin 路由配置
│   │   ├─ task_handler.go    # 任务相关接口（提交任务、获取历史）
│   │   └─ health.go          # 健康检查接口
│   ├─ task/                  # 任务管理模块
│   │   ├─ manager.go         # TaskManager，调度任务
│   │   ├─ runner.go          # Runner 接口定义
│   │   ├─ nmap_runner.go     # Nmap Runner
│   │   ├─ dirsearch_runner.go# Dirsearch Runner
│   │   └─ task_model.go      # 任务结构体、状态定义
│   ├─ queue/                 # Redis Stream / MQ 封装
│   │   ├─ redis_client.go    # Redis 客户端封装
│   │   ├─ stream.go          # Stream 生产 / 消费封装
│   │   └─ consumer.go        # 消费者组逻辑
│   ├─ ws/                    # WebSocket 服务
│   │   ├─ hub.go             # WebSocket Hub，管理客户端 & 广播
│   │   └─ handler.go         # 处理 WebSocket 连接
│   ├─ log/                   # 日志管理
│   │   └─ logger.go          # 全局日志封装
│   └─ utils/                 # 通用工具函数
│       └─ uuid.go            # 生成 taskId 等
├─ scripts/                   # 脚本文件（如 bash / python 工具调用）
├─ go.mod
└─ go.sum
```

------

## 二、模块职责说明

| 模块             | 说明                                                         |
| ---------------- | ------------------------------------------------------------ |
| `cmd/`           | 后端程序入口，main.go 初始化配置、连接 Redis、启动 API 和 WebSocket 服务 |
| `config/`        | 读取配置文件或环境变量（Redis 地址、WebSocket 端口、任务参数等） |
| `internal/api`   | REST API，供前端提交任务、查询任务状态、获取历史输出         |
| `internal/task`  | 核心任务管理逻辑，负责任务调度、执行、状态维护               |
| `internal/queue` | Redis Stream 封装，处理任务队列和日志流                      |
| `internal/ws`    | WebSocket Hub，实现实时推送任务输出给前端                    |
| `internal/log`   | 日志管理模块，记录系统日志、任务异常                         |
| `internal/utils` | 工具函数（如生成唯一 taskId）                                |
| `scripts/`       | 外部工具命令封装（nmap、dirsearch 脚本或二进制）             |

------

## 三、模块之间的交互示意

```
前端（Vue3） 
     │ REST API 提交任务
     ▼
internal/api/task_handler.go
     │
     ▼
internal/task/manager.go（TaskManager）
     │ 创建 taskId, 调度 Runner
     ▼
internal/task/nmap_runner.go
     │ 执行命令行, 将输出写入 Redis Stream
     ▼
internal/queue/stream.go
     │
     ▼
internal/ws/hub.go
     │ WebSocket 广播输出给前端
```

------

## 四、可扩展性说明

1. **新增工具**
   - 新建 `internal/task/<tool>_runner.go`
   - 实现 Runner 接口 `Run(taskId string, params map[string]string)`
   - 不需要修改其他模块，TaskManager 会动态调度。
2. **新增队列类型**
   - queue 模块抽象接口
   - Redis Stream / RabbitMQ / Kafka 都可实现接口替换
3. **多 Worker 支持**
   - 每个 Worker 实例只消费队列的任务
   - 消费者组机制保证消息只被一个 Worker 执行
4. **前端刷新恢复**
   - 任务输出写入 Redis Stream，WebSocket Hub 在客户端连接时可回放历史消息

------

## 五、目录使用建议

- `internal/task`：核心逻辑，任何新增扫描工具都放这里
- `internal/queue` + `internal/ws`：保持通用，解耦任务逻辑与传输层
- `api`：尽量保持轻量，只做请求解析与响应
- `scripts`：命令行脚本与工具调用保持独立，不耦合 Go 逻辑