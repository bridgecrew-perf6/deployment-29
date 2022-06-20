package request

type Request struct {
	Time int    `json:"time,omitempty"`
	Sign string `json:"sign,omitempty"`
}

type StatusByte struct {
	Request
	Frpc      Status `json:"frpc"`
	Udp2Raw   Status `json:"udp2raw,omitempty"`
	Oneclick  Status `json:"oneclick"`
	Wireguard Status `json:"wireguard"`
}

type Status struct {
	Status int `json:"status"`
}

type Returned struct {
	Ret int    `json:"ret"`
	Msg string `json:"msg,omitempty"`
}

type ClientBinding struct {
	Returned
	Request
	Domain         string `json:"domain"`
	Key            string `json:"key"`
	Networklaundry int    `json:"networklaundry"`
}

type PingRequest struct {
	Request
	Interval int      `json:"interval"`
	Data     []string `json:"data"`
}

type PingReturned struct {
	Ret int    `json:"ret,omitempty"`
	Msg string `json:"msg,omitempty"`
	Request
	Interval int        `json:"interval,omitempty"`
	Data     []PingData `json:"data"`
}

type PingData struct {
	Address string `json:"address"`
	Result  string `json:"result,omitempty"`
}

var daemonMap = make(map[string]string)

type show interface {
	ShowStatus()
	sendStatus(envMap map[string]int)
	SendTimingPing(name string)
	showPc()
}

type Req struct {
	Interval int      `json:"interval"`
	Data     []string `json:"data"`
	show
}
