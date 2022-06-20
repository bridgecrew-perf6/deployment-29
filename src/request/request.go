package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"oneclick/log"
	"oneclick/utlis"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

func NewReq() (r Req) {
	return r
}

func (r *Req) showPc() {
	if utlis.Exists("/root/oneclickdeployment/config/pingconfig.json") {
		viper.SetConfigName("pingconfig")
		viper.SetConfigType("json")
		viper.AddConfigPath("/root/oneclickdeployment/config")
		err := viper.ReadInConfig()
		if err != nil {
			log.Error(err.Error())
		}
		r.Interval = viper.GetInt("interval")
		r.Data = viper.GetStringSlice("data")
		return
	}
	err := ioutil.WriteFile("/root/oneclickdeployment/config/pingconfig.json", []byte(utlis.PingConfig), 0644)
	if err != nil {
		log.Error(err.Error())
	}
}

func (r Req) SendTimingPing(name string) {
	rand := utlis.RandString(6)
	daemonMap[name] = rand
	go func() {
		r.showPc()
		// 设置定时任务时间
		t := time.NewTicker(time.Second * time.Duration(r.Interval))
		var sign string
		for {
			select {
			case <-t.C:
				if daemonMap[name] != rand {
					return
				}
				tt := int(time.Now().Unix())
				pinata := PingReturned{
					Request: Request{
						Time: tt,
						Sign: sign,
					},
					Data: []PingData{},
				}
				for i := 0; i < len(r.Data); i++ {
					cmd := fmt.Sprintf("ping -c 1 '%s'", r.Data[i])
					result, _, _ := utlis.RunCommand("-c", cmd)
					// 截取ping结果的字符串
					// 以第一个换行开始 第一个---结束
					// 去掉截取字符串中的空格
					hh := utlis.GetBetweenStr(result, "\n", "---")
					str := strings.Replace(hh, "\n", "", -1)
					if str == "" {
						str = fmt.Sprintf("echo reply from %s (%s) : timeout", r.Data[i], r.Data[i])
					}
					pinata.Data = append(
						pinata.Data, PingData{
							Address: r.Data[i],
							Result:  str,
						})
				}
				res, _ := json.Marshal(pinata.Data)
				SignMap := make(map[string]string)
				SignMap["time"] = strconv.Itoa(tt)
				SignMap["data"] = string(res)
				// 生成签名
				pinata.Sign = utlis.SignatureVisitor(SignMap)
				res, err := json.MarshalIndent(pinata, "", "  ")
				if err != nil {
					log.Error(err.Error())
				}
				client := &http.Client{}
				resp, err := client.Post("http://192.168.200.124:8082/ArchivedPingsInfo/post", "application/json", bytes.NewBuffer(res))
				if err != nil {
					log.Error(err.Error())
					return
				}
				bodyText, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Error(err.Error())
				}
				fmt.Printf("%s\n", bodyText)
			}
		}
	}()
}

func (r Req) sendStatus(envMap map[string]int) {
	var sign string
	t := int(time.Now().Unix())
	data := StatusByte{
		Request: Request{
			Time: t,
			Sign: sign,
		},
		Frpc: Status{
			Status: envMap["frpc"],
		},
		Udp2Raw: Status{
			Status: envMap["udp2raw"],
		},
		Oneclick: Status{
			Status: envMap["oneclick"],
		},
		Wireguard: Status{
			Status: envMap["wireguard"],
		},
	}
	_, k := envMap["udp2raw"]
	if !k {
		data.Udp2Raw.Status = 0
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(t)
	SignMap["frpc"] = strconv.Itoa(envMap["frpc"])
	SignMap["udp2raw"] = strconv.Itoa(envMap["udp2raw"])
	SignMap["oneclick"] = strconv.Itoa(envMap["oneclick"])
	SignMap["wireguard"] = strconv.Itoa(envMap["wireguard"])
	// 生成签名
	data.Sign = utlis.SignatureVisitor(SignMap)
	res, _ := json.Marshal(data)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://192.168.200.124:8082/CheckStartup/post", strings.NewReader(string(res)))
	if err != nil {
		log.Error(err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return
	}
}

// ShowStatus 查看启动后各服务的状态并上传
func (r Req) ShowStatus() {
	time.Sleep(time.Second * 15)
	// fmt.Println(time.Now().Format("2006-01-02 15:04:05"))
	data := make(map[string]int)
	networklaundry := viper.GetInt("networklaundry")
	process := []string{0: "frpc", 1: "oneclick", 2: "udp2raw"}
	slice := process[:]
	if networklaundry == 1 {
		slice = process[0:2]
	}
	for _, v := range slice {
		cmd := fmt.Sprintf("cd /etc/supervisor;supervisorctl status %s | awk '{print $2}'", v)
		result, _, _ := utlis.RunCommand("-c", cmd)
		if strings.Contains(result, "RUNNING") {
			data[v] = 1
		} else {
			data[v] = 0
			cmd = fmt.Sprintf("supervisorctl status %s", v)
			res, _, _ := utlis.RunCommand("-c", cmd)
			log.Error(fmt.Sprintf("%s启动失败\n%s", v, res))
		}
	}
	wg, _, _ := utlis.RunCommand("-c", "wg")
	if len(wg) == 0 {
		res := "wireguard启动失败"
		data["wireguard"] = 0
		log.Error(res)
	} else {
		res := "wireguard启动成功"
		data["wireguard"] = 1
		log.Info(res)
	}
	r.sendStatus(data)
	//dataType, _ := json.Marshal(data)
	//if !(strings.Contains(string(dataType), "RUNNING") && len(wg) == 0) {
	//	r.SendStatus(data)
	//}
}
