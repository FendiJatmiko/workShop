package main
import "time"
import "fmt"
func main() {
    channel1 := make(chan string)
    channel2 := make(chan string)
    go func() {
      time.Sleep(time.Second * 4)
      channel1 <- "satu"
    }()
    go func() {
      time.Sleep(time.Second * 2)
      channel2 <- "dua"
    }()
    for i :=0; i < 2; i++{
      select {
      case msg1 := <-channel1:
        fmt.Println("menerima", msg1)
      case msg2 := <-channel2:
        fmt.Println("menerima", msg2)
      }
    }
}
