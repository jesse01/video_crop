package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"video_crop/internal/ffmpeg"

	"github.com/spf13/cobra"
)

var (
	// Version 版本号
	Version string
	// BuildTime 构建时间
	BuildTime string

	inputFile  string
	outputFile string
	overwrite  bool
	detectTime int // 检测黑边的时长（秒）
)

var rootCmd = &cobra.Command{
	Use:   "video_crop [input_file] [output_file]",
	Short: "自动裁剪视频黑边",
	Long: `视频黑边自动裁剪工具

该工具使用 ffmpeg 的 cropdetect 功能自动检测视频黑边，
并裁剪掉黑边生成新的视频文件。

示例:
  video_crop input.mp4 output.mp4
  video_crop input.mp4 output.mp4 --detect-time 10`,
	Args: cobra.MaximumNArgs(2),
	RunE: runCrop,
}

// SetVersion 设置版本信息
func SetVersion(version, buildTime string) {
	Version = version
	BuildTime = buildTime
	rootCmd.Version = fmt.Sprintf("%s (built %s)", Version, BuildTime)
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "输入视频文件")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "输出视频文件")
	rootCmd.Flags().BoolVarP(&overwrite, "overwrite", "y", false, "覆盖已存在的输出文件")
	rootCmd.Flags().IntVarP(&detectTime, "detect-time", "d", 10, "检测黑边的时长（秒），建议5-15秒")
}

func Execute() error {
	return rootCmd.Execute()
}

func runCrop(cmd *cobra.Command, args []string) error {
	// 处理参数
	if len(args) > 0 {
		inputFile = args[0]
	}
	if len(args) > 1 {
		outputFile = args[1]
	}

	// 验证输入
	if inputFile == "" {
		return fmt.Errorf("请指定输入视频文件")
	}

	if !fileExists(inputFile) {
		return fmt.Errorf("输入文件不存在: %s", inputFile)
	}

	// 自动生成输出文件名
	if outputFile == "" {
		base := filepath.Base(inputFile)
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]
		outputFile = filepath.Join(filepath.Dir(inputFile), name+"_cropped"+ext)
	}

	// 检查输出文件是否已存在
	if fileExists(outputFile) && !overwrite {
		return fmt.Errorf("输出文件已存在: %s (使用 --overwrite 覆盖)", outputFile)
	}

	fmt.Printf("输入文件: %s\n", inputFile)
	fmt.Printf("输出文件: %s\n", outputFile)
	fmt.Printf("检测时长: %d秒\n", detectTime)
	fmt.Println()

	// 初始化 ffmpeg
	ffmpegPath, err := ffmpeg.InitFFmpeg()
	if err != nil {
		return fmt.Errorf("初始化ffmpeg失败: %w", err)
	}
	fmt.Printf("使用ffmpeg: %s\n\n", ffmpegPath)

	// 检测黑边
	fmt.Println("正在检测视频黑边...")
	cropParams, err := ffmpeg.DetectCrop(inputFile, detectTime)
	if err != nil {
		return fmt.Errorf("检测黑边失败: %w", err)
	}
	fmt.Printf("检测到的裁剪参数: %s\n\n", cropParams)

	// 裁剪视频
	fmt.Println("正在裁剪视频...")
	startTime := time.Now()
	progressChan := make(chan ffmpeg.ProgressInfo, 10)

	go func() {
		for progress := range progressChan {
			fmt.Printf("\r进度: %.1f%% | 速度: %.1fx | 已用时间: %s",
				progress.Percent,
				progress.Speed,
				progress.TimeElapsed.Round(time.Second))
		}
		fmt.Println()
	}()

	err = ffmpeg.CropVideo(inputFile, outputFile, cropParams, overwrite, progressChan)
	close(progressChan)

	if err != nil {
		return fmt.Errorf("裁剪视频失败: %w", err)
	}

	elapsed := time.Since(startTime)
	fmt.Printf("\n✓ 裁剪完成！用时: %s\n", elapsed)
	fmt.Printf("✓ 输出文件: %s\n", outputFile)

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
