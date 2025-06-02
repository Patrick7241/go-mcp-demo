# go-mcp-demo

本项目基于 [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go) 实现了一个 `MCP Server` 与 `MCP Client` 通信的演示，结合 [Ollama](https://ollama.com/) 启动本地大语言模型，实现自然语言调用工具的能力。

[b站演示视频](https://www.bilibili.com/video/BV1N87Dz2Efn)

## ✨ 功能概览

- 集成 MCP Server + Client 通信框架
- 支持基于自然语言执行：
  - 读取本地 `.sql` 文件
  - 执行任意 MySQL 查询
- 使用 Ollama 本地大模型进行推理（默认端口：`http://localhost:11434`）

---

## 🚀 快速开始

### 1. 安装 Ollama

请参考官网：[https://ollama.com](https://ollama.com)  
根据操作系统下载安装 Ollama CLI 工具。

---

### 2. 拉取并运行模型

在命令行中执行以下命令，拉取你需要的模型（示例使用 [`qwen2:1.5b`](https://ollama.com/library/qwen2)）：

```bash
ollama pull qwen2:1.5b
ollama run qwen2:1.5b
````

模型将监听本地端口：`http://localhost:11434/api/generate`

---

### 3. 修改配置文件

打开并编辑 `pkg/config/config.yml`：

```yaml
ollama:
  model: qwen2:1.5b

mysql:
  host: localhost
  port: 3306
  user: root
  password: your_password
  database: your_database

sql:
  sqlFilePath: E:/your/path/to/test.sql
```

> ✅ 注意事项：
>
> * `model` 名称需与 `ollama pull` 和 `ollama run` 保持一致
> * `mysql` 信息替换为你本地或远程数据库连接参数
> * `sqlFilePath` 设置为你希望读取的 `.sql` 文件绝对路径

---

### 4. 启动主服务（Client）

进入主服务目录并运行：

```bash
cd cmd/server
go run main.go
```

默认监听地址：`http://localhost:2001`

---

### 5. 启动 MCP 服务（Server）

在新终端中启动 MCP 服务：

```bash
cd cmd/mcp-server
go run main.go
```

默认监听地址：`http://localhost:2002`
服务将自动注册以下工具：

* `read_file`: 读取指定的 `.sql` 文件内容
* `query_db`: 执行任意 SQL 并返回查询结果

---
