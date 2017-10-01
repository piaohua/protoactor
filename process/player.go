package process

import (
	"fmt"
	"protoactor/data"
	"protoactor/entity"
	"protoactor/messages"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

const (
	LOGINING = 0 //登录中
	LOGINED  = 1 //登录成功
	LOGOUT   = 2 //下线
	INGAME   = 3 //游戏中
)

//test
type testPlayer struct{ Who string }

//PlayerActor
type PlayerActor struct {
	*entity.User            //角色数据
	Code         string     //邀请码
	State        int        //登录状态
	Ctime        int64      //创建时间
	WsPid        *actor.PID //ws进程ID
	DeskPid      *actor.PID //牌桌进程ID
}

func (state *PlayerActor) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *actor.Started:
		fmt.Println("Starting, initialize actor here")
		fmt.Println("Starting, msg -> %v", msg)
	case *actor.Stopping:
		fmt.Println("Stopping, actor is about to shut down")
	case *actor.Stopped:
		fmt.Println("Stopped, actor and its children are stopped")
	case *actor.Restarting:
		fmt.Println("Restarting, actor is about to restart")
	case *actor.ReceiveTimeout:
		//超时处理(登出)
		data.SaveUser(state.User) //数据存储
		ctx.Self().Stop()
	case *testPlayer:
		fmt.Printf("self %s\n", ctx.Self().String())
		fmt.Printf("player %v\n", msg.Who)
		fmt.Printf("Ctime %d\n", state.Ctime)
		//state.Userid = msg.Who
		state.Ctime = time.Now().Unix()
		fmt.Printf("Ctime %d\n", state.Ctime)
	case string:
		fmt.Printf("msg %s\n", msg)
	case *messages.Logout:
		//下线处理
		//日志消息
		LogPID.Tell(&messages.LogLogout{
			Userid: state.User.Userid,
			Event:  1,
		})
		fmt.Printf("logout pid: %s", msg.Sender.String())
		state.WsPid = nil          //断开
		if state.State != INGAME { //不在游戏内3分钟后下线
			ctx.SetReceiveTimeout(3 * time.Minute) //timeout
		}
		state.State = LOGOUT //改变状态
		glog.Infof("Loginout: %s", ctx.Self().String())
		//TODO 下线消息
	case *messages.Login:
		ctx.SetReceiveTimeout(0) //login Successfully, timeout off
		msg1 := &messages.OnlineReq{
			Pid:      ctx.Self(),
			Sender:   msg.Sender,
			Userid:   msg.Userid,
			Phone:    msg.Phone,
			Nickname: msg.Nickname,
		}
		OnlinePID.Request(msg1, ctx.Self())
		glog.Infof("Login: %s", ctx.Self().String())
	case *messages.OnlineResp:
		exist := msg.Result
		//登录处理
		if exist && state.WsPid != nil {
			state.WsPidTell(&messages.RepeatLogin{
				Sender: ctx.Self(),
				Repeat: msg.Sender,
			})
		} else {
			// 初始化数据
			user, err := data.InitUser(msg.Userid, msg.Phone, msg.Nickname)
			if err != nil {
				glog.Errorf("Userid %s InitUser err:%s", state.Userid, err)
				//TODO logout
			}
			state.User = user
			ctx.Self().Tell(&messages.Logined{
				Sender: msg.Sender,
			})
		}
		state.State = LOGINING //登录中
	case *messages.Logined:
		//登录成功
		if state.State == LOGINING {
			state.State = LOGINED //登录成功
		}
		state.WsPid = msg.Sender //替换
		state.WsPidTell(&messages.Logined{
			Sender: ctx.Self(),
		})
		state.WsPidTell(&messages.SLogin{
			Userid: state.User.Userid,
		})
		//日志消息
		state.WsPidTell(&messages.LogLogin{
			Userid: state.User.Userid,
			Sender: LogPID,
		})
		glog.Infof("Logined: %s", state.User.Userid)
	case proto.Message:
		//消息请求
		//glog.Infof("PlayerMsg: %s, %#v", ctx.Self().String(), msg)
		state.Handler(msg, ctx)
	}
}

//初始化
func InitPlayer(userid string) *actor.PID {
	props := actor.FromInstance(&PlayerActor{}) //桌子实例
	pid, err := actor.SpawnNamed(props, userid) //启动一个进程
	if err != nil {
		fmt.Printf("init player err -> %v", err)
		return pid //TODO name exists  如果已经存在?
	}
	return pid
}

//pid.String() == "player"+角色ID
