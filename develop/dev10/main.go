package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===
Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123
Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).
При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

const (
	DEADLINETIME    = time.Millisecond * 500
	READBUFFER      = 1024
	WAITBEFORECLOSE = time.Millisecond * 500
)

type client struct {
	serverAddr  string
	timeout     time.Duration
	conn        net.Conn
	ctx         context.Context
	cancel      context.CancelFunc
	abortChan   chan bool
	stdinChan   chan string
	lastMessage string
}

func newClient(serverAddr string, timeout time.Duration) client {
	c := client{
		serverAddr: serverAddr,
		timeout:    timeout,
		abortChan:  make(chan bool),
		stdinChan:  make(chan string),
	}
	c.ctx, c.cancel = context.WithCancel(context.Background())
	return c
}

func (c *client) dial() error {
	var err error
	dialer := &net.Dialer{Timeout: c.timeout}
	c.conn, err = dialer.Dial("tcp", c.serverAddr)
	if err == nil {
		log.Printf("Connected to: %s", c.serverAddr)
	}
	return err
}

func (c *client) waitExit() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)
		sig := <-ch
		fmt.Println(sig)
		c.abortChan <- true
	}()
}

func (c *client) close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	log.Print("Exited")
	return nil
}

func (c *client) cancelReadWriteClose() error {
	c.cancel()
	time.Sleep(WAITBEFORECLOSE)
	if err := c.close(); err != nil {
		return err
	}
	return nil
}

func (c *client) readFromConn() chan bool {
	go c.readRoutine()
	return c.abortChan
}

func (c *client) readFromWriteToConn() chan bool {
	go c.readRoutine()
	go c.writeRoutine()
	return c.abortChan
}

func (c *client) readRoutine() {
	reply := make([]byte, READBUFFER)
OUTER:
	for {
		select {
		case <-c.ctx.Done():
			break OUTER
		default:
			if err := c.conn.SetReadDeadline(time.Now().Add(DEADLINETIME)); err != nil {
				log.Println(err)
			}
			n, err := c.conn.Read(reply)
			if err != nil {
				if err == io.EOF {
					c.abortChan <- true
					break OUTER
				}
				if netErr, ok := err.(net.Error); ok && !netErr.Timeout() {
					log.Println(err)
				}
			}
			if n == 0 {
				break
			}
			bs := reply[:n]
			if len(bs) != 0 {
				c.lastMessage = string(bs)
			}
			fmt.Printf(c.lastMessage)
		}
	}
}

func (c *client) writeRoutine() {
	go func(stdin chan<- string) {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					c.abortChan <- true
					return
				}
				log.Println(err)
			}
			stdin <- s
		}
	}(c.stdinChan)

OUTER:
	for {
		select {
		case <-c.ctx.Done():
			break OUTER
		default:

		STDIN:
			for {
				select {
				case stdin, ok := <-c.stdinChan:
					if !ok {
						break STDIN
					}
					if _, err := c.conn.Write([]byte(stdin)); err != nil {
						log.Println(err)
					}
					c.lastMessage = stdin
				case <-time.After(DEADLINETIME):
					break STDIN
				}
			}
		}
	}
}

type args map[string]string

func getCmdArgsMap() args {
	args := make(args)

	timeoutArg := flag.String("timeout", "10s", "timeout for connection (duration)")
	fileName := filepath.Base(os.Args[0])

	flag.Usage = func() {
		fmt.Printf("example1: %s 1.2.3.4 567\n", fileName)
		fmt.Printf("example2: %s --timeout=Ns 8.9.10.11 1213\n", fileName)
		flag.PrintDefaults()
	}

	flag.Parse()
	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(2)
	}

	args["addr"] = flag.Arg(0) + ":" + flag.Arg(1)
	args["timeout"] = *timeoutArg

	return args
}
func main() {

	args := getCmdArgsMap()
	timeout, err := time.ParseDuration(args["timeout"])
	if err != nil {
		log.Fatalln(err)
	}

	client := newClient(args["addr"], timeout)
	if err := client.dial(); err != nil {
		log.Fatalln("Cannot connect:", err)
	}

	abort := client.readFromWriteToConn()
	client.waitExit()

	<-abort

	if err := client.cancelReadWriteClose(); err != nil {
		log.Fatalln("Error close client:", err)
	}

	time.Sleep(time.Second)
}
