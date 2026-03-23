# Go MVC

快速开发的 Go MVC 项目，采用模块化架构，支持未来微服务迁移。

## 特性

- 基于 Gin 框架
- 模块化架构设计
- 懒加载初始化模式
- 支持 MySQL、Redis、JWT
- 任务队列支持（Asynq）

## 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/sy121718/go-mvc.git
cd go-mvc
```

### 2. 配置

```bash
cp config.yaml.example config.yaml
# 编辑 config.yaml，填入你的配置
```

### 3. 安装依赖

```bash
go mod download
```

### 4. 运行

```bash
go run cmd/main.go
```

## 项目结构

```
go-mvc/
├── cmd/main.go             # 启动入口
├── config/                 # 配置解析
├── internal/               # 私有应用代码
│   ├── routers/            # 主路由聚合
│   ├── middleware/         # 全局中间件
│   ├── module/             # 业务模块
│   └── task/               # 任务处理
├── pkg/                    # 通用组件
│   ├── database/           # MySQL
│   ├── cache/              # Redis
│   ├── auth/               # JWT
│   └── crypto/             # 签名
├── storage/                # 文件存储
├── docs/                   # 项目文档
└── config.yaml             # 配置文件
```

## 开发规范

详见 [CLAUDE.md](./CLAUDE.md)

## License

MIT