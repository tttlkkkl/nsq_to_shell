package service

import (
	"context"
	"fmt"
	"io"
	"nsq_exec_php/com"
	"os/exec"
	"strings"

	nnsq "github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// Service 服务定义
type Service struct{}

// Start 启动服务
func Start() error {
	// 添加处理函数,并且指定启动的 golang 处理协程数量为与服务端约定的可同时下发的消息数量
	nsqConsumer.Consumer.AddConcurrentHandlers(&Service{}, nsqConsumer.Config.MaxInFlight)
	// 通过 nsqlookup 启动消费客户端
	if err := nsqConsumer.Consumer.ConnectToNSQLookupds(nsqConsumer.LookupdAddr); err != nil {
		return err
	}
	return nil
}

// Stop 停止服务
func Stop() error {
	for _, v := range nsqConsumer.LookupdAddr {
		if err := nsqConsumer.Consumer.DisconnectFromNSQLookupd(v); err != nil {
			com.Log.Warning("移除 nsqlookup 出现一些问题:", v, err)
		}
	}
	// 停止服务
	nsqConsumer.Consumer.Stop()
	return nil
}

// HandleMessage 消息处理
// 默认情况下此方法返回错误消息将会被重新排队，否则消息被正常消费
func (s *Service) HandleMessage(message *nnsq.Message) error {
	com.Log.Infof("收到订阅消息，消息ID：%s,消息内容:%s", message.ID, string(message.Body))
	// 关闭消息自动处理
	message.DisableAutoResponse()
	// 设定执行超时时间
	ctx, cancel := context.WithTimeout(context.Background(), nsqConsumer.HandleTimeout)
	defer cancel()
	cmdDone := make(chan struct{}, 1)
	// 协程启动命令行执行
	go ExecCmd(ctx, cmdDone, string(message.Body))
	select {
	// 超时
	case <-ctx.Done():
		com.Log.Infof("消息 %s 执行超时", message.ID)
		// 未超过重排次数限制则重新排队,否则直接丢弃
		if nsqConsumer.Attempts >= int(message.Attempts) {
			message.Requeue(-1)
		} else {
			com.Log.Infof("消息 %s 超过 %d 次重排，直接丢弃", message.ID, message.Attempts)
			message.Finish()
		}
	// 执行完成
	case <-cmdDone:
		com.Log.Infof("消息 %s 已处理完毕", message.ID)
		message.Finish()
	}
	return nil
}

// ExecCmd 执行指令
func ExecCmd(ctx context.Context, cmdDone chan struct{}, messageBody string) {
	cmds := getCommand(messageBody)
	if len(cmds) == 0 {
		com.Log.Error("找不到可执行的指令，请检查配置")
		return
	}
	com.Log.Info("即将执行指令：", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		com.Log.Error("获取 shell stdout 输出流失败：", err)
		return
	}
	defer stdout.Close()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		com.Log.Error("获取 shell stderr 输出流失败：", err)
		return
	}
	defer stderr.Close()
	go asyncLog(stdout)
	go asyncLog(stderr)
	select {
	case <-ctx.Done():
		com.Log.Error("任务已执行超时，强制退出")
		return
	default:
		// 启动命令
		if err = cmd.Run(); err != nil {
			com.Log.Error("shell 指令执行出错：", err)
			return
		}
		// 执行完成，写入退出信号
		cmdDone <- struct{}{}
	}
}

// 实时输出执行日志
func asyncLog(reader io.ReadCloser) {
	cache := ""
	buf := make([]byte, 1024)
	for {
		num, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			com.Log.Debug("日志监听结束")
			return
		}
		if num > 0 {
			b := buf[:num]
			s := strings.Split(string(b), "\n")
			line := strings.Join(s[:len(s)-1], "\n")
			fmt.Printf("%s%s\n", cache, line)
			cache = s[len(s)-1]
		}
	}
}

// 获取组装执行指令
func getCommand(messageBody string) []string {
	var commands []string
	for _, v := range viper.GetStringSlice("nsq.command") {
		if strings.Count(v, "%s") == 1 {
			v = fmt.Sprintf(v, messageBody)
		}
		commands = append(commands, v)
	}
	return commands
}
