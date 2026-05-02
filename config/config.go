package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Aliyun   AliyunConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type DatabaseConfig struct {
	Driver       string
	Host         string
	Port         int
	Username     string
	Password     string
	Database     string
	Charset      string
	MaxIdleConns int
	MaxOpenConns int
}

type JWTConfig struct {
	Secret string
	Expire int
}

type AliyunConfig struct {
	AccessKeyID     string
	AccessKeySecret string
	Endpoint        string
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

var AppConfig *Config

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	AppConfig = &Config{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	log.Println("配置文件加载成功")
}
