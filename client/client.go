package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {

	var (
		port = 8000
		addr = "127.0.0.1"
		c    *websocket.Conn
		err  error
	)
	var wg sync.WaitGroup
	c, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/echo", addr, port), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	wg.Add(1)
	defer c.Close()
	go func(wg *sync.WaitGroup) {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(msg))
		wg.Done()
	}(&wg)
	if err = c.WriteMessage(websocket.TextMessage, []byte("TEST!!")); err != nil {
		fmt.Println(err)
	}
	wg.Wait()

	c, _, err = websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/cli", addr, port), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Fatal("** client websocket ** ", err)
		}
		fmt.Println(string(msg))
	}
}
