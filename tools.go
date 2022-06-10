package remote_log_aliyun_go

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"time"
)

// 记录上传失败的日志到本地的地址
var ErrorLogPath string

// 记录上传失败的日志
func FailUploadLog(log string) {
	file, _ := os.OpenFile(fmt.Sprintf("%v/error_log_%v.log", ErrorLogPath, time.Now().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(log + "\n")
	write.Flush()
}

//获取ip
func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

//获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
