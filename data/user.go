package data

import (
	"protoactor/entity"

	"github.com/golang/glog"
)

func InitUser(userid, phone, nickname string) (*entity.User, error) {
	user := new(entity.User)
	user.Userid = userid
	user.Phone = phone
	user.Nickname = nickname
	has, err := Engine.Get(user)
	if err != nil {
		glog.Infof("get user: %s, err -> %v", userid, err)
	}
	if err == nil && !has { //记录不存在
		//TODO 初始化值
		user.Coin = 50000
		user.Diamond = 50000
		user.RoomCard = 1
		user.Vip = 1
		affected, err := Engine.Insert(user)
		if err != nil {
			glog.Infof("Insert user: %s, err -> %v", userid, err)
			return nil, err
		}
		if affected == 1 {
			return user, nil
		}
	}
	return user, nil
}

func SaveUser(user *entity.User) error {
	affected, err := Engine.Id(user.Userid).AllCols().Omit("Userid").Update(user)
	if err != nil {
		glog.Infof("Update user: %s, err -> %v", user.Userid, err)
		return err
	}
	if affected == 1 {
	}
	return nil
}
