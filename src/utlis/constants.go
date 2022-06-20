package utlis

const pemKey = string(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDFRNPZ+2A/XUdIm9+VvVQQ2xWELa/TprZrAboi4R7bDDom2nUc
eIcrmI8fV+9iHfx0cVNg9YSANj7VbUOSHVzDVy5R7oo3T7ME6dS/t24I/gGvnWcz
nKHCsCek+T50siLlwfHfId2hQAKraU2gWeF2lrATprl3mC12wFdEKnz2CwIDAQAB
AoGAbXXEi+b1QBO1My/yv3bfx76ZUM+9CZcvD29U5ne+FFPTjK2ZYCPs9R7hA8Za
eTokVERxvJJfZHk1Il5PqSsLxgsHmMq/yXNvyG2Bya43sO4Vzlzx/+SvpSBuwrr9
gg150M6vJJeKsjIYFDX3/GK+TVSIuPxoL21Ho7nlBZKxNzECQQDmknu8uTCztAPF
7cpZi7nMGtg+UqdoWU2YFceHYPhqvCbhpWctfmToHJEpRpClpE8dXF1a51HmqdNE
F+N/X6J5AkEA2wYftxTomO6jVW6ntE7pCvsbBJpfyvGU0lIkyYeCRRREHR2GkkAJ
lp0pZGoi1wgNjgYA+tWlrCurhysAVrPbowJBAIamBJyxiT9oYMu1kfW5I0eOZbn/
isPlYurtzRfCCVBLkGk1rotixIrII/12uAIDcjAzQFFVxP5vLnEVgkVgFAECQQDV
962UFgEFJly6YVfEdjKEX7uNS6K5iDhzH3yAxLkm8x13tBh7V8QGN5LwXh+bImrb
jFH4ui8Xe7IeYov6J8sxAkB+LgFFbzU2UJ6FiGN/DT9Nia2dBMI28Z2GW686ZTpA
Av4fcEDNAL1Tu0gHf1VvmGb5+xNIdLA6MJF2alCMkPNk
-----END RSA PRIVATE KEY-----
`)

const Config = string(`{
  "domain": "",
  "key": "",
  "networklaundry": 0
}`)

const PingConfig = string(`{
    "interval":360,
    "data":[
        "www.hitosea.com",
        "www.json.cn",
        "qq.com",
        "baidu.com"
    ]
}`)

const oneClickDeployment = string(`#!/bin/bash

function checkOS() {
        # Check OS version
        if [[ -e /etc/debian_version ]]; then
                source /etc/os-release
                OS="${ID}" # debian or ubuntu
        elif [[ -e /etc/centos-release ]]; then
                source /etc/os-release
                OS="${ID}"
                echo $OS
        else
                echo "看起来您没有在Debian、Ubuntu、CentOS Linux系统上运行此安装程序"
                exit 1
        fi
}

function Kernelupgrade(){
        # KERNEL_VERSION=$(cut -d '.' -f1 <<<$(uname -r))
        main={{.MAIN}}

        if [[ "$main" -lt 5 ]]; then
                yn={{.KERNEL_VERSION_YN}}
                case $yn in
                        [yY])
                                if [[ ${OS} == 'ubuntu' ]]; then
                                        sudo apt-get upgrade linux-image-generic -y
                                        echo
                                        echo "重新启动使内核升级完毕"
                                        echo
                                elif [[ ${OS} == 'centos' ]]; then
                                        yum -y update
                                        rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
                                        rpm -Uvh http://www.elrepo.org/elrepo-release-7.0-3.el7.elrepo.noarch.rpm
                                        yum --enablerepo=elrepo-kernel install kernel-ml -y
                                        grub2-set-default 0
                                        sed -i '3c GRUB_DEFAULT=0' /etc/default/grub
                                        grub2-set-default 'CentOS Linux (5.15.11-1.el7.elrepo.x86_64) 7 (Core)'
                                        echo "2"
                                fi
                                local yn="y"
                                while [[ "$yn" != "y" && "$yn" != "n" ]]; do
                                        read -p "是否现在重启？否则将终止脚本 [y/n] :" yn
                                done
                                if [[ "$yn" == "y" ]]; then
                                        echo "重启中"
                                        sleep 1s
                #                       reboot
                                        exit 1
                                else
                                        echo "即将终止脚本"
                                        sleep 2s
                                        exit 1
                                fi
                                ;;
                        [nN])
                                echo "no"
                                exit 1
                esac
        fi
}






