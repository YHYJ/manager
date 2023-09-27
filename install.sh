#!/usr/bin/env bash

: << !
Name: install.sh
Author: YJ
Email: yj1516268@outlook.com
Created Time: 2023-09-27 10:21:55

Description: manager安装脚本

Attentions:
-

Depends:
-
!

set -eo pipefail

####################################################################
#+++++++++++++++++++++++++ Define Variable ++++++++++++++++++++++++#
####################################################################
#------------------------- Program Variable
# program name
name=$(basename "$0")
readonly name
# program version
readonly major_version=0.0.1
readonly minor_version=20230927
readonly rel_version=1

#------------------------- Exit Code Variable
readonly normal=0          # 一切正常
readonly err_file=1        # 文件/路径类错误
readonly err_param=2       # 参数错误
readonly err_run_command=3 # 命令执行错误

#------------------------- Parameter Variable
# description variable
readonly desc="用于安装manager"
# app variable
readonly app_name="manager"
readonly repo_url="https://git.yj1516.top/YJ/manager.git"
readonly repo_name="manager"
readonly repo_branch="main"
readonly repo_temp="/tmp/manager-build/$repo_name"
# install variable
readonly attribution="root"
readonly install_path="/usr/local/bin"

####################################################################
#+++++++++++++++++++++++++ Define Function ++++++++++++++++++++++++#
####################################################################
#------------------------- Info Function
function helpInfo() {
  echo -e ""
  echo -e "\x1b[32m$name\x1b[0m\x1b[1m$desc\x1b[0m"
  echo -e "--------------------------------------------------"
  echo -e "Usage:"
  echo -e ""
  echo -e "     $name [OPTION]"
  echo -e ""
  echo -e "Options:"
  echo -e "     -u, --update      升级manager"
  echo -e ""
  echo -e "     -h, --help        显示帮助信息"
  echo -e "     -v, --version     显示版本信息"
}

function versionInfo() {
  echo -e ""
  echo -e "\x1b[1m$name\x1b[0m version (\x1b[1m$major_version-$minor_version.$rel_version\x1b[0m)"
}

#------------------------- Feature Function
# 检测程序是否已安装
function checkProgram() {
  if [ -e "$install_path/$app_name" ]; then
    echo -e "\x1b[32mSuccess: $app_name found!\x1b[0m"
    exit $normal
  fi
}

# 克隆仓库
function cloneRepo() {
  if [ ! -e "$repo_temp" ]; then
    mkdir -p "$repo_temp"
    git clone -b "$repo_branch" "$repo_url" "$repo_temp" --depth=1
    result="$?"
    if [[ $result -eq 0 ]]; then
      echo -e "\x1b[32mSuccess: clone $repo_name to $repo_temp\x1b[0m"
    else
      echo -e "\x1b[31mError: clone $repo_name failed!\x1b[0m"
      exit $err_run_command
    fi
  else
    echo -e "\x1b[32mTrying to update $repo_name\x1b[0m"
    cd "$repo_temp" && git pull origin "$repo_branch"
    result="$?"
    if [[ $result -eq 0 ]]; then
      echo -e "\x1b[32mSuccess: update $repo_name\x1b[0m"
    else
      echo -e "\x1b[31mError: update $repo_name failed!\x1b[0m"
    fi
  fi
}

# 编译程序
function compileInstall() {
  cd "$repo_temp" || exit $err_file
  # 检测Makefile文件是否存在，存在则使用Makefile编译，不存在则使用go build编译
  if [ -e "Makefile" ]; then
    if command -v make &> /dev/null; then
      # Compile
      make
      # Install
      pkexec make install
    else
      echo -e "\x1b[31mError: 'make' not installed!\x1b[0m"
    fi
  elif [ -e "main.go" ]; then
    if command -v go &> /dev/null; then
      # Compile
      echo -e "\x1b[32m==>\x1b[0m Trying to compile project"
      go build -trimpath -ldflags="-s -w" -o "$app_name"
      echo -e "\x1b[32m[✔]\x1b[0m Successfully generated \x1b[32m$app_name\x1b[0m"
      # Install
      echo -e "\x1b[32m==>\x1b[0m Trying to install $app_name"
      pkexec install --mode=755 --owner=$attribution --group=$attribution $app_name $install_path/$app_name
      echo -e "\x1b[32m[✔]\x1b[0m Successfully installed \x1b[32m$app_name\x1b[0m"
    else
      echo -e "\x1b[31mError: 'go' not installed!\x1b[0m"
    fi
  else
    echo -e "\x1b[31mError: 'Makefile' or 'main.go' not found!\x1b[0m"
  fi
}

####################################################################
#++++++++++++++++++++++++++++++ Main ++++++++++++++++++++++++++++++#
####################################################################
ARGS=$(getopt --options "hv" --longoptions "help,version" -n "$name" -- "$@")
eval set -- "$ARGS"

if [[ ${#@} -lt 2 ]]; then
  checkProgram
  cloneRepo
  compileInstall
else
  while true; do
    case $1 in
      -h | --help)
        helpInfo
        shift 1
        ;;
      -v | --version)
        versionInfo
        shift 1
        ;;
      --)
        shift 1
        break
        ;;
      *)
        helpInfo
        exit $err_param
        ;;
    esac
  done
fi
