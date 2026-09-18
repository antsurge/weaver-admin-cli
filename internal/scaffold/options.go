package scaffold

import (
	"fmt"
	"regexp"
)

// DefaultTemplateSource 默认模板源（远程 git 仓库）
const DefaultTemplateSource = "https://github.com/antsurge/weaver-admin.git"

// Options create 命令的全部参数
type Options struct {
	// Name 新项目名（kebab-case，如 my-admin）
	Name string
	// Module 新 Go module 路径（如 github.com/you/my-admin），为空时自动推导
	Module string
	// Source 模板源：本地目录路径或 git 仓库地址，为空时使用默认值
	Source string
	// Output 输出目录，默认当前目录
	Output string
	// Init 生成后自动执行 go mod tidy
	Init bool
	// Force 目标目录已存在时强制覆盖
	Force bool
}

var namePattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// Validate 校验并补齐参数
func (o *Options) Validate() error {
	if o.Name == "" {
		return fmt.Errorf("项目名不能为空，请使用 --name 指定（kebab-case，如 my-admin）")
	}
	if !namePattern.MatchString(o.Name) {
		return fmt.Errorf("项目名 %q 格式不合法，需为 kebab-case（小写字母/数字，以 - 分隔，如 my-admin）", o.Name)
	}
	if o.Output == "" {
		o.Output = "."
	}
	if o.Source == "" {
		o.Source = DefaultTemplateSource
	}
	return nil
}
