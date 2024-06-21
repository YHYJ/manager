<h1 align="center">Manager</h1>

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-06-07 11:09:05 -->

---

<p align="center">
  <a href="https://github.com/YHYJ/manager/actions/workflows/release.yml"><img src="https://github.com/YHYJ/manager/actions/workflows/release.yml/badge.svg" alt="Go build and release by GoReleaser"></a>
</p>

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Install](#install)
  * [一键安装](#一键安装)
* [Usage](#usage)
* [Configuration](#configuration)
* [Compile](#compile)
  * [当前平台](#当前平台)
  * [交叉编译](#交叉编译)
    * [Linux](#linux)
    * [macOS](#macos)
    * [Windows](#windows)

<!-- vim-markdown-toc -->

---

<!--------------------------------------------------->
<!--  _ __ ___   __ _ _ __   __ _  __ _  ___ _ __  -->
<!-- | '_ ` _ \ / _` | '_ \ / _` |/ _` |/ _ \ '__| -->
<!-- | | | | | | (_| | | | | (_| | (_| |  __/ |    -->
<!-- |_| |_| |_|\__,_|_| |_|\__,_|\__, |\___|_|    -->
<!--                              |___/            -->
<!--------------------------------------------------->

---

自定义软件和脚本管理器，支持 Linux、macOS 和 Windows

## Install

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/YHYJ/manager/main/install.sh | sudo bash -s
```

## Usage

- `install`子命令

  该子命令用于安装/更新自开发的程序/脚本，有以下参数：

  - 'all'：安装/更新程序和脚本
  - 'go'： 安装/更新基于 go 开发的程序
    步骤：
    1. 获取远端程序版本
    2. 比较远端程序和本地程序版本
    3. 不一样则更新，一样则跳过
    4. 不一样时包含尚未安装到本地的情况，执行安装
  - 'shell'：安装/更新 shell 脚本
    步骤：
    1. 获取远端程序的 Hash 值
    2. 比较远端程序和本地程序的 Hash 值
    3. 不一样则更新，一样则跳过
    4. 不一样时包含尚未安装到本地的情况，执行安装

- `setup`子命令

  该子命令用于配置安装的程序/脚本，有以下参数：

  - 'all'：配置以下所有程序/脚本
  - 'docker'：配置 Docker Root Directory
  - 'git'：配置 git 用户信息等并生成 SSH 密钥
  - 'npm'：配置 npm 使用的镜像源
  - 'pip'：配置 pip 使用的镜像源

- `config`子命令

  该子命令用于操作配置文件，有以下参数：

  - 'create'：创建默认内容的配置文件，可以使用全局参数'--config'指定配置文件路径
  - 'force'：当指定的配置文件已存在时，使用该参数强制覆盖原文件
  - 'print'：打印配置文件内容

- `version`子命令

  查看程序版本信息

- `help`子命令

  查看程序帮助信息

## Configuration

1. 使用`config`子命令生成默认配置文件（具体使用方法执行`manager config --help`查看）
2. 参照如下说明修改配置文件：

```toml
[install]                               # 全局安装配置
  method = "release"                    # 安装方法，可选值为 [release, source]，分别为使用 Github 的 Release 和源码安装（release 方法目前仅支持基于 go 的程序）
  program_path = "/usr/local/bin"       # 程序/脚本存储路径
  release_temp = "/tmp/manager/release" # release 安装方法所下载文件的存储路径
  resources_path = "/usr/local/share"   # 资源文件存储路径
  source_temp = "/tmp/manager/source"   # source 安装方法所下载文件的存储路径

  [install.go]                                                                                             # 基于 go 的程序的安装配置
    completion_dir = ["/home/yj/.cache/oh-my-zsh/completions", "/home/yj/.oh-my-zsh/cache/completions"]    # 补全文件的存储路径
    generate_path = "build"                                                                                # 编译生成文件的存储路径
    gitea_api = "https://git.yj1516.top/api/v1"                                                            # source 安装方法 - Gitea 安装源 API 地址
    gitea_url = "https://git.yj1516.top"                                                                   # source 安装方法 - Gitea 安装源地址
    gitea_username = "YJ"                                                                                  # source 安装方法 - Gitea 安装源用户名
    github_api = "https://api.github.com"                                                                  # source 安装方法 - Github 安装源 API 地址
    github_url = "https://github.com"                                                                      # source 安装方法 - Github 安装源地址
    github_username = "YHYJ"                                                                               # source 安装方法 - Github 安装源用户名
    names = ["checker", "repos", "eniac", "kbdstage", "manager", "rolling", "scleaner", "skynet", "trash"] # 要安装的程序列表
    release_accept = "application/vnd.github+json"                                                         # release 安装方法 - API 请求头参数
    release_api = "https://api.github.com"                                                                 # release 安装方法 - API 地址

  [install.self]                                                                                           # 管理程序本身的配置
    completion_dir = ["/home/yj/.cache/oh-my-zsh/completions", "/home/yj/.oh-my-zsh/cache/completions"]    # 补全文件的存储路径
    generate_path = "build"                                                                                # 编译生成文件的存储路径
    gitea_api = "https://git.yj1516.top/api/v1"                                                            # source 安装方法 - Gitea 安装源 API 地址
    gitea_url = "https://git.yj1516.top"                                                                   # source 安装方法 - Gitea 安装源地址
    gitea_username = "YJ"                                                                                  # source 安装方法 - Gitea 安装源用户名
    github_api = "https://api.github.com"                                                                  # source 安装方法 - Github 安装源 API 地址
    github_url = "https://github.com"                                                                      # source 安装方法 - Github 安装源地址
    github_username = "YHYJ"                                                                               # source 安装方法 - Github 安装源用户名
    name = "manager"                                                                                       # 管理程序名
    release_accept = "application/vnd.github+json"                                                         # release 安装方法 - API 请求头参数
    release_api = "https://api.github.com"                                                                 # release 安装方法 - API 地址

  [install.shell]                                    # shell 脚本的安装配置
    dir = "System-Script/app"                        # shell 脚本在仓库中的路径
    gitea_api = "https://git.yj1516.top/api/v1"      # Gitea 安装源 API 地址
    gitea_branch = "ArchLinux"                       # Gitea 安装源分支名
    gitea_raw = "https://git.yj1516.top"             # Gitea 安装源文件下载地址
    gitea_username = "YJ"                            # Gitea 安装源用户名
    github_api = "https://api.github.com"            # Github 安装源 API 地址
    github_branch = "ArchLinux"                      # Github 安装源分支名
    github_raw = "https://raw.githubusercontent.com" # Github 安装源文件下载地址
    github_username = "YHYJ"                         # Github 安装源用户名
    repo = "Program"                                 # 存储脚本的仓库名
    names = [                                        # 要安装的 shell 脚本列表
      "collect-system",
      "configure-dtags",
      "py-virtualenv-tool",
      "save-docker-images",
      "sfm",
      "spacevim-update",
      "spider",
      "system-checkupdates",
      "trash-manager",
      "usb-manager",
    ]

[variable]
  http_proxy = ""  # HTTP 代理
  https_proxy = "" # HTTPS 代理
```

## Compile

### 当前平台

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/manager/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/manager/general.BuildTime=`date +%s` -X github.com/yhyj/manager/general.BuildBy=$USER" -o build/manager main.go
```

### 交叉编译

使用命令`go tool dist list`查看支持的平台

#### Linux

```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/manager/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/manager/general.BuildTime=`date +%s` -X github.com/yhyj/manager/general.BuildBy=$USER" -o build/manager main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### macOS

```bash
CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/manager/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/manager/general.BuildTime=`date +%s` -X github.com/yhyj/manager/general.BuildBy=$USER" -o build/manager main.go
```

> 使用`uname -m`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64

#### Windows

```powershell
CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -gcflags="-trimpath" -ldflags="-s -w -H windowsgui -X github.com/yhyj/manager/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/manager/general.BuildTime=`date +%s` -X github.com/yhyj/manager/general.BuildBy=$USER" -o build/manager.exe main.go
```

> 使用`echo %PROCESSOR_ARCHITECTURE%`确定硬件架构
>
> - 结果是 x86_64 则 GOARCH=amd64
> - 结果是 aarch64 则 GOARCH=arm64
