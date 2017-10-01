package login

import (
	"protoactor/data"
	"protoactor/entity"
	"protoactor/errorcode"
	"protoactor/messages"
	"protoactor/process"
	"utils"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/golang/glog"
)

func Login(ctos *messages.CLogin, pid *actor.PID) {
	stoc := &messages.SLogin{Error: 0}
	phone := ctos.GetPhone()
	password := ctos.GetPassword()
	glog.Infof("phone -> %s", phone)
	glog.Infof("password -> %s", password)

	if phone == "" || !utils.PhoneRegexp(phone) {
		glog.Errorln("phone legal")
		stoc.Error = errorcode.AccountLegal
		pid.Tell(stoc)
		return
	}

	if password == "" || len(password) != 32 {
		glog.Errorln("password legal")
		stoc.Error = errorcode.PasswordLegal
		pid.Tell(stoc)
		return
	}

	acc := new(entity.Account)
	acc.Name = phone
	has, err := data.Engine.Get(acc)
	userid := utils.String(acc.Id)
	if err != nil || !has || userid == "" {
		glog.Errorln("account login failed", phone)
		stoc.Error = errorcode.LoginFail
		pid.Tell(stoc)
		return
	}
	//验证密码
	if acc.Password != utils.Md5(password+acc.Salt) {
		glog.Errorln("account regist failed", phone)
		stoc.Error = errorcode.PasswordError
		pid.Tell(stoc)
		return
	}

	//处理重复登录,关闭掉老的连接后再接管pid
	playerPid := process.InitPlayer("player" + userid)
	playerPid.Tell(&messages.Login{
		Userid:   userid,
		Phone:    phone,
		Nickname: acc.Nickname,
		Sender:   pid,
	}) //登录消息
	glog.Infof("playerPid %v", playerPid.String())
}
