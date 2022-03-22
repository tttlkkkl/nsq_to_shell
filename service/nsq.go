package service

import (
	"errors"
	"time"

	nnsq "github.com/nsqio/go-nsq"
	"github.com/spf13/viper"
)

// 全局服务变量
var nsqConsumer NsqConsumer

// NsqConsumer nsq 消费者定义
type NsqConsumer struct {
	Topic         string
	Channel       string
	LookupdAddr   []string
	Consumer      *nnsq.Consumer
	Config        *nnsq.Config
	HandleTimeout time.Duration
	Attempts      int
}

// 初始化 nsq 订阅服务
func InitNsq() (err error) {
	c := NsqConsumer{
		Config: nnsq.NewConfig(),
	}
	// 初始化配置
	if err = c.setNsqConfig(); err != nil {
		return err
	}
	// 实例化消费者
	nc, err := nnsq.NewConsumer(c.Topic, c.Channel, c.Config)
	if err != nil {
		return err
	}
	c.Consumer = nc
	// 全局服务变量
	nsqConsumer = c
	return nil
}

// setNsqConfig 根据配置文件设置 nsq 否则取默认值
func (n *NsqConsumer) setNsqConfig() error {
	if v := viper.GetStringSlice("nsq.lookupAddress"); len(v) > 0 {
		n.LookupdAddr = v
	} else {
		return errors.New("必须指定 lookupAddress！")
	}
	if v := viper.GetString("nsq.topic"); v != "" {
		n.Topic = v
	} else {
		return errors.New("必须指定 topic")
	}
	if v := viper.GetString("nsq.channel"); v != "" {
		n.Channel = v
	} else {
		return errors.New("必须指定 channel")
	}
	if v := viper.GetString("nsq.clientName"); v != "" {
		n.Config.ClientID = v
	}
	if v := viper.GetDuration("nsq.msgTimeout"); v >= 0 {
		n.Config.MsgTimeout = v
	}
	if v := viper.GetInt("nsq.maxInFlight"); v >= 0 {
		n.Config.MaxInFlight = v
	}
	if v := viper.GetDuration("nsq.handleTimeout"); v >= 0 {
		if v > n.Config.MsgTimeout {
			return errors.New("handleTimeout 大于 msgTimeout，设置无效")
		}
		n.HandleTimeout = v
	}
	if v := viper.GetInt("nsq.attempts"); v >= 0 {
		n.Attempts = v
	}
	return nil
}
