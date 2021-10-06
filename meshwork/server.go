package meshwork

import (
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	"github.com/phuhao00/bromel/logger"
	"go.uber.org/zap"
	"net"
	"time"
)

// Server is a server for TCP connections.
type Server struct {
	Listener net.Listener

	// Packer is the message packer, will be passed to session.
	Packer Packer

	// Codec is the message codec, will be passed to session.
	Codec encoding.Codec

	// OnSessionCreate is an event hook, will be invoked when session's created.
	OnSessionCreate func(sess *Session)

	// OnSessionClose is an event hook, will be invoked when session's closed.
	OnSessionClose func(sess *Session)

	socketReadBufferSize  int
	socketWriteBufferSize int
	socketSendDelay       bool
	readTimeout           time.Duration
	writeTimeout          time.Duration
	respQueueSize         int
	router                *Router
	printRoutes           bool
	accepting             chan struct{}
	stopped               chan struct{}
	ServerAddress         string
	msgReqChan            chan *Context
	msgRspChan            chan *Context
	Logger                *logger.Logger
}

// ServerOption is the option for Server.
type ServerOption struct {
	SocketReadBufferSize  int            // sets the socket read buffer size.
	SocketWriteBufferSize int            // sets the socket write buffer size.
	SocketSendDelay       bool           // sets the socket delay or not.
	ReadTimeout           time.Duration  // sets the timeout for connection read.
	WriteTimeout          time.Duration  // sets the timeout for connection write.
	Packer                Packer         // packs and unpacks packet payload, default packer is the packet.DefaultPacker.
	Codec                 encoding.Codec // encodes and decodes the message data, can be nil.
	RespQueueSize         int            // sets the response channel size of session, 1024 will be used if < 0.
	ReqQueueSize          int            // sets the request channel size of router, 1024 will be used if < 0.
	DoNotPrintRoutes      bool           // whether to print registered route handlers to the console.
	ServerAddress         string
	ReqCtxQueue           chan *Context
	Logger                *logger.Logger
}

// NewServer creates a Server according to opt.
func NewServer(opt *ServerOption) *Server {
	if opt.Packer == nil {
		opt.Packer = &CustomPacker{}
	}
	if opt.RespQueueSize < 0 {
		opt.RespQueueSize = 1024
	}
	if opt.ReqQueueSize < 0 {
		opt.ReqQueueSize = 1024
	}
	s := &Server{
		socketReadBufferSize:  opt.SocketReadBufferSize,
		socketWriteBufferSize: opt.SocketWriteBufferSize,
		socketSendDelay:       opt.SocketSendDelay,
		respQueueSize:         opt.RespQueueSize,
		readTimeout:           opt.ReadTimeout,
		writeTimeout:          opt.WriteTimeout,
		Packer:                opt.Packer,
		Codec:                 opt.Codec,
		printRoutes:           !opt.DoNotPrintRoutes,
		router:                newRouter(),
		accepting:             make(chan struct{}),
		stopped:               make(chan struct{}),
		ServerAddress:         opt.ServerAddress,
		Logger:                opt.Logger,
	}
	s.router.reqCtxQueue = opt.ReqCtxQueue
	return s
}

// Serve starts to listen TCP and keeps accepting TCP connection in a loop.
// The loop breaks when error occurred, and the error will be returned.
func (s *Server) Serve(addr string) error {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}
	lis, err := net.ListenTCP("tcp", address)
	if err != nil {
		s.Logger.Error(err.Error())
		return err
	}
	s.Listener = lis
	if s.printRoutes {
		s.router.printHandlers(fmt.Sprintf("tcp://%s", s.Listener.Addr()))
	}
	go s.router.consumeRequest()
	return s.acceptLoop()
}

// acceptLoop accepts TCP connections in a loop, and handle connections in goroutines.
// Returns error when error occurred.
func (s *Server) acceptLoop() error {
	close(s.accepting)
	for {
		select {
		case <-s.stopped:
			s.Logger.Info("server accept loop stopped")
			return ErrServerStopped
		default:
		}

		conn, err := s.Listener.Accept()
		if err != nil {
			select {
			case <-s.stopped:
				s.Logger.Info("server accept loop stopped")
				return ErrServerStopped
			default:
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				tempDelay := time.Millisecond * 5
				s.Logger.Error("accept", zap.String("err", err.Error()))
				time.Sleep(tempDelay)
				continue
			}
			return fmt.Errorf("accept err: %s", err)
		}
		if s.socketReadBufferSize > 0 {
			if err := conn.(*net.TCPConn).SetReadBuffer(s.socketReadBufferSize); err != nil {
				return fmt.Errorf("conn set read buffer err: %s", err)
			}
		}
		if s.socketWriteBufferSize > 0 {
			if err := conn.(*net.TCPConn).SetWriteBuffer(s.socketWriteBufferSize); err != nil {
				return fmt.Errorf("conn set write buffer err: %s", err)
			}
		}
		if s.socketSendDelay {
			if err := conn.(*net.TCPConn).SetNoDelay(false); err != nil {
				return fmt.Errorf("conn set no delay err: %s", err)
			}
		}
		go s.handleConn(conn)
	}
}

// handleConn creates a new session with conn,
// handles the message through the session in different goroutines,
// and waits until the session's closed.
func (s *Server) handleConn(conn net.Conn) {
	sess := newSession(conn, &SessionOption{
		Packer:        s.Packer,
		Codec:         s.Codec,
		respQueueSize: s.respQueueSize,
		logger:        s.Logger,
	})
	Sessions().Add(sess)
	if s.OnSessionCreate != nil {
		go s.OnSessionCreate(sess)
	}

	go sess.readInbound(s.router.reqCtxQueue, s.readTimeout) // start reading message packet from connection.
	go sess.writeOutbound(s.writeTimeout)                    // start writing message packet to connection.

	<-sess.closed                // wait for session finished.
	Sessions().Remove(sess.ID()) // session has been closed, remove it.

	if s.OnSessionClose != nil {
		go s.OnSessionClose(sess)
	}
	if err := conn.Close(); err != nil {
		s.Logger.Error("connection close ", zap.String("err", err.Error()))
	}
}

// Stop stops server by closing all the TCP sessions, listener and the router.
func (s *Server) Stop() error {
	close(s.stopped)

	// close all sessions
	closedNum := 0
	Sessions().Range(func(id string, sess *Session) (next bool) {
		sess.close()
		closedNum++
		return true
	})
	s.Logger.Error("session(s) closed", zap.Int("closedNum", closedNum))

	s.router.stop()
	return s.Listener.Close()
}

// AddRoute registers message handler and middlewares to the router.
func (s *Server) AddRoute(msgID interface{}, handler HandlerFunc, middlewares ...MiddlewareFunc) {
	s.router.register(msgID, handler, middlewares...)
}

// Use registers global middlewares to the router.
func (s *Server) Use(middlewares ...MiddlewareFunc) {
	s.router.registerMiddleware(middlewares...)
}

// NotFoundHandler sets the not-found handler for router.
func (s *Server) NotFoundHandler(handler HandlerFunc) {
	s.router.setNotFoundHandler(handler)
}
