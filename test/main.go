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
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")  // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
func main() {

	r := gin.Default()
	//r.Use(LoggerToFile())
	r.Use(Cors())
	r.GET("/data", func(c *gin.Context) {
		id := c.DefaultQuery("device_id", "1")
		price := c.DefaultQuery("value","1")
		ch <- msg{Id: id, Price: price}
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/ping", ping)
	r.Run() // listen and serve on 0.0.0.0:8080
}
