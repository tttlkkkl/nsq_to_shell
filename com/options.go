package com

// InitOptions 初始化配置
func InitOptions() error {
	vp := Options{
		{
			Key:    "nsq.lookupAddress",
			EnvKey: "LH_LOOKUP_ADDRESS",
		},
		{
			Key:          "nsq.clientName",
			EnvKey:       "LH_LOOKUP_ADDRESS",
			DefaultValue: "nsq_exec_php",
		},
		{
			Key:          "nsq.msgTimeout",
			EnvKey:       "LH_LOOKUP_ADDRESS",
			DefaultValue: "10s",
		},
		{
			Key:          "nsq.maxInFlight",
			EnvKey:       "LH_LOOKUP_ADDRESS",
			DefaultValue: 3,
		},
	}
	return vp.Init()
}
