package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB      DBConfig      `mapstructure:"db"`
	Server  ServerConfig  `mapstructure:"server"`
	Cluster ClusterConfig `mapstructure:"cluster"`
}

type DBConfig struct {
	Postgres PostgresConfig `mapstructure:"postgres"`
	Minio    MinioConfig    `mapstructure:"minio"`
}

type ClusterConfig struct {
	DefaultRegion   string        `mapstructure:"default_region"`
	DefaultImage    string        `mapstructure:"default_image"`
	CreationTimeout time.Duration `mapstructure:"creation_timeout"`
}

type ServerConfig struct {
	Port           string        `mapstructure:"port"`
	RequestTimeout time.Duration `mapstructure:"request_timeout"`
	Prometheus     bool          `mapstructure:"prometheus"`
	Pprof          bool          `mapstrusture:"pprof"`
	BodyLimit      int           `mapstructure:"body_limit"`
}

type PostgresConfig struct {
	Host   string `mapstructure:"host"`
	User   string `mapstructure:"user"`
	Pwd    string `mapstructure:"pwd"`
	DBName string `mapstructure:"db_name"`
	Port   string `mapstructure:"port"`
}

type MinioConfig struct {
	Url       string `mapstructure:"url"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Port      string `mapstructure:"port"`
	Path      string `mapstructure:"path"`
}

func ViperCfg(path string) Config {
	viper.AutomaticEnv()
	viper.SetConfigName("seismographd")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
