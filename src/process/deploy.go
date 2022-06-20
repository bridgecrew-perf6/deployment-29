package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"oneclick/deploy"
	"oneclick/log"
	"oneclick/utlis"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Deploy(data BindingReturned) {
	config := BindingRequest{
		Domain:         domain,
		Key:            key,
		Networklaundry: networklaundry,
	}
	res, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Error(err.Error())
	}
	err = exec.Command("/bin/bash", "-c", "mkdir -p /etc/supervisor/conf.d /etc/wireguard").Run()
	if err != nil {
		log.Error(err.Error())
	}
	// 写入配置文件
	go func() {
		if networklaundry == 2 {
			err := ioutil.WriteFile("/etc/supervisor/conf.d/udp2raw.conf", []byte(utlis.Udp2rawConfig(strconv.Itoa(data.Data.Udp2rawEndpoint), data.Data.Udp2rawServerIp, data.Data.Udp2rawPasswd)), 0644)
			if err != nil {
				log.Error(err.Error())
			}
		}
		err := ioutil.WriteFile("/root/oneclickdeployment/config/config.json", res, 0644)
		if err != nil {
			log.Error(err.Error())
		}
		err = ioutil.WriteFile("/etc/wireguard/wg0.conf", []byte(data.Data.PeerInfo), 0644)
		if err != nil {
			log.Error(err.Error())
		}
		err = ioutil.WriteFile("/root/oneclickdeployment/frp/frpc.ini", []byte(data.Data.FrpcConfig), 0644)
		if err != nil {
			log.Error(err.Error())
		}
		err = ioutil.WriteFile("/etc/supervisor/supervisord.conf", []byte(utlis.SupervisorConfig), 0644)
		if err != nil {
			log.Error(err.Error())
		}
		err = ioutil.WriteFile("/etc/supervisor/conf.d/frpc.conf", []byte(utlis.FrpConfig), 0644)
		if err != nil {
			log.Error(err.Error())
		}
		err = ioutil.WriteFile("/etc/supervisor/conf.d/oneclick.conf", []byte(utlis.Oneclick), 0644)
		if err != nil {
			log.Error(err.Error())
		}
	}()
	ins := deploy.InitIns()
	ins.Source("install")
	//启动后台守护进程和wireguard
	//err = exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisord;supervisorctl reload;wg-quick down wg0;wg-quick up wg0").Run()
	err = exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisord;supervisorctl reload").Run()
	if err != nil {
		log.Error(err.Error())
	}
	Menu()
}

func Menu() {
	var inputStr string
	for {
		f := bufio.NewReader(os.Stdin)
		fmt.Println()
		fmt.Println("请选择你希望的操作:")
		fmt.Println("(1) 重新部署")
		fmt.Println("(2) 一键卸载")
		fmt.Println("(3) 停止使用")
		fmt.Println("(4) 重新使用")
		fmt.Println("(5) 解除密钥绑定")
		fmt.Println("(6) 退出脚本")
		fmt.Print("请选择 [1~6]:")
		inputStr, _ = f.ReadString('\n')
		inputStr = strings.TrimRight(inputStr, "\n")
		if inputStr == "1" {
			ins := deploy.InitIns()
			ins.Source("install")
		} else if inputStr == "2" {
			// 关闭后台守护进程和wireguard
			exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisorctl stop frpc udp2raw oneclick;wg-quick down wg0").Run()
			uns := deploy.InitIns()
			uns.Source("autoremove")
		} else if inputStr == "3" {
			// 关闭后台守护进程和wireguard
			err := exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisorctl stop frpc udp2raw oneclick;wg-quick down wg0").Run()
			if err != nil {
				fmt.Println(err.Error())
			}
			os.Exit(0)
		} else if inputStr == "4" {
			fmt.Println()
			err := exec.Command("/bin/bash", "-c", "cd /etc/supervisor;supervisorctl reload;wg-quick up wg0").Run()
			if err != nil {
				fmt.Println(err.Error())
			}
			fmt.Println("已为你重新启动")
		} else if inputStr == "5" {
			fmt.Println("============== 解除密钥绑定 ==============")
			os.Exit(0)
		} else if inputStr == "6" {
			os.Exit(0)
		} else {
			fmt.Println(utlis.Red("\n请输入正确的选项"))
			continue
		}
	}
}
