package communication

import (
	"net"
	"bufio"
	"strings"
	"log"
	"strconv"
	"fmt"
)

type Channel struct {
	socket net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewChannel(server string, port int) (*Channel) {
	channel := new(Channel)
	err := channel.Open(server, port)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to %s:%d. Cause: %s", server, port, err))
	}
	return channel
}

func (c *Channel) BindIOChannels(input chan string, output chan string) {
	// constantly read from socket and put into input go channel
	go func(readChannel chan string, pipe *Channel) {
		for true {
			line, err := pipe.ReadLine()
			if err != nil {
				log.Fatalf("Error when reading. Reason: [%s]", err)
			} else {
				readChannel <- line
			}
		}
	}(input, c)

	// constantly read from output go channel and write on socket
	go func(channelToReadFrom chan string, pipe *Channel) {
		for true {
			pipe.WriteLine(<-channelToReadFrom)
		}
	}(output, c)

}

func (c *Channel) Open(server string, port int) (error) {
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

func (c *Channel) Close() {
	c.WriteLine("QUIT interrupted\n")
	c.socket.Close()
}

func (c *Channel) ReadLine() (string, error) {
	line, err := c.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	line = strings.TrimSpace(line)
	//log.Printf(">>> %s\n", line)

	return line, nil
}

func (c *Channel) WriteLine(line string) (error) {
	log.Printf("<<< %s", line)
	_, err := c.writer.WriteString(line)
	if err != nil {
		return err
	}
	c.writer.Flush()
	return nil
}

