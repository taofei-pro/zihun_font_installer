package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed fonts/*.ttf
var fonts embed.FS

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// SelectDirectory 打开目录选择对话框
func (a *App) SelectDirectory() (string, error) {
	// 调用 Wails 的 runtime 方法，显示目录选择对话框
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目标目录",
	})
	if err != nil {
		return "", err
	}
	return dir, nil
}

// InstallFonts 将嵌入的字体文件复制到目标目录
func (a *App) InstallFonts(outputPath string) error {
	// 确保目标路径存在
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("目标路径不存在: %s", outputPath)
	}

	// 包含字体包名称的数组
	fontNames := []string{
		"字魂书雅宋-Regular.ttf",
		"Zihun Serif-Regular.ttf",
	}

	// 获取嵌入的字体文件列表
	fontFiles, err := fonts.ReadDir("fonts")
	if err != nil {
		return fmt.Errorf("读取字体文件列表失败: %v", err)
	}

	fontNameExist := false

	// 遍历并复制每个字体文件
	for _, fontFile := range fontFiles {
		if fontFile.IsDir() {
			continue
		}

		// 检查字体文件是否在字体名称数组中
		fontName := fontFile.Name()
		if !fontNameExist {
			for _, item := range fontNames {
				if item == fontName {
					fontNameExist = true
					break
				}
			}
		}

		// 打开嵌入的字体文件
		fontData, err := fonts.Open(filepath.Join("fonts", fontName))
		if err != nil {
			return fmt.Errorf("无法打开嵌入的字体文件 %s: %v", fontName, err)
		}
		defer fontData.Close()

		// 创建目标文件
		destPath := filepath.Join(outputPath, fontName)
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("无法创建目标文件 %s: %v", destPath, err)
		}
		defer destFile.Close()

		// 将字体数据写入目标文件
		_, err = io.Copy(destFile, fontData)
		if err != nil {
			return fmt.Errorf("复制字体文件 %s 时出错: %v", fontName, err)
		}
	}

	// 安装成功后提示
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "安装成功",
		Message: "所有字体已成功安装到目标目录！",
	})
	if fontNameExist {
		fmt.Println("字体安装成功！")
	}

	return nil
}
