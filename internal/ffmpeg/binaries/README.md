# FFmpeg 二进制文件

此目录用于存放嵌入到应用程序中的 FFmpeg 二进制文件。

## 下载 FFmpeg

### Windows

下载地址: https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip

解压后，将 `ffmpeg.exe` 重命名为 `ffmpeg-windows.exe` 并放入此目录。

### macOS

```bash
# 使用 Homebrew 安装
brew install ffmpeg

# 复制到当前目录
cp $(which ffmpeg) ffmpeg-mac
```

或从官网下载: https://evermeet.cx/ffmpeg/

### Linux

```bash
# Ubuntu/Debian
apt-get install ffmpeg

# 复制到当前目录
cp $(which ffmpeg) ffmpeg-linux
```

或从官网下载: https://johnvansickle.com/ffmpeg/

## 嵌入说明

将对应平台的二进制文件放入此目录后，Go 编译时会自动将其嵌入到可执行文件中。

注意：FFmpeg 二进制文件较大（约 50-80MB），会导致最终的可执行文件也较大。
