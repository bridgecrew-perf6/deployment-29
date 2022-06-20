package deploy

import (
	"fmt"
	"oneclick/log"
	"os"
	"os/exec"
	"strings"
	"time"
)

func (i ItsOUin) start() {
	i.Curl()
	i.Wget()
	i.Supervisor()
	i.Jq()
	i.Wireguard()
	i.deploy()
}

func (i ItsOUin) aptBak() {
	cmd := fmt.Sprintf("mv /etc/apt/sources.list \"/etc/apt/sources.list.%v.back\" 2>/dev/null", time.Now().Unix())
	err := exec.Command("/bin/bash", "-c", cmd).Run()
	if err != nil {
		log.Error("apt源备份失败")
	}
	stdout, err := exec.Command("/bin/bash", "-c", "lsb_release -cs").Output()
	if err != nil {
		log.Error(err.Error())
	}
	str := strings.Replace(aptStr, "$apt_url", aptUrl, -1)
	str = strings.Replace(str, "$i.name", strings.Replace(string(stdout), "\n", "", -1), -1)
	f, err := os.OpenFile("/etc/apt/sources.list", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error(err.Error())
	}
	defer f.Close()
	_, err = f.Write([]byte(str))
	if err != nil {
		log.Error(err.Error())
	}
}

func (i ItsOUin) aptUpdate() {
	err := exec.Command("/bin/bash", "-c", "apt-get update -y").Run()
	if err != nil {
		log.Error(err.Error())
	}
}

func (i *ItsOUin) Curl() {
	i.name = append(i.name, "curl")
}

func (i *ItsOUin) Wget() {
	i.name = append(i.name, "wget")
}

func (i *ItsOUin) Supervisor() {
	i.name = append(i.name, "supervisor")
}

func (i *ItsOUin) Jq() {
	i.name = append(i.name, "jq")
}

func (i *ItsOUin) Wireguard() {
	i.name = append(i.name, "wireguard")
}
