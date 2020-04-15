package utils

import (
	"errors"
	"flag"
	"fmt"
	"github.com/larspensjo/config"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	path, _    = GetCurrentPath()
	configFile = flag.String("configfile", path+"/../config/conf.ini", "General configuration file")
	options    = make(map[string]string)
)

// 获取配置文件内容
func GetOption(name string, category string) (string, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	cfg, err := config.ReadDefault(*configFile)
	if err != nil {
		return "", err
	}
	//set config file std End

	//Initialized topic from the configuration
	if cfg.HasSection(category) {
		section, err := cfg.SectionOptions(category)
		if err == nil {
			for _, v := range section {
				item, err := cfg.String(category, v)
				if err == nil {
					options[v] = item
				}
			}
		}
	}

	value, ok := options[name]
	if ok {
		return value, nil
	}
	return "", errors.New(fmt.Sprintf("config item '%s' not found. %+v", name))
}

//获取当前程序所在目录
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}
