package socket

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
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

type WSConn struct {
	conn      *websocket.Conn // websocket连接
	writeChan chan []byte     // 消息写入通道
	maxMsgLen uint32          // 最大消息长度
	closeFlag bool            // 连接状态
	index     int             // 包序
	pid       *actor.PID      // ws进程ID
	playerPid *actor.PID      // 角色进程ID
}

//创建连接
func newWSConn(conn *websocket.Conn, pendingWriteNum int, maxMsgLen uint32) *WSConn {
	return &WSConn{
		conn:      conn,
		closeFlag: true,
		writeChan: make(chan []byte, pendingWriteNum),
		maxMsgLen: maxMsgLen,
	}
}

//连接地址
func (ws *WSConn) localAddr() string {
	return ws.conn.LocalAddr().String()
}

func (ws *WSConn) remoteAddr() string {
	return ws.conn.RemoteAddr().String()
}

func (ws *WSConn) GetIPAddr() string {
	return strings.Split(ws.remoteAddr(), ":")[0]
}

//断开连接
func (ws *WSConn) Close() {
	ws.conn.Close()
}

//index(1byte) + proto(4byte) + msgLen(4byte) + msg
//func (ws *WSConn) msgHandler() {
//}

func (ws *WSConn) readPump() {
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
func (ws *WSConn) writePump() {
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
func (ws *WSConn) write(mt int, msg []byte) error {
	ws.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.conn.WriteMessage(mt, msg)
}
