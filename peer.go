package main

import (
	"net"
)

type Peer struct {
	conn    net.Conn
	msgChan chan []byte
}

func NewPeer(conn net.Conn, msgChan chan []byte) *Peer {

	return &Peer{
		conn:    conn,
		msgChan: msgChan,
	}
}

func (p *Peer) readLoop() error {
	buff := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buff)
		if err != nil {
			return err
		}
		p.msgChan <- buff[:n]
	}
}
