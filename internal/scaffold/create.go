package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Result create 命令的执行结果
type Result struct {
	// Dir 生成的项目目录（绝对路径）
	Dir string
	// Names 实际使用的命名映射
	Names *Names
}

// Create 执行完整的脚手架流程：
// 获取模板 -> 校验参数 -> 复制并替换 -> 清理 -> 可选初始化
func Create(opts Options) (*Result, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// 1. 获取模板
	srcDir, cleanup, err := fetchTemplate(opts.Source)
	if err != nil {
		return nil, err
	}
	defer cleanup()

	// 2. 读取模板元数据，构造命名映射
	meta, err := loadTemplateMeta(srcDir)
	if err != nil {
		return nil, err
	}
	names := BuildNames(meta, opts)

	// 3. 确定目标目录
	target := filepath.Join(opts.Output, names.Kebab)
	if exists(target) && !opts.Force {
		return nil, fmt.Errorf("目标目录 %s 已存在，请更换 --name 或使用 --force 覆盖", target)
	}
	if err := os.MkdirAll(opts.Output, 0o755); err != nil {
		return nil, err
	}
	if exists(target) {
		if err := os.RemoveAll(target); err != nil {
			return nil, fmt.Errorf("清理旧目标目录失败: %w", err)
		}
	}

	// 4. 复制模板并全局替换
	fmt.Printf("📦 正在复制模板到 %s ...\n", target)
	if err := CopyAndReplace(srcDir, target, names.Replacer()); err != nil {
		return nil, fmt.Errorf("复制模板失败: %w", err)
	}

	// 5. 清理模板残留
	if err := Cleanup(target); err != nil {
		return nil, fmt.Errorf("清理模板残留失败: %w", err)
	}

	// 6. 可选初始化
	if opts.Init {
		if err := runInit(target); err != nil {
			return nil, err
		}
	}

	return &Result{Dir: target, Names: names}, nil
}

// exists 判断路径是否存在
func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// runInit 在生成目录执行 go mod tidy
func runInit(dir string) error {
	fmt.Printf("🔧 执行 go mod tidy ...\n")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// PrintNextSteps 打印生成后的后续指引
func PrintNextSteps(res *Result) {
	name := res.Names.Kebab
	module := res.Names.Module
	fmt.Printf(`
✅ 项目创建成功！

  目录:    %s
  module:  %s

接下来:
  cd %s
  make installCli   # 安装 protoc/ent/wire 等工具（首次）
  make wire         # 生成依赖注入代码
  make ent          # 生成 ent 代码
  make buf-generate # 生成 proto 代码

然后按 docs 下的 SQL 初始化数据库，配置 configs/config.yaml，
前端进入 frontend/ 目录执行 pnpm install 即可启动。
`, filepath.Clean(name), module, name)
}
