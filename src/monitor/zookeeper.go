package monitor

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
	"utils"
)

func GetZkConnect(servers string) (*zk.Conn, error) {
	// 连接Zookeeper服务器
	zkConnect, _, err := zk.Connect(strings.Split(servers, ","), time.Second)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return nil, err
	}
	return zkConnect, err
}

func GetServicesPath(zkConnect *zk.Conn, servicesPath string) ([]string, error) {
	if servicesPath == "" {
		return []string{}, nil
	}

	// 获取存储服务下的子目录
	children, _, err := zkConnect.Children(servicesPath)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return []string{}, err
	}
	return children, err
}

func GetService(zkConnect *zk.Conn, servicesPath string, serviceName string) (string, error) {
	if serviceName == "" || servicesPath == "" {
		return "", nil
	}
	servicesPath += "/" + serviceName

	_, _, err := zkConnect.Exists(servicesPath)
	if err != nil {
		utils.Warning(fmt.Sprintf("%+v (%s)", err, servicesPath))
		return "", err
	}

	// get service content
	value, _, err := zkConnect.Get(servicesPath)
	if err != nil {
		utils.Errors(fmt.Sprintf("%+v", err))
		return "", err
	}

	return string(value[:]), err
}

func CheckService(service string) bool {
	address := jsoniter.Get([]byte(service), "Address").ToString()
	port := jsoniter.Get([]byte(service), "Port").ToString()

	if address == "" && port == "" {
		return false
	}

	return utils.CheckTCP(address + ":" + port)
}

func SetService(zkConnect *zk.Conn, servicesPath string, servicesName string, serviceValue string) (bool, error) {
	if servicesPath == "" {
		servicesPath = "/hyperf/jsonrpc/checks/" + servicesName
	} else {
		servicesPath = servicesPath + "/" + servicesName
	}

	exists, stat, err := zkConnect.Exists(servicesPath)
	if err != nil {
		utils.Warning(fmt.Sprintf("%+v(%s)", err, servicesPath))
		return false, err
	}
	if exists == false {
		_, err := MakeNodes(zkConnect, servicesPath)
		if err != nil {
			utils.Errors(fmt.Sprintf("%+v", err))
			return false, err
		}
	}

	_, error := zkConnect.Set(servicesPath, []byte(serviceValue), stat.Version)
	if error != nil {
		utils.Errors(fmt.Sprintf("%+v (%s)", error, servicesPath))
		return false, error
	}

	return true, error
}

func DeleteService(zkConnect *zk.Conn, servicesPath string, servicesName string) bool {
	if servicesName != "" {
		servicesPath += "/" + servicesName
	}

	exists, stat, err := zkConnect.Exists(servicesPath)
	if err != nil {
		utils.Warning(fmt.Sprintf("%+v(%s)", err, servicesPath))
		return true
	}
	if exists == false {
		utils.Warning(fmt.Sprintf("'%s' not exists.", servicesPath))
		return true
	}

	err = zkConnect.Delete(servicesPath, stat.Version)
	if err != nil {
		utils.Errors(fmt.Sprintf("Delete Check Jsonrpc Info, error => %+v", err))
		return false
	}
	utils.Trace("Delete Check Jsonrpc Info: success")
	return true
}

func MakeNodes(zkConnect *zk.Conn, servicesPath string) (bool, error) {
	ServicesPath := ""
	for _, path := range strings.Split(servicesPath, "/") {
		if path == "" {
			continue
		}
		ServicesPath += "/" + path
		exists, _, err := zkConnect.Exists(ServicesPath)
		if err != nil {
			utils.Errors(fmt.Sprintf("%+v (%s)", err, ServicesPath))
			return false, err
		}

		if exists == false {
			_, err := zkConnect.Create(ServicesPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
			if err != nil {
				utils.Errors(fmt.Sprintf("Make Nodes '%s' error : %+v", ServicesPath, err))
				return false, err
			}
			utils.Trace(fmt.Sprintf("Make Nodes '%s': success", ServicesPath))
		}
	}

	return true, nil
}
