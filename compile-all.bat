@echo off
set GOFILE=run.go

echo Compiling for Windows (64-bit)...
set GOOS=windows
set GOARCH=amd64
go build -o run-windows.exe %GOFILE%
upx run-windows.exe

echo Compiling for Linux (64-bit)...
set GOOS=linux
set GOARCH=amd64
go build -o run-linux %GOFILE%
upx run-linux

echo Compiling for macOS (64-bit)...
set GOOS=darwin
set GOARCH=amd64
go build -o run-mac %GOFILE%
upx run-mac


echo Compilation finished!
