package server

import (
    "fmt"
    "net"
    "time"
    "../common"
)

type Message struct {
    protoId  byte
    data []byte
}

type ClientService struct {
    uid   string
    queueIdx  int64
    conn net.Conn
    readChan chan []byte
    writeChan chan Message
    closeChan chan bool
}

func NewClientService(uid string, idx int64, conn net.Conn) *ClientService {
    return &ClientService { 
        uid: uid,
        queueIdx: idx,
        conn: conn,
        readChan: make(chan []byte),
        writeChan: make(chan Message),
        closeChan: make(chan bool),
    }
}

func (cs *ClientService) writeCoroutine() {
    for {
        select {
        case <- cs.closeChan:
            fmt.Println("close writeCoroutine goroutine")
            return

        case msg := <-cs.writeChan:
            data := common.EncodePacket(msg.protoId, msg.data)
            cs.conn.Write(data)
        }
    }
}

// readCoroutine read coroutine
func (cs *ClientService) readCoroutine() {
    for {
        data := make([]byte, 128)
        cs.conn.SetReadDeadline(time.Now().Add(2 * time.Second))
        _, err := cs.conn.Read(data)

        if err != nil { 
            if neterr, ok := err.(net.Error); ok && neterr.Timeout() {
                //timeout, go on
            } else {
                cs.conn.Close()
                time.Sleep(1000 * time.Millisecond)
            }
        }
        select {
            case <- cs.closeChan:
                fmt.Println("close readCoroutine goroutine")
                return 
            default:
                continue
        }
	}
}

// Close close connection
func (cs *ClientService) Close() {
    cs.closeChan <- true
}

func (cs *ClientService) SetQueueIdx(idx int64) {
    cs.queueIdx = idx
}