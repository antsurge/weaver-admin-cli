package cmd

import (
	"fmt"
	"os"

	"github.com/antsurge/weaver-admin-cli/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "weaver",
	Short: "Weaver-admin 脚手架 CLI",
	Long: `weaver 是 weaver-admin（Kratos + Ent + Vue3/vben-admin 中后台）的项目脚手架工具。

通过一条命令从 weaver-admin 模板创建全新的业务项目，
自动完成项目名 / Go module / 品牌标识的批量替换与多余文件清理。

用法:
  weaver create --name my-admin --module github.com/you/my-admin
`,
}

// Execute 执行根命令
func Execute() {
	rootCmd.AddCommand(newCreateCommand())
	rootCmd.SetVersionTemplate(`weaver version {{.Version}}
`)
	rootCmd.Version = version.Version

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "❌ "+err.Error())
		os.Exit(1)
	}
}
