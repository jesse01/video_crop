# GitHub Actions 自动编译指南

## 快速开始

### 第一步：创建 GitHub 仓库

```bash
# 在 GitHub 上创建一个新仓库（比如 video_crop）

# 添加远程仓库（替换 YOUR_USERNAME）
git remote add origin https://github.com/YOUR_USERNAME/video_crop.git

# 推送代码
git branch -M main
git push -u origin main
```

### 第二步：触发编译

选择以下任一方式：

#### 方式 A：创建版本发布（推荐）

```bash
# 创建版本 tag
git tag v1.0.0
git push origin v1.0.0
```

#### 方式 B：手动触发

1. 打开 GitHub 仓库页面
2. 点击 "Actions" 标签
3. 选择 "Build Release"
4. 点击 "Run workflow" → "Run workflow"

### 第三步：下载编译产物

编译完成后（大约5-10分钟）：

**从 Release 下载：**
- 访问仓库的 "Releases" 页面
- 下载 `v1.0.0` 版本中的文件

**从 Actions 下载：**
- 访问 "Actions" 页面
- 点击最新的运行记录
- 滚动到底部下载 artifacts

## 下载的文件

| 文件名 | 平台 | 说明 |
|--------|------|------|
| `VideoCrop-Windows-x64-GUI.zip` | Windows | GUI版本，解压后运行 video_crop_gui.exe |
| `video_crop-Windows-x64-CLI.zip` | Windows | 命令行版本，解压后运行 video_crop.exe |
| `VideoCrop-macOS-GUI.zip` | macOS | GUI版本，解压后双击 VideoCrop.app |
| `video_crop-macOS-CLI.zip` | macOS | 命令行版本，解压后运行 video_crop-mac-arm64 |
| `VideoCrop-Linux-GUI.tar.gz` | Linux | GUI版本，解压后运行 video_crop_gui |
| `video_crop-Linux-CLI.tar.gz` | Linux | 命令行版本，解压后运行 video_crop-linux |

## 注意事项

1. **FFmpeg 已经内嵌**在可执行文件中，无需额外安装
2. **Windows GUI** 首次运行可能需要允许防火墙访问
3. **macOS** 如果无法打开，右键点击 → "打开" → "打开"

## 故障排除

### 编译失败
- 检查 GitHub Actions 日志获取详细错误信息
- 确保 `main_gui.go` 和 `cropdetect.go` 语法正确

### 下载的文件无法运行
- Windows：可能被杀毒软件拦截，添加到白名单
- macOS：运行 `xattr -cr VideoCrop.app` 移除隔离属性
- Linux：确保有执行权限 `chmod +x video_crop_gui`
