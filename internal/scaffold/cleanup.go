package scaffold

import (
	"os"
	"path/filepath"
)

// cleanupRelPaths 生成后需要删除的模板专属文件/目录
var cleanupRelPaths = []string{
	".git",
	".codebuddy",
	".workbuddy",
	".idea",
	".DS_Store",
	".ocrrc.yml",
	"dump.rdb",
	"frontend/dump.rdb",
	"screenshot.png",
	"BUG_README.md",
	"docs/bugfix-log.md", // 修复历史，对新项目无意义
	templateMetaFile,     // 模板元数据，生成后移除
}

// Cleanup 清理生成目录中的模板残留
func Cleanup(root string) error {
	for _, p := range cleanupRelPaths {
		full := filepath.Join(root, p)
		if err := os.RemoveAll(full); err != nil {
			return err
		}
	}
	return nil
}
