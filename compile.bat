@echo off
:: Set the Go environment variables for Windows
set GOOS=windows
set GOARCH=amd64
set GO111MODULE=auto
:: Compile the Go program
go build -o run.exe run.go

:: Check if the build was successful
if %errorlevel%==0 (
    echo Compilation successful! The output is run.exe
upx run.exe
) else (
    echo Compilation failed!
    pause
)

