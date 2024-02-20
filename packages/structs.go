package packages

import (
	"net"
	"net/url"
	"time"
)

type Address struct {
	Host string
	Port int
}
type ProxyOptions struct {
	Proxy   *Address
	Target  *Address
	Timeout time.Duration
}
type Socket struct {
	Options    *ProxyOptions
	Connection net.Conn
}

type AttackInfo struct {
	Target      *url.URL
	Duration    int
	RequestRate int
	Threads     int
	Proxies     []string
}
