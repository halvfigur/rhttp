package rhttp

import (
	"context"
	"errors"
	"io"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	proxyAddr string
	proxy     net.Conn
	local     net.Conn
}

// NewClient forwards traffic to from the proxy at addr to a local port
func NewClient(ctx context.Context, addr string, port int) (client *Client, err error) {
	proxy, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	buf := make([]byte, 128)
	n, err := proxy.Read(buf)
	if err != nil {
		return
	}

	local, err := net.Dial("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(port)))
	if err != nil {
		proxy.Close()
		return
	}

	proxyAddr := string(buf[:n])

	go func() {
		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer local.Close()
			defer wg.Done()
			if _, err := io.Copy(local, proxy); err != nil {
				return
			}
		}()

		go func() {
			defer proxy.Close()
			defer wg.Done()
			if _, err := io.Copy(proxy, local); err != nil {
				return
			}
		}()

		wg.Wait()
	}()

	client = &Client{
		proxyAddr: proxyAddr,
		proxy:     proxy,
		local:     local,
	}

	return
}

func (c *Client) ProxyAddr() string {
	return c.proxyAddr
}

func (c *Client) Close() (err error) {
	if c.proxy != nil {
		if err = c.proxy.Close(); err != nil {
			return err
		}

		c.proxy = nil
	}

	if c.local != nil {
		if err = c.local.Close(); err != nil {
			return err
		}

		c.local = nil
	}

	if c.proxy == nil && c.local == nil {
		return errors.New("closed")
	}

	return
}
