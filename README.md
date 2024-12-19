<h1 align="center">Manager</h1>
<h3 align="center">自定义软件和脚本管理器</h3>

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

* [适配](#适配)
* [安装](#安装)
  * [一键安装](#一键安装)
  * [编译安装](#编译安装)
    * [当前平台](#当前平台)
    * [交叉编译](#交叉编译)
* [用法](#用法)

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

## 适配

- Linux: 适配
- macOS: 适配
- Windows: 适配

## 安装

### 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/YHYJ/manager/main/install.sh | sudo bash -s
```

也可以从 [GitHub Releases](https://github.com/YHYJ/manager/releases) 下载解压后使用

### 编译安装

#### 当前平台

如果要为当前平台编译，可以使用以下命令：

```bash
go build -gcflags="-trimpath" -ldflags="-s -w -X github.com/yhyj/manager/general.GitCommitHash=`git rev-parse HEAD` -X github.com/yhyj/manager/general.BuildTime=`date +%s` -X github.com/yhyj/manager/general.BuildBy=$USER" -o build/manager main.go
```

#### 交叉编译

> 使用命令`go tool dist list`查看支持的平台
>
> Linux 和 macOS 使用命令`uname -m`，Windows 使用命令`echo %PROCESSOR_ARCHITECTURE%` 确认系统架构
>
> - 例如 x86_64 则设 GOARCH=amd64
> - 例如 aarch64 则设 GOARCH=arm64
> - ...

设置如下系统变量后使用 [编译安装](#编译安装) 的命令即可进行交叉编译：

- CGO_ENABLED: 不使用 CGO，设为 0
- GOOS: 设为 linux, darwin 或 windows
- GOARCH: 根据当前系统架构设置

## 用法

- `install`子命令

  该子命令用于安装/更新自开发的程序/脚本，有以下参数：

  - '--all'：安装/更新程序和脚本
  - '--go'： 安装/更新基于 go 开发的程序

    步骤：

    1. 获取远端程序版本
    2. 比较远端程序和本地程序版本
    3. 版本不一样则更新，一样则跳过
    4. 版本不一样时包含尚未安装到本地的情况，执行安装

  - '--shell'：安装/更新 shell 脚本

    步骤：

    1. 获取远端程序的哈希值
    2. 比较远端程序和本地程序的哈希值
    3. 哈希值不一样则更新，一样则跳过
    4. 哈希值不一样时包含尚未安装到本地的情况，执行安装

- `setup`子命令

  配置指定程序，有以下参数：

  - '--all'：配置以下所有程序
  - '--chezmoi'：配置 chezmoi
  - '--cobra'：配置 cobra
  - '--docker'：配置 Docker
  - '--frpc'：配置 frpc
  - '--git'：配置 git 并生成 SSH 密钥
  - '--go'：配置 golang
  - '--pip'：配置 pip 使用的镜像源
  - '--update-checker'：配置系统更新检测服务

- `config`子命令

  操作配置文件，有以下参数：

  - '--create'：交互式创建配置文件
  - '--open'：使用系统默认编辑器打开配置文件
  - '--print'：打印配置文件内容

- `version`子命令

  查看程序版本信息

- `help`子命令

  查看程序帮助信息
