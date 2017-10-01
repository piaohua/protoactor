/**********************************************************
 * Author        : piaohua
 * Email         : 814004090@qq.com
 * Last modified : 2017-04-01 21:12:42
 * Filename      : robot.go
 * Description   : 机器人
 * *******************************************************/
package robots

import (
	"bytes"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second // Time allowed to write a message to the peer.
	pongWait       = 60 * time.Second // Time allowed to read the next pong message from the peer.
	pingPeriod     = 9 * time.Second  // Send pings to peer with this period. Must be less than pongWait.
	maxMessageSize = 10240            // Maximum message size allowed from peer.
	waitForLogin   = 5 * time.Second  // 连接建立后5秒内没有收到登陆请求,断开socket
)

type WebsocketConnSet map[*websocket.Conn]struct{}

// 机器人连接数据
type Robot struct {
	conn      *websocket.Conn // websocket连接
	writeChan chan []byte     // 消息写入通道
	maxMsgLen uint32          // 最大消息长度
	closeFlag bool            // 连接状态
	index     int             // 包序
	//游戏数据
	data  *user    //数据
	code  string   //邀请码
	seat  uint32   //位置
	cards []uint32 //手牌
}

// 基本数据
type user struct {
	Userid    string // 用户id
	Nickname  string // 用户昵称
	Sex       uint32 // 用户性别,男1 女0
	Phone     string // 绑定的手机号码
	Photo     string // 头像
	Coin      uint32 // 金币
	Diamond   uint32 // 钻石
	Vip       uint32 // vip
	VipExpire uint32 // vip
	RoomCard  uint32 // 房卡
}

//创建连接
func newRobot(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *Robot {
	return &Robot{
		conn:      conn,
		data:      new(user),
		closeFlag: true,
		writeChan: make(chan []byte, pendingWriteNum),
		maxMsgLen: maxMsgLen,
	}
}

//断开连接
func (ws *Robot) Close() {
	ws.conn.Close()
}

//接收
func (ws *Robot) Router(id uint32, body []byte) {
	msg, err := unpack(id, body)
	if err != nil {
		glog.Errorln("protocol unpack err:", id, err)
		return
	}
	ws.Recv(id, msg)
}

//发送消息
func (ws *Robot) Sender(pb interface{}) {
	code, body, err := packet(pb)
	if err != nil {
		glog.Errorf("packet msg err -> %#v, %v", pb, err)
		return
	}
	msg := Pack(code, body, ws.index)
	if ws.closeFlag {
		if uint32(len(msg)) > ws.maxMsgLen {
			glog.Errorf("msg too long -> %d", len(msg))
			return
		}
		if len(ws.writeChan) == cap(ws.writeChan) {
			glog.Errorf("writeChan channel full -> %d", len(ws.writeChan))
			ws.Close()
			return
		}
		ws.writeChan <- msg
	}
}

func (ws *Robot) readPump() {
	defer func() {
		ws.Close()
		ws.closeFlag = false
		close(ws.writeChan) //TODO 加锁
	}()
	ws.conn.SetReadLimit(maxMessageSize)
	ws.conn.SetReadDeadline(time.Now().Add(pongWait))
	ws.conn.SetPongHandler(func(string) error { ws.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	// 消息缓冲
	msgbuf := bytes.NewBuffer(make([]byte, 0, 10240))
	// 消息长度
	var length int = 0
	// 包序长度
	var index int = 0
	// 协议编号
	var proto uint32 = 0
	for {
		n, message, err := ws.conn.ReadMessage()
		if err != nil {
			glog.Errorf("Read error: %s, %d\n", err, n)
			break
		}
		// 数据添加到消息缓冲
		m, err := msgbuf.Write(message)
		if err != nil {
			glog.Errorf("Buffer write error: %s, %d\n", err, m)
			return
		}
		// 消息分割循环
		for {
			// 消息头
			if length == 0 && msgbuf.Len() >= 9 {
				index = int(msgbuf.Next(1)[0])             //包序
				proto = DecodeUint32(msgbuf.Next(4))       //协议号
				length = int(DecodeUint32(msgbuf.Next(4))) //消息长度
				// 检查超长消息
				if length > 10240 {
					glog.Errorf("Message too length: %d\n", length)
					return
				}
			}
			//fmt.Printf("index: %d, proto: %d, length: %d, len: %d\n", index, proto, length, msgbuf.Len())
			// 消息体
			if length > 0 && msgbuf.Len() >= length { //TODO length 空消息体
				//fmt.Printf("Client messge: %s\n", string(msgbuf.Next(length)))
				//包序验证
				ws.index++
				ws.index = ws.index % 256
				//fmt.Printf("Message index error: %d, %d\n", index, ws.index)
				if ws.index != index {
					fmt.Printf("Message index error: %d, %d\n", index, ws.index)
					glog.Errorf("Message index error: %d, %d\n", index, ws.index)
					//return
				}
				//代理,路由
				//proxyHandler2(proto, msgbuf.Next(length), ws)
				ws.Router(proto, msgbuf.Next(length))
				length = 0
			} else {
				break
			}
		}
	}
}

//消息写入 TODO write Buff
func (ws *Robot) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		ws.Close()
	}()
	for {
		select {
		case message, ok := <-ws.writeChan:
			if !ok {
				ws.write(websocket.CloseMessage, []byte{})
				return
			}
			err := ws.write(websocket.TextMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			err := ws.write(websocket.PingMessage, []byte{})
			if err != nil {
				return
			}
		}
	}
}

//写入
func (ws *Robot) write(mt int, msg []byte) error {
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, msg)
}

//Big Endian
func DecodeUint32(data []byte) uint32 {
	return (uint32(data[0]) << 24) | (uint32(data[1]) << 16) | (uint32(data[2]) << 8) | uint32(data[3])
}

//Big Endian
func EncodeUint32(n uint32) []byte {
	b := make([]byte, 4)
	b[3] = byte(n & 0xFF)
	b[2] = byte((n >> 8) & 0xFF)
	b[1] = byte((n >> 16) & 0xFF)
	b[0] = byte((n >> 24) & 0xFF)
	return b
}

//封包
func Pack(code uint32, msg []byte, index int) []byte {
	buff := make([]byte, 9+len(msg))
	msglen := uint32(len(msg))
	buff[0] = byte(index)
	copy(buff[1:5], EncodeUint32(code))
	copy(buff[5:9], EncodeUint32(msglen))
	copy(buff[9:], msg)
	return buff
}
