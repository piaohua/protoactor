package process

import (
	"protoactor/desk"
	"protoactor/entity"
	"protoactor/errorcode"
	"protoactor/messages"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/glog"
)

func (state *RoomsActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *messages.EnterRoomReq:
		state.enterRoom(msg, ctx)
	case *messages.CreateRoomReq:
		state.createRoom(msg, ctx)
	case *messages.DeskClose:
		state.deskClose(msg, ctx)
	default:
		glog.Errorf("request msg wrong -> %#v", msg)
	}
}

//生成一个牌桌邀请码,全列表中唯一
func (state *RoomsActor) genInviteCode() (s string) {
	for n := 0; n < 100; n++ { //一定次数后放弃尝试
		//生成长度为6的随机邀请码
		s = utils.RandStr(6)
		//是否已经存在
		_, exist := state.RoomsMap.Get(s)
		//如果存在,重复尝试
		if !exist {
			return s
		}
	}
	return ""
}

//关闭桌子
func (state *RoomsActor) deskClose(req interface{}, ctx actor.Context) {
	msg := req.(*messages.DeskClose)
	val, exist := state.RoomsMap.Get(msg.Code)
	if exist {
		d := val.(*entity.DeskData)
		d.Pid.Stop() //关闭
		//TODO 日志记录
	}
	state.RoomsMap.Remove(msg.Code)
}

//进入房间
func (state *RoomsActor) enterRoom(req interface{}, ctx actor.Context) {
	msg := req.(*messages.EnterRoomReq)
	stoc := &messages.SEnterRoom{Error: 0}
	val, exist := state.RoomsMap.Get(msg.Code)
	if !exist {
		stoc.Error = errorcode.RoomNotExist
		ctx.Respond(stoc)
		ctx.Respond(&messages.EnterRoomResp{})
		return
	}
	d := val.(*entity.DeskData)
	if msg.Diamond < d.Cost {
		stoc.Error = errorcode.NotEnoughDiamond
		ctx.Respond(stoc)
		ctx.Respond(&messages.EnterRoomResp{})
		return
	}
	u := &entity.DeskUser{
		Userid:   msg.Userid,
		Nickname: msg.Nickname,
		Sex:      msg.Sex,
		Photo:    msg.Photo,
		Pid:      msg.Pid, //玩家进程
	}
	var ok bool = desk.Add(u, d)
	if !ok {
		stoc.Error = errorcode.CreateRoomFail
		ctx.Respond(stoc)
		ctx.Respond(&messages.EnterRoomResp{})
		return
	}
	//发消息给桌子进程
	d.Pid.Tell(&messages.EnterDeskReq{
		Result: true,
		Userid: msg.Userid,
		Pid:    msg.Pid, //玩家进程
	})
}

//创建房间
func (state *RoomsActor) createRoom(req interface{}, ctx actor.Context) {
	msg := req.(*messages.CreateRoomReq)
	stoc := &messages.SCreateRoom{Error: 0}
	//条件验证
	glog.Infof("createRoom: %d", msg.Rtype)
	if (msg.Payment != 1 && msg.Payment != 0) ||
		(msg.Pao != 1 && msg.Pao != 0) ||
		(msg.Count != 2 && msg.Count != 4) {
		stoc.Error = errorcode.CreateRoomFail
		ctx.Respond(stoc)
		ctx.Respond(&messages.CreateRoomResp{})
		return
	}
	//msg.Pao //0,1
	//msg.Count //2,4
	//msg.Ante //2,5,10
	//msg.Round //4,12,20
	var cost uint32 = 20 //20,40,80
	if msg.Round == 12 {
		cost = 40
	}
	if msg.Round == 20 {
		cost = 80
	}
	if msg.Diamond < cost {
		stoc.Error = errorcode.NotEnoughDiamond
		ctx.Respond(stoc)
		ctx.Respond(&messages.CreateRoomResp{})
		return
	}
	//生成邀请码
	var code string = state.genInviteCode()
	if code == "" {
		stoc.Error = errorcode.CreateRoomFail
		ctx.Respond(stoc)
		ctx.Respond(&messages.CreateRoomResp{})
		return
	}
	var now uint32 = uint32(utils.Timestamp())
	var expire uint32 = now + msg.Round*600
	//初始化
	d := desk.NewDeskData(state.RoomLastID, msg.Round, expire,
		msg.Rtype, msg.Ante, cost, msg.Payment, msg.Pao, now,
		msg.Count, msg.Userid, msg.Rname, code)
	//加入房间
	u := &entity.DeskUser{
		Userid:   msg.Userid,
		Nickname: msg.Nickname,
		Sex:      msg.Sex,
		Photo:    msg.Photo,
		Pid:      msg.Pid,
	}
	var ok bool = desk.Add(u, d)
	if !ok {
		stoc.Error = errorcode.CreateRoomFail
		ctx.Respond(stoc)
		ctx.Respond(&messages.CreateRoomResp{})
		return
	}
	//返回消息
	roomdata := &messages.RoomData{
		Roomid:     state.RoomLastID,
		Rtype:      msg.Rtype,
		Rname:      msg.Rname,
		Userid:     msg.Userid,
		Expire:     expire,
		Round:      msg.Round,
		Invitecode: code,
		Count:      msg.Count,
	}
	stoc.Rdata = roomdata
	//启动一张桌子
	pid := InitDesk(d)
	d.Pid = pid //设置进程ID
	//添加到房间列表
	state.RoomsMap.Put(code, d)
	//房间ID递增
	state.RoomLastID += 1 //房间ID递增
	//打包消息, stoc可以直接分开发
	ctx.Respond(&messages.CreateRoomResp{
		Result: true,
		Pid:    pid,
	})
	ctx.Respond(stoc)
	//召唤机器人
	go callRobot(code)
}

//创建远程连接消息通道
func callRobot(code string) {
	channel := make(chan *messages.RobotMsg, 1)
	message := &messages.RobotMsg{
		Code: code,
		Num:  3,
	}
	remote := actor.NewPID("127.0.0.1:7070", "RobotMsg")
	go func() {
		for msg := range channel {
			remote.Tell(msg)
		}
	}()
	channel <- message
	utils.Sleep(3)
	close(channel)
}
