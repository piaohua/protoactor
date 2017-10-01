package process

import (
	"protoactor/errorcode"
	"protoactor/messages"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

func (state *PlayerActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *messages.CUserData:
		state.userData(msg, ctx.Self())
	case *messages.CEnterRoom:
		state.enterRoom(msg, ctx.Self())
	case *messages.CCreateRoom:
		state.create(msg, ctx.Self())
	case *messages.EnterRoomResp:
		msg1 := msg.(*messages.EnterRoomResp)
		//glog.Infof("EnterRoomResp -> %#v", msg1)
		if msg1.Result {
			state.DeskPid = msg1.Pid //设置房间进程ID
		} else { //进入失败,不在房间内
			state.DeskPid = nil //断开
		}
	case *messages.CreateRoomResp:
		msg1 := msg.(*messages.CreateRoomResp)
		//glog.Infof("CreateRoomResp -> %#v", msg1)
		if msg1.Result {
			state.DeskPid = msg1.Pid //设置房间进程ID
		}
	case *messages.CReady:
		state.ready(msg, ctx.Self())
	case *messages.CLeave:
		state.leave(msg, ctx.Self())
	case *messages.CDiscard:
		state.discard(msg, ctx.Self())
	case *messages.COperate:
		state.operate(msg, ctx.Self())
	case *messages.LeaveDesk:
		state.DeskPid = nil
		if state.State == LOGOUT {
			ctx.SetReceiveTimeout(3 * time.Minute) //timeout
		} else {
			state.State = LOGINED //登录状态
		}
	case *messages.UpdateUser:
		state.updateUser(msg, ctx.Self())
	case proto.Message:
		state.WsPidTell(msg) //响应消息
	default:
		glog.Errorf("request msg wrong -> %#v", msg)
	}
}

//响应消息
func (state *PlayerActor) WsPidTell(msg interface{}) {
	if state.WsPid == nil {
		glog.Errorf("Respond msg -> %#v", msg)
		return
	}
	state.WsPid.Tell(msg)
}

//进入房间
func (state *PlayerActor) userData(req interface{}, self *actor.PID) {
	msg := req.(*messages.CUserData)
	//TODO
	msg2 := &messages.SUserData{}
	msg2.Data = &messages.UserData{
		Userid: msg.Userid,
	}
	state.WsPidTell(msg2) //响应
}

//进入房间
func (state *PlayerActor) enterRoom(req interface{}, self *actor.PID) {
	msg := req.(*messages.CEnterRoom)
	msg2 := &messages.EnterRoomReq{
		Userid:   state.Userid,
		Nickname: state.Nickname,
		Sex:      state.Sex,
		Photo:    state.Photo,
		Diamond:  state.Diamond,
		Pid:      self, //玩家进程ID
	}
	if state.DeskPid != nil { //已经在房间
		state.DeskPid.Request(msg2, self) //进入房间
	} else {
		msg2.Code = msg.GetInvitecode()
		RoomsPID.Request(msg2, self) //进入已经创建的房间
	}
}

//创建房间
func (state *PlayerActor) create(req interface{}, self *actor.PID) {
	msg := req.(*messages.CCreateRoom)
	if state.DeskPid != nil { //已经在房间
		msg3 := &messages.EnterRoomReq{
			Userid:   state.Userid,
			Nickname: state.Nickname,
			Sex:      state.Sex,
			Photo:    state.Photo,
			Diamond:  state.Diamond,
			Pid:      self, //玩家进程ID
		}
		state.DeskPid.Request(msg3, self) //进入房间
		return
	}
	//获取参数
	var rname string = msg.GetRname()
	var rtype uint32 = msg.GetRtype()
	var ante uint32 = msg.GetAnte()       //2,5,10
	var round uint32 = msg.GetRound()     //4,12,20
	var payment uint32 = msg.GetPayment() //0一家,1=AA支付
	var pao uint32 = msg.GetPao()         //1一赖到底,0干瞪眼
	var count uint32 = msg.GetCount()     //2,4
	msg2 := &messages.CreateRoomReq{
		Userid:   state.Userid,
		Nickname: state.Nickname,
		Sex:      state.Sex,
		Photo:    state.Photo,
		Diamond:  state.Diamond,
		Rtype:    rtype,
		Rname:    rname,
		Ante:     ante,
		Round:    round,
		Payment:  payment,
		Pao:      pao,
		Count:    count,
		Pid:      self, //玩家进程ID
	}
	RoomsPID.Request(msg2, self) //开始创建房间
}

//玩家准备
func (state *PlayerActor) ready(req interface{}, self *actor.PID) {
	msg := req.(*messages.CReady)
	var ready bool = msg.GetReady()
	if state.DeskPid != nil { //已经在房间
		msg2 := &messages.ReadyRoomReq{
			Userid: state.Userid,
			Ready:  ready,
			Pid:    self, //玩家进程ID
		}
		state.DeskPid.Request(msg2, self) //进入房间
		return
	}
	stoc := &messages.SReady{Error: 0}
	stoc.Ready = ready
	stoc.Error = errorcode.NotInRoom
	state.WsPidTell(stoc) //响应
}

//离开房间
func (state *PlayerActor) leave(req interface{}, self *actor.PID) {
	//msg := req.(*messages.CLeave)
	if state.DeskPid != nil { //已经在房间
		msg2 := &messages.LeaveRoomReq{
			Userid: state.Userid,
			Pid:    self, //玩家进程ID
		}
		state.DeskPid.Request(msg2, self) //进入房间
		return
	}
	stoc := &messages.SLeave{Error: 0}
	stoc.Error = errorcode.NotInRoom
	state.WsPidTell(stoc) //响应
}

//打牌
func (state *PlayerActor) discard(req interface{}, self *actor.PID) {
	msg := req.(*messages.CDiscard)
	if state.DeskPid != nil { //已经在房间
		msg2 := &messages.DiscardReq{
			Userid: state.Userid,
			Card:   msg.GetCard(),
			Pid:    self, //玩家进程ID
		}
		state.DeskPid.Tell(msg2) //进入房间
		return
	}
	stoc := &messages.SDiscard{Error: 0}
	stoc.Error = errorcode.NotInRoom
	state.WsPidTell(stoc) //响应
}

//操作
func (state *PlayerActor) operate(req interface{}, self *actor.PID) {
	msg := req.(*messages.COperate)
	if state.DeskPid != nil { //已经在房间
		msg2 := &messages.OperateReq{
			Userid: state.Userid,
			Card:   msg.GetCard(),
			Value:  msg.GetValue(),
			Pid:    self, //玩家进程ID
		}
		state.DeskPid.Tell(msg2) //进入房间
		return
	}
	stoc := &messages.SOperate{Error: 0}
	stoc.Error = errorcode.NotInRoom
	state.WsPidTell(stoc) //响应
}

//更新数据
func (state *PlayerActor) updateUser(req interface{}, self *actor.PID) {
	msg := req.(*messages.UpdateUser)
	state.User.Win += msg.Win
	state.User.Lost += msg.Lost
	state.User.Ping += msg.Ping
	state.User.Piao += msg.Piao
	//state.WsPidTell(msg) //响应
}
