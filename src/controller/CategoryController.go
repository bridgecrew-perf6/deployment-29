package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"oneclick/log"
	"oneclick/model"
	"oneclick/request"
	"oneclick/response"
	"oneclick/utlis"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (c2 CategoryController) UpdateTP(c *gin.Context) {
	var RequestBody model.PingRequest
	var sign string
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	m, _ := json.Marshal(&RequestBody.Data)
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	SignMap["data"] = string(m)
	SignMap["interval"] = strconv.Itoa(RequestBody.Interval)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	PostBody := model.PingReturned{
		Ret:  1,
		Msg:  "请求成功",
		Sign: sign,
		Data: []model.PingResult{},
	}
	config := model.PingRequest{
		Interval: RequestBody.Interval,
	}
	for i := 0; i < len(RequestBody.Data); i++ {
		cmd := fmt.Sprintf("ping -c 1 %s", RequestBody.Data[i])
		result, _, _ := utlis.RunCommand("-c", cmd)
		// 截取ping结果的字符串
		// 以第一个换行开始 第一个---结束
		// 去掉截取字符串中的空格
		hh := utlis.GetBetweenStr(result, "\n", "---")
		str := strings.Replace(hh, "\n", "", -1)
		if str == "" {
			str = fmt.Sprintf("echo reply from %s (%s) : timeout", RequestBody.Data[i], RequestBody.Data[i])
		}
		PostBody.Data = append(
			PostBody.Data,
			model.PingResult{
				Address: RequestBody.Data[i],
				Result:  str,
			},
		)
		config.Data = append(config.Data, RequestBody.Data[i])
	}
	res, _ := json.Marshal(config)
	//res, _ = json.MarshalIndent(config, "", "  ")
	ioutil.WriteFile("/root/oneclickdeployment/config/pingconfig.json", res, 0644)
	res, _ = json.Marshal(PostBody.Data)
	SignMap["ret"] = strconv.Itoa(PostBody.Ret)
	SignMap["msg"] = PostBody.Msg
	SignMap["data"] = string(res)
	// 生成签名
	PostBody.Sign = utlis.SignatureVisitor(SignMap)
	c.JSON(http.StatusOK, PostBody)
	// 修改ping间隔
	request.NewReq().SendTimingPing("ping")
}

func (c2 CategoryController) ExecPing(c *gin.Context) {
	var RequestBody model.PingRequest
	var sign string
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	m, _ := json.Marshal(&RequestBody.Data)
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	SignMap["data"] = string(m)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	PostBody := model.PingReturned{
		Ret:  1,
		Msg:  "请求成功",
		Sign: sign,
		Data: []model.PingResult{},
	}
	log.Info("ExecPing " + string(m))
	for i := 0; i < len(RequestBody.Data); i++ {
		cmd := fmt.Sprintf("ping -c 1 '%s'", RequestBody.Data[i])
		result, _, _ := utlis.RunCommand("-c", cmd)
		// 截取ping结果的字符串
		//以第一个换行开始 第一个---结束
		// 去掉截取字符串中的空格
		hh := utlis.GetBetweenStr(result, "\n", "---")
		str := strings.Replace(hh, "\n", "", -1)
		if str == "" {
			str = fmt.Sprintf("echo reply from %s (%s) : timeout", RequestBody.Data[i], RequestBody.Data[i])
		}
		PostBody.Data = append(PostBody.Data, model.PingResult{
			Address: RequestBody.Data[i],
			Result:  str,
		})
	}
	res, _ := json.Marshal(PostBody.Data)
	SignMap["ret"] = strconv.Itoa(PostBody.Ret)
	SignMap["msg"] = PostBody.Msg
	SignMap["data"] = string(res)
	// 生成签名
	PostBody.Sign = utlis.SignatureVisitor(SignMap)
	c.JSON(http.StatusOK, PostBody)
}

func (c2 CategoryController) ShowConfig(c *gin.Context) {
	var RequestBody model.Request
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		log.Error("签名错误")
		response.Error(c, "签名错误")
		return
	}
	cmd := exec.Command("/bin/bash", "-c", "cat /etc/wireguard/wg0.conf")
	log.Info("cat /etc/wireguard/wg0.conf")
	output, _ := cmd.CombinedOutput()
	if len(output) == 0 {
		output = []byte("null")
	}
	response.Success(c, "请求成功", gin.H{
		"wireguard_config": string(output),
	})
}

