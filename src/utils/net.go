package utils

import (
	"log"
	"net"
	"strings"
)

func CheckTCP(url string) bool {
	if len(strings.Split(url, ":")) <= 1 {
		url = url + ":80"
	}
	_, err := net.Dial("tcp", url)
	if err != nil {
		log.Printf("%+v", err)
		return false
	}
	return true
}

func CheckUDP(url string) bool {
	if len(strings.Split(url, ":")) <= 1 {
		url = url + ":80"
	}
	_, err := net.Dial("udp", url)
	if err != nil {
		log.Printf("%+v", err)
		return false
	}
	return true
}
