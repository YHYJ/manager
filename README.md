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

* [Usage](#usage)
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

自开发程序管理器

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
  - 'pip'：配置 pip 使用的镜像源
  - 'npm'：配置 npm 使用的镜像源
  - 'docker'：配置 Docker Root Directory
  - 'git'：配置 git 用户信息等并生成 SSH 密钥

- `config`子命令

  该子命令用于操作配置文件，有以下参数：

  - 'create'：创建默认内容的配置文件，可以使用全局参数'--config'指定配置文件路径
  - 'force'：当指定的配置文件已存在时，使用该参数强制覆盖原文件
  - 'print'：打印配置文件内容

- `version`子命令

  查看程序版本信息

- `help`子命令

  查看程序帮助信息

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