func (c2 CategoryController) ShowStatus(c *gin.Context) {
	var RequestBody model.Request
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	output, _ := exec.Command("/bin/bash", "-c", "wg").Output()
	if len(output) == 0 {
		output = []byte("null")
	}
	response.Success(c, "请求成功", gin.H{
		"wireguard_config": string(output),
	})
}

func (c2 CategoryController) UpdateWgConfig(c *gin.Context) {
	var RequestBody model.Wireguard
	var sign string
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	SignMap["private_key"] = RequestBody.PrivateKey
	SignMap["address"] = RequestBody.Address
	SignMap["dns"] = RequestBody.DNS
	SignMap["mtu"] = RequestBody.MTU
	SignMap["public_key"] = RequestBody.PublicKey
	SignMap["allowed_i_ps"] = RequestBody.AllowedIPs
	SignMap["endpoint"] = RequestBody.Endpoint
	SignMap["persistent_keepalive"] = RequestBody.PersistentKeepalive
	SignMap["peer_info"] = RequestBody.PeerInfo
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	WgConfig := model.Wireguard{
		PrivateKey:          RequestBody.PrivateKey,
		Address:             RequestBody.Address,
		DNS:                 RequestBody.DNS,
		MTU:                 RequestBody.MTU,
		PublicKey:           RequestBody.PublicKey,
		AllowedIPs:          RequestBody.AllowedIPs,
		Endpoint:            RequestBody.Endpoint,
		PersistentKeepalive: RequestBody.PersistentKeepalive,
		PeerInfo:            RequestBody.PeerInfo,
	}
	res, _ := json.MarshalIndent(WgConfig, "", "  ")
	f, err := os.OpenFile("/etc/wireguard/wg0.conf", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	defer f.Close()
	_, err = f.Write([]byte(RequestBody.PeerInfo))
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	err = exec.Command("/bin/bash", "-c", "wg-quick down wg0;wg-quick up wg0").Run()
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	log.Info("wg-quick down wg0;wg-quick up wg0")
	output, _ := exec.Command("/bin/bash", "-c", "wg").Output()
	fmt.Println(string(output))
	PostBody := model.T{
		Ret:  1,
		Msg:  "请求成功",
		Sign: sign,
		Data: model.WireguardStatus{
			WireguardRunning: string(output),
		},
	}
	res, _ = json.Marshal(PostBody.Data)
	SignMap["ret"] = strconv.Itoa(PostBody.Ret)
	SignMap["msg"] = PostBody.Msg
	SignMap["RequestBody"] = string(res)
	// 生成签名
	PostBody.Sign = utlis.SignatureVisitor(SignMap)
	c.JSON(200, PostBody)
}

func (c2 CategoryController) Update(c *gin.Context) {
	var RequestBody model.Request
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	// 更新脚本
	// cmd := fmt.Sprintf("cd /root/oneclickdeployment/bin && git fetch && git checkout origin/master %s", data.File)
	cmd := fmt.Sprintf("cd /root/oneclickdeployment/bin && git fetch && git checkout origin/test .")
	err = exec.Command("/bin/bash", "-c", cmd).Run()
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	log.Info(cmd)
	response.Success(c, "请求成功", nil)
	//time.Sleep(time.Second *5)
	// 重启后台守护进程
	go func() {
		cmd := fmt.Sprintf("cd /etc/supervisor;supervisorctl restart time_ping update_wireguard udp2raw frpc")
		err := exec.Command("-c", cmd).Run()
		if err != nil {
			log.Error(err.Error())
		}
		log.Info(cmd)
	}()
}

func (c2 CategoryController) ShowLog(c *gin.Context) {
	var RequestBody model.GetLogger
	err := c.ShouldBindJSON(&RequestBody)
	if err != nil {
		log.Error(err.Error())
		response.Error(c, err.Error())
		return
	}
	SignMap := make(map[string]string)
	SignMap["time"] = strconv.Itoa(RequestBody.Time)
	SignMap["logtype"] = RequestBody.LogType
	SignMap["numrows"] = strconv.Itoa(RequestBody.NumRows)
	if !utlis.KuratowskiConstraint(SignMap, RequestBody.Sign) {
		response.Error(c, "签名错误")
		return
	}
	cmd := fmt.Sprintf("tail -n %d /root/oneclickdeployment/log/%s.log", RequestBody.NumRows, RequestBody.LogType)
	result, _, _ := utlis.RunCommand("-c", cmd)
	log.Info(cmd)
	response.Success(c, "请求成功", gin.H{
		"LogData": result,
	})
}

func NewCategoryController() ICategoryController {
	return CategoryController{}
}
