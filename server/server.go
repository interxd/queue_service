package server

// server
import (
    "fmt"
    "net"
    "time"
	"./queue"
    "../common"
    "github.com/Unknwon/goconfig"
)

// Server
type Server struct {
    ip   string
    port   int
    clients common.ConcurrentMap   
    clientQueue queue.LeakyBucketQueue
    maxQueueSize int
    dequeueChan chan string
    broadcastInterval int
}

func CreateServer(
    ip string,
    port int,
    cfg *goconfig.ConfigFile) *Server {
    
    maxQueueSize, _ := cfg.Int("queue", "max_queue_size")
    interval, _ := cfg.Int("queue", "interval")
    rate, _ := cfg.Int("queue", "rate")
    dequeueChan:= make(chan string, maxQueueSize)

    fmt.Println("CreateLeakyBucketQueue: ", interval, rate)

    lbq := queue.CreateLeakyBucketQueue(
        dequeueChan,
        maxQueueSize,
        interval,
        rate,
    )
    //var q queue.RateLimitQueue;
    //q = lbq 

    server := &Server{
        ip: ip,
        port: port,
        clients: common.New(),
        clientQueue: lbq,
        dequeueChan : dequeueChan,
        maxQueueSize : maxQueueSize,
        broadcastInterval : interval,
	}
	
    return server
}

func (s *Server) Start() {
    listener, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(s.ip), s.port, ""})
    if err != nil {
        fmt.Println("server start listen failed: ", err.Error())
        return
    }
    
    go s.AssignToken()

    s.clientQueue.Start();
    
    fmt.Println("start work succ, waiting for clients")

    for {
        conn, err := listener.AcceptTCP()
        if err != nil {
            fmt.Println("client connect exception", err.Error())
            continue
        }
        go s.OnNewConnection(conn)
    }   
}

func (s *Server) AssignToken() {
    totalCnt := 0
    lastBroadcastTime := time.Now().Unix()
    for {
        shouldBroadcast := false
        select {
        case uid := <- s.dequeueChan:
            totalCnt++
            fmt.Println("AssignToken: ", uid, time.Now().Unix(), totalCnt)
            shouldBroadcast = true
			if cs, ok := s.clients.Get(uid); ok {
                token := "random_token"
                client := cs.(*ClientService)
                client.writeChan <- Message {common.SC_PUSH_TOKEN, []byte(token)}
                client.Close()
			}
		default:
            time.Sleep(1 * time.Second)
        }
        
        if shouldBroadcast && (time.Now().Unix() - lastBroadcastTime >= int64(s.broadcastInterval)) {
            s.broadcastQueuePosition()
            lastBroadcastTime = time.Now().Unix()
            shouldBroadcast = false
        }
    }
}


func (s *Server) broadcastQueuePosition() {
    fmt.Println("broadcastQueuePosition#####################: ", time.Now().Unix())
    
    s.clients.IterCb(func(key string, val interface{}) {
        cs := val.(*ClientService)
        s.NotifyClientCurrentPos(cs)
	})
}

func (s *Server) OnNewConnection(conn net.Conn) {
    connFrom := conn.RemoteAddr().String()
    fmt.Println("Connection from: ", connFrom)
    
    defer func() { 
        if conn != nil {
            conn.Close()
        }
    }()

    if s.clients.Count() >= s.maxQueueSize {
        fmt.Println("reach max clients num, return", s.clients.Count(), s.maxQueueSize)
        return
    }

	ret, uid := s.HandleRegister(conn);
	//fmt.Println("HandleRegister", ret, uid)
    if !ret {
        return;
    }
    
    //create new client and enqueue
    client := NewClientService(uid, 0, conn)
    go client.readCoroutine()
    go client.writeCoroutine()

    idx := s.clientQueue.Enqueue(uid);
    client.SetQueueIdx(idx);
    
    s.clients.Set(uid, client)
    s.NotifyClientCurrentPos(client)

    select {
	case <- client.closeChan:
        fmt.Println("close:" , uid)
        s.clients.Remove(uid)
		return
	}
}

func (s *Server) HandleRegister(conn net.Conn) (bool, string) {
    data := make([]byte, 128)
    conn.Read(data)
    _, idx, content := common.DecodePacket(data)

    if idx != common.CS_REGISTER {
        return false, ""
    }

    uid := string(content)
    return true, uid
}

func (s *Server) NotifyClientCurrentPos(client *ClientService) {
    pos := int32(s.clientQueue.GetRelativePos(client.queueIdx))
	client.writeChan <- Message {common.SC_SYNC_POS, common.ConvertInt32ToBytes(pos)}
	//fmt.Println("NotifyClientCurrentPos: ", client.uid, pos)
}
