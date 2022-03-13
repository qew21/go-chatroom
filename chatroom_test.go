package main

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// 消息结构
type Message struct {
	Token   string `json:"token"`   //token用于表示发起用户
	Content string `json:"content"` //content表示消息内容，这里简化别的消息
}

type Sender struct {
	conn *websocket.Conn
	send chan Message
}

var (
	newline        = []byte{'\n'}
	space          = []byte{' '}
	pongWait       = 60 * time.Second
	maxMessageSize = 512
	maxClient      = 10
	maxMessage     = 1000
	mList          = list.New()
	cList          = list.New()
)

func Test_Guest(t *testing.T) {
	assert := assert.New(t)
	go main()
	time.Sleep(time.Second)
	clients := createClients(maxClient)
	process(clients)
	assert.Equal(mList.Len(), maxMessage-1)
	assert.Equal(cList.Len(), 10*(maxMessage-1))
	assert.Equal(removeDuplicate(cList).Len(), maxMessage-1)
}

// 创建一定数量的客户端
func createClients(amount int) []*Sender {
	var clients []*Sender
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8081", Path: "/chatroom/guest"}
	for i := 0; i < amount; i++ {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			panic(err)
		}
		sender := &Sender{
			conn: c,
			send: make(chan Message, 128),
		}
		go sender.loopSendMessage()
		go sender.loopReceiveMessage()
		clients = append(clients, sender)
	}
	return clients
}

// 随机的选择一些用户去发送消息
func process(clients []*Sender) {
	// 设置随机种子
	rand.Seed(time.Now().UnixNano())

	clientsAmount := len(clients) - 1

	for i := 0; i < maxMessage; i++ {
		time.Sleep(time.Millisecond * 5)
		go func() {
			randIndex := rand.Intn(clientsAmount)
			clients[randIndex].sendMessage(fmt.Sprintf("id:%v, cnt:%v", randIndex, strconv.Itoa(i)))
		}()
	}
}

// 写入消息
func (sender *Sender) sendMessage(str string) {
	message := Message{
		Token:   time.Now().Format(time.RFC3339),
		Content: str,
	}
	sender.send <- message
}

// 发送消息
func (sender *Sender) loopSendMessage() {
	for {
		m := <-sender.send
		if err := sender.conn.WriteJSON(m); err != nil {
			fmt.Println(err)
		}
		fmt.Println("发送消息", m)
		mList.PushFront(m)
	}
}

// 接收消息
func (sender *Sender) loopReceiveMessage() {
	sender.conn.SetReadLimit(int64(maxMessageSize))
	sender.conn.SetReadDeadline(time.Now().Add(pongWait))
	sender.conn.SetPongHandler(func(string) error { sender.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := sender.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("收到消息", string(message))
		cList.PushFront(string(message))
	}
}

// 转字符串
func Obj2Str(m Message) string {
	jsons, errs := json.Marshal(m)
	if errs != nil {
		fmt.Println(errs.Error())
	}
	return string(jsons)
}

// 去重
func removeDuplicate(l *list.List) *list.List {
	var sMap = make(map[string]bool)
	var next *list.Element
	for e := l.Front(); e != nil; e = next {
		m := e.Value.(string)
		next = e.Next()
		if sMap[m] == true {
			l.Remove(e)
		} else {
			sMap[m] = true
		}
	}
	return l
}
