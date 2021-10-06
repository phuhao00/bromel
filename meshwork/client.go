package meshwork

import (
	"errors"
	"github.com/phuhao00/bromel/logger"
	"io"
	"net"
	"time"
)

type ReceiveHandle func(entry *Entry, client *Client)

type Client struct {
	Packer        Packer
	DiaAddress    string
	ReaderNum     int
	logger        *logger.Logger
	msgCh         chan *Entry
	ReceiveHandle ReceiveHandle
}

type ClientOption struct {
	Packer        Packer
	DiaAddress    string
	ReaderNum     int
	Logger        *logger.Logger
	MsgCh         chan *Entry
	ReceiveHandle ReceiveHandle
}

func NewClient(opt *ClientOption) *Client {
	obj := &Client{
		Packer:        opt.Packer,
		DiaAddress:    opt.DiaAddress,
		logger:        opt.Logger,
		msgCh:         opt.MsgCh,
		ReceiveHandle: opt.ReceiveHandle,
	}
	return obj
}

func (c *Client) Run() {
	c.sender()
	for i := 0; i < c.ReaderNum; i++ {
		c.readerClient(i)
	}
}

func (c *Client) sender() {
	conn, err := c.establish(c.DiaAddress)
	if err != nil {
		c.logger.Error(err.Error())
	}
	go func() {
		for {
			select {
			case msg := <-c.msgCh:
				packedMsg, err := c.Packer.Pack(msg)
				if err != nil {
					c.logger.Error(err.Error())
					panic(err)
				}
				if _, err := conn.Write(packedMsg); err != nil {
					c.logger.Error(err.Error())
					panic(err)
				}
			}
		}
	}()
	go func() {
		for {
			msg, err := c.Packer.Unpack(conn)
			if err != nil {
				ne, ok := err.(net.Error)
				if errors.Is(err, io.EOF) || (ok && ne.Temporary()) {
					continue
				}
				c.logger.Error(err.Error())
				return
			}
			c.ReceiveHandle(msg, c)
		}
	}()
}

func (c *Client) readerClient(id int) {
	conn, err := c.establish(c.DiaAddress)
	if err != nil {
		c.logger.Error(err.Error())
		return
	}

	go func() {
		for {
			msg, err := c.Packer.Unpack(conn)
			if err != nil {
				ne, ok := err.(net.Error)
				if errors.Is(err, io.EOF) || (ok && ne.Temporary() && ne.Timeout()) {
					continue
				}
				return
			}
			c.ReceiveHandle(msg, c)
		}
	}()
}

func (c *Client) WriteMsg(msg *Entry) {
	c.msgCh <- msg
}

func (c *Client) establish(address string) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", address, time.Second*30)
	if err != nil {
		time.Sleep(time.Second)
		return c.establish(address)
	}
	err = conn.(*net.TCPConn).SetKeepAlive(true)

	return conn, err
}
