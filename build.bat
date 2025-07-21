@echo off
setlocal enabledelayedexpansion

:: =============================================================================
::  Windows Build Script for slowsql-analysis
::
::  This script cross-compiles the Go application for various platforms,
::  similar to the provided build.sh script.
:: =============================================================================

:: =============================================================================
::  Main execution block
:: =============================================================================
CALL :main
GOTO :EOF

:: =============================================================================
::  Subroutines (Functions)
:: =============================================================================

:main
    ECHO [INFO] Starting build for slowsql-analysis...
    ECHO.

    CALL :check_environment
    IF !ERRORLEVEL! NEQ 0 (
        ECHO [ERROR] Environment check failed. Aborting.
        EXIT /B 1
    )

    CALL :clean_build
    IF !ERRORLEVEL! NEQ 0 (
        ECHO [ERROR] Build cleanup failed. Aborting.
        EXIT /B 1
    )

    ECHO [INFO] Starting cross-compilation...
    ECHO.
    
    :: Build for different platforms
    CALL :build_binary linux amd64 slowsql-analysis-linux-amd64
    IF !ERRORLEVEL! NEQ 0 EXIT /B 1

    CALL :build_binary linux arm64 slowsql-analysis-linux-arm64
    IF !ERRORLEVEL! NEQ 0 EXIT /B 1

    CALL :build_binary windows amd64 slowsql-analysis-windows-amd64.exe
    IF !ERRORLEVEL! NEQ 0 EXIT /B 1
    
    ECHO.
    ECHO [SUCCESS] All builds completed!
    ECHO [INFO] Build artifacts are in the 'build' directory:
    dir build /B
    ECHO.

    ECHO [INFO] --- Usage Notes ---
    ECHO 1. Linux AMD64:   build\slowsql-analysis-linux-amd64
    ECHO 2. Linux ARM64:   build\slowsql-analysis-linux-arm64
    ECHO 3. Windows:       build\slowsql-analysis-windows-amd64.exe
    ECHO.
    GOTO :EOF


:check_environment
    ECHO [INFO] Checking system environment...
    
    :: Check OS Version
    ver

    :: Check for Go installation
    go version >nul 2>&1
    IF !ERRORLEVEL! NEQ 0 (
        ECHO [ERROR] Go command not found. Please install Go.
        ECHO [INFO]  You can download it from https://golang.org/dl/
        EXIT /B 1
    )
    
    :: Display Go version
    ECHO [INFO] Found Go version:
    go version
    
    :: Check CGO
    IF "%CGO_ENABLED%"=="1" (
        ECHO [WARNING] CGO_ENABLED is 1. It will be disabled during build to ensure static linking.
    )
    
    ECHO [SUCCESS] Environment check passed.
    ECHO.
    GOTO :EOF


:clean_build
    ECHO [INFO] Cleaning up previous build...
    IF EXIST build (
        ECHO [INFO] Removing existing 'build' directory.
        RD /S /Q build
    )
    MKDIR build
    ECHO [SUCCESS] 'build' directory is ready.
    ECHO.
    GOTO :EOF


:build_binary
    set "goos=%1"
    set "goarch=%2"
    set "output_name=%3"
    
    ECHO [INFO] Compiling for %goos%/%goarch%...
    
    :: Set build environment variables
    set CGO_ENABLED=0
    set GOOS=%goos%
    set GOARCH=%goarch%
    
    :: Get current timestamp for build info using WMIC for reliability
    for /f "tokens=2 delims==" %%I in ('wmic os get localdatetime /value') do set "localdatetime=%%I"
    set "build_time=!localdatetime:~0,4!-!localdatetime:~4,2!-!localdatetime:~6,2! !localdatetime:~8,2!:!localdatetime:~10,2!:!localdatetime:~12,2!"

    :: Build the command with version info and build flags
    go build -ldflags="-s -w -X main.Version=1.0.0 -X 'main.BuildTime=!build_time!'" -o "build\%output_name%"
    
    :: Check for build errors
    IF !ERRORLEVEL! NEQ 0 (
        ECHO [ERROR] Failed to compile for %goos%/%goarch%.
        EXIT /B 1
    )
    
    ECHO [SUCCESS] Successfully compiled: build\%output_name%
    
    :: Get file size
    for %%F in ("build\%output_name%") do (
        ECHO [INFO]   ^> File size: %%~zF bytes
    )
    ECHO.
    
    GOTO :EOF 
