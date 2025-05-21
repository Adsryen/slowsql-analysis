# slowsql-analysis

> This project is optimized and enhanced based on [kbnote/slowsql-analysis](https://github.com/kbnote/slowsql-analysis).

A MySQL slow query log analysis tool based on pt-query-digest, providing a user-friendly Web interface to display analysis results.

[中文版](README.md)

## Table of Contents

- [Features](#features)
- [System Requirements](#system-requirements)
- [Quick Start](#quick-start)
- [Detailed Configuration](#detailed-configuration)
- [Performance Metrics](#performance-metrics)
- [Usage Scenarios](#usage-scenarios)
- [Troubleshooting](#troubleshooting)
- [Development Guide](#development-guide)

## Features

- Support analysis of slow query logs within specified time ranges
- Generate beautiful HTML reports for analysis results
- Support web server for online report viewing
- Automatically identify and group similar SQL queries
- Provide detailed query performance metrics statistics
- Support one-click SQL statement copying
- Automatically mark different performance levels based on query time
- Support multi-platform operation (Linux/Windows/macOS)
- Built-in pt-query-digest tool, no additional installation required
- Automatic system environment and dependency detection
- Support UTF-8 encoded log files

## System Requirements

### Basic Requirements
- Operating System: Linux, Windows or macOS
- Perl Runtime Environment (5.10 or higher)
- Memory: 2GB or more recommended
- Disk Space: At least 100MB free space

### Perl Module Dependencies
- DBI
- DBD::mysql
- Digest::MD5
- Time::HiRes
- IO::Socket::SSL
- Term::ReadKey

### Installing Dependencies

CentOS/RHEL:
```bash
yum install -y perl-DBI perl-DBD-MySQL perl-Time-HiRes perl-IO-Socket-SSL perl-Digest-MD5 perl-TermReadKey
```

Ubuntu/Debian:
```bash
apt-get install -y libdbi-perl libdbd-mysql-perl libtime-hires-perl libio-socket-ssl-perl libdigest-md5-perl libterm-readkey-perl
```

## Quick Start

### 1. Download and Installation

Download the binary file for your platform from the Release page:

| Platform | Filename |
|----------|----------|
| Linux AMD64 | `slowsql-analysis-linux-amd64` |
| Linux ARM64 | `slowsql-analysis-linux-arm64` |
| Windows | `slowsql-analysis-windows-amd64.exe` |
| macOS Intel | `slowsql-analysis-darwin-amd64` |
| macOS M1/M2 | `slowsql-analysis-darwin-arm64` |

### 2. Basic Usage

```bash
# Add execution permissions (Linux/macOS)
chmod +x slowsql-analysis-*

# Analyze slow query log
./slowsql-analysis -f <slow query log path>

# Start web server (default port 6033)
./slowsql-analysis -f <slow query log path> -port 6033
```

## Detailed Configuration

### Command Line Parameters

| Parameter | Description | Required | Default | Example |
|-----------|-------------|----------|---------|---------|
| -f | Slow query log file path | Yes | - | `/var/log/mysql-slow.log` |
| -port | Web service port | No | 6033 | `8080` |
| -startTime | Start time | No | - | `2024-04-16 00:00:00` |
| -endTime | End time | No | - | `2024-04-16 23:59:59` |

## Performance Metrics

### Key Metrics Explained

1. **Query Time Metrics**
   - `Query_time`: SQL execution time
   - `95%`: 95th percentile query execution time
   - `99%`: 99th percentile query execution time
   - `Max_Query_Time`: Maximum query time

2. **Lock Metrics**
   - `Lock_time`: Lock wait time
   - `Rows_examined`: Number of rows scanned
   - `Rows_sent`: Number of rows returned

3. **Performance Levels**
   - 🟢 Good: < 1 second
   - 🟡 Warning: 1-5 seconds
   - 🔴 Critical: > 5 seconds

## Usage Scenarios

### 1. Daily Monitoring
```bash
# Analyze previous day's slow queries at midnight
0 1 * * * /path/to/slowsql-analysis -f /var/log/mysql-slow.log -startTime="$(date -d 'yesterday' +'%Y-%m-%d 00:00:00')" -endTime="$(date -d 'yesterday' +'%Y-%m-%d 23:59:59')"
```

### 2. Performance Optimization
```bash
# Analyze slow queries for a specific time period
./slowsql-analysis -f /var/log/mysql-slow.log -startTime="2024-04-16 10:00:00" -endTime="2024-04-16 12:00:00"
```

### 3. Real-time Monitoring
```bash
# Start web server for continuous monitoring
./slowsql-analysis -f /var/log/mysql-slow.log -port 6033
```

## Troubleshooting

### Common Issues

1. **Service Won't Start**
   - Check if the port is already in use
   - Verify sufficient permissions
   - Check firewall settings

2. **Empty Analysis Report**
   - Verify log file permissions
   - Check if log format is correct
   - Verify time range settings

3. **Performance Issues**
   - Recommended log file size < 1GB
   - Avoid analyzing too long time ranges
   - Consider increasing system memory

### Log Information

Program logs are located at:
- Linux/macOS: `/var/log/slowsql-analysis.log`
- Windows: `C:\ProgramData\slowsql-analysis\logs\`

## Usage Examples

1. Analyze slow queries from the last day:
```bash
./slowsql-analysis -f /var/log/mysql-slow.log -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"
```

2. Start web server to view reports:
```bash
./slowsql-analysis -f /var/log/mysql-slow.log -port 6033
```

After generating the report, access it through your browser at `http://<server-ip>:6033/<report-filename>`.

## Notes

1. Ensure pt-query-digest tool is installed on your system
2. Running directory must contain `cmd/pt-query-digest` and `template/template.html` files
3. Ensure read permissions for the slow query log file
4. In web server mode, ensure the specified port is not in use

## Dependencies

- Bootstrap 3.3.7
- jQuery 3.3.1
- clipboard.js 2.0.8

## Build Instructions

### Prerequisites

- Go 1.22 or higher
- pt-query-digest tool
- MySQL (for generating slow query logs)

### Installing Dependencies

```bash
# Install pt-query-digest
# Ubuntu/Debian
sudo apt-get install percona-toolkit

# CentOS/RHEL
sudo yum install percona-toolkit

# Install Go dependencies
go mod download
```

### Building the Project

```bash
# Compile the project
go build -o slowsql-analysis

# Run tests
go test ./...
```

### Development Mode

```bash
go run main.go -f <slow query log path>
``` 