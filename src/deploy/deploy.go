package deploy

import (
	"fmt"
	"oneclick/utlis"
	"os/exec"

	"github.com/spf13/viper"
)

func (i ItsOUin) Source(Type string) {
	i.Type = Type
	fmt.Println()
	if i.Type == "install" {
		fmt.Println("开始安装...")
	}
	if i.Type == "autoremove" {
		fmt.Println("开始卸载...")
	}
	switch i.OS {
	case "ubuntu":
		i.OS = "apt-get"
		// i.aptBak()
		// i.aptUpdate()
		i.start()
	case "centos":
		i.OS = "yum"
		i.start()
	default:
		panic("暂时不支持centos、ubuntu以外的系统")
	}
}

func (i *ItsOUin) deploy() {
	for _, v := range i.name {
		st := v
		if v == "supervisor" {
			st = "supervisorctl"
		}
		if v == "wireguard" {
			st = "wg"
		}
		fmt.Printf("%s%s  %s\r", "...    \t", i.Type, v)
		if v == "wireguard" && i.OS == "yum" {
			str := i.OS + i.Type + " -y epel-release elrepo-release"
			str2 := i.OS + i.Type + " -y kmod-wireguard wireguard-tools iptables qrencode"
			fmt.Println(str)
			fmt.Println(str2)
			exec.Command("/bin/bash", "-c", str).Run()
			exec.Command("/bin/bash", "-c", str2).Run()
		}
		cmd := "command -v " + st
		// stdout, _ := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
		str := i.OS + " " + i.Type + " " + v + " -y"
		exec.Command("/bin/bash", "-c", str).Run()
		switch i.Type {
		case "install":
			stdout, _ := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
			if string(stdout) == "" {
				fmt.Printf("%s%s%s\n", utlis.Red("×      \t"), i.Type+"  ", v)
				continue
			}
			fmt.Printf("%s%s%s\n", utlis.LightGreen("✔      \t"), i.Type+"  ", v)
		case "autoremove":
			stdout, _ := exec.Command("/bin/bash", "-c", cmd).CombinedOutput()
			if string(stdout) != "" {
				fmt.Printf("%s%s%s\n", utlis.Red("×      \t"), i.Type+"  ", v)
				continue
			}
			fmt.Printf("%s%s%s\n", utlis.LightGreen("✔      \t"), i.Type+"  ", v)
		}

	}
}

func InitIns() ItsOUin {
	viper.SetConfigName("os-release")
	viper.SetConfigType("env")
	viper.AddConfigPath("/etc")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err.Error())
	}
	return ItsOUin{
		OS: viper.GetString("ID"),
	}
}
