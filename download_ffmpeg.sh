#!/bin/bash

# FFmpeg 下载脚本
# 用于自动下载对应平台的 ffmpeg 二进制文件

set -e

BIN_DIR="internal/ffmpeg/binaries"
mkdir -p "$BIN_DIR"

detect_platform() {
    case "$(uname -s)" in
        Darwin)
            echo "macos"
            ;;
        Linux)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            echo "unknown"
            ;;
    esac
}

PLATFORM=$(detect_platform)
echo "检测到平台: $PLATFORM"

case $PLATFORM in
    macos)
        echo "下载 macOS 版本的 ffmpeg..."
        # 使用静态编译版本
        FFMPEG_URL="https://evermeet.cx/ffmpeg/getrelease/zip"
        TMP_ZIP="/tmp/ffmpeg-mac.zip"

        if command -v curl >/dev/null 2>&1; then
            curl -L -o "$TMP_ZIP" "$FFMPEG_URL"
        elif command -v wget >/dev/null 2>&1; then
            wget -O "$TMP_ZIP" "$FFMPEG_URL"
        else
            echo "错误: 需要 curl 或 wget 来下载 ffmpeg"
            exit 1
        fi

        unzip -o "$TMP_ZIP" -d "$BIN_DIR"
        mv "$BIN_DIR/ffmpeg" "$BIN_DIR/ffmpeg-mac"
        rm "$TMP_ZIP"
        chmod +x "$BIN_DIR/ffmpeg-mac"
        echo "✓ ffmpeg-mac 已下载到 $BIN_DIR/ffmpeg-mac"
        ;;

    linux)
        echo "下载 Linux 版本的 ffmpeg..."
        FFMPEG_URL="https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz"
        TMP_TAR="/tmp/ffmpeg-linux.tar.xz"

        if command -v curl >/dev/null 2>&1; then
            curl -L -o "$TMP_TAR" "$FFMPEG_URL"
        elif command -v wget >/dev/null 2>&1; then
            wget -O "$TMP_TAR" "$FFMPEG_URL"
        else
            echo "错误: 需要 curl 或 wget 来下载 ffmpeg"
            exit 1
        fi

        tar -xf "$TMP_TAR" -C /tmp
        mv /tmp/ffmpeg-git-*-*-amd64-static/ffmpeg "$BIN_DIR/ffmpeg-linux"
        rm -rf /tmp/ffmpeg-git-*-*-amd64-static "$TMP_TAR"
        chmod +x "$BIN_DIR/ffmpeg-linux"
        echo "✓ ffmpeg-linux 已下载到 $BIN_DIR/ffmpeg-linux"
        ;;

    windows)
        echo "Windows 平台请手动下载 ffmpeg:"
        echo "1. 访问: https://www.gyan.dev/ffmpeg/builds/"
        echo "2. 下载 'ffmpeg-release-essentials.zip'"
        echo "3. 解压并找到 bin/ffmpeg.exe"
        echo "4. 将其复制到 $BIN_DIR/ffmpeg-windows.exe"
        ;;

    *)
        echo "不支持的平台: $PLATFORM"
        echo "请手动下载 ffmpeg 并放入 $BIN_DIR 目录"
        exit 1
        ;;
esac

echo ""
echo "下载完成！现在可以构建项目了："
echo "  go build -o video_crop"
