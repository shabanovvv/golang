package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	ctx     context.Context
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	cancel  context.CancelFunc
}

func (t *telnetClient) Connect() error {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(t.ctx, "tcp", t.address)
	if err != nil {
		return err
	}
	t.conn = conn
	_, err = fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	if err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Send() error {
	buf := make([]byte, 1024)
	n, err := t.in.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	if n > 0 {
		_, err = t.conn.Write(buf[:n])
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *telnetClient) Receive() error {
	scanner := bufio.NewScanner(t.conn)
	for scanner.Scan() {
		select {
		case <-t.ctx.Done():
			return t.Close()
		default:
			_, err := t.out.Write(scanner.Bytes())
			if err != nil {
				return err
			}

			_, err = t.out.Write([]byte{'\n'})
			if err != nil {
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func (t *telnetClient) Close() error {
	defer t.cancel()
	if t.conn != nil {
		err := t.conn.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	return &telnetClient{
		ctx:     ctx,
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		cancel:  cancel,
	}
}
