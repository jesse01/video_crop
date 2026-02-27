# GitHub Actions 自动编译

## 使用方法

### 方式1：通过 Tag 触发（自动发布版本）

```bash
# 1. 提交所有更改
git add .
git commit -m "your changes"
git push

# 2. 创建并推送tag（触发自动编译和发布）
git tag v1.0.0
git push origin v1.0.0
```

GitHub Actions 会自动：
- 编译所有平台的 GUI 和 CLI 版本
- 创建 GitHub Release
- 上传所有编译产物

### 方式2：手动触发

1. 访问 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择 "Build Release" 工作流
4. 点击 "Run workflow" 按钮

## 下载编译产物

编译完成后，可以从以下位置下载：

### 方式A：从 Release 下载
- 访问仓库的 "Releases" 页面
- 下载对应平台的文件

### 方式B：从 Actions 下载
- 访问 "Actions" 页面
- 选择一次运行记录
- 在页面底部下载 artifacts

## 编译产物

| 文件 | 平台 | 类型 |
|------|------|------|
| `VideoCrop-Windows-x64-GUI.zip` | Windows | GUI应用 |
| `video_crop-Windows-x64-CLI.zip` | Windows | 命令行 |
| `VideoCrop-macOS-GUI.zip` | macOS | GUI应用 |
| `video_crop-macOS-CLI.zip` | macOS | 命令行 |
| `VideoCrop-Linux-GUI.tar.gz` | Linux | GUI应用 |
| `video_crop-Linux-CLI.tar.gz` | Linux | 命令行 |

## 本地测试

在推送之前，可以先在本地测试编译是否正常：

### macOS
```bash
# GUI
CGO_ENABLED=1 go build -tags gui -o test_gui .

# CLI
CGO_ENABLED=0 go build -o test_cli .
```

### Windows
```cmd
REM GUI
set CGO_ENABLED=1
go build -tags gui -o test_gui.exe .

REM CLI
set CGO_ENABLED=0
go build -o test_cli.exe .
```
