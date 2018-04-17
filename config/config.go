package config

import "github.com/jinzhu/configor"

var Config = struct {
	Server struct {
		Domain       string `default:"localhost"`
		Port         string `default:"8080"`
		DocumentRoot string `default:"/home/young/software/boo/boo-blogger"`
	}
	DB struct {
		Driver   string `default:"mysql"`
		Host     string `default:"127.0.0.1"`
		Port     string `default:"3306"`
		UserName string `default:"root"`
		Password string `default:"young159357789"`
		Database string `default:"boo"`
	}
}{}

func InitConfig(configFile string) {
	configor.Load(&Config, configFile)
}
