package cmd

import (
	"nsq_exec_php/com"
	"nsq_exec_php/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

// sCmd represents the s command
var sCmd = &cobra.Command{
	Use:   "s",
	Short: "启动订阅服务",
	Long:  `s 启动订阅服务`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// 启动订阅服务
		if args[0] == "s" {
			com.Log.Info("服务启动中...")
			service.Start()
			q := make(chan os.Signal)
			signal.Notify(q, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
			<-q
			com.Log.Info("收到退出信号，准备关闭服务...")
			if err := service.Stop(); err != nil {
				com.Log.Fatalln("服务异常关闭")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(sCmd)
}
