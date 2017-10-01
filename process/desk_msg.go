package process

import (
	"protoactor/desk"
	"protoactor/errorcode"
	"protoactor/messages"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/glog"
)

func (state *DeskActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *messages.EnterRoomReq:
		state.enter(msg, ctx)
	case *messages.EnterDeskReq:
		state.enterResp(msg, ctx)
	case *messages.ReadyRoomReq:
		state.ready(msg, ctx)
	case *messages.LeaveRoomReq:
		state.leave(msg, ctx)
	case *messages.DiscardReq:
		state.discard(msg, ctx)
	case *messages.OperateReq:
		state.operate(msg, ctx)
	case *messages.TimerDiscardReq:
		state.timerDiscard(msg, ctx)
	case *messages.TimerTurnReq:
		state.timerTurn(msg, ctx)
	case *messages.DeskClose:
		RoomsPID.Tell(msg)                     //房间关闭
		ctx.SetReceiveTimeout(5 * time.Second) //timeout
	default:
		glog.Errorf("request msg wrong -> %#v", msg)
	}
}

//进入房间
func (state *DeskActor) enterResp(req interface{}, ctx actor.Context) {
	msg := req.(*messages.EnterDeskReq)
	//glog.Infof("msg: %#v", msg)
	result, msg1 := desk.EnterMsg(msg.Userid, state.Desk)
	if !result {
		glog.Errorf("enter desk error:%#v", msg)
	}
	msg2 := &messages.EnterRoomResp{
		Result: msg.Result,
		Pid:    ctx.Self(),
	}
	//响应消息给玩家进程
	msg.Pid.Tell(msg2)
	msg.Pid.Tell(msg1)
}

//重复进入房间
func (state *DeskActor) enter(req interface{}, ctx actor.Context) {
	msg := req.(*messages.EnterRoomReq)
	//glog.Infof("msg: %#v", msg)
	result, msg1 := desk.Enter(msg.Userid, state.Desk)
	if result {
		ctx.Respond(&messages.EnterRoomResp{
			Result: true,
			Pid:    ctx.Self(), //桌子进程
		})
		ctx.Respond(msg1)
		return
	}
	stoc := &messages.SEnterRoom{
		Error: errorcode.NotInRoom,
	}
	ctx.Respond(stoc)
}

//准备
func (state *DeskActor) ready(req interface{}, ctx actor.Context) {
	msg := req.(*messages.ReadyRoomReq)
	glog.Infof("msg: %#v", msg)
	result := desk.Ready(msg.Userid, msg.Ready, state.Desk)
	switch result {
	case 1:
		ctx.Respond(&messages.SReady{
			Ready: msg.Ready,
			Error: errorcode.InTheVote,
		})
	case 2:
		ctx.Respond(&messages.SReady{
			Ready: msg.Ready,
			Error: errorcode.NotInRoom,
		})
	}
}

//离开
func (state *DeskActor) leave(req interface{}, ctx actor.Context) {
	msg := req.(*messages.LeaveRoomReq)
	glog.Infof("msg: %#v", msg)
	result := desk.Leave(msg.Userid, state.Desk)
	switch result {
	case 1:
		ctx.Respond(&messages.SLeave{Error: errorcode.NotInRoom})
	case 2:
		ctx.Respond(&messages.SLeave{Error: errorcode.InTheGame})
	}
}

//打牌
func (state *DeskActor) discard(req interface{}, ctx actor.Context) {
	msg := req.(*messages.DiscardReq)
	glog.Infof("msg: %#v", msg)
	result := desk.DiscardL(msg.Userid, msg.Card, false, state.Desk)
	switch result {
	case 1:
		ctx.Respond(&messages.SDiscard{Error: errorcode.NotInRoom})
	case 2:
		ctx.Respond(&messages.SDiscard{Error: errorcode.NotInTheGame})
	case 3:
		ctx.Respond(&messages.SDiscard{Error: errorcode.NotYourTurn})
	}
}

//准备
func (state *DeskActor) operate(req interface{}, ctx actor.Context) {
	msg := req.(*messages.OperateReq)
	glog.Infof("msg: %#v", msg)
	result := desk.OperateL(msg.Card, msg.Value, msg.Userid, state.Desk)
	switch result {
	case 1:
		ctx.Respond(&messages.SOperate{Error: errorcode.NotInRoom})
	case 2:
		ctx.Respond(&messages.SOperate{Error: errorcode.NotInTheGame})
	case 3:
		ctx.Respond(&messages.SOperate{Error: errorcode.NotYourTurn})
	}
}

//超时打牌
func (state *DeskActor) timerDiscard(req interface{}, ctx actor.Context) {
	msg := req.(*messages.TimerDiscardReq)
	glog.Infof("msg: %#v", msg)
	result := desk.DiscardL(msg.Userid, msg.Card, true, state.Desk)
	glog.Infof("result: %#v", result)
}

//超时操作
func (state *DeskActor) timerTurn(req interface{}, ctx actor.Context) {
	msg := req.(*messages.TimerTurnReq)
	glog.Infof("msg: %#v", msg)
	desk.TurnL(state.Desk)
}
