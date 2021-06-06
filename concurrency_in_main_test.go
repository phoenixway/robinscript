package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestConcurrentLogic(t *testing.T) {
	msg := make(chan string)
	reply := make(chan string)
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go wsserver(t, wg, msg)
	go ai_engine(t, wg, msg, reply)
	go executor(t, wg, reply)
	wg.Wait()
	fmt.Println("main finished")
}

func wsserver(t *testing.T, wg *sync.WaitGroup, msg chan string) {
	defer wg.Done()
	t.Logf("wsserver started\n")
	msg <- "message1"
	msg <- "message2"
	msg <- "message3"
	close(msg)
	t.Logf("wsserver finished\n")
}

func ai_engine(t *testing.T, wg *sync.WaitGroup, msg, reply chan string) {
	t.Logf("ai_engine started\n")
	defer wg.Done()
	var m string
	var ok bool
	for {
		m, ok = <-msg
		if !ok {
			wg.Done()
			fmt.Println("ai_core finished!")
			return
		}
		t.Logf("ai_engine got '%s'!\n", m)
		reply <- fmt.Sprintf("echo reply from '%s'", m)
	}
}

func executor(t *testing.T, wg *sync.WaitGroup, reply chan string) {
	t.Logf("executor started\n")
	defer wg.Done()
	for {
		t.Logf("executor got text reply '%s'!\n", <-reply)
	}
}
