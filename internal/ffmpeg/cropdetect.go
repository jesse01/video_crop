package ffmpeg

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// ProgressInfo 表示处理进度信息
type ProgressInfo struct {
	Percent     float64
	Speed       float64
	TimeElapsed time.Duration
}

// InitFFmpeg 初始化ffmpeg，返回ffmpeg可执行文件路径
func InitFFmpeg() (string, error) {
	// 如果用户环境中有ffmpeg，优先使用
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// 尝试从嵌入的ffmpeg二进制文件
	embeddedPath, err := extractEmbeddedFFmpeg()
	if err != nil {
		return "", fmt.Errorf("无法找到或提取ffmpeg: %w", err)
	}

	return embeddedPath, nil
}

// DetectCrop 检测视频黑边
func DetectCrop(videoPath string, detectSeconds int) (string, error) {
	ffmpegPath, err := InitFFmpeg()
	if err != nil {
		return "", err
	}

	// 构建ffmpeg命令
	args := []string{
		"-i", videoPath,
		"-t", fmt.Sprintf("%d", detectSeconds),
		"-vf", "cropdetect=24:16:0",
		"-f", "null",
		"-",
	}

	cmd := exec.Command(ffmpegPath, args...)

	// 捕获stderr输出（ffmpeg的日志输出到stderr）
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", fmt.Errorf("创建stderr管道失败: %w", err)
	}

	// 运行命令
	if err := cmd.Start(); err != nil {
		return "", fmt.Errorf("启动ffmpeg失败: %w", err)
	}

	// 解析输出获取crop参数
	cropParams, err := parseCropOutput(stderr)
	if err != nil {
		cmd.Wait()
		return "", fmt.Errorf("解析crop输出失败: %w", err)
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return "", fmt.Errorf("ffmpeg执行失败: %w", err)
	}

	return cropParams, nil
}

// parseCropOutput 解析ffmpeg的cropdetect输出
func parseCropOutput(pipe interface{}) (string, error) {
	scanner := bufio.NewScanner(pipe.(interface {
		Read([]byte) (int, error)
	}))

	// 正则匹配 crop 参数
	// 格式: [Parsed_cropdetect_0 @ ...] crop=1920:1080:0:0
	cropRegex := regexp.MustCompile(`crop=(\d+):(\d+):(\d+):(\d+)`)

	var bestCrop string
	maxPixels := 0

	for scanner.Scan() {
		line := scanner.Text()

		// 查找包含crop的行
		if strings.Contains(line, "crop=") {
			matches := cropRegex.FindStringSubmatch(line)
			if len(matches) >= 5 {
				width := matches[1]
				height := matches[2]
				x := matches[3]
				y := matches[4]

				cropParam := fmt.Sprintf("%s:%s:%s:%s", width, height, x, y)

				// 选择像素最多的crop结果（最准确的检测）
				// 计算像素数
				var w, h int
				fmt.Sscanf(width+":"+height, "%d:%d", &w, &h)
				pixels := w * h

				if pixels > maxPixels {
					maxPixels = pixels
					bestCrop = cropParam
				}
			}
		}
	}

	if bestCrop == "" {
		return "", fmt.Errorf("未能从ffmpeg输出中检测到crop参数")
	}

	return bestCrop, nil
}

// CropVideo 裁剪视频
func CropVideo(inputPath, outputPath, cropParams string, overwrite bool, progressChan chan<- ProgressInfo) error {
	ffmpegPath, err := InitFFmpeg()
	if err != nil {
		return err
	}

	// 确保函数结束时关闭channel
	defer close(progressChan)

	args := []string{}

	if overwrite {
		args = append(args, "-y")
	}

	args = append(args,
		"-i", inputPath,
		"-vf", fmt.Sprintf("crop=%s", cropParams),
		"-c:a", "copy", // 直接复制音频流，不重新编码
		outputPath,
	)

	cmd := exec.Command(ffmpegPath, args...)

	// 捕获stderr以解析进度
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("创建stderr管道失败: %w", err)
	}

	// 启动进度监控
	progressDone := make(chan error)
	go func() {
		progressDone <- monitorProgress(stderr, progressChan)
	}()

	// 运行命令
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动ffmpeg失败: %w", err)
	}

	// 等待进度监控完成
	if err := <-progressDone; err != nil {
		return err
	}

	// 等待命令完成
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg执行失败: %w", err)
	}

	return nil
}

// monitorProgress 监控ffmpeg处理进度
func monitorProgress(pipe interface{}, progressChan chan<- ProgressInfo) error {
	scanner := bufio.NewScanner(pipe.(interface {
		Read([]byte) (int, error)
	}))

	startTime := time.Now()

	// 正则匹配进度信息
	// 格式: frame=  123 fps= 30 q=28.0 size=    1234kB time=00:00:05.00 bitrate= 1234.5kbits/s speed=1.23x
	timeRegex := regexp.MustCompile(`time=(\d+):(\d+):(\d+\.\d+)`)
	speedRegex := regexp.MustCompile(`speed=\s*([\d.]+)x`)

	var duration float64 // 视频总时长（秒）

	for scanner.Scan() {
		line := scanner.Text()

		// 首先尝试获取视频总时长
		if duration == 0 && strings.Contains(line, "Duration:") {
			durationRegex := regexp.MustCompile(`Duration:\s*(\d+):(\d+):([\d.]+)`)
			matches := durationRegex.FindStringSubmatch(line)
			if len(matches) >= 4 {
				var h, m, s float64
				fmt.Sscanf(matches[1]+":"+matches[2]+":"+matches[3], "%f:%f:%f", &h, &m, &s)
				duration = h*3600 + m*60 + s
			}
		}

		// 解析进度
		timeMatches := timeRegex.FindStringSubmatch(line)
		if len(timeMatches) >= 4 && duration > 0 {
			var h, m, s float64
			fmt.Sscanf(timeMatches[1]+":"+timeMatches[2]+":"+timeMatches[3], "%f:%f:%f", &h, &m, &s)
			currentTime := h*3600 + m*60 + s

			percent := (currentTime / duration) * 100
			if percent > 100 {
				percent = 100
			}

			speed := 1.0
			speedMatches := speedRegex.FindStringSubmatch(line)
			if len(speedMatches) >= 2 {
				fmt.Sscanf(speedMatches[1], "%f", &speed)
			}

			progressChan <- ProgressInfo{
				Percent:     percent,
				Speed:       speed,
				TimeElapsed: time.Since(startTime),
			}
		}
	}

	return nil
}

// extractEmbeddedFFmpeg 提取嵌入的ffmpeg二进制文件
func extractEmbeddedFFmpeg() (string, error) {
	// 获取平台对应的ffmpeg二进制
	data, err := getEmbeddedFFmpeg()
	if err != nil {
		return "", err
	}

	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "video_crop")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %w", err)
	}

	// 确定ffmpeg可执行文件名
	exeName := "ffmpeg"
	if runtime.GOOS == "windows" {
		exeName = "ffmpeg.exe"
	}

	exePath := filepath.Join(tmpDir, exeName)

	// 如果文件已存在，直接返回
	if _, err := os.Stat(exePath); err == nil {
		return exePath, nil
	}

	// 写入二进制文件
	if err := os.WriteFile(exePath, data, 0755); err != nil {
		return "", fmt.Errorf("写入ffmpeg二进制文件失败: %w", err)
	}

	return exePath, nil
}
