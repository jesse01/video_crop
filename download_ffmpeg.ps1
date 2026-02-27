# FFmpeg 下载脚本 (Windows)
# 用于自动下载 Windows 版本的 ffmpeg

$ErrorActionPreference = "Stop"

$BinDir = "internal\ffmpeg\binaries"
New-Item -ItemType Directory -Force -Path $BinDir | Out-Null

Write-Host "下载 Windows 版本的 ffmpeg..." -ForegroundColor Green

$FFmpegUrl = "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip"
$TmpZip = "$env:TEMP\ffmpeg-windows.zip"
$ExtractDir = "$env:TEMP\ffmpeg-extract"

# 检查是否安装了必要的工具
if (-not (Get-Command curl -ErrorAction SilentlyContinue) -and -not (Get-Command wget -ErrorAction SilentlyContinue)) {
    Write-Host "错误: 需要 curl 或 wget 来下载 ffmpeg" -ForegroundColor Red
    Write-Host "请手动下载并解压 ffmpeg-release-essentials.zip" -ForegroundColor Yellow
    Write-Host "下载地址: $FFmpegUrl" -ForegroundColor Yellow
    exit 1
}

# 下载
if (Get-Command curl -ErrorAction SilentlyContinue) {
    curl -L -o "$TmpZip" "$FFmpegUrl"
} else {
    wget -O "$TmpZip" "$FFmpegUrl"
}

Write-Host "解压文件..." -ForegroundColor Yellow

# 检查是否有 tar 命令（Windows 10+ 内置）
if (Get-Command tar -ErrorAction SilentlyContinue) {
    # 使用 tar 解压
    New-Item -ItemType Directory -Force -Path $ExtractDir | Out-Null
    tar -xf "$TmpZip" -C "$ExtractDir"

    # 查找 ffmpeg.exe
    $FfmpegExe = Get-ChildItem -Path $ExtractDir -Recurse -Filter "ffmpeg.exe" | Select-Object -First 1

    if ($FfmpegExe) {
        Copy-Item $FfmpegExe.FullName "$BinDir\ffmpeg-windows.exe"
        Write-Host "✓ ffmpeg-windows.exe 已下载到 $BinDir\ffmpeg-windows.exe" -ForegroundColor Green
    } else {
        Write-Host "错误: 在解压的文件中找不到 ffmpeg.exe" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "错误: 需要 tar 命令来解压文件（Windows 10+ 内置）" -ForegroundColor Red
    Write-Host "请手动解压 $TmpZip 并将 bin\ffmpeg.exe 复制到 $BinDir\ffmpeg-windows.exe" -ForegroundColor Yellow
    exit 1
}

# 清理临时文件
Remove-Item "$TmpZip" -ErrorAction SilentlyContinue
Remove-Item $ExtractDir -Recurse -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "下载完成！现在可以构建项目了：" -ForegroundColor Green
Write-Host "  go build -o video_crop.exe" -ForegroundColor Yellow
