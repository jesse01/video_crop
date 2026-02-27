.PHONY: all build clean test run-windows deps download-ffmpeg gui gui-windows

# 版本信息
VERSION ?= 1.0.0
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

# 默认目标
all: build

# 构建当前平台
build:
	@echo "构建当前平台..."
	@CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o video_crop
	@echo "✓ 构建完成: video_crop"

# 构建 Windows 版本
windows:
	@echo "构建 Windows (amd64)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o video_crop.exe
	@echo "✓ 构建完成: video_crop.exe"

# 构建 macOS 版本
macos:
	@echo "构建 macOS (arm64)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o video_crop-mac-arm64
	@echo "✓ 构建完成: video_crop-mac-arm64"

# 构建 Linux 版本
linux:
	@echo "构建 Linux (amd64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o video_crop-linux
	@echo "✓ 构建完成: video_crop-linux"

# 构建所有平台
build-all:
	@echo "构建所有平台..."
	@mkdir -p build
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/video_crop-windows-amd64.exe
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/video_crop-mac-amd64
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/video_crop-mac-arm64
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/video_crop-linux-amd64
	@echo "✓ 所有平台构建完成！"
	@ls -lh build/

# 下载 FFmpeg
download-ffmpeg:
	@echo "下载 FFmpeg..."
	@if [ "$(shell uname)" = "Darwin" ]; then \
		chmod +x download_ffmpeg.sh && ./download_ffmpeg.sh; \
	elif [ "$(shell uname)" = "Linux" ]; then \
		chmod +x download_ffmpeg.sh && ./download_ffmpeg.sh; \
	else \
		echo "请手动下载 FFmpeg"; \
	fi

# 安装依赖
deps:
	@echo "安装依赖..."
	@go mod download
	@go mod tidy
	@echo "✓ 依赖安装完成"

# 运行测试
test:
	@echo "运行测试..."
	@CGO_ENABLED=0 go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf build/
	@rm -f video_crop video_crop.exe video_crop-* video_crop-*
	@echo "✓ 清理完成"

# 在 Windows 上运行（需要 Wine）
run-windows: video_crop.exe
	@echo "运行 Windows 版本..."
	@if command -v wine >/dev/null 2>&1; then \
		wine video_crop.exe --help; \
	else \
		echo "需要安装 Wine 来运行 Windows 程序"; \
	fi

# 构建 GUI 版本（需要 CGO）
gui:
	@echo "构建 GUI 版本（当前平台）..."
	@CGO_ENABLED=1 go build -tags gui -o video_crop_gui .
	@echo "✓ GUI 构建完成: video_crop_gui"

# 构建 macOS GUI App Bundle
gui-macos:
	@echo "构建 macOS GUI 版本..."
	@mkdir -p VideoCrop.app/Contents/MacOS
	@CGO_ENABLED=1 go build -tags gui -o VideoCrop.app/Contents/MacOS/videocrop .
	@echo "正在添加代码签名..."
	@codesign --force --deep --sign - VideoCrop.app
	@echo "✓ GUI 构建完成: VideoCrop.app"
	@echo ""
	@echo "验证签名:"
	@codesign -vv VideoCrop.app

# 构建 Windows GUI 版本（需要在 Windows 上编译）
gui-windows:
	@echo "Windows GUI 版本需要在 Windows 上编译"
	@echo "请使用 build_gui_windows.bat 脚本"
	@if [ "$(shell uname)" = "Darwin" ] || [ "$(shell uname)" = "Linux" ]; then \
		echo ""; \
		echo "或者使用 Docker 编译:"; \
		echo "  docker run --rm -v \"\$$PWD\":/app -w /app ttyao/golang-windows:latest go build -tags windows -o video_crop_gui.exe ."; \
	fi

# 显示帮助
help:
	@echo "可用的 make 目标:"
	@echo "  make              - 构建当前平台"
	@echo "  make build        - 构建当前平台"
	@echo "  make gui          - 构建 GUI 版本（当前平台，需要 CGO）"
	@echo "  make gui-macos    - 构建 macOS GUI App Bundle（自动签名）"
	@echo "  make gui-windows  - 显示 Windows GUI 编译说明"
	@echo "  make windows      - 构建 Windows 命令行版本"
	@echo "  make macos        - 构建 macOS 版本"
	@echo "  make linux        - 构建 Linux 版本"
	@echo "  make build-all    - 构建所有平台"
	@echo "  make download-ffmpeg - 下载 FFmpeg"
	@echo "  make deps         - 安装依赖"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理构建文件"
	@echo "  make help         - 显示此帮助"
