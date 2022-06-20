package process

type BindingRequest struct {
	Request
	Domain         string `json:"domain"`
	Key            string `json:"key"`
	Networklaundry int    `json:"networklaundry,omitempty"`
}

type Request struct {
	Time int    `json:"time,omitempty"`
	Sign string `json:"sign,omitempty"`
}

type BindingReturned struct {
	Returned
	Request
	Networklaundry int `json:"networklaundry,omitempty"`
	Binding        int `json:"binding,omitempty"`
	Data           struct {
		BindingReturnedConfig
	} `json:"data"`
}

type BindingReturnedConfig struct {
	PeerInfo        string `json:"peerinfo"`
	FrpcConfig      string `json:"frpcconfig"`
	Udp2rawEndpoint int    `json:"udp2rawendpoint,omitempty"`
	Udp2rawServerIp string `json:"udp2rawserver,omitempty"`
	Udp2rawPasswd   string `json:"udp2rawpasswd,omitempty"`
}

type Returned struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg,omitempty"`
}
