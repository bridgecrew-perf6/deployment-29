package main

import (
	"oneclick/inspect"
	"oneclick/request"
	"oneclick/route"
	// "sync"
)

func main() {
	lts := inspect.NewCheck()
	if _, bl := lts.ShowConfig(); !bl {
		panic("未检测到配置文件，是否已绑定")
	}
	go request.NewReq().ShowStatus()
	lts.ShowLog()
	req := request.NewReq()
	req.SendTimingPing("ping")
	route.Route()
}
