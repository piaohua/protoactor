// Code generated by protoc-gen-main.
// source: regist/msg.go
// DO NOT EDIT!

package robots

import (
	"errors"
	"protoactor/messages"
)

//打包消息
func packet(msg interface{}) (uint32, []byte, error) {
	switch msg.(type) {
	case *messages.CLogin:
		b, err := msg.(*messages.CLogin).Marshal()
		return 1000, b, err
	case *messages.CBuy:
		b, err := msg.(*messages.CBuy).Marshal()
		return 3000, b, err
	case *messages.CLeave:
		b, err := msg.(*messages.CLeave).Marshal()
		return 4002, b, err
	case *messages.CWxLogin:
		b, err := msg.(*messages.CWxLogin).Marshal()
		return 1004, b, err
	case *messages.CWxpayQuery:
		b, err := msg.(*messages.CWxpayQuery).Marshal()
		return 3004, b, err
	case *messages.CBuildAgent:
		b, err := msg.(*messages.CBuildAgent).Marshal()
		return 3006, b, err
	case *messages.CCreateRoom:
		b, err := msg.(*messages.CCreateRoom).Marshal()
		return 4010, b, err
	case *messages.CUserData:
		b, err := msg.(*messages.CUserData).Marshal()
		return 1036, b, err
	case *messages.CConfig:
		b, err := msg.(*messages.CConfig).Marshal()
		return 1030, b, err
	case *messages.CWxpayOrder:
		b, err := msg.(*messages.CWxpayOrder).Marshal()
		return 3002, b, err
	case *messages.CChatText:
		b, err := msg.(*messages.CChatText).Marshal()
		return 2000, b, err
	case *messages.CChatVoice:
		b, err := msg.(*messages.CChatVoice).Marshal()
		return 2004, b, err
	case *messages.CEnterRoom:
		b, err := msg.(*messages.CEnterRoom).Marshal()
		return 4000, b, err
	case *messages.CReady:
		b, err := msg.(*messages.CReady).Marshal()
		return 4004, b, err
	case *messages.CKick:
		b, err := msg.(*messages.CKick).Marshal()
		return 4014, b, err
	case *messages.CDiscard:
		b, err := msg.(*messages.CDiscard).Marshal()
		return 4029, b, err
	case *messages.CRegist:
		b, err := msg.(*messages.CRegist).Marshal()
		return 1002, b, err
	case *messages.COperate:
		b, err := msg.(*messages.COperate).Marshal()
		return 4030, b, err
	default:
		return 0, []byte{}, errors.New("msg wrong")
	}
}
