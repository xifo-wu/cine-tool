package symlink

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func CreateSymlink(sourceFile, targetPath string) error {
	fileInfo, err := os.Stat(sourceFile)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	mode := fileInfo.Mode()
	if mode.IsDir() {

		err := filepath.Walk(sourceFile, func(path string, info os.FileInfo, err error) error {
			relativePath, _ := filepath.Rel(sourceFile, path)

			linkPath := filepath.Join(targetPath, relativePath)

			if !info.IsDir() && IsExtensionMatch(linkPath) {
				createLink(path, linkPath)
				if err != nil {
					log.Printf("创建软链接失败：%s -> %s", linkPath, path)
					return err
				}

				log.Printf("软链接已成功创建：%s -> %s", linkPath, path)
			}

			return err
		})

		return err
	}

	err = createLink(sourceFile, targetPath)
	if err != nil {
		log.Printf("创建软链接失败：%s -> %s", sourceFile, targetPath)
		return err
	}

	log.Printf("软链接已成功创建：%s -> %s", sourceFile, targetPath)
	return nil
}

func createLink(path, linkPath string) error {
	os.MkdirAll(filepath.Dir(linkPath), 0755)
	return os.Symlink(path, linkPath)
}

func IsExtensionMatch(filePath string) bool {
	extensions := viper.GetStringSlice("EXTENSIONS")
	for _, ext := range extensions {
		if filepath.Ext(filePath) == ext {
			return true
		}
	}

	return false
}

func IsSubdirectory(parentPath string, childPath string) bool {
	// 获取父目录路径和子目录路径的绝对路径
	parentAbsPath, err := filepath.Abs(parentPath)
	if err != nil {
		fmt.Println("Error getting absolute path for the parent directory:", err)
		return false
	}

	childAbsPath, err := filepath.Abs(childPath)
	if err != nil {
		fmt.Println("Error getting absolute path for the child directory:", err)
		return false
	}

	// 获取子目录相对于父目录的相对路径
	relPath, err := filepath.Rel(parentAbsPath, childAbsPath)
	if err != nil {
		fmt.Println("Error getting relative path between directories:", err)
		return false
	}

	// 检查相对路径是否不包含父目录之外的内容
	return !strings.HasPrefix(relPath, "..") && !filepath.IsAbs(relPath)
}
