package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	FormID   int64  // 发送者
	TargetID int64  // 接收者
	Type     int    // 发送类型 群聊 私聊 广播
	Media    int    // 消息类型 文字 图片 音频
	Content  string // 消息内容
	Pic      string
	Url      string
	Desc     string
	Amount   string // 其他数字统计

}

func (t *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 影射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

// 需要: 发送者ID, 接收者ID, 消息类型, 发送的内容, 发送类型
func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. 获取参数, 校验token
	query := request.URL.Query()
	idStr := query.Get("userId")
	userID, _ := strconv.ParseInt(idStr, 10, 64)
	// msgType := query.Get("type")
	// targetID := query.Get("targetId")
	// context := query.Get("context")

	isvalida := true // token校验
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2. 获取conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}
	// 3. 用户关系

	// 4. userid  跟 node 绑定
	rwLocker.Lock()
	clientMap[userID] = node
	rwLocker.Unlock()

	// 5. 完成发送逻辑
	go sendProc(node)
	// 6. 完成接收逻辑
	go recvProc(node)
	sendMsg(userID, []byte("欢迎进入聊天室aa"))
}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("[ws] <<<<<< ", string(data))
	}
}

var udpSendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

func init() {
	go updSendProc()
	go updRecvPorc()
}

// 完成UPD数据发送协程
func updSendProc() {
	con, err :=
		net.DialUDP("udp", nil, &net.UDPAddr{
			IP:   net.IPv4(43, 139, 0, 255),
			Port: 3000,
		})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpSendChan:
			_, err := con.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成UPD数据接收协程
func updRecvPorc() {
	con, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer con.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		var buf [512]byte
		n, err := con.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispath(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispath(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch msg.Type {
	case 1: // 群发
		sendMsg(msg.TargetID, data)
		// case 2:
		// 	sendGroupMsg()
		// case 3 :
		// 	sendAllMsg()

	}
}

func sendMsg(userID int64, msg []byte) {
	rwLocker.Lock()
	node, ok := clientMap[userID]
	rwLocker.Unlock()
	if ok {
		node.DataQueue <- msg
	}

}
