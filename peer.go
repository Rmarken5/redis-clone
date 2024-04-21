package main

import (
	"fmt"
	"github.com/tidwall/resp"
	"io"
	"log"
	"log/slog"
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan Message
}

func (p *Peer) Send(val []byte) error {
	_, err := p.conn.Write(val)
	if err != nil {
		return err
	}
	return nil
}

func NewPeer(conn net.Conn, dataChan chan Message) *Peer {
	return &Peer{
		conn:    conn,
		msgChan: dataChan,
	}
}

func (p *Peer) readLoop() error {
	rd := resp.NewReader(p.conn)
	for {
		slog.Debug("readloop", "iter")
		v, _, err := rd.ReadValue()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		if v.Type() == resp.Array {
			for _, value := range v.Array() {
				switch value.String() {
				case "SET":
					if len(v.Array()) != 3 {
						return fmt.Errorf("invalid number of variables for SET command")
					}
					cmd := SetCommand{
						key: v.Array()[1].Bytes(),
						val: v.Array()[2].Bytes(),
					}
					p.msgChan <- Message{
						cmd:  cmd,
						peer: p,
					}
				case "GET":
					if len(v.Array()) != 2 {
						return fmt.Errorf("invalid number of variables for GET command")
					}
					cmd := GetCommand{
						key: v.Array()[1].Bytes(),
					}
					p.msgChan <- Message{
						cmd:  cmd,
						peer: p,
					}
				}
			}
		}
	}
	return nil
}
