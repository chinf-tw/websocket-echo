package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func main() {
	var wg sync.WaitGroup
	c, _, err := websocket.DefaultDialer.Dial("ws://server.cirlab:8000/echo", nil)
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

	c, _, err = websocket.DefaultDialer.Dial("ws://server.cirlab:8000/cli", nil)
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
