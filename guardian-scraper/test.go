package main

import (
	"fmt"
	"sync"
)

var token = make(chan int,1)

func say_hello(wg *sync.WaitGroup,i int){
	token <- 0
	fmt.Printf("hello %v\n",i)
	wg.Done()
}

func say_hi(wg *sync.WaitGroup,i int){
	<-token
	fmt.Printf("hi %v\n",i)
	wg.Done()
}

func main(){
    
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0;i<5;i++{
		 go say_hello(&wg,i)
		 go say_hi(&wg,i)
	}
    wg.Wait()

}