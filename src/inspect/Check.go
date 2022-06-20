package inspect

import (
	"fmt"
	"oneclick/log"
	"oneclick/utlis"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Inspect interface {
	showUser()
	showKernel()
	ShowLog()
	ShowConfig() (ck *Check, bl bool)
}

type Check struct {
	Networklaundry int
	Domain         string
	Key            string
}

// ShowLog 守护进程
// 定时查看日志大小
// 每小时检测一次
func (c Check) ShowLog() {
	go func() {
		t := time.NewTicker(60 * time.Minute)
		for {
			select {
			case <-t.C:
				opBytes, err := exec.Command("/bin/bash", "-c", "du -sh /root/oneclickdeployment/log/").Output()
				fmt.Println(string(opBytes))
				if err != nil {
					log.Error(err.Error())
				}
				if string(opBytes[2:3]) == "M" {
					num, _ := strconv.Atoi(string(opBytes[:2]))
					if num > 20 {
						err = exec.Command("/bin/bash/", "-c", "rm -rf /root/oneclickdeployment/log/*").Run()
						if err != nil {
							log.Error(err.Error())
						}
					}
				}
			}
		}
	}()
}

func (c Check) showUser() {
	u, _ := user.Current()
	if u.Username != "root" {
		fmt.Println(utlis.Red("您需要以root用户身份运行此脚本"))
		os.Exit(1)
	}
}

func (c Check) showKernel() {
	var (
		KernelVersionYn = "n"
		IsAliyunMirror  = "n"
	)
	stdout, err := exec.Command("/bin/bash", "-c", "uname -r | awk -F . '{print $1}'").Output()
	if err != nil {
		panic(err.Error())
	}
	num, _ := strconv.Atoi(string(stdout[0]))
	if num < 5 {
		fmt.Println("当前内核版本", num)
		fmt.Println("检测到内核版本太低,需要升级到 5 以上")
		fmt.Print(utlis.Red("升级内核会重启你的主机,是否升级 [y/N] :"))
		fmt.Scanln(&KernelVersionYn)
		if KernelVersionYn == "n" {
			os.Exit(1)
		}
		fmt.Print("使用阿里源 [y/N] :")
		fmt.Scanln(&IsAliyunMirror)
		utlis.WriteFile("/root/oneclickdeployment/bin/deployment.sh", utlis.BaseUtils("", KernelVersionYn, IsAliyunMirror))
		utlis.Command("sudo bash /root/oneclickdeployment/bin/deployment.sh")
	}
}

func (c *Check) ShowConfig() (ck *Check, bl bool) {
	if utlis.Exists("/root/oneclickdeployment/config/config.json") {
		viper.SetConfigName("config")
		viper.SetConfigType("json")
		viper.AddConfigPath("/root/oneclickdeployment/config")
		err := viper.ReadInConfig()
		if err != nil {
			panic(err.Error())
		}
		c.Networklaundry = viper.GetInt("networklaundry")
		c.Domain = viper.GetString("domain")
		c.Key = viper.GetString("key")
		return c, true
	}
	return nil, false
}

func NewCheck() (c Check) {
	c.showUser()
	c.showKernel()
	return c
}
