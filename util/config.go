package util

import (
	"fmt"
	"github.com/spf13/viper"
	"path"
	"path/filepath"
	"runtime"
)

var (
	ProjectRootPath = getOnCurrentPath()
)

func getOnCurrentPath() string {
	_, filename, _, _ := runtime.Caller(0)

	absoluteParentPath, err := filepath.Abs(filepath.Join(path.Dir(filename), ".."))
	if err != nil {
		panic(err)
	}
	// 将路径转换为字符串，并添加斜杠
	return absoluteParentPath + string(filepath.Separator)
}

func CreateConfig(filename string) *viper.Viper {
	config := viper.New()
	configPath := filepath.Join(ProjectRootPath, "config")
	config.AddConfigPath(configPath)
	config.SetConfigName(filename)
	config.SetConfigType("yaml")
	configFilename := filepath.Join(configPath, filename+".yaml")

	if err := config.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(fmt.Errorf("找不到配置文件：%s", configFilename))
		} else {
			panic(fmt.Errorf("解析配置文件 %s 出错: %s", configFilename, err.Error()))
		}

	}
	return config
}
