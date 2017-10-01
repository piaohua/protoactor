package entity

import "time"

//账号
type Account struct {
	Id       int64
	Name     string    `xorm:"varchar(25) notnull unique"` //账户名称
	Nickname string    `xorm:"varchar(50) notnull"`        //账户名称
	Password string    `xorm:"varchar(32) notnull"`        //账户密码
	Salt     string    `xorm:"varchar(6) notnull"`         //盐
	Ctime    time.Time `xorm:"created"`                    //regist Time
}
