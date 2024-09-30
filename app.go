package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ReadDirFileCount(dir string, includeChild bool) int {
	data := FormData{
		SourceDir:    dir,
		IncludeChild: includeChild,
	}
	fileList, err := GetFileList(data)
	if err != nil {
		log.Println("读取文件列表时出错:", err)
		return 0
	}

	fileListSize := len(fileList)

	return fileListSize
}

func (a *App) HandleFormSubmission(data FormData) {

	jsonData, _ := json.Marshal(data)
	log.Println(string(jsonData))

	// 根据 data.SourceDir 和 data.IncludeChild 读取文件列表详情

	fileList, err := GetFileList(data)
	if err != nil {
		log.Println("读取文件列表时出错:", err)
		runtime.EventsEmit(a.ctx, "task-complete", "任务已完成！")

		return
	}

	fileListSize := len(fileList)

	skipSizeInBytes := int64(data.SkipFileLessThan * 1024)

	//md5Map := make(map[string]bool) // 存储文件的 MD5 值，判断是否有重复
	fileMetaMap := make(map[string]bool) // 存储文件的元信息，判断是否有重复

	for i, file := range fileList {
		log.Println("Index:", i, "File:", file)
		runtime.EventsEmit(a.ctx, "task-progress", i+1, fileListSize)

		// process file with data struct
		// 获取文件信息
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		// 跳过小于指定大小的文件
		if info.Size() < skipSizeInBytes {
			continue
		}

		// 检查文件名是否包含跳过的内容
		if data.SkipFileContains != "" && len(data.SkipFileContains) > 0 {
			if strings.Contains(info.Name(), data.SkipFileContains) {
				continue
			}
		}

		// 检查是否为系统文件
		if isSystemFile(info) {
			log.Println("跳过系统文件:", file)
			continue
		}

		// 生成文件唯一键
		fileKey := generateFileKey(info.Size(), info.ModTime())

		// 判断是否是重复文件（基于文件唯一键）
		if _, exists := fileMetaMap[fileKey]; exists {
			log.Println("跳过重复文件:", file)
			continue
		}

		// 记录文件唯一键，防止重复文件
		fileMetaMap[fileKey] = true

		/*		// 计算文件的 MD5 值
				fileMD5, err := calculateFileMD5(file)
				if err != nil {
					log.Println("计算 MD5 失败:", file, err)
					continue
				}

				// 判断是否是重复文件
				if _, exists := md5Map[fileMD5]; exists {
					log.Println("跳过重复文件:", file)
					continue
				}

				// 记录 MD5，防止重复文件
				md5Map[fileMD5] = true
		*/
		// 读取拍摄时间
		shootTime, err := getPhotoTakenTime(file, data.NoShotTimeType)
		if err != nil {
			log.Println("无法读取拍摄时间:", err)
			continue // 继续处理其他文件
		}

		// 生成目标路径
		targetDir := generateTargetDir(data.TargetDir, shootTime, data.DirType)
		err = os.MkdirAll(targetDir, os.ModePerm)
		if err != nil {
			continue
		}

		targetPath := filepath.Join(targetDir, info.Name())

		// 移动或复制文件
		err = copyOrMoveFile(file, targetPath, data.IsMove, data.SkipSameFile)
		if err != nil {
			log.Println("文件操作失败:", err)
			continue
		}

		log.Println("文件操作完成:", file)
		// 向前端发送进度更新事件
		//log.Printf("向前端发送进度更新事件: task-progress: %+v,%+v \n", i+1, fileListSize)
	}

	// 任务完成后通知前端
	log.Printf("任务完成后通知前端: task-complete\n")

	runtime.EventsEmit(a.ctx, "task-complete", "任务已完成！")

}
