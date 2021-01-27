package internal

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

const maxInt = int(^uint(0) >> 1)

const (
	codeOffset    = 0
	codeLen       = 2
	payloadOffset = codeLen
	payloadLen    = 4
	packageLen    = codeLen + payloadLen
)

const (
	StatusOK = iota
	StatusShutdown
)

var (
	ErrMsgTooLarge = errors.New("message too large")
)

type parser struct {
	r      *bufio.Reader
	w      *bufio.Writer
	code   [2]byte
	header [4]byte
}

type message struct {
	code    uint16
	header  uint32
	payload []byte
}

func (p *parser) recvMsg() (*message, error) {
	var buf [packageLen]byte
	msg := &message{}
	if _, err := p.r.Read(buf[:]); err != nil {
		return nil, err
	}

	// parse code and the length of payload
	msg.code = binary.BigEndian.Uint16(buf[codeOffset:codeOffset + codeLen])
	length := binary.BigEndian.Uint32(buf[payloadOffset:payloadOffset + payloadLen])
	if int64(length) > int64(maxInt) {
		return nil, ErrMsgTooLarge
	}

	// parse payload
	msg.payload = make([]byte, length)
	if _, err := p.r.Read(msg.payload); err != nil {
		if err == io.EOF {
			err = io.ErrUnexpectedEOF
		}
		return nil, err
	}
	return msg, nil
}

func (p *parser) sendMsg(code uint16, b []byte) error {
	binary.BigEndian.PutUint16(p.code[:], code)
	binary.BigEndian.PutUint32(p.header[:], uint32(len(b)))

	// write header
	if _, err := p.w.Write(p.code[:]); err != nil {
		return err
	}
	if _, err := p.w.Write(p.header[:]); err != nil {
		return err
	}

	// write payload
	if _, err := p.w.Write(b); err != nil {
		return err
	}
	if err := p.w.Flush(); err != nil {
		return err
	}
	return nil
}
