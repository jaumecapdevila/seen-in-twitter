package twittervotes

import (
	"io"
	"net"
	"time"
)

var conn net.Conn
var reader io.ReadCloser

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}
	netc, err := net.DialTimeout(netw, addr, time.Second*5)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

func closseConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}
