package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"im-client/protocol"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	protocolStr := "wss"
	host := "ws-7chubuwilgxul-51135.gw002.oneitfarm.com"
	port := "443"
	route := ""
	userId := "d98815d1m1y00deef6fce5e12108018q6xwy35z4"
	platformId := 201

	url := fmt.Sprintf("%s://%s:%s/%s", protocolStr, host, port, route)
	StartClient(url, userId, uint32(platformId))

	go StartHttpServer("80")
}

func StartClient(url string, userId string, platformId uint32){
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("连接到 %s", url)

	//1、建立连接
	client, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("连接错误:", err)
	}
	defer client.Close()

	//2、连接验证
	authMsg := protocol.Message{
		Cmd: protocol.MsgAuthToken,
		Seq: 0,
		Body: &protocol.Message_AuthTokenMsg{
			AuthTokenMsg: &protocol.AuthTokenMessage{
				UserId: userId,
				PlatformId: platformId,
			},
		},
	}
	byte1, err1 := proto.Marshal(&authMsg)
	if err1 != nil{
		fmt.Println("Marshal failed!")
	}
	err2 := client.WriteMessage(websocket.BinaryMessage, byte1)
	if err2 != nil{
		fmt.Println("auth failed!")
	}
	//3、发送心跳


	done := make(chan struct{})

	//读取消息
	go func() {
		defer close(done)
		for {
			_, message, err := client.ReadMessage()
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

			err := client.WriteMessage(websocket.BinaryMessage, byte2)
			if err != nil {
				log.Println("写入错误:", err)
				return
			}
		case <-interrupt:
			log.Println("接收到中断信号，关闭连接...")
			client.Close()
			//err := client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			//if err != nil {
			//	log.Println("关闭连接错误:", err)
			//	return
			//}
			//select {
			//case <-done:
			//case <-time.After(time.Second):
			//}
			return
		}
	}
}

func StartHttpServer(port string){
	route := gin.Default()

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

	route.Run(":" + port)
}