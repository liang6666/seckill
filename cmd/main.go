package main

import (
	"bufio"
	"flag"
	"github.com/sunshibao/go-jdmt/global"
	"github.com/sunshibao/go-jdmt/logs"
	"github.com/sunshibao/go-jdmt/secKill"
	"os"
	"strings"
	"time"
)
//
//
var skuId = flag.String("sku", "100012043978", "")
var num = flag.Int("num", 2, "")
var works = flag.Int("works", 7, "")
var start = flag.String("time", "11:59:59.500", "")
var browserPath = flag.String("execPath", "", "")
var eid = flag.String("eid", "", "")
var fp = flag.String("fp", "", "")
var payPwd = flag.String("payPwd", "", "")
var isFileLog = flag.Bool("isFileLog", true, "")

func init() {
	flag.StringVar(&global.PushToken, "token", "", "")

	flag.Parse()
}

func main() {
	var err error

	if *isFileLog {
		logs.AllowFileLogs()
	}

	execPath := ""
	if *browserPath != "" {
		execPath = *browserPath
	}
RE:
	jdSecKill := secKill.NewJdSecKill(execPath, *skuId, *num, *works)
	jdSecKill.StartTime, err = global.Hour2Unix(*start)
	if err != nil {
		logs.Fatal("开始时间初始化失败", err)
	}

	jdSecKill.PayPwd = *payPwd
	if *eid != "" {
		if *fp == "" {
			logs.Fatal("请传入fp参数")
		}
		jdSecKill.SetEid(*eid)
	}

	if *fp != "" {
		if *eid == "" {
			logs.Fatal("请传入eid参数")
		}
		jdSecKill.SetFp(*fp)
	}

	if jdSecKill.StartTime.Unix() < time.Now().Unix() {
		jdSecKill.StartTime = jdSecKill.StartTime.AddDate(0, 0, 1)
	}
	jdSecKill.SyncJdTime()
	logs.PrintlnInfo("开始执行时间为：", jdSecKill.StartTime.Format(global.DateTimeFormatStr))

	err = jdSecKill.Run()
	if err != nil {
		if strings.Contains(err.Error(), "exec") {
			logs.PrintlnInfo("默认浏览器执行路径未找到，" + execPath + "  请重新输入：")
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				execPath = scanner.Text()
				if execPath != "" {
					break
				}
			}
			goto RE
		}
		logs.Fatal(err)
	}
}
