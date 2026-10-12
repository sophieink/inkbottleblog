package internal
import (
	"io"
	"os"
	"fmt"
	"path/filepath"
	"strings"

	"inkbottle/conf"
)

var (
	//continue filename
	skipfiles = map[string]bool{ 
		"README.md": true,
	}
	//continue 跳过的文件后缀
  skipsuffixes = []string { 
		".tmp", ".swp", ".log", "~",
	}
	//continue index
	skipdirs = map[string]bool {
		".git": true,
		".github": true,
		".cache": true,
	}
)

//判断某个文件名是否需要跳过
func shouldSkip(n string, isDir bool) bool {
	if isDir {
		return skipdirs[n]
	}
	if skipfiles[n] {
		return true
	}
	lower := strings.ToLower(n)
	for _, suf := range skipsuffixes {
		if suf == "" {
			continue
		}
		if strings.HasSuffix(lower, strings.ToLower(suf)) {
			return true
		}
	}
	return false
}
// 确保目录在，不在就递归创建
func EnsureDir(dir string) error {
	if err := os.MkdirAll(dir, conf.DefaultDirPerm); err != nil {
		return fmt.Errorf("mkdir %s: %w", dir, err)
	}
	return nil
}

// witefile,父目录不存在自动创建
func WriteFile(path string, data []byte) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, conf.DefaultFilePerm); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// ReadFile 读文件
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	return data, nil
}

//手动递归
func CopyDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("stat src %s: %w", src, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("src %s is not a directory", src)
	}
	return copyRecursive(src, dst)
}

func copyRecursive(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("readdir %s: %w", src, err)
	}
	for _, e := range entries {
		n := e.Name()
		if shouldSkip(n, e.IsDir()) {
			continue
		}
		srcPath := filepath.Join(src, n)
		dstPath := filepath.Join(dst, n)
		if e.IsDir() {
			if err := os.MkdirAll(dstPath, conf.DefaultDirPerm); err != nil {
				return fmt.Errorf("mkdir %s: %w", dstPath, err) 
			}
			if err := copyRecursive(srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		if !e.Type().IsRegular() {
			continue
		}
		if err := copyFile(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

// 递归复制web/static/
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open src %s: %w", src, err)	
	}
	defer in.Close()

	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, conf.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("open dst %s: %w", dst, err)
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("copy %s -> %s: %w", src, dst, err)
	}
	return nil
}
