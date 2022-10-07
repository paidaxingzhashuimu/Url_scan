package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

var waitgroup sync.WaitGroup

func getpath() []string {
	var urldata []string                             //创建一个数组，后面将读取到的文件内容，添加进去
	urlfile := flag.String("f", "", "输入url路径（相对路径）") //flag包获取扫描url路径
	flag.Parse()                                     //解析flag包接收的数据
	urlopen, err := os.Open(*urlfile)                //打开flag包获取到文件
	if err != nil {
		fmt.Println("不存在此文件，请检查路径和文件名！")
	}
	urlread := bufio.NewReader(urlopen) //利用bufio包创建缓存读取文件内容
	for {                               //for循环读取文件内容
		urldeline, _, err := urlread.ReadLine() //readline（）函数逐行读取
		if err == io.EOF {
			break //读到最后一行后，退出
		}
		urldeline_string := string(urldeline)       //将读取到的每一行转换为string类型
		urldata = append(urldata, urldeline_string) //将内容添加进数组
	}
	return urldata //函数返回数组内容
}

func jiaxin_gogogo(url string, file *os.File) {
	defer func() {
		err_logs, _ := os.OpenFile("err_logs.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		err := recover()
		if err != nil {
			fmt.Println(url + "  [*]格式存在问题，请检查")
			err_write := bufio.NewWriter(err_logs)
			err_write.WriteString(url + "\n")
			err_write.Flush()
		}
	}()
	fmt.Println("[*]正在扫描:", url)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //跳过ssl证书验证
	}
	client := &http.Client{Transport: tr, Timeout: 5 * time.Second} //
	httpresp, _ := http.NewRequest("GET", url, nil)                 //http包创建一个新的get请求
	resp, err := client.Do(httpresp)                                //resp接收请求返回的内容
	if err != nil {
		fmt.Println(url + "访问失败")
	}
	if resp != nil {
		statuscode := resp.Status
		if statuscode == "200 OK" { //判断返回结果的状态码是否为200
			file.WriteString(url + "\n") //将能访问的url结果写入results.txt
			fmt.Println("[+]存在url:", url)

		}

	}
	waitgroup.Done()
}

func main() {

	scan_result, err := os.OpenFile("results.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("输出结果出现意外")
	}
	urldata := getpath() //获取getpat（）函数返回的urldata数组
	waitgroup.Add(len(urldata))
	for _, urldata_read := range urldata { //遍历访问urldata数组内容

		go jiaxin_gogogo(urldata_read, scan_result) //并发运行函数

	}
	waitgroup.Wait()
}
