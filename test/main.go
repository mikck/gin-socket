package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	_ "github.com/sirupsen/logrus"
	"net/http"
)

type msg struct {
	Id    string
	Price string
}

var ch = make(chan msg, 100)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//webSocket请求ping 返回pong
func ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for {
		select {
		case a := <-ch:
			fmt.Print(a)
			err = ws.WriteJSON(a)
			if err != nil {
				fmt.Print(err)
				break
			}
		}
	}

	//for {
	//	//读取ws中的数据
	//	mt, message, err := ws.ReadMessage()
	//	if err != nil {
	//		break
	//	}
	//	if string(message) == "ping" {
	//		message = []byte("pong")
	//	}
	//	//写入ws数据
	//	err = ws.WriteMessage(mt, message)
	//	if err != nil {
	//		break
	//	}
	//}
}

func main() {

	r := gin.Default()
	r.Use(LoggerToFile())
	r.GET("/data", func(c *gin.Context) {
		id := c.Query("id")
		price := c.Query("price")
		ch <- msg{Id: id, Price: price}
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ping", ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
