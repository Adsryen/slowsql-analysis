#!/bin/bash

# 设置颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${BLUE}[信息] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[成功] $1${NC}"
}

print_error() {
    echo -e "${RED}[错误] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[警告] $1${NC}"
}

# 检查系统环境
check_environment() {
    print_info "检查系统环境..."
    
    # 检查操作系统
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        print_info "当前系统: $NAME $VERSION_ID"
    fi
    
    # 检查是否安装了Go
    if ! command -v go &> /dev/null; then
        print_error "未找到Go环境，请先安装Go"
        print_info "可以访问 https://golang.org/dl/ 下载安装"
        exit 1
    fi
    
    # 显示Go版本
    GO_VERSION=$(go version)
    print_info "Go版本: $GO_VERSION"
    
    # 检查GOPATH
    if [ -z "$GOPATH" ]; then
        print_warning "GOPATH未设置，将使用默认值"
    else
        print_info "GOPATH: $GOPATH"
    fi
    
    # 检查CGO
    if [ "$CGO_ENABLED" = "1" ]; then
        print_warning "CGO已启用，将被禁用以确保静态链接"
    fi
}

# 清理构建目录
clean_build() {
    print_info "清理之前的构建文件..."
    if [ -d "build" ]; then
        rm -rf build/*
        print_success "清理完成"
    else
        mkdir -p build
        print_success "创建build目录"
    fi
}

# 编译函数
build_binary() {
    local os=$1
    local arch=$2
    local output=$3
    
    print_info "正在编译 ${os}/${arch} 版本..."
    
    # 设置编译环境变量
    export CGO_ENABLED=0
    export GOOS=$os
    export GOARCH=$arch
    
    # 执行编译，添加版本信息和构建时间
    if go build -ldflags="-s -w -X main.Version=1.0.0 -X 'main.BuildTime=$(date)'" -o "build/${output}"; then
        print_success "编译成功: build/${output}"
        # 检查文件大小
        local size=$(ls -lh "build/${output}" | awk '{print $5}')
        print_info "文件大小: ${size}"
        # 添加执行权限
        chmod +x "build/${output}"
    else
        print_error "编译失败: ${os}/${arch}"
        return 1
    fi
}

# 主函数
main() {
    print_info "开始构建慢查询日志分析工具..."
    
    # 检查环境
    check_environment
    
    # 清理构建目录
    clean_build
    
    # 编译不同平台的版本
    build_binary linux amd64 "slowsql-analysis-linux-amd64" || exit 1
    build_binary linux arm64 "slowsql-analysis-linux-arm64" || exit 1
    build_binary windows amd64 "slowsql-analysis-windows-amd64.exe" || exit 1
    build_binary darwin amd64 "slowsql-analysis-darwin-amd64" || exit 1
    build_binary darwin arm64 "slowsql-analysis-darwin-arm64" || exit 1
    
    print_success "编译完成！"
    print_info "编译结果在 build 目录下:"
    ls -lh build/
    
    print_info "使用说明:"
    print_info "1. Linux AMD64 版本: build/slowsql-analysis-linux-amd64"
    print_info "2. Linux ARM64 版本: build/slowsql-analysis-linux-arm64"
    print_info "3. Windows 版本: build/slowsql-analysis-windows-amd64.exe"
    print_info "4. macOS Intel 版本: build/slowsql-analysis-darwin-amd64"
    print_info "5. macOS M1/M2 版本: build/slowsql-analysis-darwin-arm64"
    
    print_warning "注意: 如果在 CentOS 上运行遇到问题，请尝试:"
    print_info "1. 确保有执行权限: chmod +x slowsql-analysis-linux-amd64"
    print_info "2. 如果有 SELinux，临时关闭: sudo setenforce 0"
    print_info "3. 使用 strace 调试: strace ./slowsql-analysis-linux-amd64 ..."
}

# 执行主函数
main 