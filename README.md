# 视频黑边自动裁剪工具

> 自动检测并裁剪视频黑边的桌面工具，支持 Windows、macOS 和 Linux

[![Build Release](https://github.com/YOUR_USERNAME/video_crop/actions/workflows/build.yml/badge.svg)](https://github.com/YOUR_USERNAME/video_crop/actions/workflows/build.yml)

## 功能特点

- 🎬 **自动检测黑边** - 使用 FFmpeg 智能检测视频黑边区域
- 🖥️ **图形界面** - 简洁易用的 GUI，无需命令行
- 💻 **命令行支持** - 支持批量处理和脚本自动化
- 📦 **独立运行** - 内嵌 FFmpeg，无需额外安装
- 🌍 **跨平台** - 支持 Windows、macOS、Linux
- ⚡ **实时进度** - 显示处理进度、速度和预计时间

## 截图

### GUI 版本

```
┌─────────────────────────────────────┐
│    📹 视频黑边自动裁剪工具          │
├─────────────────────────────────────┤
│ 输入文件: example.mp4               │
│ [选择视频文件]                      │
│                                     │
│ 检测时长: 10 秒                      │
│ [━━━━━━━━━━━━━━━━━━━━] 5-30秒      │
│                                     │
│ [开始裁剪]                          │
│ [━━━━━━━━━━━━━━━━] 进度: 45%       │
│ 状态: 进度: 45.2% | 速度: 2.1x     │
└─────────────────────────────────────┘
```

## 下载安装

### 方式 1：下载预编译版本（推荐）

访问 [Releases](https://github.com/YOUR_USERNAME/video_crop/releases) 页面下载对应平台：

| 平台 | GUI 版本 | 命令行版本 |
|------|----------|------------|
| Windows | `VideoCrop-Windows-x64-GUI.zip` | `video_crop-Windows-x64-CLI.zip` |
| macOS | `VideoCrop-macOS-GUI.zip` | `video_crop-macOS-CLI.zip` |
| Linux | `VideoCrop-Linux-GUI.tar.gz` | `video_crop-Linux-CLI.tar.gz` |

解压后直接运行即可，无需安装。

### 方式 2：从源码编译

详见 [BUILD_GUIDE.md](AUTO_BUILD_GUIDE.md)

## 使用方法

### GUI 版本

1. 双击运行 `VideoCrop.app`（macOS）或 `video_crop_gui.exe`（Windows）
2. 点击"选择视频文件"选择要处理的视频
3. 调整检测时长（默认 10 秒）
4. 点击"开始裁剪"
5. 等待处理完成，输出文件保存在原文件同目录

### 命令行版本

#### 基本用法

```bash
# 自动生成输出文件名
video_crop input.mp4

# 指定输出文件名
video_crop input.mp4 output.mp4

# 使用长参数
video_crop --input input.mp4 --output output.mp4
```

#### 高级选项

```bash
# 覆盖已存在的输出文件
video_crop -i input.mp4 -o output.mp4 -y

# 指定黑边检测时长（秒）
video_crop -i input.mp4 -o output.mp4 -d 15

# 查看帮助
video_crop --help
```

#### 参数说明

| 参数 | 简写 | 说明 | 默认值 |
|------|------|------|--------|
| `--input` | `-i` | 输入视频文件 | 必填 |
| `--output` | `-o` | 输出视频文件 | `[input]_cropped.[ext]` |
| `--overwrite` | `-y` | 覆盖已存在的输出文件 | false |
| `--detect-time` | `-d` | 黑边检测时长（秒） | 10 |

## 支持的视频格式

- MP4, MKV, AVI, MOV, WMV, FLV, WEBM, M4V

## 工作原理

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  输入视频   │ ──> │  黑边检测    │ ──> │  裁剪视频   │
│  input.mp4  │     │  cropdetect  │     │  crop       │
└─────────────┘     └──────────────┘     └─────────────┘
                                           │
                                           ▼
                                    ┌─────────────┐
                                    │  输出视频   │
                                    │output_cropped│
                                    └─────────────┘
```

1. **检测阶段**：分析视频开头指定时长（默认 10 秒），检测黑边位置
2. **裁剪阶段**：使用检测到的参数裁剪整个视频
3. **输出**：生成去除黑边的新视频文件

## 常见问题

### Q: macOS 提示"已损坏，无法打开"？

A: 这是 macOS 的安全机制，右键点击应用 → "打开" → 再次点击"打开"

### Q: Windows 提示找不到 FFmpeg？

A: 工具已内嵌 FFmpeg，如果仍有问题，检查是否被杀毒软件拦截

### Q: 裁剪后还有黑边？

A: 尝试增加检测时长（`-d 20`），让检测更准确

### Q: 处理速度很慢？

A: 处理速度取决于视频编码和电脑性能，通常为视频时长的 1/3 到 1/2

## 开发

### 环境要求

- Go 1.21+
- GCC（GUI 版本需要）

### 编译

```bash
# 克隆仓库
git clone https://github.com/YOUR_USERNAME/video_crop.git
cd video_crop

# 编译当前平台
make build

# 编译 GUI 版本（需要 CGO）
make gui-macos  # macOS
```

### 使用 GitHub Actions 自动编译

详见 [GITHUB_ACTIONS.md](GITHUB_ACTIONS.md)

## 项目结构

```
video_crop/
├── main.go              # CLI 入口
├── main_gui.go          # GUI 入口
├── internal/
│   ├── cmd/            # CLI 逻辑
│   └── ffmpeg/         # FFmpeg 封装
│       ├── cropdetect.go
│       └── embed_*.go   # 平台相关嵌入
└── .github/workflows/   # CI/CD 配置
```

## 技术栈

- **Go 1.21+** - 主要编程语言
- **FFmpeg** - 视频处理引擎
- **Fyne v2** - GUI 框架
- **Cobra** - CLI 框架

## 许可证

MIT License - 详见 [LICENSE](LICENSE)

## 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 致谢

- [FFmpeg](https://ffmpeg.org/) - 强大的视频处理工具
- [Fyne](https://fyne.io/) - 跨平台 GUI 框架
- [Cobra](https://github.com/spf13/cobra) - CLI 框架

---

**Made with ❤️ by [Your Name]**
