package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-chatroom/pkg/db"
	"simple-chatroom/pkg/pb"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func formatMessage(jsonMsg pb.SocketMessage) []byte {
	byteMsg, _ := json.Marshal(&jsonMsg)
	return byteMsg
}

func main() {
	r := gin.Default()
	guest := melody.New()
	user := melody.New()

	r.GET("/guest", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "guest.html")
	})

	r.GET("/user", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "user.html")
	})

	r.GET("/chatroom/user", func(c *gin.Context) {
		user.HandleRequest(c.Writer, c.Request)
	})

	r.GET("/chatroom/guest", func(c *gin.Context) {
		guest.HandleRequest(c.Writer, c.Request)
	})

	user.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("收到消息", string(msg))
		var jmsg pb.SocketMessage
		if err := json.Unmarshal([]byte(msg), &jmsg); err == nil {
			switch jmsg.Action {
			case "login":
				chats := db.PullMessage(jmsg.Group)
				fmt.Println(chats)
				for _, chat := range chats {
					bmsg, _ := json.Marshal(&chat)
					user.BroadcastFilter(bmsg, func(q *melody.Session) bool {
						return s == q
					})
				}
			case "registUser":
				err := db.CreateUser(jmsg.User)
				if err == nil {
					jmsg.Content = "注册成功"
					user.BroadcastFilter(formatMessage(jmsg), func(q *melody.Session) bool {
						return s == q
					})
				}
			case "registGroup":
				err := db.CreateGroup(jmsg.Group)
				if err == nil {
					err := db.JoinGroup(jmsg.User, jmsg.Group)
					if err == nil {
						jmsg.Content = fmt.Sprintf("注册%s成功", jmsg.Group)
						user.BroadcastFilter(formatMessage(jmsg), func(q *melody.Session) bool {
							return s == q
						})
					}
				}
			case "addUser":
				err := db.JoinGroup(jmsg.User, jmsg.Group)
				if err == nil {
					jmsg.Content = fmt.Sprintf("加入%s成功", jmsg.Group)
					user.BroadcastFilter(formatMessage(jmsg), func(q *melody.Session) bool {
						return s == q
					})
					db.SaveMessage(jmsg)
				}
			case "message":
				db.SaveMessage(jmsg)
				user.BroadcastFilter(formatMessage(jmsg), func(q *melody.Session) bool {
					return s == q
				})
			default:
				user.Broadcast(msg)
			}
		} else {
			fmt.Println(err)
		}

	})

	guest.HandleMessage(func(s *melody.Session, msg []byte) {
		guest.Broadcast(msg)
	})

	r.Run(":8081")
}
