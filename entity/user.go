package entity

import "time"

//用户数据
type User struct {
	Userid    string    `xorm:"pk",orm:"",bson:"_id"`                //用户id
	Nickname  string    `xorm:"Varchar(20)",orm:"",bson:"nickname"`  //用户昵称
	Sex       uint32    `xorm:"Int default(0)",orm:"",bson:"sex"`    //用户性别,男1 女0
	Phone     string    `xorm:"Varchar(11)",orm:"",bson:"phone"`     //绑定的手机号码
	Photo     string    `xorm:"Varchar(255)",orm:"",bson:"photo"`    //头像
	Status    uint32    `xorm:"Int default(0)",orm:"",bson:"status"` //正常0  锁定1  黑名单2
	Coin      uint32    `xorm:"Int",orm:"",bson:"coin"`              //金币
	Diamond   uint32    `xorm:"Int",orm:"",bson:"diamond"`           //钻石
	RoomCard  uint32    `xorm:"Int",orm:"",bson:"room_card"`         //房卡
	Vip       uint32    `xorm:"Int",orm:"",bson:"vip"`               //vip
	VipExpire uint32    `xorm:"Int",orm:"",bson:"vip_expire"`        //vip有效期
	Wxuid     string    `xorm:"Varchar(50)",orm:"",bson:"wxuid"`     //微信uid
	Win       uint32    `xorm:"Int",orm:"",bson:"win"`               //胜局数
	Lost      uint32    `xorm:"Int",orm:"",bson:"lost"`              //败局数
	Ping      uint32    `xorm:"Int",orm:"",bson:"ping"`              //平局数
	Piao      uint32    `xorm:"Int",orm:"",bson:"piao"`              //飘赖数
	Robot     bool      `xorm:"Bool",orm:"",bson:"robot"`            //是否是机器人
	Utime     time.Time `xorm:"updated",orm:"",bson:"utime"`         //更新时间
}
