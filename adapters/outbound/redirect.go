package adapters

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"

	C "github.com/Dreamacro/clash/constant"
)

type RedirectOption struct {
	Name   string `proxy:"name"`
	Server string `proxy:"server"`
	Port   int    `proxy:"port"`
}

// RedirectAdapter is a Redirectly connected adapter
type RedirectAdapter struct {
	conn net.Conn
}

// Close is used to close connection
func (rd *RedirectAdapter) Close() {
	rd.conn.Close()
}

// Conn is used to http request
func (rd *RedirectAdapter) Conn() net.Conn {
	return rd.conn
}

type Redirect struct {
	name   string
	server string
	port   int
}

func (rd *Redirect) Name() string {
	return rd.name
}

func (rd *Redirect) Type() C.AdapterType {
	return C.Redirect
}

func (rd *Redirect) Generator(metadata *C.Metadata) (adapter C.ProxyAdapter, err error) {
	var c net.Conn

	isProxyServer := true

	c, err = net.DialTimeout("tcp", net.JoinHostPort(rd.server, strconv.Itoa(rd.port)), tcpTimeout)
	if err != nil {
		c, err = net.DialTimeout("tcp", net.JoinHostPort(metadata.String(), metadata.Port), tcpTimeout)
		if err != nil {
			return
		}
		isProxyServer = false
	}
	tcpKeepAlive(c)

	if isProxyServer {
		_, err = c.Write([]byte(fmt.Sprintf("CONNECT %s:%s HTTP/1.1\r\n\r\n", metadata.String(), metadata.Port)))
		if err != nil {
			return
		}

		resp, err := http.ReadResponse(bufio.NewReader(c), nil)
		if err != nil {
			return nil, err
		}
		resp.Body.Close()
	}

	return &RedirectAdapter{conn: c}, nil
}

func NewRedirect(option RedirectOption) *Redirect {
	return &Redirect{
		name:   option.Name,
		server: option.Server,
		port:   option.Port,
	}
}
