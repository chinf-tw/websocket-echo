package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var counter = 0

func main() {
	r := gin.Default()
	r.GET("/echo", echo)
	r.GET("/cli", cliControl)
	fmt.Println("listen and serve webscoket on 0.0.0.0:8080/echo")
	r.Run() // listen and serve on 0.0.0.0:8080
}

func cliControl(c *gin.Context) {
	command := make(chan []byte)

	upgrade := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer func() {
		log.Println("disconnect !!")
		if err = ws.Close(); err != nil {
			log.Println("ws.Close()問題\n", err)
		}
	}()
	go func(w *websocket.Conn) {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(strconv.Itoa(counter)+". client: ", string(msg))
		counter++
	}(ws)
	go func(w *websocket.Conn) {
		for {
			if err := w.WriteMessage(websocket.TextMessage, <-command); err != nil {
				log.Fatal("** websocket ** ", err)
			}
		}
	}(ws)
	for {
		var cmd string

		fmt.Print("Enter your command: ")
		if _, err := fmt.Scan(&cmd); err != nil {
			log.Fatal("** command **", err)
		}
		fmt.Println("start", string(cmd), "end")
		command <- []byte(cmd)
	}
}
func echo(c *gin.Context) {
	upgrade := &websocket.Upgrader{
		//如果有 cross domain 的需求，可加入這個，不檢查 cross domain
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	fmt.Println("ws is connected")
	defer func() {
		log.Println("disconnect !!")
		if err = ws.Close(); err != nil {
			log.Println("ws.Close()問題\n", err)
		}
	}()
	var msg []byte
	for {
		_, msg, err = ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(strconv.Itoa(counter)+". client: ", string(msg))
		msg = []byte(strconv.Itoa(counter) + ". server: " + string(msg))
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			return
		}
		counter++
	}
}
