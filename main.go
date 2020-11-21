package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PackageManager struct { // {{{
	Name    string   `json:"name"`
	First   []string `json:"first"`
	End     []string `json:"end"`
	Files   []string `json:"files"`
	Mirrors []string `json:"mirrors"`
	Format  []string `json:"format"`
} // }}}

type Config struct { // {{{
	Update string           `json:"update"`
	PM     []PackageManager `json:"list"`
} // }}}

var (
	Debug      = true
	Version    = "v0.1"
	ConfigFile = "./config.json"
)

func main() {
	if len(os.Args) == 1 {
		// mcn
		fmt.Println("功能:")
		fmt.Println("\t根据配置文件修改镜像源地址")
		fmt.Println()
		fmt.Println("版本:")
		fmt.Println("\t发行版: " + strconv.FormatBool(!Debug))
		fmt.Println("\t版本名: " + Version)
		fmt.Println("\t项目地址: " + "https://github.com/hbk01/MirrorCN")
		fmt.Println()
		fmt.Println("使用:")
		fmt.Println("\tmcn [--update] | [[--config=FILE] PackageManagerTag]")
		fmt.Println()
		fmt.Println("示例:")
		fmt.Println("\tmcn pacman apt # 修改pacman和apt的源，默认使用与mcn命令同目录的./config.json")
		fmt.Println("\tmcn --config=\"$HOME/MirrorCN.json\" pip # 使用指定的配置文件修改pip的源")
		fmt.Println("\tmcn --update # 下载最新的配置文件")
	} else if os.Args[1] == "--update" {
		// mcn --update
		fmt.Println("Update...")
		os.Exit(1)
	} else if strings.HasPrefix(os.Args[1], "--config-file=") {
		// mcn --config-file="./MirrorCN.json" apt
		ConfigFile = strings.Split(os.Args[1], "=")[1]
		fmt.Println(ConfigFile)
		pm := getAllPM(os.Args)
		fmt.Println("You Input The PM(s):", pm)
	} else {
		// mcn pacman apt
		pm := getAllPM(os.Args)
		fmt.Println("You Input The PM(s):", pm)
	}
	// config := Parse("./config.json")
}

func getAllPM(args []string) (pm []string) {
	for i, s := range args {
		if i == 0 {
			continue
		} else if strings.HasPrefix(s, "--") {
			continue
		} else {
			pm = append(pm, s)
		}
	}
	return pm
}
