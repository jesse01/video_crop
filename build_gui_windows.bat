@echo off
REM Windows GUI版本编译脚本
REM 需要先安装Go和GCC（TDM-GCC推荐）

echo ========================================
echo   视频黑边裁剪工具 - Windows GUI编译
echo ========================================
echo.

REM 检查Go是否安装
where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未找到Go，请先安装Go
    echo 下载地址: https://golang.org/dl/
    pause
    exit /b 1
)

REM 检查GCC是否安装
where gcc >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未找到GCC，GUI版本需要CGO支持
    echo.
    echo 请安装TDM-GCC: https://jmeubank.github.io/tdm-gcc/
    echo 安装时选择 "Create" 绿色部分
    pause
    exit /b 1
)

echo [1/3] 检测到编译环境...
go version
gcc --version | findstr gcc
echo.

echo [2/3] 开始编译GUI版本...
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

go build -tags gui -ldflags="-s -w -H windowsgui" -o video_crop_gui.exe .

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo [失败] 编译出错！
    pause
    exit /b 1
)

echo.
echo [3/3] 编译成功！
echo.
echo ----------------------------------------
echo   生成文件: video_crop_gui.exe
echo ----------------------------------------
echo.

REM 显示文件大小
for %%A in (video_crop_gui.exe) do echo 文件大小: %%~zA 字节
echo.

echo 是否现在运行程序？ (Y/N)
choice /c YN /n
if %ERRORLEVEL% EQU 1 (
    start video_crop_gui.exe
)

echo.
echo 编译完成！
pause
