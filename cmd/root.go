package cmd

import (
	"os"

	"nsq_exec_php/com"
	"nsq_exec_php/service"

	"github.com/spf13/cobra"
)

const (
	webAppName  = "dn.mc.web.auth"
	grpcAppName = "dn.mc.srv.auth"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{}

// Execute 执行
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		com.Log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig 初始化全局配置
func initConfig() {
	if err := com.InitOptions(); err != nil {
		com.Log.Fatalln("配置初始化出错：", err)
	}
	if err := service.InitNsq(); err != nil {
		com.Log.Fatalln("nsq 消费者初始化失败：", err)
	}
}
