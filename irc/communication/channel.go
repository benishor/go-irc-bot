package communication

import (
	"net"
	"bufio"
	"strings"
	"log"
	"strconv"
	"fmt"
)

type Channel interface {
	ReadLine() (string, error)
	WriteLine(string) (error)
}

type TcpChannel struct {
	Channel

	socket net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTcpChannel(server string, port int) (*TcpChannel) {
	channel := new(TcpChannel)
	err := channel.Open(server, port)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to %s:%d. Cause: %s", server, port, err))
	}
	return channel
}

func (c *TcpChannel) Open(server string, port int) (error) {
	address := server + ":" + strconv.Itoa(port)
	log.Println("Attempting to connect to ", address)

	var err error
	c.socket, err = net.Dial("tcp", address)
	if err != nil {
		return err
	}

	c.writer = bufio.NewWriter(c.socket)
	c.reader = bufio.NewReader(c.socket)

	return nil
}

func (c *TcpChannel) Close() {
	c.socket.Close()
}

func (c *TcpChannel) ReadLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	//log.Printf(">>> %s\n", line)

	return line, nil
}

func (c *TcpChannel) WriteLine(line string) (error) {
	log.Printf("<<< %s", line)
	_, err := c.writer.WriteString(line)
	if err != nil {
		return err
	}
	c.writer.Flush()
	return nil
}

