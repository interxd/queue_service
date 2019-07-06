package queue
    
import (
    "testing"
    "fmt"
    "time"
)

func AssignTocken(uidChan chan string) {
    for {
        select {
        case uid := <- uidChan:
            fmt.Println("AssignTocken: ", uid, time.Now().Unix())
        }
    }
}

func Test_Queue_1(t *testing.T) {
    outChan:= make(chan string, 100)
    lbq := CreateLeakyBucketQueue(
        outChan,
        10,
        2,
        5,
    )

    lbq.Start()
    go AssignTocken(outChan)

    i := 1
    for {
        lbq.Enqueue(fmt.Sprintf("%v", i))
        i += 1
        if i > 20 {
            break
        }
    }

    time.Sleep(10 * time.Second)
}


func Test_Queue_2(t *testing.T) { 
    outChan:= make(chan string, 100)
    lbq := CreateLeakyBucketQueue(
        outChan,
        10,
        2,
        5,
    )

    lbq.Start()
    go AssignTocken(outChan)

    lbq.Enqueue("1")
    lbq.Enqueue("2")
    time.Sleep(3 * time.Second)

    lbq.Enqueue("3")
    lbq.Enqueue("4")
    lbq.Enqueue("5")
    lbq.Enqueue("6")
    lbq.Enqueue("7")
    lbq.Enqueue("8")
    lbq.Enqueue("9")
    lbq.Enqueue("10")

    time.Sleep(5 * time.Second)
}
