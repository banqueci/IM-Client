package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"im-client/protocol"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"
)

var MsgChannel chan *protocol.Message

type Client struct{
	Conn 			*websocket.Conn
	SendChan 		chan *protocol.Message
	ReceiveChan 	chan *protocol.Message
	UserId			string
	PlatformId 		uint32
}

type Route struct {
	clients map[string]interface{}
}
var clientRoute Route

func main() {
	protocolStr := "wss"
	host := "ws-7chubuwilgxul-51135.gw002.oneitfarm.com"
	port := "443"
	route := ""
	userId := "d98815d1m1y00deef6fce5e12108018q6xwy35z4"
	platformId := 201

	MsgChannel = make(chan *protocol.Message, 1024)

	go StartHttpServer("81")

	url := fmt.Sprintf("%s://%s:%s/%s", protocolStr, host, port, route)

	client1 := NewClient(userId, uint32(platformId))





	client1.StartClient(url)
	clientRoute.clients[userId] = client1
}

func NewClient(userId string, platformId uint32) *Client {
	client1 := new(Client)
	client1.UserId = userId
	client1.PlatformId = platformId
	client1.SendChan = make(chan *protocol.Message, 1024)
	client1.ReceiveChan = make(chan *protocol.Message, 1024)
	return client1
}

func (client *Client) StartClient(url string){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("连接到 %s", url)

	//1、建立连接
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer conn.Close()

	//2、连接验证
	authMsg := protocol.Message{
		Cmd: protocol.MsgAuthToken,
		Seq: 0,
		Body: &protocol.Message_AuthTokenMsg{
			AuthTokenMsg: &protocol.AuthTokenMessage{
				UserId: client.UserId,
				PlatformId: client.PlatformId,
			},
		},
	}
	byte1, err1 := proto.Marshal(&authMsg)
	if err1 != nil{
		fmt.Println("Marshal failed!")
	}
	err2 := conn.WriteMessage(websocket.BinaryMessage, byte1)
	if err2 != nil{
		fmt.Println("auth failed!")
	}

	done := make(chan struct{})
	//读取消息
	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("读取错误:", err)
				return
			}
			msg := &protocol.Message{}
			proto.Unmarshal(message, msg)
			if msg.Cmd == 5{
				continue
			}
			log.Println("接收到消息: ", msg)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	//写入消息
	for {
		select {
		case <-done:
			return
		case _ = <-ticker.C:
			//定时发送心跳消息
			hearMsg := protocol.Message{
				Cmd: protocol.MsgHeartbeat,
				Seq: 0,
				Body: &protocol.Message_HeartMsg{
					HeartMsg: &protocol.HeartBeatMessage{
					},
				},
			}
			byte2, err3 := proto.Marshal(&hearMsg)
			if err3 != nil{
				log.Println("marshal failed:", err3)
			}

			err := conn.WriteMessage(websocket.BinaryMessage, byte2)
			if err != nil {
				log.Println("写入错误:", err)
				return
			}
		case msg := <-MsgChannel:
			//读取通道里面的消息
			if msg.GetImMsg().Sender != client.UserId{
				MsgChannel<-msg
				continue
			} 
			byte3, err4 := proto.Marshal(msg)
			if err4 != nil{
				log.Println("marshal failed:", err4)
				continue
			}
			err5 := conn.WriteMessage(websocket.BinaryMessage, byte3)
			if err5 != nil {
				log.Println("写入错误:", err5)
				return
			}
		case <-interrupt:
			log.Println("接收到中断信号，关闭连接...")
			conn.Close()
			return
		}
	}
}

func StartHttpServer(port string){
	route := gin.Default()

	route.GET("/debug/pprof/*any", gin.WrapF(http.DefaultServeMux.ServeHTTP))

	route.GET("/healthcheck", func(context *gin.Context) {
		var response struct {
			State uint        `json:"state"`
			Data  interface{} `json:"data"`
			Msg   string      `json:"msg"`
		}
		response.State = 200
		response.Data = nil
		response.Msg = "success"
		context.JSON(http.StatusOK, response)
	})

	route.POST("/sendMessage", func(context *gin.Context) {
		var response struct {
			State uint        `json:"state"`
			Data  interface{} `json:"data"`
			Msg   string      `json:"msg"`
		}
		sender := context.PostForm("sender")
		receiver := context.PostForm("receiver")
		msg := context.PostForm("message")

		message := &protocol.Message{
			Cmd: protocol.MsgIm,
			Seq: 0,
			Body: &protocol.Message_ImMsg{
				ImMsg: &protocol.IMMessage{
					Sender: sender,
					Receiver: receiver,
					Content: msg,
				},
			},
		}

		MsgChannel<-message

		response.State = 200
		response.Data = nil
		response.Msg = "success"
		context.JSON(http.StatusOK, response)
	})

	err := route.Run(":" + port)
	if err != nil {
		return 
	}
}