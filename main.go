//go:build !gui

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"video_crop/internal/cmd"
)

var (
	// Version 版本号
	Version = "1.0.0"
	// BuildTime 构建时间
	BuildTime = "unknown"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// 设置版本信息
	cmd.SetVersion(Version, BuildTime)

	// 确保临时目录存在
	tmpDir := filepath.Join(os.TempDir(), "video_crop")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to create temp directory: %v\n", err)
	}
}
