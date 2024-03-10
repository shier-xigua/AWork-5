package zhttp

import (
	"AWork-5/zfunc"
	"AWork-5/zvar"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// 发起一个请求  抓取工单上工单告警内容
func Zhttp2(Form []map[string]string, Payload string) [][]interface{} {
	var orderTrue [][]interface{}
	var MatchToBusiness *bool
	for _, WorkOrderSummary := range Form {
		//2.1、获取url+pid = url2
		url2 := fmt.Sprintf("http://ccops-itsm.cmecloud.cn:8686/bit-bpc-process/process/viewInfo?instanceId=%s", WorkOrderSummary["processInstanceId"])

		//2.1.2、 定义一个请求
		req2, err := http.NewRequest("GET", url2, bytes.NewBufferString(Payload))
		if err != nil {
			log.Println("Error creating request2:", err)
			return nil
		}

		//2.1.3、设置请求头参数
		for key, value := range zvar.Headers {
			req2.Header.Set(key, value)
		}

		//2.1.4、发起请求
		client := &http.Client{}
		respons, err := client.Do(req2)
		if err != nil {
			log.Println("Error sending request2:", err)
			return nil

		} else {
			log.Println("Response2 Status:", respons.Status)
		}

		//2.1.5.1、获取工单页面的body！页面字数挺多的
		body, err := io.ReadAll(respons.Body)
		if err != nil {

			log.Println("Error reading response2 body:", err)
			return nil
		}

		//2.1.5.2、抓取告警描述内容
		var matchAlarms []string
		pattern := `"fault_desc":\s*\{([\s\S]+?)\s*"sourceValue"`
		matches := zfunc.MatchString(pattern, string(body))
		for _, matchAlarms = range matches {
			log.Printf("工单: %v, 工单内容: %s\n", WorkOrderSummary["processKey"], matchAlarms[0]) //排错使用  抓取告警描述内容
		}
		respons.Body.Close()

		//匹配business.txt内容和告警内容是否匹配，匹配则True，否则False
		b, s := MatchToBusinessFunc(MatchToBusiness, matchAlarms[0]) //b=bool s=business
		//log.Println("测试 b s m:", *b, s, WorkOrderSummary, "结束")

		//如果返回值为True则写入orderTrue
		if *b {
			var tmp []interface{}
			tmp = append(tmp, *b, s, WorkOrderSummary)
			orderTrue = append(orderTrue, tmp)
		}

	}
	//如果orderTrue为0则未匹配到工单
	if len(orderTrue) == 0 {
		log.Println("未匹配到工单")
	}
	//最总返回一个*bool,string,[][]interface{}的总内容这个内容交给zhttp3处理
	return orderTrue
}

// 2.2、business和告警内容是否匹配
func MatchToBusinessFunc(matchToBusiness *bool, matchAlarms string) (*bool, string) {
	var business string
	var exclude string
	var businessAndexclude []string
	var trueFlag = true
	var falseFlag = false
	var businessContains bool
	var excludeContains bool
	var by []byte
	by, err := os.ReadFile("business.txt")
	if err != nil {
		log.Fatal("读取business.txt文件错误", err)
	}
	aRow := strings.Split(string(by), "\n")
	for _, businessSplitExclude := range aRow { //循环读取文件中一行内容
		businessAndexclude = strings.Split(businessSplitExclude, "\\") //分割匹配字段和排除字段
		business = strings.TrimSpace(businessAndexclude[0])            //匹配字段去空

		//----------------判断是否切片长度为1--------------
		if len(businessAndexclude) == 1 && business != "" {
			businessContains = strings.Contains(matchAlarms, business)
			if businessContains {
				matchToBusiness = &trueFlag
				//log.Println("切片1：", *matchToBusiness, business)
				return matchToBusiness, business
			} else {
				matchToBusiness = &falseFlag
			}

			//--判断如果business.txt有空行不匹配
		} else if len(businessAndexclude) == 1 && business == "" {
			matchToBusiness = &falseFlag
			//----------------判断是否切片长度为1--------------

			//----------------判断是否切片长度为2--------------
		} else if len(businessAndexclude) >= 2 {
			//判断businessContains = 业务和告警是否匹配 匹配为=true
			businessContains = strings.Contains(matchAlarms, business) //匹配
			//fmt.Println("pipei")
			//循环取排除字段
			for i := 1; i <= len(businessAndexclude)-1; i++ {
				//fmt.Println(i) //调试排除字段循环取值
				//fmt.Println(businessAndexclude[i])

				//去除空行
				exclude = strings.TrimSpace(businessAndexclude[i])

				//判断excludeContains = 排除字段和告警是否匹配 匹配为true
				excludeContains = strings.Contains(matchAlarms, exclude) //匹配
				if excludeContains {
					//fmt.Println("跳转")
					goto JumpLabel //如果排除字段中任意一个为true跳转
				}

			}

		JumpLabel: //如果排除字段中任意一个为true跳转
			if businessContains {
				//fmt.Println("one true")
				//判断第二个excludeContains是否为true 如果true放弃接单
				if excludeContains {
					//如果为true则matchToBusiness=false
					matchToBusiness = &falseFlag
					log.Printf("判断匹配中'%v',因有排除字段'\\%v',放弃匹配该工单\n", business, exclude)
					//log.Println("切片2 排除字段：", *matchToBusiness, business)

					//2.2、判断第二个excludeContains是否为false 如果false直接接单
				} else if !excludeContains {
					matchToBusiness = &trueFlag
					//log.Println("切片2 匹配字段 未排除：", *matchToBusiness, business)
					return matchToBusiness, business
				}
			}
			// 2.1、判断第一个businessContains是否为true
		}

	}

	return matchToBusiness, "none"
}
