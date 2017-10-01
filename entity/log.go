package entity

import "time"

//注册日志
type LogRegist struct {
	Id       int64
	Userid   string    `xorm:"varchar(20) notnull index"` //账户ID
	Name     string    `xorm:"varchar(25) notnull index"` //账户名称
	Nickname string    `xorm:"varchar(50) notnull"`       //账户名称
	Ip       string    `xorm:"char(15) notnull"`          //注册IP
	DayStamp int64     `xorm:"int(11)" index`             //regist Time Today
	DayDate  int       `xorm:"int(11)" index`             //regist day date
	Ctime    time.Time `xorm:"created"`                   //create Time
}

//登录日志
type LogLogin struct {
	Id         int64
	Userid     string    `xorm:"varchar(20) notnull unique"`      //账户ID
	Event      int       `xorm:"int(3) notnull index default(0)"` //事件：0=登录,1=正常退出,2＝系统关闭时被迫退出,3＝被动退出,4＝其它情况导致的退出
	Ip         string    `xorm:"char(15) notnull"`                //登录IP
	DayStamp   int64     `xorm:"int(11)" index`                   //login Time Today
	LoginTime  time.Time `xorm:"created"`                         //login Time
	LogoutTime time.Time `xorm:"updated"`                         //logout Time
}

//钻石日志
type LogDiamond struct {
	Id     int64
	Userid string    `xorm:"varchar(20) notnull index"` //账户ID
	Type   uint32    `xorm:"int(11) notnull index"`     //类型
	Num    int32     `xorm:"int(11) notnull"`           //数量
	Rest   uint32    `xorm:"int(11) notnull"`           //剩余数量
	Ctime  time.Time `xorm:"created"`                   //create Time
}

//金币日志
type LogCoin struct {
	Id     int64
	Userid string    `xorm:"varchar(20) notnull index"` //账户ID
	Type   uint32    `xorm:"int(11) notnull index"`     //类型
	Num    int32     `xorm:"int(11) notnull"`           //数量
	Rest   uint32    `xorm:"int(11) notnull"`           //剩余数量
	Ctime  time.Time `xorm:"created"`                   //create Time
}
