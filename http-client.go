// HTTP client with timeouts.
// Client from standard library lacks support of timeouts then we use this one.
// Code taken from https://gist.github.com/dmichael/5710968
// and also explained here http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
package main

import (
	"net"
	"net/http"
	"time"
)

type HTTPConfig struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
}

func TimeoutDialer(config *HTTPConfig) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(config.ReadWriteTimeout))
		tcp_conn := conn.(*net.TCPConn)
		tcp_conn.SetLinger(0) // because we have mass connects/disconnects
		return conn, nil
	}
}

func NewTimeoutClient(args ...interface{}) *http.Client {
	// Default configuration
	config := &HTTPConfig{
		ConnectTimeout:   1 * time.Second,
		ReadWriteTimeout: 1 * time.Second,
	}

	// merge the default with user input if there is one
	if len(args) == 1 {
		timeout := args[0].(time.Duration)
		config.ConnectTimeout = timeout
		config.ReadWriteTimeout = timeout
	}

	if len(args) == 2 {
		config.ConnectTimeout = args[0].(time.Duration)
		config.ReadWriteTimeout = args[1].(time.Duration)
	}

	return &http.Client{
		Transport: &http.Transport{
			Dial: TimeoutDialer(config),
		},
	}
}
