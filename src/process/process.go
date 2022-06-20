package process

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"oneclick/log"
	"oneclick/utlis"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func Register() {
	t := int(time.Now().Unix())
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(t)
	SignMap["domain"] = domain
	SignMap["key"] = key
	SignMap["networklaundry"] = strconv.Itoa(networklaundry)
	// 生成签名
	sign := utlis.SignatureVisitor(SignMap)
	datarq := BindingRequest{
		Request: Request{
			Time: t,
			Sign: sign,
		},
		Domain:         domain,
		Key:            key,
		Networklaundry: networklaundry,
	}
	res, err := json.Marshal(datarq)
	if err != nil {
		log.Error(err.Error())
	}
	//res, err = json.MarshalIndent(datarq, "", "  ")
	// 设置超时时间
	//client := &http.Client{Timeout: 60 * t.Second}
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(domain, "application/json", strings.NewReader(string(res)))
	// resp, err := http.Post(domain, "application/json", strings.NewReader(string(res)))
	if err != nil {
		log.Error(err.Error())
		fmt.Println(err.Error())
		fmt.Println(utlis.Red("域名或ip与端口错误。请重新输入"))
		// utlis.RunCommand("-c", "rm -rf /root/oneclickdeployment/config/config.json")
		networklaundry = 0
		domain = ""
		key = ""
		GetInput()
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		fmt.Println("1")
	}
	//初始化结构体
	var data BindingReturned
	//将body的值绑定到data,反序列化
	json.Unmarshal(body, &data)
	switch {
	case string(body) == "404 page not found":
		fmt.Println(utlis.Red("域名或ip与端口错误。请重新输入"))
		networklaundry = 0
		domain = ""
		key = ""
		GetInput()
	case data.Ret == 0:
		fmt.Println(utlis.Red(data.Msg))
		// 关闭后台守护进程和wireguard删除配置文件
		err := exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisorctl stop background frpc udp2raw;wg-quick down wg0;rm -rf /root/oneclickdeployment/config/config.json").Run()
		if err != nil {
			fmt.Println(err.Error())
		}
		os.Exit(1)
	case data.Ret == 2:
		fmt.Println(utlis.Red("此密钥已绑定设备"))
		os.Exit(1)
	case data.Ret == 1:
		Deploy(data)
	}
}
