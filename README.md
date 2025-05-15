# slowsql-analysis

基于 pt-query-digest 工具的 MySQL 慢查询日志分析工具，提供友好的 Web 界面展示分析结果。

## 功能特性

- 支持分析指定时间范围的慢查询日志
- 提供美观的 HTML 报告展示分析结果
- 支持启动 Web 服务器在线查看报告
- 自动识别并分组相似的 SQL 查询
- 提供详细的查询性能指标统计
- 支持 SQL 语句的一键复制
- 根据查询时间自动标记不同性能等级

## 分析指标

- 查询执行次数统计
- 查询时间分析（最大、最小、平均、95%）
- 扫描行数统计
- 锁等待时间分析
- 涉及数据表统计
- 来源用户和主机信息

## 使用方法

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
