package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// fetchTemplate 获取模板目录：
//   - source 为本地已存在目录时直接使用
//   - 否则视为 git 仓库地址，clone 到临时目录（用完清理）
//
// 返回模板目录与清理函数。
func fetchTemplate(source string) (dir string, cleanup func(), err error) {
	if info, statErr := os.Stat(source); statErr == nil && info.IsDir() {
		abs, absErr := filepath.Abs(source)
		if absErr != nil {
			return "", nil, absErr
		}
		return abs, func() {}, nil
	}

	fmt.Printf("🌐 正在从 %s 拉取模板 ...\n", source)
	tmp, tmpErr := os.MkdirTemp("", "weaver-template-*")
	if tmpErr != nil {
		return "", nil, tmpErr
	}
	cmd := exec.Command("git", "clone", "--depth", "1", source, tmp)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if runErr := cmd.Run(); runErr != nil {
		_ = os.RemoveAll(tmp)
		return "", nil, fmt.Errorf("拉取模板失败: %w", runErr)
	}
	return tmp, func() { _ = os.RemoveAll(tmp) }, nil
}
