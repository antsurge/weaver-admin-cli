package scaffold

import "strings"

// Names 一个项目名的全部命名变体，用于全局批量替换
type Names struct {
	// 模板侧的原始标识
	TemplateKebab  string // weaver-admin
	TemplatePascal string // WeaverAdmin
	TemplateCamel  string // weaverAdmin
	TemplateSnake  string // weaver_admin
	TemplateModule string // github.com/antsurge/weaver-admin

	// 新项目侧的目标标识
	Kebab  string // my-admin
	Pascal string // MyAdmin
	Camel  string // myAdmin
	Snake  string // my_admin
	Module string // github.com/you/my-admin
}

// BuildNames 基于模板元数据与用户参数构造命名映射
func BuildNames(meta TemplateMeta, opts Options) *Names {
	n := &Names{
		TemplateKebab:  meta.Name,
		TemplatePascal: toPascal(meta.Name),
		TemplateCamel:  toCamel(meta.Name),
		TemplateSnake:  toSnake(meta.Name),
		TemplateModule: meta.Module,

		Kebab:  opts.Name,
		Pascal: toPascal(opts.Name),
		Camel:  toCamel(opts.Name),
		Snake:  toSnake(opts.Name),
		Module: opts.Module,
	}
	// module 为空时，从模板 module 推导：保留组织前缀 + 新项目名
	if n.Module == "" {
		n.Module = defaultModule(n.TemplateModule, n.Kebab)
	}
	return n
}

// Replacer 生成批量替换器。
// strings.Replacer 按"最左最长"匹配且对替换结果不重复替换，因此
// 先替换精确的 module 前缀，再替换短标识是安全的。
func (n *Names) Replacer() *strings.Replacer {
	return strings.NewReplacer(
		n.TemplateModule, n.Module,
		n.TemplatePascal, n.Pascal,
		n.TemplateCamel, n.Camel,
		n.TemplateSnake, n.Snake,
		n.TemplateKebab, n.Kebab,
	)
}

// defaultModule 推导默认 module：模板 module 的组织前缀 + 新项目名
func defaultModule(templateModule, name string) string {
	idx := strings.LastIndex(templateModule, "/")
	if idx < 0 {
		return name
	}
	return templateModule[:idx+1] + name
}

// toPascal 将 kebab-case 转 PascalCase：my-admin -> MyAdmin
func toPascal(kebab string) string {
	parts := strings.Split(kebab, "-")
	for i, p := range parts {
		if p == "" {
			continue
		}
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	return strings.Join(parts, "")
}

// toCamel 将 kebab-case 转 camelCase：my-admin -> myAdmin
func toCamel(kebab string) string {
	p := toPascal(kebab)
	if p == "" {
		return p
	}
	return strings.ToLower(p[:1]) + p[1:]
}

// toSnake 将 kebab-case 转 snake_case：my-admin -> my_admin
func toSnake(kebab string) string {
	return strings.ReplaceAll(kebab, "-", "_")
}
