package client

import (
	"fmt"
	"strconv"
	"../common"
)

type ClientContainer struct {
	clients common.ConcurrentMap
	globalDch  chan string     // done channel
	clientsNum int             
}

func CreateClientContariner(
	clientsNum int) *ClientContainer {
	clientContainer := &ClientContainer{
		clients: common.New(),
		globalDch: make(chan string),
		clientsNum: clientsNum,
	}

	return clientContainer
}


func (cc *ClientContainer) genUid(idx int) string {
	return "client_" + strconv.Itoa(idx)
}

func (cc *ClientContainer) Start(serverIp string, serverPort int) {
	for i := 0; i < cc.clientsNum; i++ {
		uid := cc.genUid(i) 
		cc.createClient(serverIp, serverPort, uid)
	}

	// if reportStatistics {
	// 	go clientCluster.report()
	// }
	
	for {
		select {
			case uid := <- cc.globalDch:
				fmt.Println(uid, "ClientContainer connection closed")
				cc.clients.Remove(uid)

				if cc.clients.Count() == 0 {
					goto ForEnd
			}
		}
	}
	ForEnd:
		fmt.Println("done")
		//log.Println("done")
  
}

func (cc *ClientContainer) createClient(serverIp string, serverPort int, uid string) {
	fmt.Println("createClient: ", uid)
	client := &Client{
		uid : uid,
		readChan : make(chan []byte),
		writeChan : make(chan []byte),
		closeChan : make(chan bool),
	}

	go client.Start(cc, serverIp, serverPort)
}

func (cc *ClientContainer) OnCreateClient(c *Client) {
	cc.clients.Set(c.uid, c)
}



// func (cl *ClientCluster) report() {
// 	// report interval 10 seconds
// 	ticker := time.NewTicker(10 * time.Second)
// 	for {
// 		select {
// 		case <- cl.globalDch:
// 			return 

// 		case <-ticker.C:
// 			activeClients := cl.cMap.Count()
// 			fmt.Printf("report statistics: %d clients total, %d connections waiting, %d connections received token \n", cl.concurrentClients, 
// 						activeClients, cl.concurrentClients - activeClients)
// 		}
// 	}
// }