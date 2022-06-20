package response

import (
	"encoding/json"
	"net/http"
	"oneclick/utlis"
	"strconv"

	"github.com/gin-gonic/gin"
)

var sign string

func Response(c *gin.Context, httpStatus int, ret int, msg string, data gin.H) {
	m, _ := json.Marshal(&data)
	SignMap := make(map[string]string)
	SignMap["ret"] = strconv.Itoa(ret)
	SignMap["msg"] = msg
	SignMap["data"] = string(m)
	sign = utlis.SignatureVisitor(SignMap)
	c.JSON(httpStatus, gin.H{
		"ret":  ret,
		"msg":  msg,
		"sign": sign,
		"data": data,
	})
}

func Responses(c *gin.Context, httpStatus int, ret int, msg, data string) {
	SignMap := make(map[string]string)
	SignMap["ret"] = strconv.Itoa(ret)
	SignMap["msg"] = msg
	SignMap["data"] = data
	sign = utlis.SignatureVisitor(SignMap)
	c.JSON(httpStatus, gin.H{
		"ret":  ret,
		"msg":  msg,
		"sign": sign,
		"data": data,
	})
}

func Success(c *gin.Context, msg string, data gin.H) {
	Response(c, http.StatusOK, 1, msg, data)
}

func Successtest(c *gin.Context, msg, data string) {
	Responses(c, http.StatusOK, 1, msg, data)
}

func Fail(c *gin.Context, msg string, data gin.H) {
	Response(c, http.StatusOK, 0, msg, data)
}

func Error(c *gin.Context, msg string) {
	//返回签名错误
	c.JSON(http.StatusBadRequest, ReturnErrors(0, msg))
}

func ReturnErrors(ret int, msg string) (res Errors) {
	res = Errors{
		Ret:  ret,
		Msg:  msg,
		Sign: sign,
	}
	SignMap := make(map[string]string)
	SignMap["ret"] = strconv.Itoa(res.Ret)
	SignMap["msg"] = res.Msg
	// 加签
	sign = utlis.SignatureVisitor(SignMap)
	res.Sign = sign
	return res
}

type Errors struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Sign string `json:"sign"`
}
