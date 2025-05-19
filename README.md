# slowsql-analysis

> 本项目基于 [kbnote/slowsql-analysis](https://github.com/kbnote/slowsql-analysis) 进行优化和功能扩展。

基于 pt-query-digest 工具的 MySQL 慢查询日志分析工具，提供友好的 Web 界面展示分析结果。

[English Version](README_EN.md)

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

## 分析指标

### 基本要求
- 操作系统：Linux、Windows 或 macOS
- Perl 运行环境（5.10 或更高版本）

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

## 使用方法

### 下载和安装

1. 从 Release 页面下载对应平台的二进制文件：
   - Linux AMD64: `slowsql-analysis-linux-amd64`
   - Linux ARM64: `slowsql-analysis-linux-arm64`
   - Windows: `slowsql-analysis-windows-amd64.exe`
   - macOS Intel: `slowsql-analysis-darwin-amd64`
   - macOS M1/M2: `slowsql-analysis-darwin-arm64`

2. 添加执行权限（Linux/macOS）：
```bash
chmod +x slowsql-analysis-*
```

### 基本用法

```bash
./slowsql-analysis -f <慢查询日志路径>
```

### 指定时间范围分析

```bash
./slowsql-analysis -f <慢查询日志路径> -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"
```

### 启动 Web 服务

```bash
./slowsql-analysis -f <慢查询日志路径> -port 6033
```

### 完整参数说明

```
参数:
    -f          慢查询日志文件路径（必需）
    -port       Web服务端口，设置后可通过浏览器访问报告（可选）
    -startTime  开始时间，格式：yyyy-mm-dd HH:mm:ss（可选）
    -endTime    结束时间，格式：yyyy-mm-dd HH:mm:ss（可选）
```

## 报告说明

生成的 HTML 报告包含以下内容：

1. 慢查询概览表格
   - 按 95% 执行时间降序排列
   - 根据执行时间自动标记性能等级
   - 支持查看详细 SQL 信息

2. SQL 详情弹窗
   - 完整 SQL 语句（支持一键复制）
   - 查询执行统计信息
   - 涉及数据表列表
   - 性能指标详细统计

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
