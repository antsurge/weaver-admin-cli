package cmd

import (
	"fmt"

	"github.com/antsurge/weaver-admin-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

func newCreateCommand() *cobra.Command {
	var opts scaffold.Options

	cmd := &cobra.Command{
		Use:   "create",
		Short: "从 weaver-admin 模板创建新项目",
		Long: `从 weaver-admin 模板创建全新业务项目。

自动完成以下工作:
  - 复制完整模板（后端 + 前端）
  - 全局替换项目名 / Go module / 品牌标识（大小写变体）
  - 清理 .git、演示数据、修复历史等模板残留
  - 可选自动执行 go mod tidy

示例:
  weaver create --name my-admin --module github.com/you/my-admin
  weaver create --name my-admin --source ../weaver-admin          # 本地模板
  weaver create --name my-admin --init                            # 生成后自动 go mod tidy
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			res, err := scaffold.Create(opts)
			if err != nil {
				return err
			}
			scaffold.PrintNextSteps(res)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Name, "name", "n", "", "新项目名（kebab-case，如 my-admin），必填")
	flags.StringVar(&opts.Module, "module", "", "新 Go module 路径（如 github.com/you/my-admin），为空时保留模板组织前缀")
	flags.StringVarP(&opts.Source, "source", "s", "", fmt.Sprintf("模板源：本地目录或 git 仓库地址（默认 %s）", scaffold.DefaultTemplateSource))
	flags.StringVarP(&opts.Output, "output", "o", ".", "输出目录（项目将生成在 <output>/<name> 下）")
	flags.BoolVar(&opts.Init, "init", false, "生成后自动执行 go mod tidy")
	flags.BoolVarP(&opts.Force, "force", "f", false, "目标目录已存在时强制覆盖")

	_ = cmd.MarkFlagRequired("name")
	return cmd
}