function install_software(){
        if [[ "$OS" == "ubuntu" ]]; then
                install_ubuntu_software
        else
                install_centos_software
        fi
        res=$(command -v supervisorctl)
        if [[ "$res" == "" ]]; then
                echo "supervisor 安装失败, 请联系管理员"
                exit 1
        else
                echo "supervisor 安装完成"
        fi
        echo "启动 supervisor"
        if [[ "$OS" == "centos" ]]; then
                systemctl start supervisord # centos 是 supervisord
        else
                systemctl start supervisor
        fi
        echo
        echo "安装完毕"
        echo
        exit 0
}

function install_ubuntu_software(){
        echo "安装 Ubuntu $DISTRIBUTION_VERSION 环境"
    local apt_url="http://archive.ubuntu.com/ubuntu/"
        IS_ALIYUN_MIRROR={{.IS_ALIYUN_MIRROR}}
    if [[ "${IS_ALIYUN_MIRROR}" == "y" ]]; then
        echo "使用阿里 apt 源"
        apt_url="http://mirrors.aliyun.com/ubuntu/"
    fi
    local timestamp=$(date +%s)
    mv /etc/apt/sources.list "/etc/apt/sources.list.$timestamp.back" 2>/dev/null
    # name=$(lsb_release -c | cut -f2)
    local name=$(lsb_release -cs)
    echo "deb $apt_url $name main restricted universe multiverse
deb-src $apt_url $name main restricted universe multiverse
deb $apt_url $name-security main restricted universe multiverse
deb-src $apt_url $name-security main restricted universe multiverse
deb $apt_url $name-updates main restricted universe multiverse
deb-src $apt_url $name-updates main restricted universe multiverse
deb $apt_url $name-proposed main restricted universe multiverse
deb-src $apt_url $name-proposed main restricted universe multiverse
deb $apt_url $name-backports main restricted universe multiverse
deb-src $apt_url $name-backports main restricted universe multiverse
" >>/etc/apt/sources.list
    apt-get update -y
    if [[ ! $(command -v curl) ]]; then
        apt-get install -y curl
    fi
    apt-get install -y apt-transport-https ca-certificates gnupg-agent software-properties-common
    if [[ ! $(command -v git) ]]; then
        apt-get install -y git
    fi
    if [[ ! $(command -v wget) ]]; then
        apt-get install -y wget
    fi
    if [[ ! $(command -v supervisorctl) ]]; then
        apt-get install -y supervisor
    fi
    if [[ ! $(command -v jq) ]]; then
        apt-get install -y jq
    fi
    if [[ ! $(command -v wg) ]]; then
        apt-get install -y wireguard iptables resolvconf qrencode
    fi
}

function install_centos_software() {
    IS_ALIYUN_MIRROR={{.IS_ALIYUN_MIRROR}}
    if [[ ! $(command -v git) ]]; then
        yum install -y git
    fi
    if [[ ! $(command -v wget) ]]; then
        yum install -y wget
    fi
    if [[ ! $(command -v supervisorctl) ]]; then
        # 如果是阿里源
        if [[ "${IS_ALIYUN_MIRROR}" == "n" ]]; then
            yum install -y supervisor
        else
            echo "使用阿里源 yum 源"
            rm -f supervisor-3.4.0-1.el7.noarch.rpm 2>/dev/null
            #wget http://www.rpmfind.net/linux/epel/7/ppc64le/Packages/s/supervisor-3.4.0-1.el7.noarch.rpm
            wget -O /etc/yum.repos.d/epel-7.repo  http://mirrors.aliyun.com/repo/epel-7.repo
            #yum localinstall -y supervisor-3.4.0-1.el7.noarch.rpm
            yum install supervisor -y
        fi
    fi
    if [[ ! $(command -v jq) ]]; then
        yum install -y jq
    fi
    if [[ ! $(command -v wg) ]]; then
            yum -y install epel-release elrepo-release
                if [[ ${VERSION_ID} -eq 7 ]]; then
                        yum -y install yum-plugin-elrepo
                fi
			yum -y install kmod-wireguard wireguard-tools iptables qrencode
    fi
}

main() {
        checkOS
        Kernelupgrade
        install_software
}

