package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	//"github.com/rwcarlsen/goexif/exif"
)

type FormData struct {
	SourceDir        string `json:"sourceDir"`
	IncludeChild     bool   `json:"includeChild"`
	TargetDir        string `json:"targetDir"`
	DirType          string `json:"dirType"`
	IsMove           bool   `json:"isMove"`         // false copy true move
	NoShotTimeType   int    `json:"noShotTimeType"` // 0 跳过 1 根据修改时间整理 2 根据创建时间整理
	SkipSameFile     bool   `json:"skipSameFile"`
	SkipFileLessThan int    `json:"skipFileLessThan"` // KB
	SkipFileContains string `json:"skipFileContains"`
}

// GetFileList 获取文件列表，排除文件夹
func GetFileList(data FormData) ([]string, error) {
	var files []string

	// 如果需要递归子目录
	if data.IncludeChild {
		// 使用 WalkDir 遍历所有文件和文件夹
		err := filepath.WalkDir(data.SourceDir, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// 检查是否为文件（排除文件夹）
			if !d.IsDir() {
				files = append(files, path)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	} else {
		// 只读取指定目录下的文件
		entries, err := os.ReadDir(data.SourceDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			// 检查是否为文件（排除文件夹）
			if !entry.IsDir() {
				fullPath := filepath.Join(data.SourceDir, entry.Name())
				files = append(files, fullPath)
			}
		}
	}

	return files, nil
}

// extractCreateDate 从 exiftool 的输出中提取创建日期
func extractDate(exifOutput string) (string, error) {
	// 使用 strings.Split 拆分字符串
	parts := strings.Split(exifOutput, ": ")
	if len(parts) < 2 {
		return "", fmt.Errorf("无法提取创建日期")
	}

	// 去除前后空格并返回日期部分
	return strings.TrimSpace(parts[1]), nil
}

// formatDate 格式化创建日期字符串
func formatDate(createDate string) (time.Time, error) {
	layout := "2006:01:02 15:04:05" // ExifTool 日期格式
	parsedTime, err := time.Parse(layout, strings.TrimSpace(createDate))
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}

func getExifDateTime(filePath string, dateTimeType int) (time.Time, error) {
	// 获取 ExifTool 的绝对路径
	exifToolPath, err := filepath.Abs("assets/exiftool.exe")
	if err != nil {
		return time.Time{}, err
	}

	// 调用 exiftool 命令

	var arg string

	switch dateTimeType {
	case 0:
		arg = "-CreateDate"
	case 1:
		arg = "-ModifyDate"
	case 2:
		arg = "-FileCreateDate"
	default:
		arg = "-CreateDate"
	}

	cmd := exec.Command(exifToolPath, arg, filePath)
	//cmd := exec.Command(exifToolPath, "-ModifyDate", filePath)

	// 创建一个字节缓冲区来捕获标准输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令并获取错误
	if err := cmd.Run(); err != nil {
		return time.Time{}, err
	}
	// 提取创建日期
	createDate, err := extractDate(out.String())
	if err != nil {
		fmt.Println("提取创建日期时出错:", err)

	}

	// 格式化创建日期
	formattedDate, err := formatDate(createDate)
	if err != nil {
		fmt.Println("日期格式化出错:", err)

	}

	fmt.Println("格式化后的创建日期:", formattedDate)
	return formattedDate, err
}

// getPhotoTakenTime 读取照片的拍摄时间或文件修改时间
func getPhotoTakenTime(filePath string, noShotTimeType int) (time.Time, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return time.Time{}, err
	}
	defer file.Close()

	// 获取创建时间
	dateTime, err := getExifDateTime(filePath, 0)

	fmt.Println(dateTime)

	if err != nil {
		// 如果读取 Exif 失败，根据 NoShotTimeType 返回修改时间或创建时间
		switch noShotTimeType { // 0 跳过 1 根据修改时间整理 2 根据创建时间整理
		case 1:
			return getExifDateTime(filePath, 1)
		case 2:

			// 获取创建时间
			return getExifDateTime(filePath, 2)
		default:
			return time.Time{}, fmt.Errorf("无法获取拍摄时间")
		}
	}
	//
	//// 获取拍摄时间
	//tm, err := x.DateTime()
	//if err != nil {
	//	return time.Time{}, err
	//}

	return dateTime, nil
}

// generateTargetDir 生成目标目录路径
func generateTargetDir(baseDir string, shootTime time.Time, dirType string) string {
	switch dirType {
	case "/2006/01/": // `/2020/03/`
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("01"))
	case "/2006/1/": // `/2020/3/`
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("1"))
	case "/2006/200601/": // `/2020/202003/`
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("200601"))
	case "/200601/": // `/202003/`
		return filepath.Join(baseDir, shootTime.Format("200601"))
	case "/2006/01/02/": // `/2020/03/01/`
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("01"), shootTime.Format("02"))
	case "/2006/0102/": // `/2020/0301/`
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("0102"))
	case "/20060102/": // `/20200301/`
		return filepath.Join(baseDir, shootTime.Format("20060102"))
	case "/2006/": // `/2020/`
		return filepath.Join(baseDir, shootTime.Format("2006"))
	default:
		// 默认情况下返回年和月格式
		return filepath.Join(baseDir, shootTime.Format("2006"), shootTime.Format("01"))
	}
}

// copyOrMoveFile 拷贝或移动文件，避免文件名重复
func copyOrMoveFile(src, dst string, isMove bool, skipSameFile bool) error {
	// 检查目标文件是否存在
	if infoDst, err := os.Stat(dst); err == nil {
		if skipSameFile {

			return fmt.Errorf("存在重复文件")
		} else {

			// 判断两个文件是否一致
			infoSrc, _ := os.Stat(src)

			// 生成文件唯一键
			fileKeySrc := generateFileKey(infoSrc.Size(), infoSrc.ModTime())
			fileKeyDst := generateFileKey(infoDst.Size(), infoDst.ModTime())

			if fileKeySrc == fileKeyDst {
				return fmt.Errorf("存在重复文件")
			}

			// 文件存在，生成新的文件名
			dst = generateUniqueFileName(dst)
		}
	}

	// 拷贝或移动文件
	var err error
	if isMove {
		err = os.Rename(src, dst) // 移动文件
	} else {
		err = copyFile(src, dst) // 拷贝文件
	}
	return err
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// generateUniqueFileName 生成唯一的文件名
func generateUniqueFileName(filePath string) string {
	dir, fileName := filepath.Split(filePath)
	ext := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, ext)

	var i int
	for {
		i++
		newFileName := fmt.Sprintf("%s (%d)%s", baseName, i, ext)
		newFilePath := filepath.Join(dir, newFileName)
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			return newFilePath
		}
	}
}

// 生成唯一键：文件大小 + 修改时间
func generateFileKey(size int64, modTime time.Time) string {
	return fmt.Sprintf("%d_%d", size, modTime.UnixNano())
}

// calculateFileMD5 计算文件的 MD5 值
func calculateFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// isSystemFile 检查文件是否为系统文件或隐藏文件
func isSystemFile(fileInfo os.FileInfo) bool {
	// 跳过已知的系统文件
	systemFiles := []string{"desktop.ini", "thumbs.db"}
	for _, systemFile := range systemFiles {
		if strings.EqualFold(fileInfo.Name(), systemFile) {
			return true
		}
	}

	// 检查是否是隐藏文件（Windows）
	return fileInfo.Name()[0] == '.'
}
