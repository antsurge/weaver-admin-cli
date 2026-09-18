package scaffold

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// copyExcludes 复制模板时需要整体跳过的目录（含其全部子内容）
var copyExcludes = map[string]bool{
	".git":         true,
	".codebuddy":   true,
	".workbuddy":   true,
	".idea":        true,
	"node_modules": true,
	"dist":         true, // 前端构建产物
	".turbo":       true, // turborepo 缓存
}

// copyFileExcludes 按文件名跳过的文件
var copyFileExcludes = map[string]bool{
	".DS_Store": true,
	"dump.rdb":  true, // 误提交的 redis dump
}

// CopyAndReplace 将 src 目录递归复制到 dst 目录，
// 文本文件内容执行 replacer 替换，二进制文件原样复制。
func CopyAndReplace(src, dst string, replacer *strings.Replacer) error {
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, relErr := filepath.Rel(src, path)
		if relErr != nil {
			return relErr
		}
		if rel == "." {
			return nil
		}

		base := filepath.Base(path)
		if copyExcludes[base] {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if copyFileExcludes[base] || strings.HasSuffix(base, ".log") {
			return nil
		}

		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		if !isBinary(data) {
			data = []byte(replacer.Replace(string(data)))
		}
		return os.WriteFile(target, data, info.Mode())
	})
}

// isBinary 通过是否包含 NUL 字节判断二进制文件
func isBinary(data []byte) bool {
	return bytes.IndexByte(data, 0) >= 0
}
