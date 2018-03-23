package config

import "github.com/jinzhu/configor"

var Config = struct {
	Server struct {
		Domain       string `default:"localhost"`
		Port         string `default:"0418"`
		DocumentRoot string `default:"/var/www/http/boohttp"`
	}
	DB struct {
		Driver   string `default:"mysql"`
		Host     string `default:"localhost"`
		Port     string `default:"3306"`
		UserName string `required:"true"`
		Password string
		Database string `default:"boo"`
	}
}{}

func InitConfig(configFile string) {
	configor.Load(&Config, configFile)
}
