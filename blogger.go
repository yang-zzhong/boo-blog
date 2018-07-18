package blog

import (
	"boo-blog/cache"
	"boo-blog/http"
	"boo-blog/model"
	"errors"
	"github.com/go-ini/ini"
	. "github.com/yang-zzhong/go-model"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	BLOGGER_APP_NAME = "boo"
	BLOGGER_VERSION  = "0.0.1"
)

type Blogger struct {
	Version       string
	Name          string
	configFile    string
	config        *ini.File
	serverRunning bool
}

func NewBlogger(configFile string) *Blogger {
	blogger := new(Blogger)
	blogger.Version = BLOGGER_VERSION
	blogger.Name = BLOGGER_APP_NAME
	blogger.serverRunning = false
	blogger.SetConfig(configFile)

	return blogger
}

func (blogger *Blogger) StartHttp() error {
	if blogger.serverRunning {
		return errors.New("server is running")
	}
	return blogger.RestartHttp()
}

func (blogger *Blogger) RestartHttp() error {
	blogger.initModel()
	blogger.initCache()
	blogger.initHttp()
	model.OpenDB()
	if err := http.Start(); err != nil {
		return err
	}

	blogger.serverRunning = true
	return nil
}

func (blogger *Blogger) HandleQuitEvent() {
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				log.Println("退出blogger")
				blogger.Exit()
			}
		}
	}()
}

func (blogger *Blogger) Exit() {
	model.CloseDB()
	os.Exit(0)
}

func (blogger *Blogger) Config() *ini.File {
	return blogger.config
}

func (blogger *Blogger) SetConfig(configFile string) error {
	var err error
	if blogger.config, err = ini.Load(configFile); err != nil {
		return err
	}
	blogger.configFile = configFile
	if blogger.serverRunning {
		return blogger.RestartHttp()
	}

	return nil
}

func (blogger *Blogger) CreateTable() error {
	blogger.initModel()
	if err := model.OpenDB(); err != nil {
		return err
	} else {
		defer model.CloseDB()
		repos := []*Repo{
			model.NewBlog().Repo(),
			model.NewVote().Repo(),
			model.NewCate().Repo(),
			model.NewImage().Repo(),
			model.NewTheme().Repo(),
			model.NewTag().Repo(),
			model.NewUserImage().Repo(),
			model.NewUser().Repo(),
			model.NewComment().Repo(),
		}
		for _, repo := range repos {
			err := repo.CreateRepo()
			if err != nil {
				log.Print(err)
				continue
			}
			s := "ALTER TABLE " + repo.QuotedTableName() + " CONVERT TO CHARACTER SET utf8mb4"
			if _, err := Conn.Exec(s); err != nil {
				return err
			}
		}
		Conn.Exec("SET @@global.sql_mode='NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'")
	}
	return nil
}

func (blogger *Blogger) initModel() {
	config := blogger.config.Section("database")
	dc := &model.Config{
		Driver:   config.Key("driver").String(),
		Host:     config.Key("host").String(),
		Port:     config.Key("port").String(),
		Username: config.Key("username").String(),
		Password: config.Key("password").String(),
		Database: config.Key("database").String(),
		ImageDir: config.Key("image_dir").String(),
		BlogDir:  config.Key("blog_dir").String(),
	}

	model.InitDriver(dc)
}

func (blogger *Blogger) initCache() {
	config := blogger.config.Section("redis")
	cache.InitRedis(config.Key("addr").String(), config.Key("password").String())
}

func (blogger *Blogger) initHttp() {
	config := blogger.config.Section("server")
	http.InitHttp(
		config.Key("doc_root").String(),
		config.Key("port").String(),
		config.Key("session_secret").String())
}
