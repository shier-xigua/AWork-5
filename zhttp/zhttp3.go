package zhttp

import (
	"AWork-5/zvar"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

func Zhttp3(lineto int, sliceInterface [][]interface{}) {
	for i := 0; i < len(sliceInterface); i++ {
		//log.Println("slicce长度", len(sliceInterface))
		if OrderSummary, ok := sliceInterface[i][2].(map[string]string); ok {
			//URL
			url3 := "http://ccops-itsm.cmecloud.cn:8686/bit-bpc-process/task/batch/claim"

			//3.1、设置请求体参数
			var matchTaskId = OrderSummary["taskId"] //获取taskid
			var payload = fmt.Sprintf(`{"taskIds":"%s"}`, matchTaskId)
			//3.2、定义一个请求
			req3, err := http.NewRequest("POST", url3, bytes.NewBufferString(payload))
			if err != nil {
				fmt.Println("Error creating request3:", err)
			}
			//3.3、设置请求头
			for key, value := range zvar.Headers {
				req3.Header.Set(key, value)
			}
			//3.4、发起一个请求
			client3 := &http.Client{}
			resp3, err := client3.Do(req3)
			if err != nil {
				log.Println("Error sending request3:", err)
				return
			} else {
				log.Println("Response3 Status:", resp3.Status)
				if business, ok := sliceInterface[i][1].(string); ok {
					log.Printf(" 接单成功,匹配字段为：%v，工单号：%v 工单标题：%v ,taskid：%v, ProcessInstanceId:%v\n", business, OrderSummary["processKey"], OrderSummary["processTitle"], OrderSummary["taskId"], OrderSummary["processInstanceId"])
					var text string
					if OrderSummary["processKey"][:3] == "001" {
						text = fmt.Sprintf("AWork匹配到字段为：%v,故障工单尾号：%v", business, OrderSummary["processKey"][13:])
						// 使用 PowerShell 播放语音

					} else if OrderSummary["processKey"][:3] == "042" {
						text = fmt.Sprintf("AWork匹配到字段为：%v,事件工单尾号：%v", business, OrderSummary["processKey"][13:])
					}
					for i := 0; i < lineto; i++ {
						// 使用 PowerShell 播放语音
						cmd := exec.Command("PowerShell", "-Command", fmt.Sprintf("Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('%s')", text))
						err := cmd.Run()
						if err != nil {
							log.Println("Error playing speech:", err)
						}

						resp3.Body.Close()
					}
				}

			}
		}
	}
}
