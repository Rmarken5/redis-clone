package main

import (
	"fmt"
	"log"
	"log/slog"
	"net"
)

const defaultAddress = ":5001"

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgChan   chan []byte
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
		msgChan:   make(chan []byte),
	}
}

type Client struct {
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

func (s *Server) handleRawMessage(rawMsg []byte) error {
	fmt.Println("Msg: " + string(rawMsg))
	cmd, err := parseCommand(rawMsg)
	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case SetCommand:
		slog.Info("someone wants to set a key to the hash table", "key", v.key, "val", v.val)
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgChan:
			if err := s.handleRawMessage(rawMsg); err != nil {
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
	server := NewServer(Config{})
	log.Fatal(server.Start())

}
