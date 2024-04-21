package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultAddress = ":5001"

type Command interface {
}

type Message struct {
	cmd  Command
	peer *Peer
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgChan   chan Message
	kv        *KV
}

type Config struct {
	ListenerAddress string
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenerAddress) == 0 {
		cfg = Config{ListenerAddress: defaultAddress}
	}
	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgChan:   make(chan Message),
		kv:        NewKV(),
	}
}
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenerAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()

	slog.Info("Server Running", "listenerAddr", s.ListenerAddress)
	return s.acceptLoop()

}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case SetCommand:
		return s.kv.Set(v.key, v.val)
	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key %s not in KV", v.key)
		}
		err := msg.peer.Send(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgChan:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("raw message err", "err", err)
			}
		case peer := <-s.addPeerCh:
			s.peers[peer] = true
		case <-s.quitCh:
			return
		}

	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.HandleConnection(conn)
	}

}

func (s *Server) HandleConnection(conn net.Conn) {
	peer := NewPeer(conn, s.msgChan)
	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("error in read loop", "err", err)
	}

}

func main() {
	listenAddr := flag.String("listenAddr", defaultAddress, "listen address of the goredis server")
	flag.Parse()
	server := NewServer(Config{
		ListenerAddress: *listenAddr,
	})
	log.Fatal(server.Start())
}
