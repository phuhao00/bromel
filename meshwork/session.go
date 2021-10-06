package meshwork

import (
	"errors"
	"fmt"
	"github.com/phuhao00/bromel/encoding"
	"github.com/phuhao00/bromel/logger"
	"go.uber.org/zap"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Session represents a TCP session.
type Session struct {
	id        string         // session's ID. it's a UUID
	conn      net.Conn       // tcp connection
	closed    chan struct{}  // to close()
	respQueue chan *Entry    // response queue channel, pushed in SendResp() and popped in writeOutbound()
	packer    Packer         // to pack and unpack message
	codec     encoding.Codec // encode/decode message data
	ctxPool   sync.Pool
	logger    *logger.Logger
}

// SessionOption is the extra options for Session.
type SessionOption struct {
	Packer        Packer
	Codec         encoding.Codec
	respQueueSize int
	logger        *logger.Logger
}

// newSession creates a new Session.
// Parameter conn is the TCP connection,
// opt includes packer, codec, and channel size.
// Returns a Session pointer.
func newSession(conn net.Conn, opt *SessionOption) *Session {
	id := uuid.NewString()
	return &Session{
		id:        id,
		conn:      conn,
		closed:    make(chan struct{}),
		respQueue: make(chan *Entry, opt.respQueueSize),
		packer:    opt.Packer,
		codec:     opt.Codec,
		ctxPool:   sync.Pool{New: func() interface{} { return new(Context) }},
		logger:    opt.logger,
	}
}

// ID returns the session's ID.
func (s *Session) ID() string {
	return s.id
}

// SendResp pushes response message entry to respQueue.
// Returns error if session is closed.
func (s *Session) SendResp(respMsg *Entry) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("sessions is closed")
		}
	}()

	select {
	case s.respQueue <- respMsg:
	case <-s.closed:
		close(s.respQueue)
		err = errors.New("sessions is closed")
	}

	return
}

// close closes the session.
func (s *Session) close() {
	defer func() { _ = recover() }()
	close(s.closed)
}

// readInbound reads message packet from connection in a loop.
// And send unpacked message to reqQueue, which will be consumed in router.
// The loop breaks if errors occurred or the session is closed.
func (s *Session) readInbound(reqQueue chan<- *Context, timeout time.Duration) {
	for {
		if timeout > 0 {
			if err := s.conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
				s.logger.Error(err.Error())
				break
			}
		}
		entry, err := s.packer.Unpack(s.conn)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() && netErr.Temporary() {
				continue
			}
			s.logger.Error(err.Error(), zap.String("s.conn.RemoteAddr:", s.conn.RemoteAddr().String()))
			break
		}
		if entry == nil {
			continue
		}

		ctx := s.ctxPool.Get().(*Context)
		ctx.session = s
		ctx.reqMsgEntry = entry
		ctx.storage = nil // reset storage
		select {
		case reqQueue <- ctx:
		case <-s.closed:
			s.logger.Error(fmt.Sprintf("session %s readInbound exit because session is closed", s.id))
			return
		}
	}
	s.logger.Error(fmt.Sprintf("session %s readInbound exit because of error", s.id))
	s.close()
}

// writeOutbound fetches message from respQueue channel and writes to TCP connection in a loop.
// Parameter writeTimeout specified the connection writing timeout.
// The loop breaks if errors occurred, or the session is closed.
func (s *Session) writeOutbound(writeTimeout time.Duration) {
FOR:
	for {
		select {
		case <-s.closed:
			s.logger.Error(fmt.Sprintf("session %s writeOutbound exit because session is closed", s.id))
			return
		case respMsg, ok := <-s.respQueue:
			if !ok {
				s.logger.Error(fmt.Sprintf("session %s writeOutbound exit because session is closed", s.id))
				return
			}
			// pack message
			outboundMsg, err := s.packer.Pack(respMsg)
			if err != nil {
				s.logger.Error(fmt.Sprintf("session %s pack outbound message err: %s", s.id, err))
				continue
			}
			if outboundMsg == nil {
				continue
			}
			if writeTimeout > 0 {
				if err := s.conn.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
					s.logger.Error(fmt.Sprintf("session %s set write deadline err: %s", s.id, err))

					break FOR
				}
			}
			if _, err := s.conn.Write(outboundMsg); err != nil {
				s.logger.Error(fmt.Sprintf("session %s conn write err: %s", s.id, err))
				break FOR
			}
		}
	}
	s.close()
	s.logger.Error(fmt.Sprintf("session %s writeOutbound exit because of error", s.id))

}
