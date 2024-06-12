/*
File: define_message.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2024-05-29 15:51:14

Description: 定义输出信息及其格式
*/

package general

var (
	LatestVersionMessage       = "is already the latest version"                   // 输出文本 - 已安装的程序和脚本为最新版
	UnableToCompileMessage     = "Makefile or main.go file does not exist"         // 输出文本 - 缺失编译文件无法完成编译
	AcsInstallSuccessMessage   = "auto-completion script installed successfully"   // 输出文本 - 自动补全脚本安装成功
	AcsInstallFailedMessage    = "auto-completion script installation failed"      // 输出文本 - 自动补全脚本安装失败
	AcsUninstallSuccessMessage = "auto-completion script uninstalled successfully" // 输出文本 - 自动补全脚本卸载成功
	AcsUninstallFailedMessage  = "auto-completion script uninstallation failed"    // 输出文本 - 自动补全脚本卸载失败
)

var (
	GoLatestReleaseTagApiFormat      = "%s/repos/%s/%s/releases/latest" // API 和下载地址 - 请求远端仓库最新 Tag 的 API - Release
	GoLatestSourceTagApiFormat       = "%s/repos/%s/%s/tags"            // API 和下载地址 - 请求远端仓库最新 Tag 的 API - Source
	ShellLatestHashApiFormat         = "%s/repos/%s/%s/contents/%s/%s"  // API 和下载地址 - 请求远端仓库最新脚本的 Hash 值的 API
	ShellGithubBaseDownloadUrlFormat = "%s/%s/%s/%s"                    // API 和下载地址 - 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - GitHub 格式
	ShellGiteaBaseDownloadUrlFormat  = "%s/%s/%s/raw/branch/%s"         // API 和下载地址 - 远端仓库脚本基础下载地址（不包括在仓库路中的路径） - Gitea 格式
)

var (
	MultiSelectTips  = "Please select from the %s below (multi-select)\n"  // 提示词 - 多选
	SingleSelectTips = "Please select from the %s below (single-select)\n" // 提示词 - 单选
	QuietTips        = "Press '%s' to quit\n"                              // 提示词 - 退出
	SelectOneTips    = "Select %s"                                         // 提示词 - 单选
	SelectAllTips    = "Select All"                                        // 提示词 - 全选
	UninstallTips    = "Do you want to uninstall these software?"          // 提示词 - 卸载软件
)

var (
	OverWriteTips = "%s file already exists, do you want to overwrite it?" // 提示词 - 文件已存在是否覆写
)

var (
	InstallTips = "Please install %s first" // 提示词 - 需要安装
	InputTips   = "Please input '%s' value" // 提示词 - 输入
)
