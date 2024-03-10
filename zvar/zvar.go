package zvar

import (
	"AWork-5/zfunc"
	"fmt"
)

var token, auth, Lineto = zfunc.GetTokenAuthor()
var URL = "http://ccops-itsm.cmecloud.cn:8686/bit-bpc-process/task/list-claim"
var Headers = map[string]string{
	"Accept":           "application/json, text/plain, */*",
	"Accept-Encoding":  "gzip, deflate",
	"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
	"Bpc-AccessSource": "web",
	"Bpc-Version":      "20220509",
	"Cache-Control":    "no-cache",
	"Content-Length":   "26",
	"Content-Type":     "application/json;charset=UTF-8",
	"Cookie":           fmt.Sprintf("blueking_language=zh-cn; ultra_msa_sso=yes; Authorization=%s; bk_token=%s", auth, token),
	"Host":             "ccops-itsm.cmecloud.cn:8686",
	"Origin":           "http://ccops-itsm.cmecloud.cn:8686",
	"Proxy-Connection": "keep-alive",
	"Referer":          "http://ccops-itsm.cmecloud.cn:8686/bpc/",
	"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0",
	"X-Requested-With": "XMLHttpRequest",
}
var Payload = `{"pageNo":"1","pageSize":"10"}`
