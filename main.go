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
	fmt.Println("listen and serve webscoket on 0.0.0.0:8080/echo")
	r.Run() // listen and serve on 0.0.0.0:8080
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
	defer func() {
		log.Println("disconnect !!")
		if err = ws.Close(); err != nil {
			log.Println("ws.Close()問題\n", err)
		}
	}()
	var msg []byte
	for {
		counter++
		_, msg, err = ws.ReadMessage()
		fmt.Println(strconv.Itoa(counter)+"client: ", string(msg))
		if err != nil {
			log.Println(err)
			return
		}
		msg = []byte(strconv.Itoa(counter) + ". server: " + string(msg))
		if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Println(err)
			return
		}
	}
}
