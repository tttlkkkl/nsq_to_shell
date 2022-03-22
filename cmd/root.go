package cmd

import (
	"os"
	"os/signal"
	"syscall"

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
var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "nsq shell 执行服务",
	Long:  `启动一个消费服务，消费 nsq 中的队列消息，并传到 shell 中执行`,
	Run: func(cmd *cobra.Command, args []string) {
		// 启动订阅服务
		com.Log.Info("服务启动中...")
		service.Start()
		q := make(chan os.Signal)
		signal.Notify(q, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
		<-q
		com.Log.Info("收到退出信号，准备关闭服务...")
		if err := service.Stop(); err != nil {
			com.Log.Fatalln("服务异常关闭")
		}
	},
}

// Execute 执行
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		com.Log.Error(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./local.toml)")
}

// initConfig 初始化全局配置
func initConfig() {
	if err := com.InitOptions(cfgFile); err != nil {
		com.Log.Fatalln("配置初始化出错：", err)
	}
	if err := service.InitNsq(); err != nil {
		com.Log.Fatalln("nsq 消费者初始化失败：", err)
	}
}
