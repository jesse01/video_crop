//go:build gui

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"video_crop/internal/ffmpeg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var (
	version = "1.0.0"
)

func main() {
	myApp := app.NewWithID("com.videocrop.app")
	myWindow := myApp.NewWindow("视频黑边裁剪工具 " + version)

	// 设置关闭回调
	myWindow.SetCloseIntercept(func() {
		myApp.Quit()
	})

	// 确保临时目录存在
	tmpDir := filepath.Join(os.TempDir(), "video_crop")
	os.MkdirAll(tmpDir, 0755)

	makeUI(myWindow)

	myWindow.Resize(fyne.NewSize(500, 350))
	myWindow.CenterOnScreen()
	myWindow.SetMaster()
	myWindow.ShowAndRun()
}

func makeUI(w fyne.Window) {
	// 存储选中的文件路径
	var selectedFilePath string

	// 显示选中文件的Label（文字清晰）
	fileLabel := widget.NewLabel("未选择文件")
	fileLabel.TextStyle = fyne.TextStyle{Bold: true}

	// 选择文件按钮
	selectBtn := widget.NewButton("选择视频文件", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer reader.Close()

			// 获取文件路径
			selectedFilePath = storage.NewFileURI(reader.URI().Path()).Path()
			// 只显示文件名，路径太长显示不下
			fileName := filepath.Base(selectedFilePath)
			fileLabel.SetText(fileName)
		}, w)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".mp4", ".mkv", ".avi", ".mov", ".wmv", ".flv", ".webm", ".m4v"}))
		fd.Show()
	})
	selectBtn.Importance = widget.MediumImportance

	// 检测时长滑块
	detectTime := widget.NewSlider(5, 30)
	detectTime.Value = 10
	detectTime.Step = 1

	detectLabel := widget.NewLabel("检测时长: 10 秒")
	detectTime.OnChanged = func(value float64) {
		detectLabel.SetText(fmt.Sprintf("检测时长: %.0f 秒", value))
	}

	// 进度条
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// 状态文本
	statusText := widget.NewLabel("请选择视频文件")
	statusText.TextStyle = fyne.TextStyle{Monospace: true}

	// 开始按钮
	var startBtn *widget.Button
	startBtn = widget.NewButton("开始裁剪", func() {
		if selectedFilePath == "" {
			dialog.ShowError(fmt.Errorf("请先选择视频文件"), w)
			return
		}

		if _, err := os.Stat(selectedFilePath); os.IsNotExist(err) {
			dialog.ShowError(fmt.Errorf("文件不存在: %s", selectedFilePath), w)
			return
		}

		// 生成输出文件路径
		base := filepath.Base(selectedFilePath)
		ext := filepath.Ext(base)
		name := base[:len(base)-len(ext)]
		outputFile := filepath.Join(filepath.Dir(selectedFilePath), name+"_cropped"+ext)

		// 禁用按钮
		startBtn.Disable()
		selectBtn.Disable()
		progressBar.Show()
		statusText.SetText("正在检测黑边...")

		// 在goroutine中处理
		go processVideo(w, selectedFilePath, outputFile, int(detectTime.Value),
			progressBar, statusText, startBtn, selectBtn)
	})
	startBtn.Importance = widget.HighImportance

	// 布局 - 文件选择区域
	fileSelectBox := container.NewBorder(nil, nil, nil, selectBtn,
		container.NewVBox(
			widget.NewLabel("输入文件:"),
			fileLabel,
		))

	content := container.NewVBox(
		widget.NewLabel("📹 视频黑边自动裁剪工具"),
		widget.NewSeparator(),
		fileSelectBox,
		container.NewBorder(nil, nil, widget.NewLabel("检测时长:"), nil, detectLabel),
		detectTime,
		widget.NewSeparator(),
		container.NewVBox(
			startBtn,
			progressBar,
			statusText,
		),
	)

	w.SetContent(container.NewPadded(content))
}

func processVideo(w fyne.Window, inputFile, outputFile string, detectSeconds int,
	progressBar *widget.ProgressBar, statusText *widget.Label,
	startBtn, selectBtn *widget.Button) {

	// 使用defer确保按钮总是会被重新启用
	defer func() {
		startBtn.Enable()
		selectBtn.Enable()
	}()

	// 初始化ffmpeg
	statusText.SetText("正在初始化FFmpeg...")
	_, err := ffmpeg.InitFFmpeg()
	if err != nil {
		dialog.ShowError(fmt.Errorf("初始化FFmpeg失败: %w", err), w)
		statusText.SetText("初始化失败")
		return
	}

	// 检测黑边
	statusText.SetText("正在检测视频黑边...")
	cropParams, err := ffmpeg.DetectCrop(inputFile, detectSeconds)
	if err != nil {
		dialog.ShowError(fmt.Errorf("检测黑边失败: %w", err), w)
		statusText.SetText("检测失败")
		return
	}

	statusText.SetText("开始裁剪...")

	// 裁剪视频 - 使用WaitGroup等待完成
	var wg sync.WaitGroup
	var cropErr error

	progressChan := make(chan ffmpeg.ProgressInfo, 100)

	wg.Add(1)
	go func() {
		defer wg.Done()
		cropErr = ffmpeg.CropVideo(inputFile, outputFile, cropParams, true, progressChan)
	}()

	startTime := time.Now()
	lastUpdate := startTime

	// 处理进度更新
	for {
		select {
		case progress, ok := <-progressChan:
			if !ok {
				// channel已关闭，退出循环
				goto done
			}

			// 每200ms更新一次UI
			if time.Since(lastUpdate) >= 200*time.Millisecond {
				lastUpdate = time.Now()
				progressBar.SetValue(progress.Percent / 100)
				elapsed := time.Since(startTime).Round(time.Second)
				statusText.SetText(fmt.Sprintf("进度: %.1f%% | 速度: %.1fx | 已用: %s",
					progress.Percent, progress.Speed, elapsed))
			}
		}
	}

done:
	wg.Wait() // 确保goroutine完成

	if cropErr != nil {
		dialog.ShowError(fmt.Errorf("裁剪失败: %w", cropErr), w)
		statusText.SetText("裁剪失败")
		return
	}

	elapsed := time.Since(startTime).Round(time.Second)
	progressBar.SetValue(1)
	statusText.SetText(fmt.Sprintf("✓ 完成！用时: %s\n输出文件: %s", elapsed, outputFile))

	dialog.ShowInformation("完成", fmt.Sprintf(
		"视频裁剪完成！\n\n用时: %s\n输出文件: %s",
		elapsed, outputFile), w)
}
