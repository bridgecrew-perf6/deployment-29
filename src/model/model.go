package model

type PstructList struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Sign string `json:"sign"`
}

type Request struct {
	Time int    `json:"time,omitempty"`
	Sign string `json:"sign,omitempty"`
}

type Wireguard struct {
	Request
	PrivateKey          string `json:"private_key,omitempty" binding:"required"`
	Address             string `json:"address,omitempty" binding:"required"`
	DNS                 string `json:"dns,omitempty" binding:"required"`
	MTU                 string `json:"mtu,omitempty" binding:"required"`
	PublicKey           string `json:"public_key,omitempty" binding:"required"`
	AllowedIPs          string `json:"allowed_i_ps,omitempty" binding:"required"`
	Endpoint            string `json:"endpoint,omitempty" binding:"required"`
	PersistentKeepalive string `json:"persistent_keepalive,omitempty" binding:"required"`
	PeerInfo            string `json:"peer_info,omitempty" binding:"required"`
}

type PingRequest struct {
	Request
	Interval int      `json:"interval"`
	Data     []string `json:"data"`
}

type PingReturned struct {
	Ret      int          `json:"ret,omitempty"`
	Msg      string       `json:"msg,omitempty"`
	Sign     string       `json:"sign,omitempty"`
	Interval int          `json:"interval,omitempty"`
	Data     []PingResult `json:"data"`
}

type PingResult struct {
	Address string `json:"address"`
	Result  string `json:"result,omitempty"`
}

type GetLogger struct {
	Time    int    `json:"time"`
	Sign    string `json:"sign"`
	LogType string `json:"logtype"`
	NumRows int    `json:"numrows"`
}

type T struct {
	Ret  int             `json:"ret"`
	Msg  string          `json:"msg"`
	Sign string          `json:"sign"`
	Data WireguardStatus `json:"data"`
}

type WireguardStatus struct {
	WireguardConfig  string `json:"wireguard_config,omitempty"`
	WireguardRunning string `json:"wireguard_running,omitempty"`
}
