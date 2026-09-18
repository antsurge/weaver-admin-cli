package version

var (
	// Version 通过 -ldflags 注入
	Version = "1.0.0"
	// Commit 通过 -ldflags 注入
	Commit = "dev"
)
