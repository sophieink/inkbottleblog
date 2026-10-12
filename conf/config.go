package conf

import (
	"path/filepath"
	"os"
)

const (
	DefaultFilePerm os.FileMode = 0o644 //默认文件权限
	DefaultDirPerm  os.FileMode = 0o755 //默认目录权限
)

var (
	RootDir="."
	PostsDir=filepath.Join(RootDir, "posts")
	// *.md
	DocsDir=filepath.Join(RootDir, "docs")
	// *.html
	TemplateDir=filepath.Join(RootDir, "web", "templates")
	StaticDir=filepath.Join(RootDir, "web", "static")
	DocsStaticDir=filepath.Join(DocsDir, "static")
)
