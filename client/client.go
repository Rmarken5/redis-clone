package client

import (
	"bytes"
	"context"
	"github.com/tidwall/resp"
	"io"
	"log"
	"net"
)

type Client struct {
	addr string
	conn net.Conn
}

func New(addr string) *Client {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		addr: addr,
		conn: conn,
	}
}

func (c *Client) Set(ctx context.Context, key string, val string) error {
	var buf = &bytes.Buffer{}
	wr := resp.NewWriter(buf)

	err := wr.WriteArray([]resp.Value{resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(val),
	})
	if err != nil {
		return err
	}
	buf.WriteRune('\x00')

	_, err = io.Copy(c.conn, buf)
	return err
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	var buf = &bytes.Buffer{}
	wr := resp.NewWriter(buf)

	err := wr.WriteArray([]resp.Value{resp.StringValue("GET"),
		resp.StringValue(key),
	})
	if err != nil {
		return "", err
	}
	buf.WriteRune('\x00')

	_, err = io.Copy(c.conn, buf)

	var readBuf = make([]byte, 1024)
	n, err := c.conn.Read(readBuf)
	if err != nil {
		return "", err
	}
	return string(readBuf[:n]), nil
}
