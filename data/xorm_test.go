package data

import (
	"protoactor/entity"
	"testing"
)

//[3 1 wang3 3000 100 1 1479033336 3]
func Test_Xorm(t *testing.T) {
	acc := new(entity.Account)
	acc.Name = "wang3"   //昵称
	acc.Password = "xxx" //初始化金币
	acc.Salt = "xxx"     //初始化钻石
	affected, err := Engine.Insert(acc)
	t.Log(affected, err)
}

//[3 1 wang3 3000 100 1 1479033336 3]
func TestUser(t *testing.T) {
	u := new(entity.User)
	u.Userid = "1111" //
	has, err := Engine.Get(u)
	t.Log(has, err)
}

func BenchmarkXorm(b *testing.B) {
}
