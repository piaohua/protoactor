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

func Regist(ctos *messages.CRegist, pid *actor.PID) {
	stoc := &messages.SRegist{Error: 0}
	nickname := ctos.GetNickname()
	phone := ctos.GetPhone()
	password := ctos.GetPassword()
	glog.Infof("nickname -> %s", nickname)
	glog.Infof("phone -> %s", phone)
	glog.Infof("password -> %s", password)

	if nickname == "" || !utils.LegalName(nickname, 7) {
		glog.Errorln("nickname legal")
		stoc.Error = errorcode.NicknameLegal
		pid.Tell(stoc)
		return
	}

	if phone == "" || !utils.PhoneRegexp(phone) {
		glog.Errorln("phone legal")
		stoc.Error = errorcode.AccountLegal
		pid.Tell(stoc)
		return
	}

	if password == "" || len(password) != 32 {
		glog.Errorln("password legal", phone)
		stoc.Error = errorcode.PasswordLegal
		pid.Tell(stoc)
		return
	}

	var salt string = utils.GetSalt()
	password = utils.Md5(password + salt)
	acc := new(entity.Account)
	acc.Name = phone
	acc.Nickname = nickname
	acc.Password = password
	acc.Salt = salt
	affected, err := data.Engine.Insert(acc)
	if err != nil {
		glog.Errorln("account already exist", phone)
		stoc.Error = errorcode.AccountExist
		pid.Tell(stoc)
		return
	}
	if affected == 1 { //注册成功
		has, err := data.Engine.Get(acc)
		userid := utils.String(acc.Id) //字符串
		if err != nil || !has || userid == "" {
			glog.Errorln("account regist failed", phone)
			stoc.Error = errorcode.RegistFail
			pid.Tell(stoc)
			return
		}
		stoc.Userid = userid
		pid.Tell(stoc)
		//日志消息
		pid.Tell(&messages.LogRegist{
			Userid:   userid,
			Name:     phone,
			Nickname: nickname,
			Sender:   process.LogPID,
		})
		glog.Infof("Userid -> %s", userid)
		//TODO 直接登录
	}
}
