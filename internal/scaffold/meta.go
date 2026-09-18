package scaffold

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// TemplateMeta 模板根目录的 .weaver-template.yaml 元数据，
// 声明模板原始标识，CLI 据此构建替换规则。
type TemplateMeta struct {
	// Name 模板项目名（kebab-case）
	Name string `yaml:"name"`
	// Module 模板 Go module 路径
	Module string `yaml:"module"`
}

// templateMetaFile 模板元数据文件名
const templateMetaFile = ".weaver-template.yaml"

// loadTemplateMeta 读取模板元数据
func loadTemplateMeta(dir string) (TemplateMeta, error) {
	var meta TemplateMeta
	data, err := os.ReadFile(filepath.Join(dir, templateMetaFile))
	if err != nil {
		return meta, fmt.Errorf("模板源缺少 %s，不是有效的 weaver-admin 模板: %w", templateMetaFile, err)
	}
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return meta, fmt.Errorf("解析 %s 失败: %w", templateMetaFile, err)
	}
	if meta.Name == "" || meta.Module == "" {
		return meta, fmt.Errorf("%s 中 name/module 字段不能为空", templateMetaFile)
	}
	return meta, nil
}
