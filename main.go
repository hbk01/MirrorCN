package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// PackageManager in the config.json is a list
type PackageManager struct {
	Name    string   `json:"name"`    // 镜像源名称
	Title   string   `json:"title"`   // 标题只会输出一次，在镜像源前输出。
	Backup  string   `json:"backup"`  // 备份到哪里，默认备份到当前目录，以.bak结尾
	First   []string `json:"first"`   // 先执行的操作
	End     []string `json:"end"`     // 替换结束后执行的操作
	Files   []string `json:"files"`   // 添加到哪个文件
	Mirrors []string `json:"mirrors"` // 镜像源地址
	Format  []string `json:"format"`  // 镜像格式
}

// Config the config.json struct
type Config struct {
	Update string           `json:"update"` // 配置文件更新地址
	Debug  bool             `json:"debug"`  // 是否开启Debug模式，该模式输出更详细的信息
	PM     []PackageManager `json:"list"`   // 源列表
}

var (
	// ConfigFile is default config file location.
	ConfigFile = "config.json"
	config     = Config{
		Update: "https://raw.githubusercontent.com/hbk01/MirrorCN/master/config.json",
		Debug:  true,
		PM:     []PackageManager{},
	}

	help = `
功能:
	根据配置文件修改镜像源地址

版本:
	发行版: 
	版本名: 
	项目地址: https://github.com/hbk01/MirrorCN

使用:
	mcn [--update] | [[--config=FILE] PackageManagerTag]

示例:
	mcn pacman apt # 修改pacman和apt的源，默认使用与mcn命令同目录的./config.json
	mcn --config="$HOME/MirrorCN.json\" pip # 使用指定的配置文件修改pip的源
	mcn --update # 下载最新的配置文件
`
)

func main() {
	// TODO 新增了一个 title 字段
	// TODO 新增了一个 backup 字段
	log(1, "Starting parse arguements", "")
	pms := parseArgs(os.Args)
	log(2, "Will change", strings.Join(pms, ", "))
	log(1, "Starting parse config from", ConfigFile)
	config = parseConfigJSON(ConfigFile)
	logd(2, "update", config.Update)
	logd(2, "debug", strconv.FormatBool(config.Debug))
	for _, pm := range config.PM {
		logd(2, "PM", pm.Name)
		logd(3, "First", strings.Join(pm.First, ", "))
		logd(3, "End", strings.Join(pm.End, ", "))
		logd(3, "Files", strings.Join(pm.Files, ", "))
		logd(3, "Mirrors", strings.Join(pm.Mirrors, ", "))
	}

	log(1, "Starting change files", "")
	changed := false
	for _, pm := range config.PM {
		for _, inputPmName := range pms {
			if pm.Name == inputPmName {
				changed = true
				log(1, "Starting first hooks for", pm.Name)
				runCommands(pm.First)
				changePm(pm)
				log(1, "Starting end hooks for", pm.Name)
				runCommands(pm.End)
			}
		}
	}
	if !changed {
		log(2, "Nothing to change.", "")
	}

	log(1, "End all task.", "")
}

func runCommands(commands []string) {
	for _, command := range commands {
		logd(2, "Run command", command)
		c := strings.Split(command, " ")[0]
		cmd := exec.Command(c, strings.Join(strings.Split(command, " ")[1:], " "))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Start()
		if err != nil {
			log(3, "Command returns", err.Error())
		}
		err = cmd.Wait()
		if err != nil {
			log(3, "Command returns", err.Error())
		}
		log(3, "Command output", out.String())
	}
}

// 解析JSON配置文件
func parseConfigJSON(path string) (config Config) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
	return config
}

// 解析命令参数
func parseArgs(args []string) (pm []string) {
	if len(os.Args) == 1 {
		// mcn
		fmt.Println(help)
		os.Exit(1)
	} else if os.Args[1] == "--update" {
		// mcn --update
		log(1, "Update...", "")
	} else if strings.HasPrefix(os.Args[1], "--config-file=") {
		// mcn --config-file="./MirrorCN.json" apt
		ConfigFile = strings.Split(os.Args[1], "=")[1]
		log(2, "Use config file", ConfigFile)
		pm = getAllPM(os.Args)
	} else {
		// mcn pacman apt
		pm = getAllPM(os.Args)
	}
	return pm
}

// 修改一个源
func changePm(pm PackageManager) {
	for _, file := range pm.Files {
		log(2, "Change '"+pm.Name+"' on ", file)
		f, err := os.OpenFile(file, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0655)
		if err != nil {
			log(2, "Open file error", file)
			return
		}
		// 备份文件
		var backupPath string
		if pm.Backup == "" {
			backupPath = file + ".bak"
		} else {
			backupPath = pm.Backup
		}
		srcBytes, err := ioutil.ReadFile(file)
		if err != nil {
			log(2, "Write backup file error", err.Error())
			return
		}

		err = ioutil.WriteFile(backupPath, srcBytes, 0655)
		if err != nil {
			log(2, "Write backup file error", err.Error())
			return
		}

		// 写入 title
		_, err = f.WriteString(pm.Title)
		if err != nil {
			log(2, "Write title error", pm.Title)
		}
		// 写入源
		for _, mirror := range pm.Mirrors {
			for _, format := range pm.Format {
				s := replaceFormat(format, map[string]string{
					"mirror": mirror,
					"file":   file,
					"name":   pm.Name,
				}) + "\n"
				_, err := f.WriteString(s)
				if err != nil {
					log(2, "Write "+s+" into "+file+" error", err.Error())
				}
			}
		}
	}
}

// replaceFormat 替换变量
// $mirror  将被替换为具体的镜像源地址
// $file 	将被替换为当前的文件
// $name     将被替换为源名称
// $date    将被替换为 yyyy-MM-dd 格式的日期，如 2020-11-21
// $time    将被替换为 HH:mm:ss 格式的时间，如 15:41:25
// $app    将被替换为本程序名称
// $url     将被替换为本程序开源地址
func replaceFormat(format string, arg map[string]string) string {
	format = strings.Replace(format, "$time", time.Now().Format("15:04:05"), -1)
	format = strings.Replace(format, "$date", time.Now().Format("2006-01-02"), -1)
	format = strings.Replace(format, "$mirror", arg["mirror"], -1)
	format = strings.Replace(format, "$file", arg["file"], -1)
	format = strings.Replace(format, "$name", arg["name"], -1)
	format = strings.Replace(format, "$app", "MirrorCN", -1)
	format = strings.Replace(format, "$url", "https://github.com/hbk01/MirrorCN", -1)
	return format
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

func logd(level int, k string, v string) {
	if config.Debug {
		switch level {
		case 1:
			fmt.Println(":: " + k + " " + v)
		case 2:
			fmt.Println("  - " + k + ": " + v)
		case 3:
			fmt.Println("    - " + k + ": " + v)
		}
	}
}

func log(level int, k string, v string) {
	switch level {
	case 1:
		fmt.Println(":: " + k + " " + v)
	case 2:
		fmt.Println("  - " + k + ": " + v)
	case 3:
		fmt.Println("    - " + k + ": " + v)
	}
}
