package com

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// ViperConf viper 配置初始化
// 根据环境标识加载运当前运行目录不同的配置文件
// Key 是必须的其余不强制
type ViperConf struct {
	Key          string
	EnvKey       string
	DefaultValue interface{}
}

// Options 应用选项
type Options []ViperConf

// Env 环境信息
type Env string

const (
	// EnvLocal 本地开发环境
	EnvLocal Env = "local"
	// EnvDev 开发环境
	EnvDev Env = "dev"
	// EnvTest 测试环境
	EnvTest Env = "test"
	// EnvPre 预发布环境
	EnvPre Env = "pre"
	// EnvPd 生产环境
	EnvPd Env = "pd"
)

// GetEnv 获取当前环境
// 默认情况或者环境设置出错的情况下返回本地开发环境
func GetEnv() Env {
	s := Env(os.Getenv("GO_MICRO_ENV"))
	if !s.IsSupport() {
		return EnvLocal
	}
	return s
}
func (e *Env) String() string {
	return string(*e)
}

// IsSupport 是否是有效的,受到支持的环境
func (e *Env) IsSupport() bool {
	v := *e
	return v == EnvLocal || v == EnvDev || v == EnvTest || v == EnvPre || v == EnvPd
}

// Init 执行初始化
func (o Options) Init() error {
	var err error
	env := GetEnv()
	// 尝试从环境变量读取配置文件路径
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "."
	}
	configType := os.Getenv("CONFIG_TYPE")
	if configType == "" {
		configType = "toml"
	}
	configPrefix := os.Getenv("CONFIG_PREFIX")
	if configPrefix == "" {
		configPrefix = "LH"
	}
	log.Infof("loading config of %s environment ", env)
	viper.SetConfigName(env.String())
	viper.AddConfigPath(path)
	viper.SetConfigType(configType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(configPrefix)
	viper.AllowEmptyEnv(true)
	for _, v := range o {
		if v.EnvKey != "" {
			err = viper.BindEnv(v.Key, v.EnvKey)
		} else {
			err = viper.BindEnv(v.Key)
		}
		if err != nil {
			return err
		}
		if v.DefaultValue != nil {
			viper.SetDefault(v.Key, v.DefaultValue)
		}
	}
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigParseError); ok {
			return err
		}
		log.Error("read config error", err)
	}
	return nil
}
