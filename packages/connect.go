package packages

import (
	"crypto/tls"
	"fmt"
	"net"
)

func (options *ProxyOptions) connectSOCKS5Proxy() (*Socket, error) {
	socket, err := net.DialTimeout("tcp", options.Proxy.Host+":"+ToStr(options.Proxy.Port), options.Timeout)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			socket.Close()
		}
	}()
	greeting := []byte{0x05, 0x01, 0x00}
	if _, err := socket.Write(greeting); err != nil {
		return nil, err
	}
	greetingResponse := make([]byte, 2)
	_, err = socket.Read(greetingResponse)
	if err != nil {
		return nil, err
	} else if len(greetingResponse) != 2 {
		return nil, fmt.Errorf("error: invalid response from proxy server")
	} else if greetingResponse[0] != 0x05 {
		return nil, fmt.Errorf("error: proxy server does not support socks5")
	} else if greetingResponse[1] != 0x00 {
		return nil, fmt.Errorf("error: can not perform greeting to socks5 server")
	}
	hostLength := len(options.Target.Host)
	buffer := []byte{
		0x05,
		0x01,
		0x00,
		0x03,
		byte(hostLength),
	}
	buffer = append(buffer, []byte(options.Target.Host)...)
	portBigByte := byte(options.Target.Port >> 8)
	portByte := byte(options.Target.Port)
	buffer = append(buffer, portBigByte)
	buffer = append(buffer, portByte)
	if _, err := socket.Write(buffer); err != nil {
		return nil, err
	}
	handshakeResponse := make([]byte, 10)
	_, err = socket.Read(handshakeResponse)
	if err != nil {
		return nil, err
	} else if len(handshakeResponse) != 10 {
		return nil, fmt.Errorf("error: invalid response from proxy server")
	} else if handshakeResponse[1] != 0x00 {
		return nil, fmt.Errorf("error: can not perform handshake with socks5 server")
	}
	return &Socket{
		Connection: socket,
		Options:    options,
	}, nil
}

func (socket Socket) configureTLS() (net.Conn, error) {
	config := &tls.Config{
		ServerName:         socket.Options.Target.Host,
		InsecureSkipVerify: true,
	}
	secureClient := tls.Client(socket.Connection, config)
	if err := secureClient.Handshake(); err != nil {
		return nil, err
	}
	return secureClient, nil
}

func GetClientConnection(options *ProxyOptions) (net.Conn, error) {
	socket, err := options.connectSOCKS5Proxy()
	if err != nil {
		return nil, err
	}
	secureClient, err := socket.configureTLS()
	if err != nil {
		return nil, err
	}
	return secureClient, nil
}
