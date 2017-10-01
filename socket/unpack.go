// Code generated by protoc-gen-main.
// source: regist/msg.go
// DO NOT EDIT!

package socket

import (
	"errors"
	"protoactor/messages"
)

//解包消息
func unpack(id uint32, b []byte) (interface{}, error) {
	switch id {
	case 4010:
		msg := &messages.CCreateRoom{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1036:
		msg := &messages.CUserData{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1004:
		msg := &messages.CWxLogin{}
		err := msg.Unmarshal(b)
		return msg, err
	case 3004:
		msg := &messages.CWxpayQuery{}
		err := msg.Unmarshal(b)
		return msg, err
	case 3006:
		msg := &messages.CBuildAgent{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4000:
		msg := &messages.CEnterRoom{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4004:
		msg := &messages.CReady{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4014:
		msg := &messages.CKick{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4029:
		msg := &messages.CDiscard{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1030:
		msg := &messages.CConfig{}
		err := msg.Unmarshal(b)
		return msg, err
	case 3002:
		msg := &messages.CWxpayOrder{}
		err := msg.Unmarshal(b)
		return msg, err
	case 2000:
		msg := &messages.CChatText{}
		err := msg.Unmarshal(b)
		return msg, err
	case 2004:
		msg := &messages.CChatVoice{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1002:
		msg := &messages.CRegist{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4030:
		msg := &messages.COperate{}
		err := msg.Unmarshal(b)
		return msg, err
	case 1000:
		msg := &messages.CLogin{}
		err := msg.Unmarshal(b)
		return msg, err
	case 3000:
		msg := &messages.CBuy{}
		err := msg.Unmarshal(b)
		return msg, err
	case 4002:
		msg := &messages.CLeave{}
		err := msg.Unmarshal(b)
		return msg, err
	default:
		return nil, errors.New("msg id wrong")
	}
}