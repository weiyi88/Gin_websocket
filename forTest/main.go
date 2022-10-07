package main

import (
	"fmt"
	"time"
)

func recv(c chan int) {
	ret := <-c
	fmt.Println("接收成功", ret)
}

func main() {
	ch := make(chan int, 1)
	fmt.Println("111111")
	go recv(ch)
	fmt.Println("222222")
	ch <- 11
	fmt.Println("3333333")
	time.Sleep(time.Second)
	fmt.Println("发送成功")
}
