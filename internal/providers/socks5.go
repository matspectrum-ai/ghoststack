package providers

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"time"
)

type socks5Provider struct {
	mu     sync.RWMutex
	state  ProviderState
	listen string
	ln     net.Listener
}

func newSocks5Provider(config map[string]any) (Provider, error) {
	listen := "127.0.0.1:1080"
	if v, ok := config["listen"]; ok {
		listen = fmt.Sprintf("%v", v)
	}

	return &socks5Provider{
		state:  ProviderStopped,
		listen: listen,
	}, nil
}

func (p *socks5Provider) Name() string {
	return "socks5"
}

func (p *socks5Provider) State() ProviderState {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.state
}

func (p *socks5Provider) Start(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state == ProviderRunning {
		return nil
	}

	ln, err := net.Listen("tcp", p.listen)
	if err != nil {
		return fmt.Errorf("listen socks5: %w", err)
	}

	go p.accept(ctx, ln)

	p.ln = ln
	p.state = ProviderRunning
	return nil
}

func (p *socks5Provider) accept(ctx context.Context, ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go p.handleConn(ctx, conn)
	}
}

func (p *socks5Provider) handleConn(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(30 * time.Second))

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil || n < 3 {
		return
	}

	if buf[0] != 0x05 {
		return
	}

	conn.Write([]byte{0x05, 0x00})

	n, err = conn.Read(buf)
	if err != nil || n < 7 {
		return
	}

	atyp := buf[3]
	var addr string

	switch atyp {
	case 1: // IPv4
		if n < 10 {
			return
		}
		ip := net.IP(buf[4:8])
		port := binary.BigEndian.Uint16(buf[8:10])
		addr = net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	case 3: // Domain name
		domainLen := int(buf[4])
		if n < 5+domainLen+2 {
			return
		}
		domain := string(buf[5 : 5+domainLen])
		port := binary.BigEndian.Uint16(buf[5+domainLen : 7+domainLen])
		addr = net.JoinHostPort(domain, strconv.Itoa(int(port)))
	case 4: // IPv6
		if n < 22 {
			return
		}
		ip := net.IP(buf[4:20])
		port := binary.BigEndian.Uint16(buf[20:22])
		addr = net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	default:
		return
	}

	target, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return
	}
	defer target.Close()

	conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		io.Copy(target, conn)
		target.Close()
	}()
	go func() {
		defer wg.Done()
		io.Copy(conn, target)
		conn.Close()
	}()
	wg.Wait()
}

func (p *socks5Provider) Stop(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.state != ProviderRunning {
		return ErrProviderNotStarted
	}

	if p.ln != nil {
		p.ln.Close()
	}

	p.state = ProviderStopped
	return nil
}
