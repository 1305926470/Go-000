package internal

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func Client() {
	conn ,err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	p := parser{r: r, w:w}

	for {
		sndMsg := []byte("hello world!")
		if err := p.sendMsg(StatusOK, sndMsg); err != nil {
			fmt.Printf("%v\n", err)
		}
		msg, err := p.recvMsg()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("received: %s, code :%d, now: %d\n",
			msg.payload, msg.code, time.Now().UnixNano())

		switch msg.code {
		case StatusShutdown:
			fmt.Println("Close the connection after 3 seconds")
			time.Sleep(3 * time.Second)
			conn.Close()
			return
		case StatusOK:
		default:
		}

		time.Sleep(time.Millisecond)
	}
}





