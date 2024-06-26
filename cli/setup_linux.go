//go:build linux

/*
File: setup_linux.go
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-06-09 14:13:47

Description: 子命令 'setup' 的实现
*/

package cli

import (
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gookit/color"
	"github.com/yhyj/manager/general"
)

// 预定义输出格式
var (
	subjectMinorNameFormat = "%*s%s %s\n"
	descriptorFormat       = "%*s%s %s %s: %s %s %s\n"
	targetFileFormat       = "%*s%s %s %s: %s\n"
	askItemTitleFormat     = "%*s%s %s:\n"
	askItemsFormat         = "%*s%s "
	errorFormat            = "%*s%s %s: %s\n"
	statusFormat           = "%*s%s %s: %s\n"
	yesResultFormat        = "%*s%s %s\n"
	noResultFormat         = "%*s%s %s %s\n"
)

// ProgramConfigurator 程序配置器
// 参数：
//   - flags: 系统信息各部分的开关
func ProgramConfigurator(flags map[string]bool) {
	// 预定义配置参数变量
	var (
		hostname    = general.GetHostname()              // 主机名
		userName, _ = general.GetUserName()              // 用户真实名称
		home        = general.UserInfo.HomeDir           // 用户家目录
		sep         = strings.Repeat(" ", 4)             // 分隔符
		goBin       = filepath.Join(home, ".go", "bin")  // `go install` 命令安装目录
		noProxy     = "localhost,127.0.0.1,.example.com" // 默认不代理的 URL
	)

	var (
		subjectName      string       // 配置主题，应和程序/服务名称一致
		subjectMinorName string       // 配置子主题，在配置服务时识别是 timer 还是 service
		descriptorText   string       // 配置项描述
		writeMode        string = "t" // 配置文件写入模式，t 表示覆盖写入，a 表示追加写入
	)

	// 配置 chezmoi
	if flags["chezmoiFlag"] {
		// 提示
		subjectName = "chezmoi"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			ChezmoiDependencies = "chezmoi"                                                 // 主程序
			ChezmoiConfigFile   = filepath.Join(home, ".config", "chezmoi", "chezmoi.toml") // 配置文件
			// chezmoi 配置
			chezmoiConfigFormat = "sourceDir = %s\n[git]\n%sautoCommit = %v\n%sautoPush = %v\n"
			chezmoiSourceDir    = `"~/Documents/Repos/System/Profile"`
			chezmoiAutoCommit   = "false"
			chezmoiAutoPush     = "false"
		)

		// 检测
		if _, err := exec.LookPath(ChezmoiDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(ChezmoiConfigFile))
			// 创建配置文件
			if err := general.CreateFile(ChezmoiConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				chezmoiSourceDir, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "sourceDir")), chezmoiSourceDir)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				chezmoiAutoCommit, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "autoCommit")), chezmoiAutoCommit)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				chezmoiAutoPush, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "autoPush")), chezmoiAutoPush)

				// 配置
				ChezmoiConfigContent := color.Sprintf(chezmoiConfigFormat, chezmoiSourceDir, sep, chezmoiAutoCommit, sep, chezmoiAutoPush)
				if err := general.WriteFile(ChezmoiConfigFile, ChezmoiConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 cobra
	if flags["cobraFlag"] {
		// 提示
		subjectName = "cobra-cli"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// cobra 的依赖
			CobraDependencies = "cobra-cli"                        // 主程序
			CobraConfigFile   = filepath.Join(home, ".cobra.yaml") // 配置文件
			// cobra 配置
			cobraConfigFormat = "author: %s <%s>\nlicense: %s\nuseViper: %v\n"
			cobraAuthorName   = userName
			cobraAuthorEmail  = "email@example.com"
			cobraLicense      = "GPLv3"
			cobraUseViper     = "false"
		)

		// 检测
		if _, err := exec.LookPath(CobraDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(CobraConfigFile))
			// 创建配置文件
			if err := general.CreateFile(CobraConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraAuthorName, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "authorName")), cobraAuthorName)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraAuthorEmail, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "authorEmail")), cobraAuthorEmail)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraLicense, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "license")), cobraLicense)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				cobraUseViper, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "useViper")), cobraUseViper)

				// 配置
				CobraConfigContent := color.Sprintf(cobraConfigFormat, cobraAuthorName, cobraAuthorEmail, cobraLicense, cobraUseViper)
				if err := general.WriteFile(CobraConfigFile, CobraConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 docker
	if flags["dockerFlag"] {
		// 提示
		subjectName = "docker"
		descriptorText = "daemon configuration"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// docker 的依赖
			DockerDependencies      = "dockerd"                                            // 主程序
			DockerServiceConfigFile = "/etc/systemd/system/docker.service.d/override.conf" // 配置文件
			// docker 配置
			dockerServiceConfigFormat = "[Service]\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nEnvironment=\"%s\"\nExecStart=\nExecStart=%s --data-root=%s -H fd://\n"
			dockerServiceDataRoot     = filepath.Join(home, "Documents", "Docker", "Root")
		)

		// 检测
		if dockerdAbsPath, err := exec.LookPath(DockerDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(DockerServiceConfigFile))
			// 创建配置文件
			if err := general.CreateFile(DockerServiceConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				dockerServiceDataRoot, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "--data-root")), dockerServiceDataRoot)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				general.HttpProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTP_PROXY")), general.HttpProxy)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				general.HttpsProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "HTTPS_PROXY")), general.HttpProxy)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				noProxy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "NO_PROXY")), noProxy)

				// 需要获取交互结果的配置项
				dockerHttpProxy := color.Sprintf("HTTP_PROXY=%s", general.HttpProxy)
				dockerHttpsProxy := color.Sprintf("HTTPS_PROXY=%s", general.HttpsProxy)
				dockerNoProxy := color.Sprintf("NO_PROXY=%s", noProxy)

				// 配置
				DockerServiceConfigContent := color.Sprintf(dockerServiceConfigFormat, dockerHttpProxy, dockerHttpsProxy, dockerNoProxy, dockerdAbsPath, dockerServiceDataRoot)
				if err := general.WriteFile(DockerServiceConfigFile, DockerServiceConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}

				// 重载服务
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Rebirth"))
				rebirth(subjectName, 2, 2)
			}
		}
	}

	// 配置 frpc
	if flags["frpcFlag"] {
		// 提示
		subjectName = "frpc"
		descriptorText = "restart timing"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// frpc 的依赖
			FrpcDependencies = "frpc"                                             // 主程序
			FrpcConfigFile   = "/etc/systemd/system/frpc.service.d/override.conf" // 配置文件
			// frpc 配置
			frpcConfigFormat = "[Service]\nRestart=\nRestart=%s\n"
			frpcRestart      = "always"
		)

		// 检测
		if _, err := exec.LookPath(FrpcDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(FrpcConfigFile))
			// 创建配置文件
			if err := general.CreateFile(FrpcConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				frpcRestart, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "Restart")), frpcRestart)

				// 配置
				FrpcConfigContent := color.Sprintf(frpcConfigFormat, frpcRestart)
				if err := general.WriteFile(FrpcConfigFile, FrpcConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}

				// 重载服务
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Rebirth"))
				rebirth(subjectName, 2, 2)
			}
		}
	}

	// 配置 git
	if flags["gitFlag"] {
		// 提示
		subjectName = "git"
		descriptorText = "configuration file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// git 的依赖
			GitDependencies = "git"                             // 主程序
			GitConfigFile   = filepath.Join(home, ".gitconfig") // 配置文件
			// git 配置
			gitConfigFormat      = "[user]\n%sname = %s\n%semail = %s\n[core]\n%seditor = %s\n%sautocrlf = %s\n[merge]\n%stool = %s\n[color]\n%sui = %v\n[pull]\n%srebase = %v\n[filter \"lfs\"]\n%sclean = %s\n%ssmudge = %s\n%sprocess = %s\n%srequired = %v\n"
			gitUserEmail         = "email@example.com"
			gitCoreEditor        = "vim"
			gitCoreAutoCRLF      = "input"
			gitMergeTool         = "vimdiff"
			gitColorUI           = "true"
			gitPullRebase        = "false"
			gitFilterLfsClean    = "git-lfs clean -- %f"
			gitFilterLfsSmudge   = "git-lfs smudge -- %f"
			gitFilterLfsProcess  = "git-lfs filter-process"
			gitFilterLfsRequired = "true"
		)

		// 检测
		if _, err := exec.LookPath(GitDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(GitConfigFile))
			// 创建配置文件
			if err := general.CreateFile(GitConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				hostname, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[user].name")), hostname)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitUserEmail, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[user].email")), gitUserEmail)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitCoreEditor, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[core].editor")), gitCoreEditor)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitCoreAutoCRLF, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[core].autocrlf")), gitCoreAutoCRLF)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitMergeTool, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[merge].tool")), gitMergeTool)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitColorUI, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[color].ui")), gitColorUI)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitPullRebase, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[pull].rebase")), gitPullRebase)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsClean, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].clean")), gitFilterLfsClean)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsSmudge, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].smudge")), gitFilterLfsSmudge)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsProcess, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].process")), gitFilterLfsProcess)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				gitFilterLfsRequired, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[filter \"lfs\"].required")), gitFilterLfsRequired)

				// 配置
				GitConfigContent := color.Sprintf(gitConfigFormat, sep, hostname, sep, gitUserEmail, sep, gitCoreEditor, sep, gitCoreAutoCRLF, sep, gitMergeTool, sep, gitColorUI, sep, gitPullRebase, sep, gitFilterLfsClean, sep, gitFilterLfsSmudge, sep, gitFilterLfsProcess, sep, gitFilterLfsRequired)
				if err := general.WriteFile(GitConfigFile, GitConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 golang
	if flags["goFlag"] {
		// 提示
		subjectName = "go"
		descriptorText = "environment file"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// go 的依赖
			GolangDependencies = "go"                                        // 主程序
			GolangConfigFile   = filepath.Join(home, ".config", "go", "env") // 配置文件
			// go 配置
			golangConfigFormat = "GO111MODULE=%s\nGOBIN=%s\nGOPATH=%s\nGOCACHE=%s\nGOMODCACHE=%s\n"
			golangGO111MODULE  = "on"
			golangGOPATH       = filepath.Join(home, ".go")
			golangGOCACHE      = filepath.Join(home, ".cache", "go", "go-build")
			golangGOMODCACHE   = filepath.Join(home, ".cache", "go", "pkg", "mod")
		)

		// 检测
		if _, err := exec.LookPath(GolangDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(GolangConfigFile))
			// 创建配置文件
			if err := general.CreateFile(GolangConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				golangGO111MODULE, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "GO111MODULE")), golangGO111MODULE)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				golangGOPATH, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "GOPATH")), golangGOPATH)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				golangGOCACHE, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "GOCACHE")), golangGOCACHE)
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				golangGOMODCACHE, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "GOMODCACHE")), golangGOMODCACHE)

				// 配置
				GolangConfigContent := color.Sprintf(golangConfigFormat, golangGO111MODULE, goBin, golangGOPATH, golangGOCACHE, golangGOMODCACHE)
				if err := general.WriteFile(GolangConfigFile, GolangConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 pip
	if flags["pipFlag"] {
		// 提示
		subjectName = "pip"
		descriptorText = "mirrors"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))
		color.Printf(descriptorFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))

		// 配置项
		var (
			// pip 的依赖
			PipDependencies = "pip"                                             // 主程序
			PipConfigFile   = filepath.Join(home, ".config", "pip", "pip.conf") // 配置文件
			// pip 配置
			pipConfigFormat = "[global]\nindex-url = %s\ntrusted-host = %s\n"
			pipIndexUrl     = "https://mirrors.aliyun.com/pypi/simple"
		)

		// 检测
		if _, err := exec.LookPath(PipDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			color.Printf(targetFileFormat, 2, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(PipConfigFile))
			// 创建配置文件
			if err := general.CreateFile(PipConfigFile); err != nil {
				color.Printf(errorFormat, 2, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 2, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 4, " ", general.SuccessText("-"))
				pipIndexUrl, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "index-url")), pipIndexUrl)

				// 需要获取交互结果的配置项
				pipTrustedHost, _ := general.GetUrlHost(pipIndexUrl)

				// 配置
				PipConfigContent := color.Sprintf(pipConfigFormat, pipIndexUrl, pipTrustedHost)
				if err := general.WriteFile(PipConfigFile, PipConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 4, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 4, " ", general.SuccessText("-"), general.SuccessFlag)
				}
			}
		}
	}

	// 配置 system-checkupdates
	if flags["systemcheckupdatesFlag"] {
		// 提示
		subjectName = "system-checkupdates"
		color.Printf("%s %s\n", general.SuccessText("==>"), general.FgBlueText(subjectName))

		// 配置项
		var (
			// system-checkupdates Timer 和 Service 的依赖
			SystemCheckupdatesDependencies = "system-checkupdates"                             // 主程序，需要版本 >= 3.0.0-20230313.1
			timerConfigFile                = "/etc/systemd/system/system-checkupdates.timer"   // Timer 配置文件
			serviceConfigFile              = "/etc/systemd/system/system-checkupdates.service" // Service 配置文件
			// system-checkupdates 配置 - Timer
			timerConfigFormat      = "[Unit]\nDescription=%s\n\n[Timer]\nOnBootSec=%s\nOnUnitInactiveSec=%s\nAccuracySec=%s\nPersistent=%v\n\n[Install]\nWantedBy=%s\n"
			timerDescription       = "Timer for system-checkupdates"
			timerOnBootSec         = "10min"
			timerOnUnitInactiveSec = "2h"
			timerAccuracySec       = "30min"
			timerPersistent        = "true"
			timerWantedBy          = "timers.target"
			// system-checkupdates 配置 - Service
			serviceConfigFormat = "[Unit]\nDescription=%s\nAfter=%s\nWants=%s\n\n[Service]\nType=%s\nExecStart=%s\n"
			serviceDescription  = "System checkupdates"
			serviceAfter        = "network.target"
			serviceWants        = "network.target"
			serviceType         = "oneshot"
			serviceExecStart    = "/usr/local/bin/system-checkupdates --check"
		)

		// 检测
		if _, err := exec.LookPath(SystemCheckupdatesDependencies); err != nil {
			color.Printf(statusFormat, 2, " ", general.SuccessText("-"), general.LightText("Status"), general.NoticeText(color.Sprintf(general.InstallTips, subjectName)))
		} else {
			// ---------- Timer
			subjectMinorName = "timer"
			descriptorText = "timer"
			color.Printf(subjectMinorNameFormat, 2, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
			color.Printf(descriptorFormat, 4, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))
			color.Printf(targetFileFormat, 4, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(timerConfigFile))
			// 创建配置文件
			if err := general.CreateFile(timerConfigFile); err != nil {
				color.Printf(errorFormat, 4, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 4, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerDescription, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Unit].Description")), timerDescription)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerOnBootSec, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Timer].OnBootSec")), timerOnBootSec)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerOnUnitInactiveSec, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Timer].OnUnitInactiveSec")), timerOnUnitInactiveSec)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerAccuracySec, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Timer].AccuracySec")), timerAccuracySec)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerPersistent, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Timer].Persistent")), timerPersistent)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				timerWantedBy, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Install].WantedBy")), timerWantedBy)

				// 配置
				SystemCheckupdatesTimerConfigContent := color.Sprintf(timerConfigFormat, timerDescription, timerOnBootSec, timerOnUnitInactiveSec, timerAccuracySec, timerPersistent, timerWantedBy)
				if err := general.WriteFile(timerConfigFile, SystemCheckupdatesTimerConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 6, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 6, " ", general.SuccessText("-"), general.SuccessFlag)
				}

				// 重载服务
				color.Printf(askItemTitleFormat, 4, " ", general.SuccessText("-"), general.LightText("Rebirth"))
				name := color.Sprintf("%s.%s", subjectName, subjectMinorName)
				rebirth(name, 2, 3)
			}

			// ---------- Service
			subjectMinorName = "service"
			descriptorText = "service"
			color.Printf(subjectMinorNameFormat, 2, " ", general.SuccessText("-"), general.FgBlueText(subjectMinorName))
			color.Printf(descriptorFormat, 4, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Descriptor"), general.SecondaryText("Set up"), general.SecondaryText(subjectName), general.SecondaryText(descriptorText))
			color.Printf(targetFileFormat, 4, " ", general.SuccessText("-"), general.InfoText("INFO:"), general.LightText("Target file"), general.CommentText(serviceConfigFile))
			// 创建配置文件
			if err := general.CreateFile(serviceConfigFile); err != nil {
				color.Printf(errorFormat, 4, " ", general.SuccessText("-"), general.LightText("Error"), general.DangerText(err))
			} else {
				// 交互
				color.Printf(askItemTitleFormat, 4, " ", general.SuccessText("-"), general.LightText("Config"))
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				serviceDescription, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Unit].Description")), serviceDescription)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				serviceAfter, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Unit].After")), serviceAfter)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				serviceWants, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Unit].Wants")), serviceWants)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				serviceType, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Service].Type")), serviceType)
				color.Printf(askItemsFormat, 6, " ", general.SuccessText("-"))
				serviceExecStart, _ = general.GetInput(general.QuestionText(color.Sprintf(general.InputTips, "[Service].ExecStart")), serviceExecStart)

				// 配置
				SystemCheckupdatesServiceConfigContent := color.Sprintf(serviceConfigFormat, serviceDescription, serviceAfter, serviceWants, serviceType, serviceExecStart)
				if err := general.WriteFile(serviceConfigFile, SystemCheckupdatesServiceConfigContent, writeMode); err != nil {
					color.Printf(noResultFormat, 6, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
				} else {
					color.Printf(yesResultFormat, 6, " ", general.SuccessText("-"), general.SuccessFlag)
				}

				// 重载服务
				color.Printf(askItemTitleFormat, 4, " ", general.SuccessText("-"), general.LightText("Rebirth"))
				name := color.Sprintf("%s.%s", subjectName, subjectMinorName)
				rebirth(name, 2, 3)
			}
		}
	}
}

