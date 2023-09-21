# README

<!-- File: README.md -->
<!-- Author: YJ -->
<!-- Email: yj1516268@outlook.com -->
<!-- Created Time: 2023-06-07 11:09:05 -->

---

## Table of Contents

<!-- vim-markdown-toc GFM -->

* [Usage](#usage)
* [Compile](#compile)

<!-- vim-markdown-toc -->

---

<!-- Object info -->

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
    1. 获取远端程序的Hash值
    2. 比较远端程序和本地程序的Hash值
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

- `help`

  查看程序帮助信息

## Compile

- 编译当前平台可执行文件：

```bash
go build main.go
```

- **交叉编译**指定平台可执行文件：

```bash
# 适用于Linux AArch64平台
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build main.go
```

```bash
# 适用于macOS amd64平台
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
```

```bash
# 适用于Windows amd64平台
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```
