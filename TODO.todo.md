- [X] 优化输出信息 (2023-06-16 13:40)
- [X] 完成install --script部分 (2023-06-16 16:36)
- [X] 默认配置项增加'go.fallbacksource'和'shell.fallbacksource' (2023-07-12 16:18)
- [X] 将基于文件内容不同判断是否更新程序的逻辑改为根据版本信息 (2023-09-15 16:27)
- [X] 将shell更新分支更改为使用`git hash-object`结果和`curl https://git.yj1516.top/api/v1/repos/<用户名>/<仓库名>/contents/<文件路径>`结果中的'sha'值进行比较来判断脚本是否需要更新 (2023-09-20 15:35) (2023-09-20 15:36)
- [X] 将shell更新分支中下载Program仓库更新为使用'wget https://git.yj1516.top/<用户名>/<仓库名>/raw/branch/main/<文件路径>'下载 (2023-09-20 15:38)
- [X] 'fallback*'系列参数已定义但还未使用 (2023-09-19 15:43)
- [X] 自动补全脚本位置在变量$fpath中 (2023-10-10 19:38)
- [X] 代理默认为空 (2023-10-18 09:00)
- [X] 自动补全脚本位置完善（给定一个列表依次判断是否存在） (2023-10-18 09:01)
- [X] 使`install`子命令适配Windows，`setup`子命令提示不支持Windows (2023-11-13 16:23)
- [X] 优先安装Github Release里打包好的程序，其次使用Github仓库代码下载编译安装，Gitea作为最后的备用 (2023-11-22 10:11)
