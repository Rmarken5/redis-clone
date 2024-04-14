package main

import (
	"bufio"
	"net"
)

const bufferDelim = '\x00'

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
	reader := bufio.NewReader(p.conn)
	for {
		b, err := reader.ReadBytes(bufferDelim)
		if err != nil {
			return err
		}
		p.msgChan <- Message{
			peer: p,
			data: b,
		}
	}
}
