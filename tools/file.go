package tools

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/toolkits/file"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// 创建不存在的目录
func CreatePathIfNotExists(p string) (bool, error) {
	if !file.IsExist(p) {
		err := os.MkdirAll(p, os.FileMode(0755))
		if err != nil {
			return false, fmt.Errorf("创建不存在的目录失败: %v", err)
		}
	}
	return true, nil
}

// 创建不存在文件
func CreateFileIfNotExists(filename string) (bool, error) {
	_, err := CreatePathIfNotExists(filepath.Dir(filename))
	if err != nil {
		return false, err
	}
	if !file.IsFile(filename) {
		// 创建新文件
		f, err := os.Create(filename)
		defer f.Close()
		// 创建文件异常
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// 获取指定路径的子目录下的文件夹和文件信息
func GetSubDir(dirname string, ignore string, showHiddenDir bool) []map[string]interface{} {
	// 获取正确合法不带末尾下划线的目录
	dir := filepath.Dir(dirname + string(os.PathSeparator))
	var files []map[string]interface{}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		logs.Error(fmt.Sprintf("获取子目录.读取目录%s失败, 异常信息: %s", dir, err))
		return files
	}
	for _, fi := range rd {
		// 不显示忽略目录
		if ignore != "" && strings.Contains(fi.Name(), ignore) {
			continue
		}
		// 不显示隐藏文件夹
		if !showHiddenDir && strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		f := make(map[string]interface{})
		f["name"] = fi.Name()
		f["dir"] = fi.IsDir()
		files = append(files, f)
	}
	return files
}

// 文件排序
type SortFile struct {
	// 文件地址
	Abs string
	// 文件大小
	Size int64
	// 最后访问时间
	CreateTimestamp int64
}

// 文件排序
type SortFiles []SortFile

// 实现排序方法
func (sf SortFiles) Len() int {
	return len(sf)
}
func (sf SortFiles) Less(i, j int) bool {
	// 根据访问时间排序
	return sf[i].CreateTimestamp < sf[j].CreateTimestamp
}
func (sf SortFiles) Swap(i, j int) {
	sf[i], sf[j] = sf[j], sf[i]
}

// 文件信息按修改日期排序, 并且获取总数和大小
func GetFilesByModifyTimeSort(dirs []string) (int, int64, int64, SortFiles) {
	var fileInfos []string
	for _, dir := range dirs {
		vo, err := GetAllFiles(dir)
		if err != nil {
			logs.Error("[文件信息按日期排序]", "获取目录下所有信息失败", "目录", dir, "异常", err)
			continue
		}
		fileInfos = append(fileInfos, vo.Fs...)
	}

	var files SortFiles
	if len(fileInfos) == 0 {
		logs.Warn("[文件信息按日期排序]", "文件总数为空, 无须排序")
		return 0, 0, 0, files
	}
	// 总大小
	var size int64
	for _, f := range fileInfos {
		// 判断文件是否处于占用状态
		busy, fileInfo := IsFileBusy(f)
		if busy {
			logs.Error("[文件信息按日期排序]", "文件被占用, 跳过", "文件地址", f)
			continue
		}
		files = append(files, SortFile{
			Abs:             f,
			Size:            fileInfo.Size,
			CreateTimestamp: fileInfo.CreateTimestamp,
		})
		size += fileInfo.Size
	}
	// 文件信息排序
	sort.Sort(files)
	total := len(files)
	var firstModifyTime int64
	// 获取到最早日期
	if total > 0 {
		firstModifyTime = files[0].CreateTimestamp
	}
	return total, size, firstModifyTime, files
}

// 获取包含target目录的所有子目录
func GetAllTargetDirs(dir string, target string) ([]string, error) {
	rd, err := ioutil.ReadDir(dir)
	var dirs []string
	if err != nil {
		return dirs, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			if strings.Contains(fi.Name(), target) {
				dirs = append(dirs, filepath.Join(dir, fi.Name()))
			} else {
				// 递归获取子目录
				children, err := GetAllTargetDirs(filepath.Join(dir, fi.Name()), target)
				if err != nil {
					logs.Error(err)
					continue
				}
				dirs = append(dirs, children...)
			}
		}
	}
	return dirs, nil
}

type FilesVo struct {
	Fs []string `json:"fs"`
}

// 遍历获取指定目录下的所有文件
func GetAllFiles(dir string) (*FilesVo, error) {
	return GetAllFilesBySize(dir, -1, -1)
}

// 遍历获取指定目录下的所有文件
func GetAllFilesBySize(dir string, minSize int64, maxSize int64) (*FilesVo, error) {
	vo := &FilesVo{}
	err := getAllFiles(dir, vo, minSize, maxSize)
	return vo, err
}

func getAllFiles(dir string, vo *FilesVo, minSize int64, maxSize int64) error {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			getAllFiles(path.Join(dir, fi.Name()), vo, minSize, maxSize)
		} else {
			if minSize != -1 {
				if fi.Size() < minSize {
					continue
				}
			}
			if maxSize != -1 {
				if fi.Size() > maxSize {
					continue
				}
			}
			vo.Fs = append(vo.Fs, path.Join(dir, fi.Name()))
		}
	}
	return nil
}
