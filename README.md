# weaver-admin-cli

`weaver` 是 [weaver-admin](https://github.com/antsurge/weaver-admin)（Kratos + Ent + Vue3/vben-admin 中后台）的项目脚手架命令行工具。

通过一条命令从 weaver-admin 模板创建全新业务项目，自动完成：

- 复制完整模板（后端 + 前端）
- 全局替换项目名 / Go module / 品牌标识（完整覆盖 `weaver-admin`、`WeaverAdmin`、`weaverAdmin`、`weaver_admin` 四种变体）
- 清理 `.git`、`.idea`、构建产物、修复历史等模板残留
- 可选自动执行 `go mod tidy`

## 安装

```bash
go install github.com/antsurge/weaver-admin-cli/cmd/weaver@latest
```

## 使用

```bash
# 从远程模板仓库创建（默认即官方模板）
weaver create --name my-admin --module github.com/you/my-admin

# 从本地模板目录创建（开发/验证时更快）
weaver create --name my-admin --source ../weaver-admin

# 生成后自动 go mod tidy
weaver create --name my-admin --init

# 目标目录已存在时强制覆盖
weaver create --name my-admin --force
```

### 参数

| 参数 | 说明 | 默认值 |
|---|---|---|
| `--name, -n` | 新项目名（kebab-case，如 `my-admin`），必填 | - |
| `--module` | 新 Go module 路径 | 模板组织前缀 + 新项目名 |
| `--source, -s` | 模板源：本地目录或 git 仓库地址 | `https://github.com/antsurge/weaver-admin.git` |
| `--output, -o` | 输出目录，项目生成在 `<output>/<name>` | 当前目录 |
| `--init` | 生成后自动执行 `go mod tidy` | 否 |
| `--force, -f` | 覆盖已存在的目标目录 | 否 |

## 生成后的项目

```bash
cd my-admin
make installCli   # 安装 protoc/ent/wire 等开发工具（首次）
make wire         # 生成依赖注入代码
make ent          # 生成 ent 代码
make buf-generate # 生成 proto 代码
```

数据库初始化与配置详见生成项目中 `docs/` 与 `configs/config.yaml`。

## 模板机制

`weaver-admin` 仓库根目录的 `.weaver-template.yaml` 声明模板原始标识（`name` / `module`），
CLI 读取该文件后基于 `strings.Replacer`（最左最长匹配、对结果不重复替换）完成全局替换。

## 开发

```bash
go build ./cmd/weaver
go vet ./...
```

### 本地快速验证

```bash
go build -o /tmp/weaver ./cmd/weaver
/tmp/weaver create --name demo-admin --source /path/to/weaver-admin
```