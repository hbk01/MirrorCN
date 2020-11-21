package main

import (
	"fmt"
	"os"
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

func main() {
	if len(os.Args) == 1 {
		// mcn
		fmt.Println("功能：根据配置文件修改镜像源地址")
		fmt.Println("USAGE:")
		fmt.Println("\tmcn [ [ --update ] | [ [--config=FILE] PackageManagerTag ] ]")
		fmt.Println("")
		fmt.Println("EXAMPLE:")
		fmt.Println("\tmcn go # 修改go的源，默认使用与mcn命令同目录的配置文件")
		fmt.Println("\tmcn --config=\"$HOME/MirrorCN.json\" pacman # 使用指定的配置文件修改pacman的源")
		fmt.Println("\tmcn --update # 下载最新的配置文件")
	} else if os.Args[1] == "--update" {
		// mcn --update
		fmt.Println("Update...")
		os.Exit(1)
	} else if strings.HasPrefix(os.Args[1], "--config-file=") {
		// mcn --config-file="hahaha" apt
		file := strings.Split(os.Args[1], "=")[1]
		fmt.Println(file)
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
