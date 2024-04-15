# cLink

`New` generation logging solution.

Based on kafka's natural concurrency advantages and Go portability, the developed log synchronizer is simpler than Flink

- No messages lost
- Simple little configuration
- Quick deployment and operation

# Dependency requirements

- `github.com/ClickHouse/clickhouse-go`
- `github.com/fatih/color`
- `github.com/IBM/sarama`

Minimum Supported Golang Version is 1.20.

# Getting started

**Download package**
[latest version 1.0.1](https://github.com/xxscloud5722/cLink/releases)

**Program compilation**

- Goalng 1.20

```bash
# windows OR linux
go env -w GOOS=linux
go mod tidy
go build -o ./dist/log_sync src/main.go
```

# Configuration instructions

```yaml
kafka:
  server: kafka.example.com
  port: 9094
  sasl:
    username: your username
    password: your password
  consumer:
    group-id: your group
    topic:
      - your topic
clickhouse:
  server: your clickhouse server
  port: 9000
  username: your username
  password: your password
  database: your database
  table: your table
  fields:
    system_code: system_code
    date: date
    env: env
    service_code: service_code
    level: level
    thread: thread
    class: class
    message: message
    trace_id: trace_id
pattern: '(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+) (\w+) \[([^]]+)\] \[([^]]+)\] ([^:]+) : (.+[\s\S]*)'
pattern-fields:
  - date
  - level
  - thread
  - trace_id
  - class
  - message

```

# Docker

```bash
# pull image
docker pull xxscloud5722/cLink:1.0.1
```

# Contributors

Thanks for your contributions!

- [@xiaoliuya](https://github.com/xxscloud5722/)


# Case-CN

- [Clickhouse 日志采集存储方案 （代号：明月长老）语雀](https://www.yuque.com/mcat/uggxu0/mfgouabbgg5rs8rq#OHjNq)

# Zen

Don't let what you cannot do interfere with what you can do.
