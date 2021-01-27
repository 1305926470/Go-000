package internal

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

const (
	size         = 100
	closeTimeout = 10  // seconds
)

var (
	queue = make(chan net.Conn, size)
)

type Server struct {
	closed bool
	wg     sync.WaitGroup
	mu     sync.Mutex
	conns  map[io.Closer]struct{}
}

func (s *Server) Start() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic("can not listening at port :8080")
	}
	go s.accept()
	fmt.Println("listening on :8080...")
	for {
		rawConn, err := lis.Accept()
		if err != nil {
			panic(err)
		}
		if s.closed {
			rawConn.Close()
		}
		s.addConn(rawConn)
		queue <- rawConn
	}
}

func (s *Server) accept() {
	for {
		select {
		case conn := <-queue:
			s.wg.Add(1)
			go func() {
				s.handleRawConn(conn)
				s.removeConn(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleRawConn(rawConn net.Conn) {
	r := bufio.NewReader(rawConn)
	w := bufio.NewWriter(rawConn)
	p := parser{r: r, w:w}
	defer rawConn.Close()
	for {
		msg, err := p.recvMsg()
		if err == io.EOF {
			fmt.Println("The connection was disconnected by the peer ")
			return
		}
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		fmt.Printf("received: %s, remoteAddr: %s, now: %d\n",
			msg.payload, rawConn.RemoteAddr().String(), time.Now().UnixNano())

		if err := p.sendMsg(StatusOK, []byte("server msg")); err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		if s.closed {
			p.sendMsg(StatusShutdown, []byte("the service is to being shutdown"))
			time.Sleep(3 * time.Second)
			return
		}
	}
}

func (s *Server) addConn(c io.Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns == nil {
		s.conns = make(map[io.Closer]struct{})
	}
	s.conns[c] = struct{}{}
}

func (s *Server) removeConn(c io.Closer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conns != nil {
		delete(s.conns, c)
	}
}

func (s *Server) GracefulStop() {
	s.closed = true
	fmt.Println("Incoming service shutdown")

	t := time.NewTimer(time.Second * closeTimeout)
	go func() {
		select {
		case <-t.C:
			s.mu.Lock()
			defer s.mu.Unlock()
			for c := range s.conns {
				c.Close()
			}
			os.Exit(0)
		}
	}()

	s.wg.Wait()
	os.Exit(0)
}

func NewServer() *Server {
	return &Server{
		closed: false,
		wg:     sync.WaitGroup{},
		mu:     sync.Mutex{},
		conns:  nil,
	}
}