main
`)

const Oneclick = string(`[program:oneclick]
directory=/root/oneclickdeployment/bin ; 执行前要不要先cd到目录去
command=/root/oneclickdeployment/bin/./oneclick      ; 被监控的进程路径
autostart = true        ; 随着supervisord的启动而启动
autorestart = true              ; 自动重启
startsecs = 10        ; 启动 10 秒后没有异常退出，就当作已经正常启动了
startretries = 10               ; 启动失败时的最多重试次数
exitcodes=0                   ; 正常退出代码（是说退出代码是这个时就不再重启了吗？待确定）
stopsignal=KILL               ; 用来杀死进程的信号
redirect_stderr = true  ; 把 stderr 重定向到 stdout，默认 false
stdout_logfile=/root/oneclickdeployment/log/oneclick.log
stdout_logfile_maxbytes=10MB
user=root
loglevel=info
[supervisord] ; 必须配置
[supervisorctl] ; 必须配置
`)

const SupervisorConfig = string(`; supervisor config file

[unix_http_server]
file=/var/run/supervisor.sock   ; (the path to the socket file)
chmod=0700                       ; sockef file mode (default 0700)

[supervisord]
logfile=/var/log/supervisor/supervisord.log ; (main log file;default /supervisord.log)
pidfile=/var/run/supervisord.pid ; (supervisord pidfile;default supervisord.pid)
childlogdir=/var/log/supervisor            ; ('AUTO' child log dir, default )

; the below section must remain in the config file for RPC
; (supervisorctl/web interface) to work, additional interfaces may be
; added by defining them in separate rpcinterface: sections
[rpcinterface:supervisor]
supervisor.rpcinterface_factory = supervisor.rpcinterface:make_main_rpcinterface

[supervisorctl]
serverurl=unix:///var/run/supervisor.sock ; use a unix:// URL  for a unix socket

; The [include] section can just contain the files setting.  This
; setting can list multiple files (separated by whitespace or
; newlines).  It can also contain wildcards.  The filenames are
; interpreted as relative to this file.  Included files *cannot*
; include files themselves.

[include]
files = /etc/supervisor/conf.d/*.conf
`)

const Udp2rawConfigs = string(`[program:udp2raw]
directory=/root/oneclickdeployment/udp2raw ; 执行前要不要先cd到目录去
command=/root/oneclickdeployment/udp2raw/./udp2raw_amd64 -c -l 0.0.0.0:{{.Udp2rawEndpoint}} -r {{.Udp2rawServerIp}} -a -k {{.Udp2rawPasswd}} --raw-mode faketcp --cipher-mode xor  ; 被监控的进程启动命令
autostart = true        ; 随着supervisord的启动而启动
autorestart = true              ; 自动重启
startsecs = 10        ; 启动 10 秒后没有异常退出，就当作已经正常启动了
startretries = 10               ; 启动失败时的最多重试次数
exitcodes=0                   ; 正常退出代码（是说退出代码是这个时就不再重启了吗？待确定）
stopsignal=KILL               ; 用来杀死进程的信号
redirect_stderr = true  ; 把 stderr 重定向到 stdout，默认 false
stdout_logfile=/root/oneclickdeployment/log/udp2raw.log
stdout_logfile_maxbytes=10MB
user=root
loglevel=info
[supervisord] ; 必须配置
[supervisorctl] ; 必须配置
`)

const FrpConfig = string(`[program:frpc]
directory=/root/oneclickdeployment/frp ; 执行前要不要先cd到目录去
command=/root/oneclickdeployment/frp/frpc -c /root/oneclickdeployment/frp/frpc.ini      ; 被监控的进程路径
autostart = true        ; 随着supervisord的启动而启动
autorestart = true              ; 自动重启
startsecs = 10        ; 启动 10 秒后没有异常退出，就当作已经正常启动了
startretries = 10               ; 启动失败时的最多重试次数
exitcodes=0                   ; 正常退出代码（是说退出代码是这个时就不再重启了吗？待确定）
stopsignal=KILL               ; 用来杀死进程的信号
redirect_stderr = true  ; 把 stderr 重定向到 stdout，默认 false
stdout_logfile=/root/oneclickdeployment/log/frpc.log
stdout_logfile_maxbytes=10MB
user=root
loglevel=info
[supervisord] ; 必须配置
[supervisorctl] ; 必须配置
`)
