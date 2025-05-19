# slowsql-analysis

> This project is optimized and enhanced based on [kbnote/slowsql-analysis](https://github.com/kbnote/slowsql-analysis).

A MySQL slow query log analysis tool based on pt-query-digest, providing a user-friendly Web interface to display analysis results.

[中文版](README.md)

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

## Analysis Metrics

### Basic Requirements
- Operating System: Linux, Windows or macOS
- Perl Runtime Environment (5.10 or higher)

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

## Usage

### Download and Installation

1. Download the binary file for your platform from the Release page:
   - Linux AMD64: `slowsql-analysis-linux-amd64`
   - Linux ARM64: `slowsql-analysis-linux-arm64`
   - Windows: `slowsql-analysis-windows-amd64.exe`
   - macOS Intel: `slowsql-analysis-darwin-amd64`
   - macOS M1/M2: `slowsql-analysis-darwin-arm64`

2. Add execution permissions (Linux/macOS):
```bash
chmod +x slowsql-analysis-*
```

### Basic Usage

```bash
./slowsql-analysis -f <slow query log path>
```

### Analyze Specific Time Range

```bash
./slowsql-analysis -f <slow query log path> -startTime="2024-04-16 00:00:00" -endTime="2024-04-16 23:59:59"
```

### Start Web Server

```bash
./slowsql-analysis -f <slow query log path> -port 6033
```

### Complete Parameter Description

```
Parameters:
    -f          Slow query log file path (required)
    -port       Web service port for browser access to reports (optional)
    -startTime  Start time, format: yyyy-mm-dd HH:mm:ss (optional)
    -endTime    End time, format: yyyy-mm-dd HH:mm:ss (optional)
```

## Report Description

The generated HTML report includes:

1. Slow Query Overview Table
   - Sorted by 95th percentile execution time
   - Performance levels automatically marked based on execution time
   - Detailed SQL information available

2. SQL Details Modal
   - Complete SQL statement (with one-click copy)
   - Query execution statistics
   - List of involved tables
   - Detailed performance metrics

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