// rebirth 服务配置完成后的重载/启用/重启
//
// 参数：
//   - name: 服务名称
//   - margin: 对齐时的边距
//   - coefficient: 边距应乘的系数
func rebirth(name string, margin, coefficient int) {
	cMargin := margin * coefficient

	// 重新加载 systemd 管理器配置
	reloadArgs := []string{"daemon-reload"} // 重载服务配置
	if _, _, err := general.RunCommandToBuffer("systemctl", reloadArgs); err != nil {
		color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(err))
	}

	// 询问是否需要启用/重启服务
	checkStatusArgs := []string{"is-enabled", name} // 检测服务启用状态
	status, _, _ := general.RunCommandToBuffer("systemctl", checkStatusArgs)
	switch status {
	case "enabled":
		color.Printf(askItemsFormat, cMargin, " ", general.SuccessText("-"))

		question := color.Sprintf(general.RestartServiceTips, name)
		restart, err := general.AskUser(general.QuestionText(question), []string{"y", "N"})
		if err != nil {
			color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText("Unable to get answers: ", err))
		}
		switch restart {
		case "y":
			// 重启服务
			restartArgs := []string{"restart", name}
			if _, stderr, err := general.RunCommandToBuffer("systemctl", restartArgs); err != nil {
				color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(stderr))
			} else {
				color.Printf(yesResultFormat, cMargin, " ", general.SuccessText("-"), general.SuccessFlag)
			}
		case "n":
			return
		default:
			color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.WarningFlag, general.WarnText("Unexpected answer: ", restart))
		}
	case "disabled":
		color.Printf(askItemsFormat, cMargin, " ", general.SuccessText("-"))

		question := color.Sprintf(general.EnableServiceTips, name)
		enable, err := general.AskUser(general.QuestionText(question), []string{"y", "N"})
		if err != nil {
			color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText("Unable to get answers: ", err))
		}
		switch enable {
		case "y":
			// 启用服务（并立即运行）
			enableArgs := []string{"enable", "--now", name}
			if _, stderr, err := general.RunCommandToBuffer("systemctl", enableArgs); err != nil {
				color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(stderr))
			} else {
				color.Printf(yesResultFormat, cMargin, " ", general.SuccessText("-"), general.SuccessFlag)
			}
		case "n":
			return
		default:
			color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.WarningFlag, general.WarnText("Unexpected answer: ", enable))
		}
	case "static":
		color.Printf(yesResultFormat, cMargin, " ", general.SuccessText("-"), general.SuccessFlag)
	case "not-found":
		notFound := color.Sprintf(general.NotFoundServiceTips, name)
		color.Printf(noResultFormat, cMargin, " ", general.SuccessText("-"), general.ErrorFlag, general.DangerText(notFound))
	}
}
