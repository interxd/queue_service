package client

import (
	"fmt"
	"net"
	"time"
	"strconv"
	"../common"
)

type Client struct {
	uid string
	readChan chan []byte
	writeChan chan []byte
	closeChan chan bool
}

func (c *Client) Start(cc *ClientContainer, serverIp string, serverPort int) {
	addr, err := net.ResolveTCPAddr("tcp", serverIp + ":" + strconv.Itoa(serverPort))
	conn, err := net.DialTCP("tcp", nil, addr)
	defer func() { 
		if conn != nil {
			conn.Close()
		}
	}()

	if err != nil {
		fmt.Println("connect to server failed:", err.Error())
		return
	}

	cc.OnCreateClient(c)

	conn.Write(common.EncodePacket(common.CS_REGISTER, []byte(c.uid)))
	go c.readCoroutine(conn)

	select {
		case <- c.closeChan:
			fmt.Println(c.uid, "connection closed")
			cc.globalDch <- c.uid
	}
}

func (c *Client) readCoroutine(conn *net.TCPConn) {
	for {
		//fmt.Println("process rHandler", uid)
		data := make([]byte, 128)

		conn.SetReadDeadline(time.Now().Add(2 * time.Second))
		cnt, err := conn.Read(data)

		if err != nil {	
			if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
				//timeout, go on
				continue
			} else {
				// network error
				fmt.Println("network error", c.uid, err)
				c.closeChan <- true
				return 
			}
		}

		if (cnt == 0) {
			continue
		}

		for {
			var index uint8
			var content []byte
			data, index, content = common.DecodePacket(data)

			if index == 0 {
				break
			}

			if index == common.SC_SYNC_POS {
				fmt.Println("report current positon", c.uid, common.ConvertBytesToInt32(content))
			} else if index == common.SC_PUSH_TOKEN {
				fmt.Println("get token", c.uid, string(content), time.Now())
			}
		}
	}
}
