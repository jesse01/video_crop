//go:build darwin

package ffmpeg

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed binaries
var ffmpegFS embed.FS

func getEmbeddedFFmpeg() ([]byte, error) {
	// 尝试从嵌入的文件系统读取
	data, err := ffmpegFS.ReadFile("binaries/ffmpeg-mac")
	if err != nil {
		// 尝试从本地文件系统加载（用于开发环境）
		localPath := filepath.Join("internal", "ffmpeg", "binaries", "ffmpeg-mac")
		if fileData, err := os.ReadFile(localPath); err == nil && len(fileData) > 1024*1024 {
			return fileData, nil
		}
		// 返回错误，让调用者使用系统的ffmpeg
		return nil, fs.ErrNotExist
	}

	// 检查是否是有效的可执行文件（至少1MB）
	if len(data) > 1024*1024 {
		return data, nil
	}

	// 嵌入的文件太小，不是有效的ffmpeg
	return nil, fs.ErrNotExist
}
