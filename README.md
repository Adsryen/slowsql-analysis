# slowsql-analysis

> 本项目基于 [kbnote/slowsql-analysis](https://github.com/kbnote/slowsql-analysis) 进行优化和功能扩展。

基于 pt-query-digest 工具的 MySQL 慢查询日志分析工具，提供友好的 Web 界面展示分析结果。

[English Version](README_EN.md)

## 目录

- [功能特性](#功能特性)
- [系统要求](#系统要求)
- [快速开始](#快速开始)
- [详细配置](#详细配置)
- [性能指标说明](#性能指标说明)
- [使用场景](#使用场景)
- [故障排除](#故障排除)
- [开发指南](#开发指南)

## 功能特性

- 支持分析指定时间范围的慢查询日志
- 提供美观的 HTML 报告展示分析结果
- 支持启动 Web 服务器在线查看报告
- 自动识别并分组相似的 SQL 查询
- 提供详细的查询性能指标统计
- 支持 SQL 语句的一键复制
- 根据查询时间自动标记不同性能等级
- 支持多平台运行（Linux/Windows/macOS）
- 内置 pt-query-digest 工具，无需额外安装
- 自动检测系统环境和依赖
- 支持 UTF-8 编码的日志文件

## 系统要求

### 基本环境
- 操作系统：Linux、Windows 或 macOS
- Perl 运行环境（5.10 或更高版本）
- 内存：建议 2GB 以上
- 磁盘空间：至少 100MB 可用空间

### Perl 模块依赖
- DBI
- DBD::mysql
- Digest::MD5
- Time::HiRes
- IO::Socket::SSL
- Term::ReadKey

### 依赖安装

CentOS/RHEL:
```bash
yum install -y perl-DBI perl-DBD-MySQL perl-Time-HiRes perl-IO-Socket-SSL perl-Digest-MD5 perl-TermReadKey
```

Ubuntu/Debian:
```bash
apt-get install -y libdbi-perl libdbd-mysql-perl libtime-hires-perl libio-socket-ssl-perl libdigest-md5-perl libterm-readkey-perl
```

## 快速开始

### 1. 下载安装

从 build文件夹 下载对应平台的二进制文件：

| 平台 | 文件名 |
|------|--------|
| Linux AMD64 | `slowsql-analysis-linux-amd64` |
| Linux ARM64 | `slowsql-analysis-linux-arm64` |
| Windows | `slowsql-analysis-windows-amd64.exe` |
| macOS Intel | `slowsql-analysis-darwin-amd64` |
| macOS M1/M2 | `slowsql-analysis-darwin-arm64` |

### 2. 基本使用

```bash
# Linux/macOS 添加执行权限
chmod +x slowsql-analysis-*

# 分析慢查询日志
./slowsql-analysis -f <慢查询日志路径>

# 启动Web服务（默认端口6033）
./slowsql-analysis -f <慢查询日志路径> -port 6033
```

## 详细配置

### 命令行参数

| 参数 | 说明 | 是否必需 | 默认值 | 示例 |
|------|------|----------|--------|------|
| -f | 慢查询日志文件路径 | 是 | - | `/var/log/mysql-slow.log` |
| -port | Web服务端口 | 否 | 6033 | `8080` |
| -startTime | 开始时间 | 否 | - | `2024-04-16 00:00:00` |
| -endTime | 结束时间 | 否 | - | `2024-04-16 23:59:59` |

## 性能指标说明

### 关键指标解释

1. **Query Time Metrics**
   - `Query_time`: SQL执行时间
   - `95%`: 95%的查询执行时间
   - `99%`: 99%的查询执行时间
   - `Max_Query_Time`: 最长查询时间

2. **Lock Metrics**
   - `Lock_time`: 锁等待时间
   - `Rows_examined`: 扫描行数
   - `Rows_sent`: 返回行数

3. **性能等级**
   - 🟢 良好：< 1秒
   - 🟡 警告：1-5秒
   - 🔴 严重：> 5秒

## 使用场景

### 1. 日常监控
```bash
# 每天凌晨分析前一天的慢查询
0 1 * * * /path/to/slowsql-analysis -f /var/log/mysql-slow.log -startTime="$(date -d 'yesterday' +'%Y-%m-%d 00:00:00')" -endTime="$(date -d 'yesterday' +'%Y-%m-%d 23:59:59')"
```

### 2. 性能优化
```bash
# 分析特定时间段的慢查询
./slowsql-analysis -f /var/log/mysql-slow.log -startTime="2024-04-16 10:00:00" -endTime="2024-04-16 12:00:00"
```

### 3. 实时监控
```bash
# 启动Web服务持续监控
./slowsql-analysis -f /var/log/mysql-slow.log -port 6033
```

## 故障排除

### 常见问题

1. **无法启动服务**
   - 检查端口是否被占用
   - 确认是否有足够权限
   - 验证防火墙设置

2. **分析报告为空**
   - 确认日志文件权限
   - 验证日志格式是否正确
   - 检查时间范围设置

3. **性能问题**
   - 建议日志文件大小 < 1GB
   - 避免分析过长时间范围
   - 考虑增加系统内存

### 日志说明

程序运行日志位于：
- Linux/macOS: `/var/log/slowsql-analysis.log`
- Windows: `C:\ProgramData\slowsql-analysis\logs\`

## 使用示例

1. 分析最近一天的慢查询：
```bash
./slowsql-analysis -f /var/log/mysql-slow.log -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"
```

2. 启动 Web 服务查看报告：
```bash
./slowsql-analysis -f /var/log/mysql-slow.log -port 6033
```

生成报告后，可以通过浏览器访问 `http://<服务器IP>:6033/<报告文件名>` 查看分析结果。

## 注意事项

1. 确保系统中已安装 pt-query-digest 工具
2. 运行目录需包含 `cmd/pt-query-digest` 和 `template/template.html` 文件
3. 确保对慢查询日志文件有读取权限
4. Web 服务模式下需确保指定端口未被占用

## 依赖说明

- Bootstrap 3.3.7
- jQuery 3.3.1
- clipboard.js 2.0.8

## 构建说明

### 环境要求

- Go 1.22 或更高版本
- pt-query-digest 工具
- MySQL（用于生成慢查询日志）

### 安装依赖

```bash
# 安装 pt-query-digest
# Ubuntu/Debian
sudo apt-get install percona-toolkit

# CentOS/RHEL
sudo yum install percona-toolkit

# 安装 Go 依赖
go mod download
```

### 构建项目

```bash
# 编译项目
go build -o slowsql-analysis

# 运行测试
go test ./...
```

### 开发模式运行

```bash
go run main.go -f <慢查询日志路径>
```
