package zhttp

import (
	"AWork-5/zfunc"
	"AWork-5/zvar"
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
)

// 获取body
func Zhttp1(method, url, payload string) (string, error, int) {
	//第一阶段，定义请求
	request, err := http.NewRequest(method, url, bytes.NewBufferString(payload))
	if err != nil {
		log.Println("requests", err)
		return "", err, 0
	}
	//1.2设置请求头
	for key, value := range zvar.Headers {
		request.Header.Set(key, value)
	}
	//1.3 发起请求
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("client", err)
		*zfunc.ErrorCount++
		//fmt.Println(*zfunc.ErrorCount)
		return "", err, response.StatusCode
	} else {
		if response.StatusCode == 401 {
			log.Printf("Response1 Status: %v, 认证失败！请更换Token和Authorization后重启AWork接单工具！！！", response.Status)
		} else {
			log.Println("Response1 Status:", response.Status)
		}
		*zfunc.ErrorCount = 0
		//fmt.Println(*zfunc.ErrorCount)
	}

	//1.4读取工单系统工单摘要内容
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	client.CloseIdleConnections()
	//1.5 返回值，工单系统摘要body,这个内容交给
	return string(body), err, response.StatusCode
}

// 将body内容分解
func InfoMap(body string) []map[string]string {
	newbody := strings.Split(body, "{\"assignee\":")
	var form []map[string]string
	var entry map[string]string
	for i := 1; i < len(newbody); i++ {
		//fmt.Println(i)
		//fmt.Println(newbody[i])

		pattern := `"taskId":"([a-zA-Z0-9]+)","processInstanceId":"([a-zA-Z0-9]+)".*?"processTitle":"(.*?)","processKey":"(.*?)"`
		// 使用正则表达式查找匹配的内容
		//matches内容: [[taskID,processInstanceID, processTitle, processKey],[taskID,processInstanceID, processTitle, processKey]......]
		matches := zfunc.MatchString(pattern, newbody[i])
		for _, match := range matches {
			entry = map[string]string{
				"taskId":            match[1],
				"processInstanceId": match[2],
				"processTitle":      match[3],
				"processKey":        match[4],
			}
		}
		//将循环的每一次map数据存入到切片中
		if entry["processKey"][:3] == "001" || entry["processKey"][:3] == "042" {
			form = append(form, entry)
		}
	}
	return form
}
