package queue

import (
    "time"
    "fmt"
)

type RateLimitQueue interface {
    Start() bool
    Enqueue(uid string) int64
    GetRelativePos(idx int64) int64      //获取处在队列中的相对位置
 }

 type LeakyBucketQueue struct {
    dequeueChan chan string
    queueChan chan string
    tickInterval int
    capacity int
    rate int
    headPos int64
    tailPos int64
 }

 func CreateLeakyBucketQueue(
    deqChan chan string,
    capacity int,
    interval int,
    rate int) LeakyBucketQueue {
    queue := LeakyBucketQueue{
        dequeueChan: deqChan, 
        queueChan: make(chan string, capacity),
        capacity: capacity, 
        tickInterval: interval, 
        rate: rate,
    }

    return queue
}

func (q *LeakyBucketQueue) Start() bool {
    go q.Dispatch()
    return true
}

func (q *LeakyBucketQueue) Enqueue(uid string) int64 {
    q.queueChan <- uid
    q.tailPos++
     fmt.Println("Enqueue: ", uid, "curr queue size: ", q.tailPos - q.headPos)
    return q.tailPos
}

func (q *LeakyBucketQueue) GetRelativePos(idx int64) int64 {
    return idx - q.headPos
}

func (q *LeakyBucketQueue) Dispatch() {
    for {
        //logger.Debug("call dispatch queue");
        //fmt.Println("Dispatch!!!: ", time.Now().Unix())
        startTime := time.Now().Unix()
        consumeTime := int64(0)
        cnt := 0
        
        //dequeue 
        for {
            select {
            case uid := <- q.queueChan:
                q.dequeueChan <- uid
                q.headPos++
                cnt++
                //fmt.Println("Dispatch!!!: select", time.Now().Unix(), uid, cnt)
            default:
                //fmt.Println("Dispatch!!!: sleep_1", time.Now().Unix())
                time.Sleep(1 * time.Second)
            }
            consumeTime = time.Now().Unix() - startTime
            if (consumeTime > int64(q.tickInterval) || cnt >= q.rate) {
                //fmt.Println("Dispatch!!!: break", time.Now().Unix(), consumeTime, cnt)
                break
            }
        }
        
        //sleep
        if (consumeTime < int64(q.tickInterval)) {
            //fmt.Println("Dispatch!!!: sleep_2", time.Now().Unix(), q.tickInterval,  consumeTime, int64(q.tickInterval) - consumeTime)
            time.Sleep(time.Duration(int64(q.tickInterval) - consumeTime) * time.Second) 
        }
    }
}

