package main

import (
	"flag"
	"fmt"
	"monitor"
	"runtime"
	"strconv"
	"time"
	"utils"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//  获取监控时间间隔
	Interval, err := utils.GetOption("CheckInterval", "zookeeper")
	if err != nil {
		utils.Warning(fmt.Sprintf("%+v", err))
		Interval = "30"
	}
	CheckInterval, err := strconv.Atoi(Interval)
	if err != nil {
		utils.Warning(fmt.Sprintf("%+v", err))
		CheckInterval = 30
	}

	start := 1
	for range time.Tick(time.Second) {
		if start < CheckInterval {
			start++
			continue
		}
		start = 1
		monitor.JsonrpcMonitor()
		utils.Info("任务结束")
	}
}
