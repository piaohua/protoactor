package process

import (
	"protoactor/data"
	"protoactor/entity"
	"protoactor/messages"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/glog"
)

func (state *LoggerActor) Handler(msg interface{}, ctx actor.Context) {
	switch msg.(type) {
	case *messages.LogLogout:
		state.logout(msg)
	case *messages.LogLogin:
		state.login(msg)
	case *messages.LogRegist:
		state.regist(msg)
	default:
		_, err := data.Engine.Insert(msg)
		if err != nil {
			glog.Errorf("logger msg wrong -> %#v", msg)
		}
	}
}

//登出日志
func (state *LoggerActor) logout(req interface{}) {
	msg := req.(*messages.LogLogout)
	log := &entity.LogLogin{
		Event: int(msg.Event),
	}
	_, err := data.Engine.Where("Userid=?", msg.Userid).Cols("Event").Update(log)
	if err != nil {
		glog.Errorln("logger err:", err)
	}
}

//登录日志
func (state *LoggerActor) login(req interface{}) {
	msg := req.(*messages.LogLogin)
	log := &entity.LogLogin{
		Userid:   msg.Userid,
		Event:    int(msg.Event),
		Ip:       msg.Ip,
		DayStamp: utils.TimestampToday(),
	}
	_, err := data.Engine.Insert(log)
	if err != nil {
		glog.Errorln("logger err:", err)
	}
}

//注册日志
func (state *LoggerActor) regist(req interface{}) {
	msg := req.(*messages.LogRegist)
	log := &entity.LogRegist{
		Userid:   msg.Userid,
		Name:     msg.Name,
		Nickname: msg.Nickname,
		Ip:       msg.Ip,
		DayStamp: utils.TimestampToday(),
		DayDate:  utils.DayDate(),
	}
	_, err := data.Engine.Insert(log)
	if err != nil {
		glog.Errorln("logger err:", err)
	}
}
