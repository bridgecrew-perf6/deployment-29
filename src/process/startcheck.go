package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"oneclick/inspect"
	"oneclick/log"
	"oneclick/utlis"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetInput() {
	for {
		f := bufio.NewReader(os.Stdin)
		var input string
		fmt.Println()
		fmt.Println("请选择网络协议:")
		fmt.Println("(1) tcp ")
		fmt.Println("(2) tcp+udp ")
		fmt.Print("(默认: 1) :")
		input, _ = f.ReadString('\n')
		input = strings.TrimRight(input, "\n")
		if len(input) == 0 {
			networklaundry = 1
			break
		}
		if find := strings.Contains(string(input), " "); find {
			fmt.Println(utlis.Red("\n网络协议不能有空格"))
			continue
		}
		if input == "1" {
			networklaundry = 1
			break
		} else if input == "2" {
			networklaundry = 2
			break
		} else {
			fmt.Println(utlis.Red("\n输入无效,请重新输入"))
			continue
		}
	}
	for {
		f := bufio.NewReader(os.Stdin)
		fmt.Println("请输入你要访问的域名或ip与端口:")
		domain, _ = f.ReadString('\n')
		domain = strings.TrimRight(domain, "\n")
		if len(domain) == 0 {
			fmt.Println(utlis.Red("域名或ip不能为空"))
			continue
		}
		if find := strings.Contains(string(domain), " "); find {
			fmt.Println(utlis.Red("\n域名或ip不能有空格"))
			domain = ""
			continue
		}
		if len(domain) != 0 {
			break
		}
	}
	for {
		f := bufio.NewReader(os.Stdin)
		fmt.Println("请输入密钥:")
		key, _ = f.ReadString('\n')
		key = strings.TrimRight(key, "\n")
		if len(key) == 0 {
			fmt.Println(utlis.Red("密钥不能为空"))
			continue
		}
		if find := strings.Contains(key, " "); find {
			fmt.Println(utlis.Red("\n密钥不能有空格"))
			key = ""
			continue
		}
		if len(key) != 0 {
			break
		}

	}
	Register()
}

func Verification(ins *inspect.Check) {
	t := int(time.Now().Unix())
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(t)
	SignMap["key"] = ins.Key
	// 生成签名
	sign := utlis.SignatureVisitor(SignMap)
	data := BindingRequest{
		Request: Request{
			Time: t,
			Sign: sign,
		},
		Key: ins.Key,
	}
	res, _ := json.Marshal(data)
	//res, _ := json.MarshalIndent(data, "", "  ")
	//client := &http.Client{Timeout: 60 * t.Second}
	resp, err := http.Post("http://192.168.200.124:8082/CheckKey/post", "application/json", strings.NewReader(string(res)))
	if err == nil {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err.Error())
		}
		var data BindingReturned
		json.Unmarshal(body, &data)
		switch data.Ret {
		case 0:
			log.Error(data.Msg)
			//关闭后台守护进程和wireguard删除配置文件
			//pkg.RunCommand("-c", "cd /etc/supervisor;supervisorctl stop background frpc udp2raw;wg-quick down wg0;rm -rf /root/oneclickdeployment/config/config.json")
			//等待一秒
			//t.Sleep(t.Second *1)
			os.Exit(1)
		case 1:
			Menu()
		case 2:
			log.Error(data.Msg)
			os.Exit(1)
		default:
			fmt.Println(utlis.Red("请求失败，请重试"))
			os.Exit(1)
		}
	}
	log.Error(err.Error())
	fmt.Println(utlis.Red("域名或ip与端口错误。请重新输入"))
	// utlis.RunCommand("-c", "rm -rf /root/oneclickdeployment/config/config.json")
	domain = ""
	GetInput()
}

var (
	domain, key    string
	networklaundry int
)
