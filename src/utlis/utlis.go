package utlis

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"text/template"
	"time"
)

// BaseUtils 一键部署脚本
func BaseUtils(CertificateAuto string, KernelVersionYn string, IsAliyunMirror string) string {
	var sb strings.Builder
	sb.Write([]byte(oneClickDeployment))
	var envMap = make(map[string]interface{})
	envMap["CERTIFICATE_AUTO"] = CertificateAuto
	envMap["KERNEL_VERSION_YN"] = KernelVersionYn
	envMap["IS_ALIYUN_MIRROR"] = IsAliyunMirror
	envMap["MAIN"] = "`uname -r | awk -F . '{print $1}'`"
	return FromTemplateContent(sb.String(), envMap)
}

// Udp2rawConfig upd2raw supervisorctl配置文件
func Udp2rawConfig(Udp2rawEndpoint string, Udp2rawServerIp string, Udp2rawPasswd string) string {
	var sb strings.Builder
	sb.Write([]byte(Udp2rawConfigs))
	var envMap = make(map[string]interface{})
	envMap["Udp2rawEndpoint"] = Udp2rawEndpoint
	envMap["Udp2rawServerIp"] = Udp2rawServerIp
	envMap["Udp2rawPasswd"] = Udp2rawPasswd
	return FromTemplateContent(sb.String(), envMap)
}

// FromTemplateContent 替换字符串里的关键字
func FromTemplateContent(templateContent string, envMap map[string]interface{}) string {
	tmpl, err := template.New("text").Parse(templateContent)
	defer func() {
		if r := recover(); r != nil {
			//logger.Error("Template parse failed:", err)
		}
	}()
	if err != nil {
		panic(1)
	}
	var buffer bytes.Buffer
	_ = tmpl.Execute(&buffer, envMap)
	return string(buffer.Bytes())
}

// GetBetweenStr 截取字符串
func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

// RunCommand 执行linux命令
func RunCommand(arg ...string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/sh", arg...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func Command(cmd string) (string, error) {
	c := exec.Command("/bin/bash", "-c", cmd)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return "", err
	}
	var wg sync.WaitGroup
	var res string
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			fmt.Print(readString)
			res = fmt.Sprintf("%s", readString)
		}
	}()
	err = c.Start()
	wg.Wait()
	return res, err
}

// WriteFile 保存文件
func WriteFile(path string, content string) {
	var fileByte = []byte(content)
	err := ioutil.WriteFile(path, fileByte, 0644)
	if err != nil {
		panic(err)
	}
}

// ReadFile 读取文件
func ReadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return string(content)
}

// // ReadLines 按行读取文件
// func ReadLines(filename string) ([]string, error) {
// 	f, err := os.Open(filename)
// 	if err != nil {
// 		return []string{""}, err
// 	}
// 	defer func(f *os.File) {
// 		_ = f.Close()
// 	}(f)
// 	var ret []string
// 	r := bufio.NewReader(f)
// 	for {
// 		line, readErr := r.ReadString('\n')
// 		if readErr != nil {
// 			break
// 		}
// 		ret = append(ret, strings.Trim(line, "\n"))
// 	}
// 	return ret, nil
// }

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// RandString 生成随机字符串
func RandString(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return string(bs)
}

// RandNum 生成随机数
func RandNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// RemoveIpPort 去除IP中的端口
func RemoveIpPort(ip string) string {
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// RandomString 随机字符串
func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

// StringMd5 MD5
func StringMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// ReadLines 按行读取文件
func ReadLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	var ret []string
	r := bufio.NewReader(f)
	for {
		line, readErr := r.ReadString('\n')
		if readErr != nil {
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret, nil
}

// StringsContains 数组是否包含
func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// InArray 元素是否存在数组中
func InArray(item string, items []string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// CheckFileIsExist 判断文件是否存在  存在返回 true 不存在返回false
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}
