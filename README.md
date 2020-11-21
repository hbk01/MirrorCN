# MirrorCN


根据配置文件更改镜像源地址，您可以定义自己的配置文件。

**注意：仅限 Linux! 我没有在其他平台进行测试**

# 运行

```bash
USAGE:
    mcn [ [ --update ] | [ [--config=FILE] PackageManagerTag ] ]

EXAMPLE:
    mcn go # 修改go的源，默认使用与mcn命令同目录的配置文件
    mcn --config="$HOME/MirrorCN.json" pacman # 使用指定的配置文件修改pacman的源
    mcn --update # 下载最新的配置文件
```

# 配置文件

```json
{
    // 如果 raw.githubusercontent.com 可以访问的话，用 Visual Studio Code 打开配置文件会有提示，也可以使用本地的文件
    "$schema": "https://raw.githubusercontent.com/hbk01/MirrorCN/master/schema.json",
    // 执行 mcn update 将会从此处下载最新的配置文件
    "update": "https://raw.githubusercontent.com/hbk01/MirrorCN/master/config.json",
    // 所有可以更换的源
    "list": [
        {
            // pacman 就是命令中必填的 PackageManagerTag
            "pacman": {
                // 在添加源之前执行，出错自动停止替换，用于检查时候安装相关软件
                "first": [
                    "pacman --version"
                ],
                // 在添加完成后执行，用于刷新数据库
                "end": [
                    "pacman -Syy"
                ],
                // 需要修改哪些文件
                "files": [
                    "/etc/pacman.d/mirrorlist"
                ],
                // 要添加什么源
                "mirrors": [
                    "https://mirrors.tuna.tsinghua.edu.cn/manjaro/stable/$repo/$arch"
                ],
                // 什么样的格式添加进去
                // $mirrors 将被替换为具体的镜像源地址
                // $tag     将被替换为 PackageManagerTag
                // $date    将被替换为 yyyy-MM-dd 格式的日期，如 2020-11-21
                // $time    将被替换为 HH:mm:ss 格式的时间，如 15:41:25
                // $name    将被替换为本程序名称
                // $url     将被替换为本程序开源地址
                "format": [
                    "## China Server - Added by $name at $date $time",
                    "## $name project is open-source on github: $url",
                    "Server = $mirrors"
                ]
            }
        }
    ]
}
```
