package zfunc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func GetTokenAuthor() (token, authorization, line string) {

	file, err := os.Open("cookie.txt")
	if err != nil {
		log.Fatal("cookie.txt文件不存在", err)

	}
	//创建一个新的buffer对象
	buff := bufio.NewReader(file)
	//1、读取第一行文件内容
	token, _ = buff.ReadString('\n')
	token = strings.TrimSpace(token)
	//2、读取第二行文件内容
	authorization, _ = buff.ReadString('\n')
	authorization = strings.TrimSpace(authorization)
	//3、读取第三行文件内容
	line, _ = buff.ReadString('\n')
	line = strings.TrimSpace(line)

	return token, authorization, line
}

func MatchString(pattern string, body string) [][]string {
	reg := regexp.MustCompile(pattern)
	res := reg.FindAllStringSubmatch(body, -1)
	return res
}

func PrintName() {
	fmt.Println("| ------------------------------------------------------------------------------|")
	fmt.Println("|                          欢迎使用自动接单工具AWork                            |")
	fmt.Println("|请cookie.txt第一行放入token,第二行放入authorization;第三行可以设置接单提示次数;|")
	fmt.Println("|请business.txt文件中定义需要匹配的告警描述信息,建议使用业务名,使用换行分割;    |")
	fmt.Println("|                                                            developer：linzheng|")
	fmt.Println("| ------------------------------------------------------------------------------|")
}

func Voice(sound string) int {
	var lineto, _ = strconv.Atoi(sound)
	if lineto == 0 {
		lineto = 1
	}
	fmt.Printf("当前接单成功后语言播放次数设置为：%v； 默认为：1\n", lineto)
	return lineto
}

func LogFunc() *os.File {
	var fileTimeFormat string = "2006-01-02"

	err := os.Mkdir("log", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		log.Fatal("创建目录失败", err)
	}
	fileNmae := fmt.Sprintf("log/%v-AWork日志.txt", time.Now().Format(fileTimeFormat))
	logFile, err := os.OpenFile(fileNmae, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println("文件写入失败", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)
	return logFile
}

var ErrorCount = new(int)

func PrintErrNet(count int) {
	if *ErrorCount >= count {
		text := "请求异常，请检查网络连接！"
		// 使用 PowerShell 播放语音
		cmd := exec.Command("PowerShell", "-Command", fmt.Sprintf("Add-Type -AssemblyName System.Speech; (New-Object System.Speech.Synthesis.SpeechSynthesizer).Speak('%s')", text))
		cmd.Run()
	}
}
