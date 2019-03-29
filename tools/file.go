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

// 读取某个目录下的所有文件绝对路径, 包含子目录
func GetFilesAbsByDir(dirname string) ([]string, error) {
	// 获取正确合法不带末尾下划线的目录
	dir := filepath.Dir(dirname) + string(os.PathSeparator)
	allFiles, err := filepath.Glob(fmt.Sprintf("%s%s*", string(os.PathSeparator), dir))
	var files []string
	if err != nil {
		return files, fmt.Errorf("遍历目录%s发生异常: %s", dir, err.Error())
	}
	for _, file := range allFiles {
		// 文件数量太大可能会导致系统崩溃, 这里设置超出一定值则不再读取
		if len(files) > 10000 {
			logs.Warn(fmt.Sprintf("读取目录.文件数总数过大, 只保留10000, 建议切换目录层级"))
			break
		}
		f, err := os.Stat(file)
		if err != nil {
			logs.Error(fmt.Sprintf("读取目录.获取文件信息异常, 文件位置%s, 异常详情%s", file, err.Error()))
			continue
		}
		// 递归获取子目录
		if f.IsDir() {
			childFiles, err := GetFilesAbsByDir(file)
			if err != nil {
				continue
			}
			files = append(files, childFiles...)
			continue
		}
		files = append(files, file)
	}
	return files, nil
}

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
		logs.Error(fmt.Sprintf("获取子目录.读取目录%s失败, 异常信息: %s", dir, err.Error()))
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
	vo := &FilesVo{}
	err := getAllFiles(dir, vo)
	return vo, err
}

func getAllFiles(dir string, vo *FilesVo) error {
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			getAllFiles(path.Join(dir, fi.Name()), vo)
		} else {
			vo.Fs = append(vo.Fs, path.Join(dir, fi.Name()))
		}
	}
	return nil
}
