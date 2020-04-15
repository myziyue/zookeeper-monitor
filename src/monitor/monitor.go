package monitor

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samuel/go-zookeeper/zk"
	"strconv"
	"strings"
	"time"
	"utils"
)

var (
	ServicesCheckPath, _     = utils.GetOption("ServicesCheckPath", "zookeeper")
	ServicesPath, _          = utils.GetOption("ServicesPath", "zookeeper")
	DefaultServicesCheckPath = "/hyperf/jsonrpc/checks"
	TimeFormat               = "2006-01-02 15:04:05"
)

type JsonRpc struct {
	ID                             string
	Name                           string
	Service                        string
	Address                        string
	Port                           int
	DeregisterCriticalServiceAfter string
	Interval                       string
	CheckTime                      string
}

func JsonrpcMonitor() {
	servers, err := utils.GetOption("Servers", "zookeeper")
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return
	}

	// 检测Zookeeper服务器可用性
	for _, server := range strings.Split(servers, ",") {
		if utils.CheckTCP(server) == false {
			utils.Errors("Zookeeper 服务器不可用")
			return
		}
	}

	// 获取zk连接
	zkConn, err := GetZkConnect(servers)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return
	}

	// 获取存储服务的目录
	servicesPath, err := utils.GetOption("ServicesPath", "zookeeper")
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return
	}

	servicesChildren, err := GetServicesPath(zkConn, servicesPath)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return
	}

	for _, serviceName := range servicesChildren {
		// Check JsonRpc Alive
		utils.Trace("------------------")
		utils.Trace(fmt.Sprintf("Start Check Service '%s' ", serviceName))
		service, err := GetService(zkConn, servicesPath, serviceName)

		if err != nil {
			utils.Errors(fmt.Sprintf("%+v", err))
			continue
		}

		// Check JsonRpc Alive
		CheckJsonrpcAlive(zkConn, service)
		utils.Trace("===============")
	}
}

func CheckJsonrpcAlive(zkConnect *zk.Conn, service string) bool {
	// Jsonrpc Service Check Zookeeper Config Path
	if ServicesCheckPath == "" {
		ServicesCheckPath = DefaultServicesCheckPath
	}

	// Check Service Alive
	status := CheckService(service)

	// Service Json Info
	Check := jsoniter.Get([]byte(service), "Check").ToString()
	ServiceID := jsoniter.Get([]byte(service), "ID").ToString()

	utils.Trace(fmt.Sprintf("Check Jsonrpc Service '%s' Alive : %+v", ServiceID, status))

	// Check Jsonrpc Service Data Struct
	JsonRpcStruct := JsonRpc{
		ID:                             ServiceID,
		Name:                           jsoniter.Get([]byte(service), "Name").ToString(),
		Service:                        jsoniter.Get([]byte(service), "Service").ToString(),
		Address:                        jsoniter.Get([]byte(service), "Address").ToString(),
		Port:                           jsoniter.Get([]byte(service), "Port").ToInt(),
		DeregisterCriticalServiceAfter: jsoniter.Get([]byte(Check), "DeregisterCriticalServiceAfter").ToString(),
		Interval:                       jsoniter.Get([]byte(Check), "Interval").ToString(),
		CheckTime:                      time.Now().Format(TimeFormat),
	}
	jsonRpc, err := jsoniter.Marshal(JsonRpcStruct)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
	}

	// Jsonrcp Alived
	if status {
		// Delete Jsonrpc Service Info
		DeleteService(zkConnect, ServicesCheckPath, ServiceID)

		// Save Check Service Info
		stat, err := SetService(zkConnect, ServicesCheckPath, ServiceID, string(jsonRpc[:]))
		utils.Trace(fmt.Sprintf("Check Service Info Pass. Save Check Jsonrpc Info, state => %+v, error => %+v", stat, err))
		if err != nil {
			utils.Errors(fmt.Sprintf("%+v", err))
		}
		return stat
	}

	// Check Interval Time
	ServicesCheck, err := GetService(zkConnect, ServicesCheckPath, ServiceID)

	// Check Service Info not Exists, Save Check Service Info.
	if err != nil {
		// Save Check Service Info
		stat, err := SetService(zkConnect, ServicesCheckPath, ServiceID, string(jsonRpc[:]))
		utils.Trace(fmt.Sprintf(" Check Service Info Not Exits. Save Check Service Info. state => %+v, error => %+v", stat, err))
		if err != nil {
			utils.Errors(fmt.Sprintf("%+v", err))
		}
		return false
	}

	//  Check Service Info is Exists, Check Deregister Critical Service After.
	deregisterCriticalServiceAfterStr := jsoniter.Get([]byte(ServicesCheck), "DeregisterCriticalServiceAfter").ToString()
	checkTimeStr, err := time.ParseInLocation(TimeFormat, jsoniter.Get([]byte(ServicesCheck), "CheckTime").ToString(), time.Local)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return false
	}
	// 获取最大监控时间和单位
	unitType := deregisterCriticalServiceAfterStr[len(deregisterCriticalServiceAfterStr)-1:]
	deregisterCriticalServiceAfter, err := strconv.Atoi(strings.ReplaceAll(deregisterCriticalServiceAfterStr, unitType, ""))
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return false
	}
	intervalTime := int(time.Now().Unix() - checkTimeStr.Unix())
	state := false

	switch unitType {
	case "s":
		state = intervalTime >= deregisterCriticalServiceAfter
		break
	case "m":
		state = intervalTime/60 >= deregisterCriticalServiceAfter
		break
	case "h":
		state = intervalTime/3600 >= deregisterCriticalServiceAfter
		break
	}
	utils.Trace(fmt.Sprintf("Jsonrpc Service Outline, Interval Time Status => %+v, Unit => %s", state, unitType))

	// Delete Check service info
	if state == true || ServicesCheck == "" {
		checkStatus := DeleteService(zkConnect, ServicesCheckPath, ServiceID)
		serviceStatus := DeleteService(zkConnect, ServicesPath, ServiceID)
		utils.Trace(fmt.Sprintf("Delete Jsonrpc Service Outline, Delete Check Info Status => %+v, Delete Service Info Status => %+v", checkStatus, serviceStatus))
	}

	return false
}
