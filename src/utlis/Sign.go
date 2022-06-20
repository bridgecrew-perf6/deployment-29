package utlis

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Md5 计算md5的值
func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func SignatureVisitor(data map[string]string) string {
	data["privatekey"] = "erwswrywrywrywreybsefhdg"
	SignatureData := []byte(JoinStringsInASCII(data, "&", false, false))
	sum := fmt.Sprintf("%x", md5.Sum(SignatureData))
	return sum
}

// KuratowskiConstraint 解签
func KuratowskiConstraint(data map[string]string, Sign string) bool {
	data["privatekey"] = "erwswrywrywrywreybsefhdg"
	SignatureData := []byte(JoinStringsInASCII(data, "&", false, false))
	sum := fmt.Sprintf("%x", md5.Sum(SignatureData))
	//res := Md5(string(datas))
	if sum == Sign {
		return true
	}
	return false
}

// JoinStringsInASCII 按照规则，参数名ASCII码从小到大排序后拼接
//data 待拼接的数据
//sep 连接符
//onlyValues 是否只包含参数值，true则不包含参数名，否则参数名和参数值均有
//includeEmpty 是否包含空值，true则包含空值，否则不包含，注意此参数不影响参数名的存在
//exceptKeys 被排除的参数名，不参与排序及拼接
func JoinStringsInASCII(data map[string]string, sep string, onlyValues, includeEmpty bool, exceptKeys ...string) string {
	var list []string
	var keyList []string
	m := make(map[string]int)
	if len(exceptKeys) > 0 {
		for _, except := range exceptKeys {
			m[except] = 1
		}
	}
	for k := range data {
		if _, ok := m[k]; ok {
			continue
		}
		value := data[k]
		if !includeEmpty && value == "" {
			continue
		}
		if onlyValues {
			keyList = append(keyList, k)
		} else {
			list = append(list, fmt.Sprintf("%s=%s", k, value))
		}
	}
	if onlyValues {
		sort.Strings(keyList)
		for _, v := range keyList {
			list = append(list, data[v])
		}
	} else {
		sort.Strings(list)
	}
	return strings.Join(list, sep)
}
