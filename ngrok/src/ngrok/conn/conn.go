package conn

import (
	"bufio"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	vhost "github.com/inconshreveable/go-vhost"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"ngrok/log"
	"sync"
)

type Conn interface {
	net.Conn
	log.Logger
	Id() string
	SetType(string)
	CloseRead() error
}

type loggedConn struct {
	tcp *net.TCPConn
	net.Conn
	log.Logger
	id  int32
	typ string
}

type Listener struct {
	net.Addr
	Conns chan *loggedConn
}

func wrapConn(conn net.Conn, typ string) *loggedConn {
	switch c := conn.(type) {
	case *vhost.HTTPConn:
		wrapped := c.Conn.(*loggedConn)
		return &loggedConn{wrapped.tcp, conn, wrapped.Logger, wrapped.id, wrapped.typ}
	case *loggedConn:
		return c
	case *net.TCPConn:
		wrapped := &loggedConn{c, conn, log.NewPrefixLogger(), rand.Int31(), typ}
		wrapped.AddLogPrefix(wrapped.Id())
		return wrapped
	}

	return nil
}

func Listen(addr, typ string, tlsCfg *tls.Config) (l *Listener, err error) {
	// listen for incoming connections
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	l = &Listener{
		Addr:  listener.Addr(),
		Conns: make(chan *loggedConn),
	}

	go func() {
		for {
			rawConn, err := listener.Accept()
			if err != nil {
				log.Error("不能接受新的TCP连接类型 %s: %v", typ, err)
				continue
			}

			c := wrapConn(rawConn, typ)
			if tlsCfg != nil {
				c.Conn = tls.Server(c.Conn, tlsCfg)
			}
			c.Info("新的连接 %v", c.RemoteAddr())
			l.Conns <- c
		}
	}()
	return
}

func Wrap(conn net.Conn, typ string) *loggedConn {
	return wrapConn(conn, typ)
}

func Dial(addr, typ string, tlsCfg *tls.Config) (conn *loggedConn, err error) {
	var rawConn net.Conn
	if rawConn, err = net.Dial("tcp", addr); err != nil {
		return
	}

	conn = wrapConn(rawConn, typ)
	conn.Debug("新的连接: %v", rawConn.RemoteAddr())

	if tlsCfg != nil {
		conn.StartTLS(tlsCfg)
	}

	return
}

func DialHttpProxy(proxyUrl, addr, typ string, tlsCfg *tls.Config) (conn *loggedConn, err error) {
	// parse the proxy address
	var parsedUrl *url.URL
	if parsedUrl, err = url.Parse(proxyUrl); err != nil {
		return
	}

	var proxyAuth string
	if parsedUrl.User != nil {
		proxyAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(parsedUrl.User.String()))
	}

	var proxyTlsConfig *tls.Config
	switch parsedUrl.Scheme {
	case "http":
		proxyTlsConfig = nil
	case "https":
		proxyTlsConfig = new(tls.Config)
	default:
		err = fmt.Errorf("代理URL方案必须是http或https，已获得: %s", parsedUrl.Scheme)
		return
	}

	// dial the proxy
	if conn, err = Dial(parsedUrl.Host, typ, proxyTlsConfig); err != nil {
		return
	}

	// send an HTTP proxy CONNECT message
	req, err := http.NewRequest("CONNECT", "https://"+addr, nil)
	if err != nil {
		return
	}

	if proxyAuth != "" {
		req.Header.Set("Proxy-Authorization", proxyAuth)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ngrok)")
	req.Write(conn)

	// read the proxy's response
	resp, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		return
	}
	resp.Body.Close()

	if resp.StatusCode != 200 {
		err = fmt.Errorf("非200响应来自代理服务器: %s", resp.Status)
		return
	}

	// upgrade to TLS
	conn.StartTLS(tlsCfg)

	return
}

func (c *loggedConn) StartTLS(tlsCfg *tls.Config) {
	c.Conn = tls.Client(c.Conn, tlsCfg)
}

func (c *loggedConn) Close() (err error) {
	if err := c.Conn.Close(); err == nil {
		c.Debug("Closing")
	}
	return
}

func (c *loggedConn) Id() string {
	return fmt.Sprintf("%s:%x", c.typ, c.id)
}

func (c *loggedConn) SetType(typ string) {
	oldId := c.Id()
	c.typ = typ
	c.ClearLogPrefixes()
	c.AddLogPrefix(c.Id())
	c.Info("重命名连接 %s", oldId)
}

func (c *loggedConn) CloseRead() error {
	// XXX: use CloseRead() in Conn.Join() and in Control.shutdown() for cleaner
	// connection termination. Unfortunately, when I've tried that, I've observed
	// failures where the connection was closed *before* flushing its write buffer,
	// set with SetLinger() set properly (which it is by default).
	return c.tcp.CloseRead()
}

func Join(c Conn, c2 Conn) (int64, int64) {
	var wait sync.WaitGroup

	pipe := func(to Conn, from Conn, bytesCopied *int64) {
		defer to.Close()
		defer from.Close()
		defer wait.Done()

		var err error
		*bytesCopied, err = io.Copy(to, from)
		if err != nil {
			from.Warn("复制 %d 字节到 %s 在出现错误 %v 前", *bytesCopied, to.Id(), err)
		} else {
			from.Debug("复制 %d 字节到 %s", *bytesCopied, to.Id())
		}
	}

	wait.Add(2)
	var fromBytes, toBytes int64
	go pipe(c, c2, &fromBytes)
	go pipe(c2, c, &toBytes)
	c.Info("已加入连接 %s", c2.Id())
	wait.Wait()
	return fromBytes, toBytes
}
