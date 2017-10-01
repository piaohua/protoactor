package socket

import (
	"github.com/golang/glog"
)

//1000 登录协议
//1002 注册协议
//1004 微信登录协议
func (ws *WSConn) Router(id uint32, body []byte) {
	msg, err := unpack(id, body)
	if err != nil {
		glog.Errorln("protocol unpack err:", id, err)
		return
	}
	switch id {
	case 1000, 1002, 1004:
		ws.pid.Tell(msg)
	default:
		ws.playerPid.Tell(msg)
	}
}

//发送消息
func (ws *WSConn) Send(pb interface{}) {
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

//解包
//func Unpack() {
//}
